package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrintMap(t *testing.T) {
	hashmap := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	// Redirect stdout to capture the output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printMap(hashmap)

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	r.WriteTo(&buf)

	expected := "\nkey1: value1\nkey2: value2\n\n"
	if buf.String() != expected {
		t.Errorf("Expected %s, got %s", expected, buf.String())
	}
}

func TestPassCmd(t *testing.T) {
	tests := []struct {
		args    []string
		want    []string
		wantErr bool
	}{
		{[]string{"ft", "ls"}, []string{"ls"}, false},
		{[]string{"ft", "help"}, []string{"help"}, false},
		{[]string{"ft", "rn", "key", "newKey"}, []string{"rn", "key", "newKey"}, false},
		{[]string{"ft", "set", "key"}, []string{"set", "key"}, false},
		{[]string{"ft", "invalid"}, nil, true},
		{[]string{"ft", "rn"}, nil, true},
		{[]string{"ft", "set"}, nil, true},
	}

	for _, tt := range tests {
		got, err := passCmd(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("passCmd() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !tt.wantErr && !equalSlices(got, tt.want) {
			t.Errorf("passCmd() = %v, want %v", got, tt.want)
		}
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestEnsureData(t *testing.T) {
	// Create a temporary directory for testing
	tmpdir, err := os.MkdirTemp("", "testdata")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	testFilePath := filepath.Join(tmpdir, "test.bin")
	file := ensureData(testFilePath)
	defer file.Close()

	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Fatalf("File was not created")
	}
}

func TestChangeDirectory(t *testing.T) {
	data := cmdArgs{
		cmd: []string{"to", "testKey"},
		allPaths: map[string]string{
			"testKey": "C:\\Users\\Test\\Documents",
		},
		file: nil,
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	changeDirectory(data)

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	r.WriteTo(&buf)

	expected := "/mnt/c/Users/Test/Documents\n"
	if buf.String() != expected {
		t.Errorf("Expected %s, got %s", expected, buf.String())
	}
}

func TestSetDirectoryVar(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testdata.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	data := cmdArgs{
		cmd:      []string{"set", "testKey"},
		allPaths: make(map[string]string),
		file:     tmpfile,
	}

	setDirectoryVar(data)

	expected, _ := os.Getwd()
	if data.allPaths["testKey"] != expected {
		t.Errorf("Expected key 'testKey' to have value %s, got %s", expected, data.allPaths["testKey"])
	}

	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	result := readMap(file)
	if result["testKey"] != expected {
		t.Errorf("Expected file to have key 'testKey' with value %s, got %s", expected, result["testKey"])
	}
}

func TestDisplayAllPaths(t *testing.T) {
	data := cmdArgs{
		allPaths: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	displayAllPaths(data)

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	r.WriteTo(&buf)

	expected := "\nkey1: value1\nkey2: value2\n\n"
	if buf.String() != expected {
		t.Errorf("Expected %s, got %s", expected, buf.String())
	}
}

func TestRemoveKey(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testdata.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	data := cmdArgs{
		cmd: []string{"rm", "key1"},
		allPaths: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		file: tmpfile,
	}

	removeKey(data)

	if _, ok := data.allPaths["key1"]; ok {
		t.Errorf("Expected key 'key1' to be removed")
	}

	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	result := readMap(file)
	if _, ok := result["key1"]; ok {
		t.Errorf("Expected file to not have key 'key1'")
	}
}

func TestRenameKey(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testdata.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	data := cmdArgs{
		cmd: []string{"rn", "key1", "newKey"},
		allPaths: map[string]string{
			"key1": "value1",
		},
		file: tmpfile,
	}

	renameKey(data)

	if _, ok := data.allPaths["key1"]; ok {
		t.Errorf("Expected key 'key1' to be renamed")
	}
	if data.allPaths["newKey"] != "value1" {
		t.Errorf("Expected key 'newKey' to have value 'value1', got %s", data.allPaths["newKey"])
	}

	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	result := readMap(file)
	if _, ok := result["key1"]; ok {
		t.Errorf("Expected file to not have key 'key1'")
	}
	if result["newKey"] != "value1" {
		t.Errorf("Expected file to have key 'newKey' with value 'value1', got %s", result["newKey"])
	}
}

func TestShowHelp(t *testing.T) {
	data := cmdArgs{}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	showHelp(data)

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	r.WriteTo(&buf)

	expected := "\nhelp: you are here :)\nls: display all current key value pairs - Usage: ft ls\nrm: deletes provided key - Usage: ft rm [key]\nrn: renames key to new key - Usage: ft rn [key] [new key]\nset: set current directory path to provided key - Usage: ft set [key]\nto: change directory to provided key's path - Usage: ft to [key]\n\n"
	if buf.String() != expected {
		t.Errorf("Expected %s, got %s", expected, buf.String())
	}
}

func TestMainFunc(t *testing.T) {
	// Create a temporary executable path
	tmpdir, err := os.MkdirTemp("", "testdata")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	// Create a dummy binary file
	exePath := filepath.Join(tmpdir, "fastTravel.bin")
	file, err := os.Create(exePath)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	// Set the executable path for testing
	os.Args = []string{"ft", "help"}

	oldGetwd := os.Getwd
	defer func() { os.Getwd = oldGetwd }()
	os.Getwd = func() (string, error) {
		return tmpdir, nil
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	var buf strings.Builder
	r.WriteTo(&buf)

	expected := "\nhelp: you are here :)\nls: display all current key value pairs - Usage: ft ls\nrm: deletes provided key - Usage: ft rm [key]\nrn: renames key to new key - Usage: ft rn [key] [new key]\nset: set current directory path to provided key - Usage: ft set [key]\nto: change directory to provided key's path - Usage: ft to [key]\n\n"
	if buf.String() != expected {
		t.Errorf("Expected %s, got %s", expected, buf.String())
	}
}
