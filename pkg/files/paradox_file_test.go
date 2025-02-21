package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewParadoxFile(t *testing.T) {
	t.Run("panics when file doesn't exist", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for non-existent file")
			}
		}()
		NewParadoxTxtFile("non_existent_file.txt", Vanilla)
	})

	t.Run("returns a valid ParadoxTxtFile", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test_file_*.txt")
		if err != nil {
			t.Fatalf("Failed to create temporary file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		paradoxFile := NewParadoxTxtFile(tmpFile.Name(), Vanilla)
		if paradoxFile == nil {
			t.Error("Expected non-nil ParadoxTxtFile instance")
		}
	})
}

func TestParadoxTxtFile_Methods(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_file-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	paradoxFile := NewParadoxTxtFile(tmpFile.Name(), Vanilla)

	t.Run("FileName returns correct base name", func(t *testing.T) {
		got := paradoxFile.FileName()
		want := filepath.Base(tmpFile.Name())
		if got != want {
			t.Errorf("FileName mismatch: got %s, want %s", got, want)
		}
	})

	t.Run("FullPath returns initialized path", func(t *testing.T) {
		if paradoxFile.FullPath() != tmpFile.Name() {
			t.Errorf("FullPath mismatch: got %s, want %s", paradoxFile.FullPath(), tmpFile.Name())
		}
	})

	t.Run("Kind returns initialized kind", func(t *testing.T) {
		if paradoxFile.Kind() != Vanilla {
			t.Errorf("Kind mismatch: got %v, want %v", paradoxFile.Kind(), Vanilla)
		}
	})
}
