package main

import (
	"bytes"
	"fmt"
	"io"
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
	r, w, err := os.Pipe()
    if err != nil {
        t.Fatalf("Failed to create pipe for testing printMap()")
    }
	os.Stdout = w

	printMap(hashmap)
   
    // Use go routine so printing doesn't block program
    outChan := make(chan string)
    go func() {
        var buf bytes.Buffer
        io.Copy(&buf, r)
        outChan <- buf.String()
            
    }()

	w.Close()
	os.Stdout = old
    actual := <- outChan

	expected := "\nkey1: value1\nkey2: value2\n\n"
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	} else {
        fmt.Println("PrintMap: Success")
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
	} else {
        fmt.Println("ensureData: Success")
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

    // Use go routine so printing doesn't block program
    outChan := make(chan string)
    go func() {
        var buf bytes.Buffer
        io.Copy(&buf, r)
        outChan <- buf.String()
            
    }()

	w.Close()
	os.Stdout = old
    actual := <- outChan

	expected := "/mnt/c/Users/Test/Documents\n"
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

	data := cmdArgs{
		cmd:      []string{"set", "testKey"},
		allPaths: make(map[string]string),
		file:     tmpfile,
	}

    stdout := os.Stdout
    _,w,_ := os.Pipe()
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

	result := readMap(file)
	if result["testKey"] != expected {
		t.Errorf("Expected file to have key 'testKey' with value %s, got %s", expected, result["testKey"])
        fail = true
    }

    if !fail {
        fmt.Println("setDirectoryVar: Success")
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

    // Use go routine so printing doesn't block program
    outChan := make(chan string)
    go func() {
        var buf bytes.Buffer
        io.Copy(&buf, r)
        outChan <- buf.String()
            
    }()

	w.Close()
	os.Stdout = old
    actual := <- outChan

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

	data := cmdArgs{
		cmd: []string{"rm", "key1"},
		allPaths: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		file: tmpfile,
	}
    stdout := os.Stdout
    _,w,_ := os.Pipe()
    os.Stdout = w
	removeKey(data)
    os.Stdout = stdout
    // fmt.Println("")
	if _, ok := data.allPaths["key1"]; ok {
		t.Errorf("Expected key 'key1' to be removed")
        fail = true
	}

	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	result := readMap(file)
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

	data := cmdArgs{
		cmd: []string{"rn", "key1", "newKey"},
		allPaths: map[string]string{
			"key1": "value1",
		},
		file: tmpfile,
	}

    stdout := os.Stdout
    _,w,_ := os.Pipe()
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

	result := readMap(file)
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
	data := cmdArgs{}

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
    actual := <- outChan

    expected := "\nhelp: you are here :) - Usage: ft help\nls: display all current key value pairs - Usage: ft ls\nrm: deletes provided key - Usage: ft rm [key]\nrn: renames key to new key - Usage: ft rn [key] [new key]\nset: set current directory path to provided key - Usage: ft set [key]\nto: change directory to provided key's path - Usage: ft to [key]\n\n"
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	} else {
        fmt.Println("showHelp: Success")
    }
}

func TestMainFunc(t *testing.T) {
	// Create a temp dir to run tests in
	tmpdir, err := os.MkdirTemp("", "testdata")
    if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
    tmpdir = strings.Trim(tmpdir, " ")
	defer os.RemoveAll(tmpdir)

	// Create a dummy binary file
	exePath := filepath.Join(tmpdir, "fastTravel.bin")
	file, err := os.Create(exePath)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

    // Move working dir to tempdir 
	oldGetwd, err := os.Getwd()
    if err != nil {
        t.Fatalf("Unable to establish working directory")
    }
	defer func() { 
        err = os.Chdir(oldGetwd)
        if err != nil {
            t.Fatalf("Failed to navigate to original working directory")
        }
    }()
    err = os.Chdir(tmpdir)
    if err != nil {
        t.Fatalf("Failed to navigate to temp directory")
    }
	
    tests := []struct {
		args    []string
		expected    string
	}{
		{
            []string{"ft", "help"}, 
            "\nhelp: you are here :) - Usage: ft help\nls: display all current key value pairs - Usage: ft ls\nrm: deletes provided key - Usage: ft rm [key]\nrn: renames key to new key - Usage: ft rn [key] [new key]\nset: set current directory path to provided key - Usage: ft set [key]\nto: change directory to provided key's path - Usage: ft to [key]\n\n", 
        },
		{
            []string{"ft", "set", "key"}, 
            "Added destination key", 
        },
		{
            []string{"ft", "to", "key"}, 
            fmt.Sprintf("%v\n",sanitizeDir(tmpdir)), 
        },
		{
            []string{"ft", "ls"},
            fmt.Sprintf("\nkey: %v\n\n", tmpdir), 
        },
		{
	        []string{"ft", "rn", "key", "key2"},
            "key renamed to key2", 
        },
		{
	        []string{"ft", "rm", "key2"},
            "Removed 'key2' destination", 
        },
	}
	

    for _, tt := range tests {
        
        // Create pipe for capturing output
        stdout := os.Stdout
        r, w, err := os.Pipe()
        if err != nil {
            t.Fatalf("Failed to create pipe for testing main()")
        }
        os.Stdout = w
        
        os.Args = tt.args
        main()

        // Move output to a buffer so it can be passed to a queue
        // Go routine and chan so data transfer doesn't block 
        outChan := make(chan string)
        go func() {
            var buf bytes.Buffer
            io.Copy(&buf, r)
            outChan <- buf.String()
                
        }()
        
        // Capture queue contents for comparison against expected output
        w.Close()
        os.Stdout = stdout
        actual := <- outChan
        

        if actual != tt.expected {
            t.Errorf("-> ARGS: %v\nExpected --> %v\n____________\nGot --> %v", tt.args, tt.expected, actual)
        } 
        
        
    }

    if !t.Failed() {
        fmt.Println("main: Success")
    }



}






