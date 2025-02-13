package ft

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
	actual := <-outChan

	expected := "\nkey1: value1\nkey2: value2\n\n"
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestVerifyInput(t *testing.T) {
	tests := []struct {
		expected bool
		data     CmdArgs
		wantErr  bool
	}{
		{
			true,
			CmdArgs{
				cmd:      []string{},
				allPaths: map[string]string{},
				file:     nil,
				rdr:      strings.NewReader("y"),
			},
			false,
		},
		{
			false,
			CmdArgs{
				cmd:      []string{},
				allPaths: map[string]string{},
				file:     nil,
				rdr:      strings.NewReader("n"),
			},
			false,
		},
		{
			false,
			CmdArgs{
				cmd:      []string{},
				allPaths: map[string]string{},
				file:     nil,
				rdr:      strings.NewReader("x\n"),
			},
			true,
		},
	}

	for _, tt := range tests {
		var res string
		fmt.Fscan(tt.data.rdr, &res)
		// TODO
		// - add force true to tests
		actual, err := verifyInput(res, false)
		if !tt.wantErr && err != nil {
			fmt.Println("Error:", err)
			t.Fatal(err)
		}
		if actual != tt.expected {
			t.Errorf("-> Expected --> %v\n____________\nGot --> %v", tt.expected, actual)
		}

	}

}

func TestPipeArgs(t *testing.T) {
	tests := []struct {
		name        string
		initialArgs []string
		input       string
		expected    []string
	}{
		{
			name:        "1. Pipe key name to args.",
			initialArgs: []string{"ft"},
			input:       "keyname",
			expected:    []string{"ft", "keyname"},
		},
		{
			name:        "2. Pipe in multiple args.",
			initialArgs: []string{"ft"},
			input:       "-set keyname",
			expected:    []string{"ft", "-set", "keyname"},
		},
	}
	for _, tt := range tests {
		t.Log(tt.name)

		r, w, err := os.Pipe()
		if err != nil {
			fmt.Println(tt.name)
			t.Error("Error establishing pipe.")
		}
		stdin := os.Stdin
		os.Stdin = r

		defer func() {
			os.Stdin = stdin
			r.Close()
			w.Close()
		}()

		_, err = io.WriteString(w, tt.input)
		if err != nil {
			t.Errorf("WriteString failed: %v", err)
		}
		w.Close()

		args := tt.initialArgs
		err = PipeArgs(&args)
		if err != nil {
			fmt.Println(tt.name)
			t.Error("PipeArgs produced an error: ", err)
		}

		equal := true
		if len(args) == len(tt.expected) {
			for i, _ := range args {
				if args[i] == tt.expected[i] {
					continue
				}
				equal = false
				break
			}
			if equal {
				continue
			}
		}
		t.Errorf("Expected: %q length %v, got: %q length %v", tt.expected, len(tt.expected), args, len(args))
	}

}
