package ft

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureData(t *testing.T) {
	// Create a temporary directory for testing
	tmpdir, err := os.MkdirTemp("", "testdata")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	testFilePath := filepath.Join(tmpdir, "test.bin")
	file, err := EnsureData(testFilePath)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Error("File was not created")
	} else if err != nil {
		t.Error(err)
	}
}
