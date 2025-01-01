package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadFileWithUTF8BOM(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// 1. Create a file WITHOUT BOM
	noBOMFilePath := filepath.Join(tmpDir, "no_bom.txt")
	noBOMContent := []byte("Hello, no BOM here.")
	if err := os.WriteFile(noBOMFilePath, noBOMContent, 0644); err != nil {
		t.Fatalf("Failed to create test file without BOM: %v", err)
	}

	// 2. Create a file WITH BOM
	withBOMFilePath := filepath.Join(tmpDir, "with_bom.txt")
	bom := []byte{0xEF, 0xBB, 0xBF}
	withBOMContent := append(bom, []byte("Hello with BOM!")...)
	if err := os.WriteFile(withBOMFilePath, withBOMContent, 0644); err != nil {
		t.Fatalf("Failed to create test file with BOM: %v", err)
	}

	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "File without BOM",
			args: args{
				filePath: noBOMFilePath,
			},
			want:    noBOMContent, // entire content should be returned
			wantErr: false,
		},
		{
			name: "File with BOM",
			args: args{
				filePath: withBOMFilePath,
			},
			// expected content is original minus the first 3 BOM bytes
			want:    withBOMContent[3:],
			wantErr: false,
		},
		{
			name: "Non-existent file",
			args: args{
				filePath: filepath.Join(tmpDir, "does_not_exist.txt"),
			},
			want:    nil, // we expect an error
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFileWithUTF8BOM(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFileWithUTF8BOM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFileWithUTF8BOM()\n got  = %v\n want = %v", got, tt.want)
			}
		})
	}
}

func Test_hasUTF8BOM(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty slice",
			args: args{
				data: []byte{},
			},
			want: false,
		},
		{
			name: "Short slice (less than 3 bytes)",
			args: args{
				data: []byte{0xEF, 0xBB},
			},
			want: false,
		},
		{
			name: "Exact BOM",
			args: args{
				data: []byte{0xEF, 0xBB, 0xBF},
			},
			want: true,
		},
		{
			name: "BOM plus content",
			args: args{
				data: []byte{0xEF, 0xBB, 0xBF, 'H', 'i'},
			},
			want: true,
		},
		{
			name: "Different bytes",
			args: args{
				data: []byte{0x12, 0x34, 0x56},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasUTF8BOM(tt.args.data); got != tt.want {
				t.Errorf("hasUTF8BOM() = %v, want %v", got, tt.want)
			}
		})
	}
}
