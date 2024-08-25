package ft

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// test helpers
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

// tests
func TestPassCmd(t *testing.T) {
	tests := []struct {
		args    []string
		want    []string
		wantErr bool
	}{
		{[]string{"ft", "mypath/dir"}, []string{"_", "mypath/dir"}, false},
		{[]string{"ft", "-ls"}, []string{"-ls"}, false},
		{[]string{"ft", "-help"}, []string{"-help"}, false},
		{[]string{"ft", "-rn", "key", "newKey"}, []string{"-rn", "key", "newKey"}, false},
		{[]string{"ft", "-set", "key"}, []string{"-set", "key"}, false},
		{[]string{"ft", "-invalid"}, nil, true},
		{[]string{"ft", "-rn"}, nil, true},
		{[]string{"ft", "-set"}, nil, true},
	}

	for _, tt := range tests {
		got, err := PassCmd(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("passCmd err %v, want err: %v", err, tt.wantErr)
			return
		}
		if !tt.wantErr && !equalSlices(got, tt.want) {
			t.Errorf("passCmd Args: %v\nexpected: %v\ngot:%v\n_________\n", tt.args, tt.want, got)
			return
		}

	}

	fmt.Println("passCmd: Success")

}

func TestChangeDirectory(t *testing.T) {
	data := NewCmdArgs(
		[]string{"_", "testKey"},
		map[string]string{
			"testKey": "\\usr\\Test\\Documents",
		},
		nil,
		nil,
	)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	changeDirectory(data)

	// Use go routine so printing doesn't block program
	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outChan <- buf.String()

	}()

	w.Close()
	os.Stdout = old
	actual := <-outChan

	expected := "\\usr\\Test\\Documents\n"
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	} else {
		fmt.Println("ChangeDirectory: Success")
	}
}

func TestSetDirectoryVar(t *testing.T) {
	fail := false
	tmpfile, err := os.CreateTemp("", "testdata.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	data := NewCmdArgs(
		[]string{"set", "testKey"},
		make(map[string]string),
		tmpfile,
		nil,
	)

	stdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	setDirectoryVar(data)
	os.Stdout = stdout

	expected, _ := os.Getwd()
	if data.allPaths["testKey"] != expected {
		t.Errorf("Expected key 'testKey' to have value %s, got %s", expected, data.allPaths["testKey"])
		fail = true
	}

	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	result := ReadMap(file)
	if result["testKey"] != expected {
		t.Errorf("Expected file to have key 'testKey' with value %s, got %s", expected, result["testKey"])
		fail = true
	}

	if !fail {
		fmt.Println("setDirectoryVar: Success")
	}

}

func TestDisplayAllPaths(t *testing.T) {
	data := NewCmdArgs(
		[]string{"ls"},
		map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		nil,
		nil,
	)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	displayAllPaths(data)

	// Use go routine so printing doesn't block program
	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outChan <- buf.String()

	}()

	w.Close()
	os.Stdout = old
	actual := <-outChan

	expected := "\nkey1: value1\nkey2: value2\n\n"
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	} else {
		fmt.Println("DisplayAllPaths: Success")
	}
}

func TestRemoveKey(t *testing.T) {
	fail := false
	tmpfile, err := os.CreateTemp("", "testdata.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	input := "y"
	data := &CmdArgs{
		[]string{"rm", "key1"},
		map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		tmpfile,
		strings.NewReader(input),
	}
	stdout := os.Stdout
	_, w, err := os.Pipe()
	if err != nil {
		t.Errorf("Error establishing Pipe: %v", err)
	}

	os.Stdout = w
	removeKey(data)
	os.Stdout = stdout
	if _, ok := data.allPaths["key1"]; ok {
		t.Errorf("Expected key 'key1' to be removed")
		fail = true
	}

	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	result := ReadMap(file)
	if _, ok := result["key1"]; ok {
		t.Errorf("Expected file to not have key 'key1'")
		fail = true
	}

	if !fail {
		fmt.Println("removeKey: Success")
	}
}

func TestRenameKey(t *testing.T) {
	fail := false
	tmpfile, err := os.CreateTemp("", "testdata.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	input := "y"
	data := NewCmdArgs(
		[]string{"rn", "key1", "newKey"},
		map[string]string{
			"key1": "value1",
		},
		tmpfile,
		strings.NewReader(input),
	)

	stdout := os.Stdout
	_, w, err := os.Pipe()
	if err != nil {
		t.Errorf("Error establishing Pipe: %v", err)
	}

	os.Stdout = w
	renameKey(data)
	os.Stdout = stdout

	if _, ok := data.allPaths["key1"]; ok {
		t.Errorf("Expected key 'key1' to be renamed")
		fail = true
	}
	if data.allPaths["newKey"] != "value1" {
		t.Errorf("Expected key 'newKey' to have value 'value1', got %s", data.allPaths["newKey"])
		fail = true
	}

	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	result := ReadMap(file)
	if _, ok := result["key1"]; ok {
		t.Errorf("Expected file to not have key 'key1'")
		fail = true
	}
	if result["newKey"] != "value1" {
		t.Errorf("Expected file to have key 'newKey' with value 'value1', got %s", result["newKey"])
		fail = true
	}

	if !fail {
		fmt.Println("renameKey: Success")
	}
}

func TestShowHelp(t *testing.T) {
	data := NewCmdArgs([]string{"help"}, map[string]string{}, nil, nil)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	showHelp(data)

	// Use go routine so printing doesn't block program
	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outChan <- buf.String()

	}()

	w.Close()
	os.Stdout = old
	actual := <-outChan

	// really just need confirmation something printed out
	if !(len(actual) > 1) {
		t.Errorf("Expected a string of len > 1, got %s", actual)
	} else {
		fmt.Println("showHelp: Success")
	}
}
