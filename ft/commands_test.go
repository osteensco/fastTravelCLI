package ft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
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
		{"9. Pass in stack navigation.", []string{"ft", "]"}, []string{"-]"}, false},
	}

	for _, tt := range tests {
		got, err := PassCmd(tt.args)
		if (err != nil) != tt.wantErr {
			fmt.Println(tt.name)
			t.Errorf("passCmd err %v, want err: %v", err, tt.wantErr)
		}
		if !tt.wantErr && !equalSlices(got, tt.want) {
			fmt.Println(tt.name)
			t.Errorf("passCmd Args: %v\nexpected: %v\ngot:%v\n_________\n", tt.args, tt.want, got)
		}

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
	}{
		{
			name:     "1. Valid key provided, standalone.",
			command:  []string{"_", "testKey"},
			expected: fmt.Sprintln(tmpdir),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
			file: nil,
			rdr:  nil,
		},
		{
			name:     "2. Valid key provided, evaluate path.",
			command:  []string{"_", "testKey/subdir"},
			expected: fmt.Sprintln(tmpdir2),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
			file: nil,
			rdr:  nil,
		},
		{
			name:     "3. Invalid key provided.",
			command:  []string{"_", "testKye"},
			expected: fmt.Sprintf(UnrecognizedKeyMsg, "testKye", "testKye", "testKye"),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
			file: nil,
			rdr:  nil,
		},
		{
			name:     "4. Invalid key provided, evaluate path.",
			command:  []string{"_", "testKye/subdir"},
			expected: fmt.Sprintf(InvalidDirectoryMsg, "testKye/subdir", "testKye/subdir"),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
			file: nil,
			rdr:  nil,
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		data := NewCmdArgs(
			tmpdir,
			tt.command,
			tt.allPaths,
			tt.file,
			tt.rdr,
		)

		stdout := os.Stdout
		stderr := os.Stderr
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Println(tt.name)
			t.Error("Error establishing pipe.")
		}
		os.Stdout = w
		os.Stderr = w
		err = changeDirectory(data)
		if err != nil {
			fmt.Println(tt.name)
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

		if actual != tt.expected {
			t.Errorf("Expected: %q, got: %q", tt.expected, actual)
		}
	}
}

func TestShowDirectoryVar(t *testing.T) {

	CWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		command  []string
		paths    map[string]string
		expected string
		wanterr  bool
	}{
		{
			name:     "1. Check the directory var using -is",
			command:  []string{"-is"},
			paths:    map[string]string{"testKey": CWD},
			expected: fmt.Sprintf(IsKeyMsg, CWD, "testKey"),
			wanterr:  false,
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		data := NewCmdArgs(
			CWD,
			tt.command,
			tt.paths,
			nil,
			nil,
		)

		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		os.Stdout = w

		err = showDirectoryVar(data)
		if err != nil {
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
		os.Stdout = old
		actual := <-outChan

		if actual != tt.expected {
			t.Errorf("Expected message: %q but got message: %q", tt.expected, actual)
		}
	}
}

func TestSetDirectoryVar(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testdata.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	workdir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name     string
		command  []string
		expected string
		wanterr  bool
	}{
		{
			name:     "1. Set key that doesn't exist.",
			command:  []string{"-set", "testKey"},
			expected: workdir,
			wanterr:  false,
		},
	}

	for _, tt := range tests {
		data := NewCmdArgs(
			workdir,
			tt.command,
			make(map[string]string),
			tmpfile,
			nil,
		)

		stdout := os.Stdout
		_, w, _ := os.Pipe()
		os.Stdout = w
		err := setDirectoryVar(data)
		if err != nil {
			t.Error(err)
		}

		os.Stdout = stdout

		if data.allPaths["testKey"] != tt.expected {
			t.Errorf("Expected key 'testKey' to have value %s, got %s", tt.expected, data.allPaths["testKey"])
		}

		file, err := os.Open(tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to open temp file: %v", err)
		}
		defer file.Close()

		result, err := ReadMap(file)
		if err != nil {
			t.Error(err)
		}
		if result["testKey"] != tt.expected {
			t.Errorf("Expected file to have key 'testKey' with value %s, got %s", tt.expected, result["testKey"])
		}
	}
}

func TestDisplayAllPaths(t *testing.T) {
	data := NewCmdArgs(
		"",
		[]string{"-ls"},
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

	err := displayAllPaths(data)
	if err != nil {
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
	os.Stdout = old
	actual := <-outChan

	expected := "\nkey1: value1\nkey2: value2\n\n"
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestRemoveKey(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testdata.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	tests := []struct {
		name     string
		command  []string
		allPaths map[string]string
		input    string
		wanterr  bool
	}{
		{
			name:    "1. Remove key that exists, confirm removal.",
			command: []string{"-rm", "key1"},
			allPaths: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			input:   "y",
			wanterr: false,
		},
	}

	for _, tt := range tests {
		data := NewCmdArgs(
			"",
			tt.command,
			tt.allPaths,
			tmpfile,
			strings.NewReader(tt.input),
		)
		stdout := os.Stdout
		_, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Error establishing Pipe: %v", err)
		}

		os.Stdout = w
		err = removeKey(data)
		if err != nil {
			t.Error(err)
		}
		os.Stdout = stdout
		if _, ok := data.allPaths["key1"]; ok {
			fmt.Println(tt.name)
			t.Errorf("Expected key 'key1' to be removed")

		}

		file, err := os.Open(tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to open temp file: %v", err)
		}
		defer file.Close()

		result, err := ReadMap(file)
		if err != nil {
			t.Error(err)
		}

		if _, ok := result["key1"]; ok {
			fmt.Println(tt.name)
			t.Errorf("Expected file to not have key 'key1'")

		}
	}
}

func TestRenameKey(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testdata.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	tests := []struct {
		name     string
		command  []string
		allPaths map[string]string
		input    string
		wanterr  bool
	}{
		{
			name:    "1. Rename key that exists, confirm rename.",
			command: []string{"-rn", "key1", "newKey"},
			allPaths: map[string]string{
				"key1": "value1",
			},
			input:   "y",
			wanterr: false,
		},
	}

	for _, tt := range tests {
		data := NewCmdArgs(
			"",
			tt.command,
			tt.allPaths,
			tmpfile,
			strings.NewReader(tt.input),
		)

		stdout := os.Stdout
		_, w, err := os.Pipe()
		if err != nil {
			t.Errorf("Error establishing Pipe: %v", err)
		}

		os.Stdout = w
		err = renameKey(data)
		if err != nil {
			t.Error(err)
		}

		os.Stdout = stdout

		if _, ok := data.allPaths["key1"]; ok {
			fmt.Println(tt.name)
			t.Errorf("Expected key 'key1' to be renamed")

		}
		if data.allPaths["newKey"] != "value1" {
			t.Fatalf("Expected key 'newKey' to have value 'value1', got %s", data.allPaths["newKey"])
		}

		file, err := os.Open(tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to open temp file: %v", err)
		}
		defer file.Close()

		result, err := ReadMap(file)
		if err != nil {
			t.Error(err)
		}
		if _, ok := result["key1"]; ok {
			fmt.Println(tt.name)
			t.Errorf("Expected file to not have key 'key1'")

		}
		if result["newKey"] != "value1" {
			fmt.Println(tt.name)
			t.Errorf("Expected file to have key 'newKey' with value 'value1', got %s", result["newKey"])

		}
	}
}

func TestShowVersion(t *testing.T) {
	data := NewCmdArgs("", []string{"-version"}, map[string]string{}, nil, nil)

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w

	err = showVersion(data)
	if err != nil {
		t.Error(err)
	}

	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outChan <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	actual := <-outChan
	expected := fmt.Sprintf("%sversion:\t %s\n", Logo, Version)

	if !(actual == expected) {
		t.Errorf("Expected %q got %q", expected, actual)
	}
}

func TestShowHelp(t *testing.T) {
	data := NewCmdArgs("", []string{"-help"}, map[string]string{}, nil, nil)

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w

	err = showHelp(data)
	if err != nil {
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
	os.Stdout = old
	actual := <-outChan

	// really just need confirmation something printed out
	if !(len(actual) > 1) {
		t.Errorf("Expected a string of len > 1, got %s", actual)
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
			command:  []string{"-["},
			expected: "[\n",
			wanterr:  false,
		},
		{
			name:     "2. Move up stack.",
			command:  []string{"-]"},
			expected: "]\n",
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

		data := NewCmdArgs("", tt.command, map[string]string{}, nil, nil)

		stdout := os.Stdout
		stderr := os.Stderr
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}
		os.Stdout = w
		os.Stderr = w

		err = passToShell(data)

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

		if tt.wanterr {
			if err == nil {
				fmt.Println(tt.name)
				t.Error("Expected error, got nil")
			}
		} else {
			if err != nil {
				fmt.Println(tt.name)
				t.Errorf("Unexpected error: %s", err)
			}
		}
		if actual != tt.expected {
			fmt.Println(tt.name)
			t.Errorf("Expected %s, got %s", tt.expected, actual)
		}
	}
}

func TestUpdateFT(t *testing.T) {
	// Mock server to simulate GitHub API responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag := strings.TrimPrefix(r.URL.Path, "/repos/osteensco/fastTravelCLI/")
		if tag == "releases/latest" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"tag_name": "v.0.2.0"})
		} else if tag == "tags/v.0.1.3" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"tag_name": "v.0.1.3"})

		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	}))
	defer server.Close()

	cwd, err := os.Getwd()
	testdir, err := os.MkdirTemp(cwd, "tmp")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(testdir)

	var defaultVersion string

	// Override endpoint URL to use the test server
	EndpointGH = server.URL + "/repos/osteensco/fastTravelCLI/tags/%s"
	EndpointLatestGH = server.URL + "/repos/osteensco/fastTravelCLI/releases/latest"
	// Override Git related constants for testing
	GitCloneCMD = []string{"echo", "'", "mocking", "version", "", "git", "clone", "dev", "'"}
	GitCloneDir = strings.TrimSuffix(cwd, "/ft")

	tests := []struct {
		name       string
		cmd        []string
		wantError  bool
		setVersion func()
	}{
		{
			name:      "1. Update with no version specified.",
			cmd:       []string{"update"},
			wantError: false,
		},
		{
			name:      "2. Update with specific version.",
			cmd:       []string{"update", "v.0.1.3"},
			wantError: false,
		},
		{
			name:       "3. Already up-to-date version.",
			cmd:        []string{"update", "latest"},
			wantError:  false,
			setVersion: func() { Version = "v.0.2.0" },
		},
		{
			name:      "4. Nonexistent version.",
			cmd:       []string{"update", "nonexistentversionnumber"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)

			if tt.setVersion != nil {
				defaultVersion = Version
				tt.setVersion()
			}

			stdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Error(err)
			}
			os.Stdout = w

			data := &CmdArgs{wkDir: testdir, cmd: tt.cmd}

			// Run the function
			err = updateFT(data)

			// Use go routine so printing doesn't block program
			outChan := make(chan string)
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outChan <- buf.String()
			}()

			w.Close()
			os.Stdout = stdout
			actual := <-outChan

			t.Log(actual)
			// Verify errors for cases where an error is expected
			if (err != nil) != tt.wantError {
				errwkdir, _ := os.Getwd()
				t.Errorf("updateFT() error inside of directory %v -> %v, wantError %v", errwkdir, err, tt.wantError)
			}
			if tt.setVersion != nil {
				Version = defaultVersion
			}

		})
	}
}
