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
	tmpdir, err = filepath.EvalSymlinks(tmpdir)
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

	// CDPATH setup
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to retrieve home directory: %v", err)
	}
	cdpathdir, err := os.MkdirTemp(homeDir, "cdpathdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer os.RemoveAll(cdpathdir)

	cdpathtest, err := os.MkdirTemp(cdpathdir, "cdpathtest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	cdpathtestkey := filepath.Base(cdpathtest)

	cdpath := os.Getenv("CDPATH")
	err = os.Setenv("CDPATH", fmt.Sprint(cdpathdir, ":", cdpath))

	// tests
	tests := []struct {
		name       string
		args       []string
		pipedInput string
		expected   string
		wantErr    bool
	}{
		{
			name: "1. Check help command.",
			args: []string{"ft", "-help"},
			expected: fmt.Sprintf(
				"\n-help: %s\n-hist: %s\n-is: %s\n-ls: %s\n-rm: %s\n-rn: %s\n-set: %s\n-update: %s\n-version: %s\n[: %s\n]: %s\nkey: %s\n\n",
				ft.CmdDesc["-help"],
				ft.CmdDesc["-hist"],
				ft.CmdDesc["-is"],
				ft.CmdDesc["-ls"],
				ft.CmdDesc["-rm"],
				ft.CmdDesc["-rn"],
				ft.CmdDesc["-set"],
				ft.CmdDesc["-update"],
				ft.CmdDesc["-version"],
				ft.CmdDesc["["],
				ft.CmdDesc["]"],
				ft.CmdDesc["key"],
			),
		},
		{
			name:     "2. Check set command.",
			args:     []string{"ft", "-set", "key"},
			expected: fmt.Sprintf(ft.AddKeyMsg, "key", tmpdir),
			wantErr:  false,
		},
		{
			name:     "3. Check cd command.",
			args:     []string{"ft", "key"},
			expected: fmt.Sprintf("%v\n", tmpdir),
			wantErr:  false,
		},
		{
			name:     "4. Check cd command with bad key.",
			args:     []string{"ft", "badkey"},
			expected: fmt.Sprintf(ft.UnrecognizedKeyMsg, "badkey"),
			wantErr:  false,
		},
		{
			name:     "5. Check ls command.",
			args:     []string{"ft", "-ls"},
			expected: fmt.Sprintf("\nkey: %v\n\n", tmpdir),
			wantErr:  false,
		},
		{
			name:     "6. Check navigate stack.",
			args:     []string{"ft", "["},
			expected: "[\n",
			wantErr:  false,
		},
		{
			name:     "7. Check cd command with CDPATH.",
			args:     []string{"ft", cdpathtestkey},
			expected: fmt.Sprintf("%v\n", cdpathtest),
			wantErr:  false,
		},
		{
			name:       "8. Check command args piped to ft.",
			args:       []string{"ft"},
			pipedInput: "key",
			expected:   fmt.Sprintf("%v\n", tmpdir),
			wantErr:    false,
		},

		// TODO
		{
			name:       "9. Check set command with multiple args piped in.",
			args:       []string{"ft", "-set"},
			pipedInput: fmt.Sprintf("pipekey1=%v/pipedtest/one\npipekey2=%vpipedtest/two pipekey3='%vpipe test/three'", tmpdir, tmpdir, tmpdir),
			expected:   fmt.Sprintf("%v\n", tmpdir),
			wantErr:    false,
		},
	}

	for _, tt := range tests {

		// Create pipe for capturing output
		stdout := os.Stdout
		stderr := os.Stderr
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create pipe for testing main()")
		}
		er, ew, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create pipe for error testing main()")
		}

		os.Stdout = w
		os.Stderr = ew
		os.Args = tt.args

		// Move output to a buffer so it can be passed to a queue
		// Go routine and chan so data transfer doesn't block
		outChan := make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outChan <- buf.String()

		}()

		errOutChan := make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, er)
			errOutChan <- buf.String()

		}()

		// Mock Stdin
		stdinReader, stdinWriter, err := os.Pipe()
		if err != nil {
			fmt.Println(tt.name)
			t.Error("Error establishing pipe.")
		}
		stdin := os.Stdin
		os.Stdin = stdinReader

		defer func() {
			os.Stdin = stdin
			stdinReader.Close()
			stdinWriter.Close()
		}()

		_, err = io.WriteString(stdinWriter, tt.pipedInput)
		if err != nil {
			t.Errorf("WriteString failed: %v", err)
		}
		stdinWriter.Close()

		main()

		// Capture queue contents for comparison against expected output
		w.Close()
		if err != nil {
			t.Fatal("Failed to close Stdout pipe.")
		}
		ew.Close()
		if err != nil {
			t.Fatal("Failed to close Stderr pipe.")
		}
		os.Stdout = stdout
		os.Stderr = stderr
		actual := <-outChan
		errActual := <-errOutChan

		if tt.wantErr && errActual == "" {
			t.Log(tt.name)
			t.Errorf("-> ARGS: %v\nExpected Error\nGot -> %v", tt.args, actual)
		} else if actual != tt.expected {
			t.Log(tt.name)
			t.Errorf("-> ARGS: %v\nExpected -> %q\nGot -> %q", tt.args, tt.expected, actual)
		}

	}

}
