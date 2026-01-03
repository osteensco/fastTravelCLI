package data

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureData(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "ft_test_ensure_data")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	testFilePath := filepath.Join(tmpdir, "test.bin")

	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{
			name:     "1. Create new file",
			filepath: testFilePath,
			wantErr:  false,
		},
		{
			name:     "2. Open existing file",
			filepath: testFilePath,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := EnsureData(tt.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if file != nil {
				file.Close()
			}
			if _, err := os.Stat(tt.filepath); os.IsNotExist(err) {
				t.Error("File was not created")
			}
		})
	}
}

func TestReadData(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "ft_test_read_data")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	testFilePath := filepath.Join(tmpdir, "test.bin")

	tests := []struct {
		name     string
		setup    func(f *os.File)
		expected map[string]string
		wantErr  bool
	}{
		{
			name:     "1. Empty file",
			setup:    func(f *os.File) {},
			expected: map[string]string{},
			wantErr:  false,
		},
		{
			name: "2. Single entry",
			setup: func(f *os.File) {
				key := "home"
				val := "/home/user"
				f.Write([]byte{byte(len(key))})
				f.Write([]byte(key))
				var valLen [2]byte
				binary.LittleEndian.PutUint16(valLen[:], uint16(len(val)))
				f.Write(valLen[:])
				f.Write([]byte(val))
			},
			expected: map[string]string{"home": "/home/user"},
			wantErr:  false,
		},
		{
			name: "3. Multiple entries",
			setup: func(f *os.File) {
				data := map[string]string{
					"a": "b",
					"c": "d",
				}
				// Note: map iteration order is random, but ReadData should handle it
				for k, v := range data {
					f.Write([]byte{byte(len(k))})
					f.Write([]byte(k))
					var valLen [2]byte
					binary.LittleEndian.PutUint16(valLen[:], uint16(len(v)))
					f.Write(valLen[:])
					f.Write([]byte(v))
				}
			},
			expected: map[string]string{"a": "b", "c": "d"},
			wantErr:  false,
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

			got, err := ReadData(file)
			file.Close()

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.expected) {
				t.Errorf("ReadData() length got = %v, want %v", len(got), len(tt.expected))
			}
			for k, v := range tt.expected {
				if got[k] != v {
					t.Errorf("ReadData() key %v: got %v, want %v", k, got[k], v)
				}
			}
		})
	}
}

func TestDataUpdate(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "ft_test_data_update")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	testFilePath := filepath.Join(tmpdir, "test.bin")
	file, err := os.OpenFile(testFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	tests := []struct {
		name    string
		data    map[string]string
		wantErr bool
	}{
		{
			name:    "1. Initial write empty map",
			data:    map[string]string{},
			wantErr: false,
		},
		{
			name:    "2. Write single entry",
			data:    map[string]string{"home": "/home/user"},
			wantErr: false,
		},
		{
			name:    "3. Write multiple entries",
			data:    map[string]string{"a": "b", "c": "d"},
			wantErr: false,
		},
		{
			name:    "4. Truncate with smaller map",
			data:    map[string]string{"x": "y"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DataUpdate(tt.data, file)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify by reading back
			file.Seek(0, 0)
			got, err := ReadData(file)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(tt.data) {
				t.Errorf("DataUpdate() verify length got = %v, want %v", len(got), len(tt.data))
			}
			for k, v := range tt.data {
				if got[k] != v {
					t.Errorf("DataUpdate() verify key %v: got %v, want %v", k, got[k], v)
				}
			}
		})
	}
}