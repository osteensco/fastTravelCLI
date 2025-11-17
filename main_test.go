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

	// Force set test dir
	forcedir, err := os.MkdirTemp(tmpdir, "force")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
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
		hook       func()
		pipedInput string
		expected   string
		wantErr    bool
	}{
		{
			name:     "1. Check help command.",
			args:     []string{"ft", "-help"},
			expected: ft.CreateHelpOutput(),
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
			expected: "\nfastTravelCLI returned an error: " + fmt.Sprintf(ft.UnrecognizedKeyMsg, "badkey") + "\n",
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
		{
			name: "9. Check set command with multiple args piped in using various seperators.",
			args: []string{"ft", "-set"},
			hook: func() {
				err := os.Mkdir(fmt.Sprintf("%v/pipedtest", tmpdir), 0777)
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				err = os.Mkdir(fmt.Sprintf("%v/pipedtest/one", tmpdir), 0777)
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				err = os.Mkdir(fmt.Sprintf("%v/pipedtest/two", tmpdir), 0777)
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				err = os.Mkdir(fmt.Sprintf("%v/pipedtest/three", tmpdir), 0777)
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
			},
			pipedInput: fmt.Sprintf("pipekey1=%v/pipedtest/one\npipekey2=%v/pipedtest/two pipekey3='%v/pipedtest/three'", tmpdir, tmpdir, tmpdir),
			expected: fmt.Sprintf(
				"Added destination 'pipekey1': '%v/pipedtest/one'. \nAdded destination 'pipekey2': '%v/pipedtest/two'. \nAdded destination 'pipekey3': '%v/pipedtest/three'. \n",
				tmpdir,
				tmpdir,
				tmpdir,
			),
			wantErr: false,
		},
		{
			name:     "10. Check force set command.",
			args:     []string{"ft", "-set", "-y", fmt.Sprintf("key=%v", forcedir)},
			expected: fmt.Sprintf(ft.AddKeyMsg, "key", forcedir),
			wantErr:  false,
		},
		{
			name:     "11. Check detailed help -key.",
			args:     []string{"ft", "-h", "somekey"},
			expected: ft.DisplayDetailedHelp("_") + "\n",
			wantErr:  false,
		},
		{
			name:     "12. Check detailed help -set.",
			args:     []string{"ft", "-h", "-set"},
			expected: ft.DisplayDetailedHelp("-set") + "\n",
			wantErr:  false,
		},
		{
			name:     "13. Check detailed help -ls.",
			args:     []string{"ft", "-h", "-ls"},
			expected: ft.DisplayDetailedHelp("-ls") + "\n",
			wantErr:  false,
		},
		{
			name:     "14. Check detailed help -rm.",
			args:     []string{"ft", "-help", "-rm"},
			expected: ft.DisplayDetailedHelp("-rm") + "\n",
			wantErr:  false,
		},
		{
			name:     "15. Check detailed help -rn.",
			args:     []string{"ft", "-help", "-rn"},
			expected: ft.DisplayDetailedHelp("-rn") + "\n",
			wantErr:  false,
		},
		{
			name:     "16. Check detailed help -edit.",
			args:     []string{"ft", "-help", "-edit"},
			expected: ft.DisplayDetailedHelp("-edit") + "\n",
			wantErr:  false,
		},
		{
			name:     "17. Check detailed help -].",
			args:     []string{"ft", "-]", "-h"},
			expected: ft.DisplayDetailedHelp("-]") + "\n",
			wantErr:  false,
		},
		{
			name:     "18. Check detailed help -[.",
			args:     []string{"ft", "-[", "-h"},
			expected: ft.DisplayDetailedHelp("-[") + "\n",
			wantErr:  false,
		},
		{
			name:     "19. Check detailed help -hist.",
			args:     []string{"ft", "-hist", "-h"},
			expected: ft.DisplayDetailedHelp("-hist") + "\n",
			wantErr:  false,
		},
		{
			name:     "20. Check detailed help -version.",
			args:     []string{"ft", "-version", "-help"},
			expected: ft.DisplayDetailedHelp("-version") + "\n",
			wantErr:  false,
		},
		{
			name:     "21. Check detailed help -v.",
			args:     []string{"ft", "-v", "-help"},
			expected: ft.DisplayDetailedHelp("-v") + "\n",
			wantErr:  false,
		},
		{
			name:     "22. Check detailed help -is.",
			args:     []string{"ft", "-is", "-help"},
			expected: ft.DisplayDetailedHelp("-is") + "\n",
			wantErr:  false,
		},
		{
			name:     "23. Check detailed help -update.",
			args:     []string{"ft", "-update", "-help"},
			expected: ft.DisplayDetailedHelp("-update") + "\n",
			wantErr:  false,
		},
		{
			name:     "24. Check detailed help -u.",
			args:     []string{"ft", "-u", "-help"},
			expected: ft.DisplayDetailedHelp("-u") + "\n",
			wantErr:  false,
		},
		{
			name:     "25. Check detailed help --.",
			args:     []string{"ft", "-", "-help"},
			expected: ft.DisplayDetailedHelp("--") + "\n",
			wantErr:  false,
		},
		{
			name:     "26. Check detailed help -..",
			args:     []string{"ft", "..", "-help"},
			expected: ft.DisplayDetailedHelp("-..") + "\n",
			wantErr:  false,
		},
		{
			name:       "27. Check cd command with key piped in.", //simulate echo "key" | ft
			args:       []string{"ft"},
			expected:   fmt.Sprintf("%v\n", forcedir),
			wantErr:    false,
			pipedInput: "key",
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

		if tt.hook != nil {
			tt.hook()
		}

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
