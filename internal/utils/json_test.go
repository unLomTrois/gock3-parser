package utils_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/unLomTrois/gock3/internal/utils"
)

func TestSaveJSON(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		data     interface{}
		filename string
		wantErr  bool
	}{
		{
			name:     "successful save",
			data:     map[string]interface{}{"test": 123},
			filename: "valid.json",
			wantErr:  false,
		},
		{
			name:     "invalid path",
			data:     map[string]interface{}{},
			filename: filepath.Join("nonexistent", "file.json"),
			wantErr:  true,
		},
		{
			name:     "invalid data",
			data:     make(chan int), // channels can't be serialized
			filename: "invalid.json",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath := filepath.Join(tmpDir, tt.filename)
			err := utils.SaveJSON(tt.data, fullPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file contents
				file, err := os.Open(fullPath)
				if err != nil {
					t.Fatalf("Failed to open output file: %v", err)
				}
				defer file.Close()

				var decoded interface{}
				if err := json.NewDecoder(file).Decode(&decoded); err != nil {
					t.Errorf("Failed to decode saved JSON: %v", err)
				}

				// Verify formatting
				file.Seek(0, 0)
				enc := json.NewEncoder(file)
				enc.SetIndent("", "\t")
			}
		})
	}
}

func TestSaveJSON_Formatting(t *testing.T) {
	tmpDir := t.TempDir()
	fullPath := filepath.Join(tmpDir, "formatting_test.json")

	data := struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
		HTML  string `json:"html"`
	}{
		Name:  "Test",
		Value: 42,
		HTML:  "<div>content</div>",
	}

	err := utils.SaveJSON(data, fullPath)
	if err != nil {
		t.Fatalf("SaveJSON failed: %v", err)
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expected := `{
	"name": "Test",
	"value": 42,
	"html": "<div>content</div>"
}
`
	if string(content) != expected {
		t.Errorf("Unexpected file content:\nGot:\n%s\nWant:\n%s", content, expected)
	}
}
