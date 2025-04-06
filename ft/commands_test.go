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
	"path/filepath"
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

func equalCmd(got, expected *Cmd) bool {
	if len(got.Args) != len(expected.Args) {
		return false
	}
	for i, v := range expected.Args {
		if v != got.Args[i] {
			return false
		}
	}
	if got.Cmd != expected.Cmd {
		return false
	}
	if got.Flags.Y != expected.Flags.Y {
		return false
	}

	return true
}

// tests
func TestPassCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *Cmd
		wantErr bool
	}{
		{"1. Pass in a path.", []string{"ft", "mypath/dir"}, &Cmd{Cmd: "_", Args: []string{"mypath/dir"}}, false},
		{"2. Pass in -ls.", []string{"ft", "-ls"}, &Cmd{Cmd: "-ls"}, false},
		{"3. Pass in -help.", []string{"ft", "-help"}, &Cmd{Cmd: "-help"}, false},
		{"4. Pass in -rn.", []string{"ft", "-rn", "key", "newKey"}, &Cmd{Cmd: "-rn", Args: []string{"key", "newKey"}}, false},
		{"5. pass in -set.", []string{"ft", "-set", "key"}, &Cmd{Cmd: "-set", Args: []string{"key"}}, false},
		{"6. pass in invalid command.", []string{"ft", "-invalid"}, nil, true},
		{"7. pass in not enough arguments for -rn.", []string{"ft", "-rn"}, nil, true},
		{"8. pass in not enough arguments for -set.", []string{"ft", "-set"}, nil, true},
		{"9. pass in stack navigation.", []string{"ft", "]"}, &Cmd{Cmd: "-]"}, false},
	}

	for _, tt := range tests {
		got, err := PassCmd(tt.args)
		if (err != nil) != tt.wantErr {
			t.Log(tt.name)
			t.Errorf("passCmd err %v, want err: %v", err, tt.wantErr)
		}
		if !tt.wantErr && !equalCmd(tt.want, got) {
			t.Log(tt.name)
			t.Errorf("passCmd Args: %v\nexpected: %v\ngot: %v\n_________\n", tt.args, tt.want, got)
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
		command  Cmd
		expected string
		allPaths map[string]string
		file     *os.File
		rdr      io.Reader
	}{
		{
			name:     "1. Valid key provided, standalone.",
			command:  Cmd{Cmd: "_", Args: []string{"testKey"}},
			expected: fmt.Sprintln(tmpdir),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
			file: nil,
			rdr:  nil,
		},
		{
			name:     "2. Valid key provided, evaluate path.",
			command:  Cmd{Cmd: "_", Args: []string{"testKey/subdir"}},
			expected: fmt.Sprintln(tmpdir2),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
			file: nil,
			rdr:  nil,
		},
		{
			name:     "3. Invalid key provided.",
			command:  Cmd{Cmd: "_", Args: []string{"testKye"}},
			expected: fmt.Sprintf(UnrecognizedKeyMsg, "testKye"),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
			file: nil,
			rdr:  nil,
		},
		{
			name:     "4. Invalid key provided, evaluate path.",
			command:  Cmd{Cmd: "_", Args: []string{"testKye/subdir"}},
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
		data := NewCmdAPI(
			tmpdir,
			&tt.command,
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
		command  Cmd
		paths    map[string]string
		expected string
		wanterr  bool
	}{
		{
			name:     "1. Check the directory var using -is",
			command:  Cmd{Cmd: "-is"},
			paths:    map[string]string{"testKey": CWD},
			expected: fmt.Sprintf(IsKeyMsg, CWD, "testKey"),
			wanterr:  false,
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		data := NewCmdAPI(
			CWD,
			&tt.command,
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
	// tmpfile for temporary data persistence
	tmpfile, err := os.CreateTemp("", "testdata.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// tmpdir for directories to test with
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
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	workdir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		command  *Cmd
		key      string
		_map     map[string]string
		input    string
		expected string
	}{
		{
			name:     "1. Set key that doesn't exist.",
			command:  &Cmd{Cmd: "-set", Args: []string{"testKey1"}},
			key:      "testKey1",
			_map:     map[string]string{},
			expected: workdir,
		},
		{
			name:     "2. Set key that already exist, don't overwrite.",
			command:  &Cmd{Cmd: "-set", Args: []string{fmt.Sprintf("testKey2=%v", tmpdir)}},
			key:      "testKey2",
			_map:     map[string]string{"testKey2": workdir},
			input:    "n",
			expected: workdir,
		},
		{
			name:     "3. Set key that already exist, overwrite.",
			command:  &Cmd{Cmd: "-set", Args: []string{fmt.Sprintf("testKey3=%v", tmpdir)}},
			key:      "testKey3",
			_map:     map[string]string{"testKey3": workdir},
			input:    "y",
			expected: tmpdir,
		},
		{
			name:     "4. Attempt to set key for a path that is already saved to a key, don't overwrite.",
			command:  &Cmd{Cmd: "-set", Args: []string{"newTestKey1"}},
			key:      "newTestKey1",
			_map:     map[string]string{"testKey4": workdir},
			input:    "n",
			expected: "",
		},
		{
			name:     "5. Attempt to set key for a path that is already saved to a key, overwrite.",
			command:  &Cmd{Cmd: "-set", Args: []string{"newTestKey2"}},
			key:      "newTestKey2",
			_map:     map[string]string{"testKey5": workdir},
			input:    "y",
			expected: workdir,
		},
		{
			name:     "6. Set key for a path with a space in it.",
			command:  &Cmd{Cmd: "-set", Args: []string{fmt.Sprintf("testKey6=%v/some dir", tmpdir)}},
			key:      "testKey6",
			_map:     map[string]string{},
			expected: fmt.Sprintf("%v/some dir", tmpdir),
		},
		{
			name:     "7. Force set key for a path that is already saved to a key.",
			command:  &Cmd{Cmd: "-set", Args: []string{"newTestKey7"}, Flags: CmdFlags{Y: true}},
			key:      "newTestKey7",
			_map:     map[string]string{"testKey7": workdir},
			expected: workdir,
		},
		{
			name:     "8. Force set key that already exists.",
			command:  &Cmd{Cmd: "-set", Args: []string{fmt.Sprintf("testKey8=%v", tmpdir)}, Flags: CmdFlags{Y: true}},
			key:      "testKey8",
			_map:     map[string]string{"testKey8": workdir},
			expected: tmpdir,
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		pathMap := make(map[string]string)
		if tt._map != nil {
			pathMap = tt._map
		}
		data := NewCmdAPI(
			workdir,
			tt.command,
			pathMap,
			tmpfile,
			nil,
		)

		if len(tt._map) > 0 {
			dataUpdate(data.allPaths, tmpfile)
		}

		stdin := os.Stdin
		if tt.input != "" {
			r, w, _ := os.Pipe()
			w.WriteString(tt.input)
			data.rdr = r
			w.Close()
		}

		stdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := setDirectoryVar(data)
		if err != nil {
			t.Error(err)
		}
		w.Close()

		os.Stdin = stdin
		os.Stdout = stdout

		var output bytes.Buffer
		_, err = io.Copy(&output, r)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(output.String())

		if data.allPaths[tt.key] != tt.expected {
			t.Errorf("Expected key 'testKey' to have value %s, got %s", tt.expected, data.allPaths[tt.key])
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
		if result[tt.key] != tt.expected {
			t.Errorf("Expected file to have key 'testKey' with value %s, got %s", tt.expected, result[tt.key])
		}
	}
}

func TestDisplayAllPaths(t *testing.T) {
	data := NewCmdAPI(
		"",
		&Cmd{Cmd: "-ls"},
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
		command  Cmd
		allPaths map[string]string
		input    string
		wanterr  bool
	}{
		{
			name:    "1. Remove key that exists, confirm removal.",
			command: Cmd{Cmd: "-rm", Args: []string{"key1"}},
			allPaths: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			input:   "y",
			wanterr: false,
		},
	}

	for _, tt := range tests {
		data := NewCmdAPI(
			"",
			&tt.command,
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
		command  Cmd
		allPaths map[string]string
		input    string
		wanterr  bool
	}{
		{
			name:    "1. Rename key that exists, confirm rename.",
			command: Cmd{Cmd: "-rn", Args: []string{"key1", "newKey"}},
			allPaths: map[string]string{
				"key1": "value1",
			},
			input:   "y",
			wanterr: false,
		},
	}

	for _, tt := range tests {
		data := NewCmdAPI(
			"",
			&tt.command,
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
	data := NewCmdAPI("", &Cmd{Cmd: "-version"}, map[string]string{}, nil, nil)

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
	data := NewCmdAPI("", &Cmd{Cmd: "-help"}, map[string]string{}, nil, nil)

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
		command  *Cmd
		expected string
		wanterr  bool
	}{
		{
			name:     "1. Move down stack.",
			command:  &Cmd{Cmd: "-["},
			expected: "[\n",
			wanterr:  false,
		},
		{
			name:     "2. Move up stack.",
			command:  &Cmd{Cmd: "-]"},
			expected: "]\n",
			wanterr:  false,
		},
		{
			name:     "3. Not a stack navigation.",
			command:  &Cmd{Cmd: "_"},
			expected: "",
			wanterr:  true,
		},
	}

	for _, tt := range tests {

		data := NewCmdAPI("", tt.command, map[string]string{}, nil, nil)

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
	UPDATEMOCK = true

	tests := []struct {
		name       string
		args       []string
		wantError  bool
		setVersion func()
	}{
		{
			name:      "1. Update with no version specified.",
			args:      []string{},
			wantError: false,
		},
		{
			name:      "2. Update with specific version.",
			args:      []string{"update", "v.0.1.3"},
			wantError: false,
		},
		{
			name:       "3. Already up-to-date version.",
			args:       []string{"update", "latest"},
			wantError:  false,
			setVersion: func() { Version = "v.0.2.0" },
		},
		{
			name:      "4. Nonexistent version.",
			args:      []string{"update", "nonexistentversionnumber"},
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

			data := &CmdAPI{wkDir: testdir, cmd: &Cmd{Cmd: "-update", Args: tt.args}}

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

func TestEditPath(t *testing.T) {
	// tmpfile for temporary data persistence
	tmpfile, err := os.CreateTemp("", "testdata.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// tmpdir for directories to test with
	tmpdir_relative, err := os.MkdirTemp(".", "testdata")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	tmpdir_relative, err = filepath.EvalSymlinks(tmpdir_relative)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	tmpdir_relative = strings.Trim(tmpdir_relative, " ")
	tmpdir, err := filepath.Abs(tmpdir_relative)
	if err != nil {
		t.Fatalf("Failed to get tmpdir relative path: %v", err)
	}

	defer os.RemoveAll(tmpdir)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	workdir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	workdirParent := filepath.Dir(workdir)

	tests := []struct {
		name     string
		command  *Cmd
		_map     map[string]string
		expected map[string]string
	}{
		{
			name:     "1. Update a single key using full path.",
			command:  &Cmd{Cmd: "-edit", Args: []string{workdir, "somethingelse"}},
			_map:     map[string]string{"testKey1": workdir},
			expected: map[string]string{"testKey1": fmt.Sprintf("%s/somethingelse", workdirParent)},
		},
		{
			name:     "2. Update a single key using key.",
			command:  &Cmd{Cmd: "-edit", Args: []string{"testKey2", "somethingelse"}},
			_map:     map[string]string{"testKey2": workdir},
			expected: map[string]string{"testKey2": fmt.Sprintf("%s/somethingelse", workdirParent)},
		},
		{
			name:     "3. Update a single key using relative path.",
			command:  &Cmd{Cmd: "-edit", Args: []string{tmpdir_relative, "somethingelse"}},
			_map:     map[string]string{"testKey3": tmpdir},
			expected: map[string]string{"testKey3": fmt.Sprintf("%s/somethingelse", workdir)},
		},
		{
			name:     "4. Update a single key, directory name already updated.",
			command:  &Cmd{Cmd: "-edit", Args: []string{"testKey4", tmpdir_relative}},
			_map:     map[string]string{"testKey4": fmt.Sprintf("%s/somethingelse", workdir)},
			expected: map[string]string{"testKey4": tmpdir},
		},
		{
			name:    "5. Update multiple keys using full path.",
			command: &Cmd{Cmd: "-edit", Args: []string{workdir, "somethingelse"}},
			_map: map[string]string{
				"testKey5":   workdir,
				"testKey5_0": tmpdir,
				"testKey5_1": workdirParent,
				"testKey5_2": "home/some/other/dir",
			},
			expected: map[string]string{
				"testKey5":   fmt.Sprintf("%s/somethingelse", workdirParent),
				"testKey5_0": fmt.Sprintf("%s/somethingelse/%s", workdirParent, tmpdir_relative),
				"testKey5_1": workdirParent,
				"testKey5_2": "home/some/other/dir",
			},
		},
		{
			name:    "6. Update multiple keys using a key.",
			command: &Cmd{Cmd: "-edit", Args: []string{"testKey6", "somethingelse"}},
			_map: map[string]string{
				"testKey6":   workdir,
				"testKey6_0": tmpdir,
				"testKey6_1": workdirParent,
				"testKey6_2": "home/some/other/dir",
			},
			expected: map[string]string{
				"testKey6":   fmt.Sprintf("%s/somethingelse", workdirParent),
				"testKey6_0": fmt.Sprintf("%s/somethingelse/%s", workdirParent, tmpdir_relative),
				"testKey6_1": workdirParent,
				"testKey6_2": "home/some/other/dir",
			},
		},
	}

	for _, tt := range tests {
		pathMap := make(map[string]string)
		if tt._map != nil {
			pathMap = tt._map
		}
		data := NewCmdAPI(
			workdir,
			tt.command,
			pathMap,
			tmpfile,
			nil,
		)

		if len(tt._map) > 0 {
			dataUpdate(data.allPaths, tmpfile)
		}

		stdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}
		os.Stdout = w

		err = editPath(data)
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
		os.Stdout = stdout
		output := <-outChan

		for k, v := range data.allPaths {
			ev, ok := tt.expected[k]
			if !ok {
				t.Log(tt.name)
				t.Log(output)
				t.Errorf("Key '%s' not provided in expected field.", k)
			} else if v != ev {
				t.Log(tt.name)
				t.Log(output)
				t.Errorf("Expected key %s to have value %s, got %s", k, tt.expected[k], v)
			}
		}

	}
}

func TestEvalPath(t *testing.T) {

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
		command  *Cmd
		expected string
		allPaths map[string]string
	}{
		{
			name:     "1. Valid key provided, standalone.",
			command:  &Cmd{Cmd: "_", Args: []string{"testKey"}},
			expected: fmt.Sprintln(tmpdir),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
		},
		{
			name:     "2. Valid key provided, evaluate path.",
			command:  &Cmd{Cmd: "_", Args: []string{"testKey/subdir"}},
			expected: fmt.Sprintln(tmpdir2),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
		},
		{
			name:     "3. Invalid key provided.",
			command:  &Cmd{Cmd: "_", Args: []string{"testKye"}},
			expected: fmt.Sprintf(UnrecognizedKeyMsg, "testKye"),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
		},
		{
			name:     "4. Invalid key provided, evaluate path.",
			command:  &Cmd{Cmd: "_", Args: []string{"testKye/subdir"}},
			expected: fmt.Sprintf(InvalidDirectoryMsg, "testKye/subdir", "testKye/subdir"),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		data := NewCmdAPI(
			tmpdir,
			tt.command,
			tt.allPaths,
			nil,
			nil,
		)

		actual, err := evalPath(data)
		if err != nil {
			fmt.Println(tt.name)
			t.Error(err)
		}

		if actual != tt.expected {
			t.Errorf("Expected: %q, got: %q", tt.expected, actual)
		}
	}
}
