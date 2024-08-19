package ft

import (
	"fmt"
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
	file := EnsureData(testFilePath)
	defer file.Close()

	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Fatalf("File was not created")
	} else {
		fmt.Println("ensureData: Success")
	}
}
