package ft

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	ftdata "github.com/osteensco/fastTravelCLI/data"
)

func TestEvalBookmark(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "ft_test_bookmark")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	subdir := filepath.Join(tmpdir, "subdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatal(err)
	}

	allPaths := map[string]string{
		"home": "/home/user",
		"tmp":  tmpdir,
	}

	tests := []struct {
		name     string
		path     string
		expected string
		wantErr  bool
	}{
		{
			name:     "1. Existing bookmark",
			path:     "home",
			expected: "/home/user",
			wantErr:  false,
		},
		{
			name:     "2. Invalid bookmark",
			path:     "nowhere",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "3. Existing Bookmark subdirectory",
			path:     fmt.Sprintf("tmp/subdir"),
			expected: subdir,
			wantErr:  false,
		},
		{
			name:     "4. Invalid Bookmark subdirectory",
			path:     fmt.Sprintf("invalid/nonexistent"),
			expected: "",
			wantErr:  true,
		},
		{
			name:     "5. Valid Bookmark, invalid subdirectory",
			path:     fmt.Sprintf("tmp/nonexistent"),
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evalBookmark(allPaths, tt.path)
			if (err != nil) != tt.wantErr {
				t.Log(tt.name)
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.expected {
				t.Log(tt.name)
				t.Errorf("got %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestEvalCDPATH(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "ft_test_cdpath")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	dir1 := filepath.Join(tmpdir, "dir1")
	dir2 := filepath.Join(tmpdir, "dir2")
	target1 := filepath.Join(dir1, "target1")
	target2 := filepath.Join(dir2, "target2")

	if err := os.MkdirAll(target1, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(target2, 0755); err != nil {
		t.Fatal(err)
	}

	oldCDPATH := os.Getenv("CDPATH")
	defer os.Setenv("CDPATH", oldCDPATH)

	os.Setenv("CDPATH", fmt.Sprintf("%s:%s", dir1, dir2))

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "1. Find in first CDPATH entry",
			path:     "target1",
			expected: target1,
		},
		{
			name:     "2. Find in second CDPATH entry",
			path:     "target2",
			expected: target2,
		},
		{
			name:     "3. Not found in CDPATH",
			path:     "nonexistent",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evalCDPATH(nil, tt.path)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("got %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestEvalRelative(t *testing.T) {
	wd, err := os.Getwd()
	absPath, err := os.MkdirTemp(wd, "testing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(absPath)

	tmpDir := filepath.Base(absPath)

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "1. Valid relative path",
			path:     tmpDir,
			expected: absPath,
		},
		{
			name:     "2. Invalid relative path",
			path:     "nonexistent_rel_path_12345",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evalRelative(nil, tt.path)
			if err != nil {
				t.Log(tt.name)
				t.Errorf("Unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Log(tt.name)
				t.Errorf("got %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestEvalPath(t *testing.T) {

	settings := ftdata.GenerateDefaultSettings()

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
		wantErr  bool
		err      string
		allPaths map[string]string
	}{
		{
			name:     "1. Valid key provided, standalone.",
			command:  &Cmd{Cmd: "_", Args: []string{"testKey"}},
			expected: tmpdir,
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
		},
		{
			name:     "2. Valid key provided, evaluate path.",
			command:  &Cmd{Cmd: "_", Args: []string{"testKey/subdir"}},
			expected: tmpdir2,
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
		},
		{
			name:     "3. Invalid key provided.",
			command:  &Cmd{Cmd: "_", Args: []string{"testKye"}},
			expected: "",
			wantErr:  true,
			err:      fmt.Sprintf(UnrecognizedKeyMsg, "testKye"),
			allPaths: map[string]string{
				"testKey": tmpdir,
			},
		},
		{
			name:     "4. Invalid key provided, evaluate path.",
			command:  &Cmd{Cmd: "_", Args: []string{"testKye/subdir"}},
			expected: "",
			wantErr:  true,
			err:      fmt.Sprintf(InvalidDirectoryMsg, "testKye/subdir", "testKye/subdir"),
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
			settings,
			nil,
			nil,
			nil,
		)

		actual, err := evalPath(data, &data.cmd.Args[0])
		if tt.wantErr {
			if err == nil {
				t.Error("expected the following error but did not get it - ", tt.err)
			}
			if tt.err != err.Error() {
				t.Errorf("Expected error: %q, got: %q", tt.err, err)
			}
		} else if err != nil {
			t.Error(err)
		}

		if actual != tt.expected {
			t.Errorf("Expected: %q, got: %q", tt.expected, actual)
		}
	}
}
