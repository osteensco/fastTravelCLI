package ft

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
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
		name    string
		args    []string
		want    []string
		wantErr bool
	}{
		{"1. Pass in a path.", []string{"ft", "mypath/dir"}, []string{"_", "mypath/dir"}, false},
		{"2. Pass in -ls.", []string{"ft", "-ls"}, []string{"-ls"}, false},
		{"3. Pass in -help.", []string{"ft", "-help"}, []string{"-help"}, false},
		{"4. Pass in -rn.", []string{"ft", "-rn", "key", "newKey"}, []string{"-rn", "key", "newKey"}, false},
		{"5. Pass in -set.", []string{"ft", "-set", "key"}, []string{"-set", "key"}, false},
		{"6. Pass in invalid command.", []string{"ft", "-invalid"}, nil, true},
		{"7. Pass in not enough arguments for -rn.", []string{"ft", "-rn"}, nil, true},
		{"8. Pass in not enough arguments for -set.", []string{"ft", "-set"}, nil, true},
		{"9. Pass in stack navigation.", []string{"ft", ">"}, nil, false},
	}

	for _, tt := range tests {
		got, err := PassCmd(tt.args)
		fmt.Println(tt.name)
		if (err != nil) != tt.wantErr {
			t.Errorf("passCmd err %v, want err: %v", err, tt.wantErr)
			continue
		}
		if !tt.wantErr && !equalSlices(got, tt.want) {
			t.Errorf("passCmd Args: %v\nexpected: %v\ngot:%v\n_________\n", tt.args, tt.want, got)
			continue
		}

		fmt.Println("Success")
	}

}

func TestChangeDirectory(t *testing.T) {

	tmpdir, err := os.MkdirTemp("", "testing")
	if err != nil {
		t.Fatal(err)
	}
	tmpdir2 := tmpdir + "/subdir"
	err = os.Mkdir(tmpdir2, fs.ModeDir)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(tmpdir)

	tests := []struct {
		name     string
		command  []string
		expected string
		allPaths map[string]string
		file     *os.File
		rdr      io.Reader
		wanterr  bool
	}{
		{
			name:     "1. Valid key provided, standalone.",
			command:  []string{"_", "testKey"},
			expected: tmpdir,
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
			file:    nil,
			rdr:     nil,
			wanterr: false,
		},
		{
			name:     "2. Valid key provided, evaluate path.",
			command:  []string{"_", "testKey/subdir"},
			expected: tmpdir2,
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
			file:    nil,
			rdr:     nil,
			wanterr: false,
		},
		{
			name:     "3. Invalid key provided.",
			command:  []string{"_", "testKye"},
			expected: "",
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
			file:    nil,
			rdr:     nil,
			wanterr: true,
		},
		{
			name:     "4. Invalid key provided, evaluate path.",
			command:  []string{"_", "testKye/subdir"},
			expected: "",
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
			file:    nil,
			rdr:     nil,
			wanterr: true,
		},
	}

	for _, tt := range tests {

		fmt.Println(tt.name)
		data := NewCmdArgs(
			tt.command,
			tt.allPaths,
			tt.file,
			tt.rdr,
		)

		stdout := os.Stdout
		stderr := os.Stderr
		r, w, err := os.Pipe()
		if err != nil {
			t.Error("Error establishing pipe.")
		}
		os.Stdout = w
		os.Stderr = w
		err = changeDirectory(data)
		if tt.wanterr {
			if err == nil {
				t.Error(fmt.Sprintf("%s wanted error, go none.", tt.name))
			}
		} else if err != nil {
			t.Error(err)
		}
		// Use go routine so printing doesn't block program
		outChan := make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outChan <- buf.String()

		}()

		w.Close()
		os.Stdout = stdout
		os.Stderr = stderr
		actual := <-outChan
		actual = strings.Trim(actual, "\n")

		fmt.Println(tt.name)
		if actual != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, actual)
		} else {
			fmt.Printf("%s == %s", actual, tt.expected)
			fmt.Println("    --Success")
		}
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

func TestNavStack(t *testing.T) {

	tests := []struct {
		name     string
		command  []string
		expected string
		wanterr  bool
	}{
		{
			name:     "1. Move down stack.",
			command:  []string{"<"},
			expected: "<",
			wanterr:  false,
		},
		{
			name:     "2. Move up stack.",
			command:  []string{">"},
			expected: ">",
			wanterr:  false,
		},
		{
			name:     "3. Not a stack navigation.",
			command:  []string{"_"},
			expected: "",
			wanterr:  true,
		},
	}

	for _, tt := range tests {

		data := NewCmdArgs(tt.command, map[string]string{}, nil, nil)

		stdout := os.Stdout
		stderr := os.Stderr
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}
		os.Stdout = w
		os.Stderr = w

		err = navStack(data)

		// Use go routine so printing doesn't block program
		outChan := make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outChan <- buf.String()

		}()

		w.Close()
		os.Stdout = stdout
		os.Stderr = stderr
		actual := <-outChan

		fmt.Println(tt.name)
		if tt.wanterr {
			if err == nil {
				t.Error("Expected error, got nil")
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}
		}
		if actual != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, actual)
		} else if !tt.wanterr && err == nil {
			fmt.Println("Success")
		}
	}

}
