package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/osteensco/fastTravelCLI/ft"
)












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
            fmt.Sprintf("%v\n",ft.SanitizeDir(tmpdir)), 
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






