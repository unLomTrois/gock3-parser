package cli_test

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/unLomTrois/gock3/internal/app/cli"
	// If you need to mock or intercept calls in pdxfile/utils,
	// import them and use a test double approach.
)

// --------------------
// Unit-test style tests
// --------------------

func TestNewParseCommand(t *testing.T) {
	tests := []struct {
		name string
		want *cli.ParseCommand
	}{
		{
			name: "Create new parse command",
			want: &cli.ParseCommand{}, // We can’t compare the FlagSet directly, so we mainly check type here
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cli.NewParseCommand()
			// We can’t do a direct DeepEqual on *flag.FlagSet,
			// so we do a simpler check on type for demonstration:
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewParseCommand() = %T, want %T", got, tt.want)
			}
		})
	}
}

func TestParseCommand_Name(t *testing.T) {
	cmd := cli.NewParseCommand()
	expectedName := "parse"
	if cmd.Name() != expectedName {
		t.Errorf("ParseCommand.Name() = %v, want %v", cmd.Name(), expectedName)
	}
}

func TestParseCommand_Description(t *testing.T) {
	cmd := cli.NewParseCommand()
	expectedDesc := "Parse a file and generate the output files"
	if cmd.Description() != expectedDesc {
		t.Errorf("ParseCommand.Description() = %v, want %v", cmd.Description(), expectedDesc)
	}
}

// --------------------------
// Integration-style tests
// --------------------------

func TestParseCommand_MissingArguments(t *testing.T) {
	// If no arguments are passed, we expect an error
	cmd := cli.NewParseCommand()
	err := cmd.Run([]string{})
	if err == nil {
		t.Errorf("expected error for missing arguments, got nil")
	}
}

func TestParseCommand_InvalidFlag(t *testing.T) {
	// Pass an invalid flag to trigger flag parsing error
	cmd := cli.NewParseCommand()
	err := cmd.Run([]string{"somefile.txt", "--unknown-flag"})
	if err == nil {
		t.Errorf("expected error for invalid flag, got nil")
	}
	if err != nil && err.Error() != "failed to parse flags: flag provided but not defined: -unknown-flag" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParseCommand_NonExistingFile(t *testing.T) {
	// If the file doesn’t exist, we expect an error from FileExists
	cmd := cli.NewParseCommand()
	err := cmd.Run([]string{"this_file_does_not_exist.txt"})
	if err == nil {
		t.Errorf("expected error for non-existing file, got nil")
	}
}

func TestParseCommand_ValidFile_NoAstFlag(t *testing.T) {
	// Create a temporary file to parse
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // clean up

	// Write something minimal to the file
	_, err = tmpFile.WriteString("hello = world")
	if err != nil {
		t.Fatal(err)
	}
	// Close before we parse it
	tmpFile.Close()

	cmd := cli.NewParseCommand()
	err = cmd.Run([]string{tmpFile.Name()})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestParseCommand_ValidFile_WithAstFlag(t *testing.T) {
	// Create a temporary file to parse
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("hello = world")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Create a path where we want to save our AST
	tmpDir := t.TempDir()
	astPath := filepath.Join(tmpDir, "ast.json")

	cmd := cli.NewParseCommand()
	// Provide the file and the --save-ast flag
	err = cmd.Run([]string{tmpFile.Name(), "--save-ast", astPath})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify that the AST file was created
	if _, statErr := os.Stat(astPath); os.IsNotExist(statErr) {
		t.Errorf("expected AST file to exist at %s, but it does not", astPath)
	}
}

// --------------------------
// Additional tests for error coverage
// --------------------------

// TestParseCommand_FileParseFailure simulates a failure inside parseFile.
// One way to do this is by providing a file that triggers a parse error.
// You may need to modify the parse logic or use a mock to force an error.
func TestParseCommand_FileParseFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod 0000 doesn't work on windows")
	}

	// Here we create an unreadable file (permissions 000)
	// so it fails during parseFile -> pdxfile.ParseFile
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = tmpFile.WriteString("hello \\= world")
	tmpFile.Close()

	// Make the file unreadable
	os.Chmod(tmpFile.Name(), 0000)
	defer func() {
		os.Chmod(tmpFile.Name(), 0644) // restore so we can remove it
		os.Remove(tmpFile.Name())
	}()

	cmd := cli.NewParseCommand()
	err = cmd.Run([]string{tmpFile.Name()})
	if err == nil {
		t.Errorf("expected parse error, got nil")
	}
}

// TestParseCommand_SaveASTFailure simulates a failure in handleAST.
// We make the output path unwritable or invalid.
func TestParseCommand_SaveASTFailure(t *testing.T) {
	// Create a valid input file
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = tmpFile.WriteString("hello = world")
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Use a directory path instead of a file path to cause save error
	// Or a path in a directory we don't have permissions for
	// e.g., on Unix systems: "/root/ast.json"
	// For demonstration, use a directory as the ast.json "file".
	tmpDir := t.TempDir()
	astPath := tmpDir // This is a directory, not a file

	cmd := cli.NewParseCommand()
	err = cmd.Run([]string{tmpFile.Name(), "--save-ast", astPath})
	if err == nil {
		t.Errorf("expected error when trying to save AST to a directory, got nil")
	}
	// Optionally, test the exact error message:
	// e.g., if err != nil && !strings.Contains(err.Error(), "failed to save AST") { ... }
}

// If you want to precisely confirm the returned error messages,
// you can do so with string checks or by using `errors.As/Is`.
