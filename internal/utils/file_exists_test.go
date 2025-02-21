package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileExists(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "testfile")
		if err := os.WriteFile(tmpFile, nil, 0644); err != nil {
			t.Fatal("Failed to create test file")
		}

		absPath, err := FileExists(tmpFile)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if absPath != tmpFile {
			t.Fatalf("Expected path %q, got %q", tmpFile, absPath)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		tmpDir := t.TempDir()
		nonExistent := filepath.Join(tmpDir, "missing")

		absPath, err := FileExists(nonExistent)
		if err == nil {
			t.Fatal("Expected error but got nil")
		}
		if absPath != nonExistent {
			t.Fatalf("Expected path %q, got %q", nonExistent, absPath)
		}
		if err.Error() != "file does not exist: "+nonExistent {
			t.Fatalf("Unexpected error message: %v", err)
		}
	})

	t.Run("existing directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		absPath, err := FileExists(tmpDir)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if absPath != tmpDir {
			t.Fatalf("Expected path %q, got %q", tmpDir, absPath)
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		invalidPath := string([]byte{0}) // Null character, invalid path
		absPath, err := FileExists(invalidPath)

		if err == nil {
			t.Fatal("Expected error but got nil")
		}
		if absPath != "" {
			t.Fatalf("Expected empty path, got %q", absPath)
		}
		expectedErr := "unable to get absolute path:"
		if !strings.Contains(err.Error(), expectedErr) {
			t.Fatalf("Expected error containing %q, got %q", expectedErr, err.Error())
		}
	})
}
