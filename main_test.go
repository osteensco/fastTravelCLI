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
		name     string
		args     []string
		expected string
	}{
		{
			name: "1. Check help command.",
			args: []string{"ft", "-help"},
			expected: fmt.Sprintf(
				"\n-help: %s\n-ls: %s\n-rm: %s\n-rn: %s\n-set: %s\n[: %s\n]: %s\nkey: %s\n\n",
				ft.CmdDesc["-help"],
				ft.CmdDesc["-ls"],
				ft.CmdDesc["-rm"],
				ft.CmdDesc["-rn"],
				ft.CmdDesc["-set"],
				ft.CmdDesc["["],
				ft.CmdDesc["]"],
				ft.CmdDesc["key"],
			),
		},
		{
			name:     "2. Check set command.",
			args:     []string{"ft", "-set", "key"},
			expected: "Added destination key",
		},
		{
			name:     "3. Check cd command.",
			args:     []string{"ft", "key"},
			expected: fmt.Sprintf("%v\n", tmpdir),
		},
		{
			name:     "4. Check ls command.",
			args:     []string{"ft", "-ls"},
			expected: fmt.Sprintf("\nkey: %v\n\n", tmpdir),
		},
		{
			name:     "5. Check navigate stack.",
			args:     []string{"ft", "["},
			expected: "[",
		},
		// {
		// 	        []string{"ft", "rn", "key", "key2"},
		//             "key renamed to key2",
		//         },
		// {
		// 	        []string{"ft", "rm", "key2"},
		//             "Removed 'key2' destination",
		//         },
	}

	for _, tt := range tests {

		// Create pipe for capturing output
		stdout := os.Stdout
		stderr := os.Stderr
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create pipe for testing main()")
		}
		os.Stdout = w
		os.Stderr = w
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
		os.Stderr = stderr
		actual := <-outChan

		if actual != tt.expected {
			fmt.Println(tt.name)
			t.Errorf("-> ARGS: %v\nExpected -> %v\n____________\nGot -> %v", tt.args, tt.expected, actual)
		}

	}

}
