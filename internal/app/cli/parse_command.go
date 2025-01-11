package cli

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/internal/app/pdxfile"
	"github.com/unLomTrois/gock3/internal/app/utils"
)

type ParseCommand struct {
	flagset     *flag.FlagSet
	astFilepath string
}

// NewParseCommand initializes a new ParseCommand with the appropriate flags.
func NewParseCommand() *ParseCommand {
	pc := &ParseCommand{
		flagset: flag.NewFlagSet("parse", flag.ContinueOnError),
	}

	// CLI usage example:
	//   gock3 parse file.txt --save-ast ast.json
	pc.flagset.StringVar(
		&pc.astFilepath,
		"save-ast",
		"",
		"Save the AST to a file\nExample: --save-ast ast.json",
	)

	return pc
}

// Name returns the name of the command.
func (pc *ParseCommand) Name() string {
	return pc.flagset.Name()
}

// Description returns a short description of what the command does.
func (pc *ParseCommand) Description() string {
	return "Parse a file and generate the output files"
}

// Run is the entry point for the 'parse' command. It parses the arguments,
// validates them, reads the file, and handles the resulting AST.
func (pc *ParseCommand) Run(args []string) error {
	// 1. Parse the CLI arguments
	if err := pc.parseArgs(args); err != nil {
		return err
	}

	filePath := args[0]
	fullpath, err := utils.FileExists(filePath)
	if err != nil {
		return err
	}

	// 2. Parse the file to get the AST
	ast, err := pc.parseFile(fullpath)
	if err != nil {
		return err
	}

	// 3. Handle the AST (save to file if needed)
	if err := pc.handleAST(ast); err != nil {
		return err
	}

	return nil
}

// parseArgs validates and parses the incoming arguments using the command's flagset.
func (pc *ParseCommand) parseArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("not enough arguments (no file specified)")
	}

	// Skip args[0] (the file path) when parsing flags
	if err := pc.flagset.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	return nil
}

// parseFile reads and parses the specified file into an AST structure.
func (pc *ParseCommand) parseFile(fullpath string) (*ast.AST, error) {
	fileEntry := files.NewFileEntry(fullpath, files.FileKind(files.Mod))

	ast, err := pdxfile.ParseFile(fileEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return ast, nil
}

// handleAST handles the logic for the parsed AST, such as saving it to disk.
func (pc *ParseCommand) handleAST(ast *ast.AST) error {
	// If no --save-ast path is provided, nothing more to do
	if pc.astFilepath == "" {
		return nil
	}

	// Otherwise, save the AST to the specified file
	if err := utils.SaveJSON(ast, pc.astFilepath); err != nil {
		return fmt.Errorf("failed to save AST: %w", err)
	}

	absPath, err := filepath.Abs(pc.astFilepath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	log.Println("Saved parse tree to", absPath)
	return nil
}
