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
	} else {
		fmt.Println("PrintMap: Success")
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
		actual, err := verifyInput(res)
		if !tt.wantErr && err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if actual != tt.expected {
			t.Errorf("-> Expected --> %v\n____________\nGot --> %v", tt.expected, actual)
		}

	}

}
