package data

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGenerateDefaultSettings(t *testing.T) {
	got := GenerateDefaultSettings()
	want := &Settings{
		QueryOrder:   []string{"bookmark", "CDPATH", "relative"},
		UpdatePrompt: false,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GenerateDefaultSettings() got = %v, want %v", got, want)
	}
}

func TestReadSettings(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "ft_test_read_settings")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	testFilePath := filepath.Join(tmpdir, "settings.json")

	tests := []struct {
		name     string
		setup    func(f *os.File)
		expected *Settings
		wantErr  bool
	}{
		{
			name:  "1. Empty file returns defaults",
			setup: func(f *os.File) {},
			expected: &Settings{
				QueryOrder:   []string{"bookmark", "CDPATH", "relative"},
				UpdatePrompt: false,
			},
			wantErr: false,
		},
		{
			name: "2. Valid JSON file",
			setup: func(f *os.File) {
				s := &Settings{
					QueryOrder:   []string{"relative", "bookmark"},
					UpdatePrompt: true,
				}
				json.NewEncoder(f).Encode(s)
			},
			expected: &Settings{
				QueryOrder:   []string{"relative", "bookmark"},
				UpdatePrompt: true,
			},
			wantErr: false,
		},
		{
			name: "3. Invalid JSON file",
			setup: func(f *os.File) {
				f.WriteString("{invalid json")
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.OpenFile(testFilePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
			if err != nil {
				t.Fatal(err)
			}
			tt.setup(file)
			file.Seek(0, 0)

			got, err := ReadSettings(file)
			file.Close()

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ReadSettings() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestWriteSettings(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "ft_test_write_settings")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	t.Run("1. Write to new file", func(t *testing.T) {
		fPath := filepath.Join(tmpdir, "new_settings.json")
		f, err := os.OpenFile(fPath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		settings := &Settings{
			QueryOrder:   []string{"a", "b"},
			UpdatePrompt: true,
		}

		if err := WriteSettings(f, settings); err != nil {
			t.Fatalf("WriteSettings failed: %v", err)
		}

		f.Seek(0, 0)
		var got Settings
		if err := json.NewDecoder(f).Decode(&got); err != nil {
			t.Fatalf("Failed to decode written settings: %v", err)
		}
		if !reflect.DeepEqual(&got, settings) {
			t.Errorf("Got %v, want %v", got, settings)
		}
	})

	t.Run("2. Overwrite existing file", func(t *testing.T) {
		fPath := filepath.Join(tmpdir, "existing_settings.json")
		// initialContent is longer than what we are writing to so we ensure file is properly truncated
		initialContent := "{ \"garbage\": \"data\", \"more\": \"stuff\" }"
		if err := os.WriteFile(fPath, []byte(initialContent), 0644); err != nil {
			t.Fatal(err)
		}

		f, err := os.OpenFile(fPath, os.O_RDWR, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		settings := &Settings{
			QueryOrder:   []string{"c"},
			UpdatePrompt: false,
		}

		if err := WriteSettings(f, settings); err != nil {
			t.Fatalf("WriteSettings failed: %v", err)
		}

		f.Seek(0, 0)
		var got Settings
		if err := json.NewDecoder(f).Decode(&got); err != nil {
			t.Fatalf("Failed to decode overwritten settings: %v", err)
		}
		if !reflect.DeepEqual(&got, settings) {
			t.Errorf("Got %v, want %v", got, settings)
		}
	})

	t.Run("3. Write error (closed file)", func(t *testing.T) {
		fPath := filepath.Join(tmpdir, "error_file.json")
		f, err := os.Create(fPath)
		if err != nil {
			t.Fatal(err)
		}
		f.Close() // Close immediately to trigger error

		settings := &Settings{}
		err = WriteSettings(f, settings)
		if err == nil {
			t.Error("Expected error writing to closed file, got nil")
		}
	})
}
