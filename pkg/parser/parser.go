package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/unLomTrois/gock3/internal/utils"
	"github.com/unLomTrois/gock3/pkg/cache"
	"github.com/unLomTrois/gock3/pkg/files"
	"github.com/unLomTrois/gock3/pkg/lexer"
	"github.com/unLomTrois/gock3/pkg/lexer/tokens"
	"github.com/unLomTrois/gock3/pkg/parser/ast"
	"github.com/unLomTrois/gock3/pkg/report"
)

// Parser represents the parser with its current state and error manager.
type Parser struct {
	tokenstream  *tokens.TokenStream
	currentToken *tokens.Token
	lookahead    *tokens.Token
	loc          *tokens.Loc
	*report.ErrorManager
}

// NewParser creates a new Parser instance.
func NewParser(tokenstream *tokens.TokenStream) *Parser {
	p := &Parser{
		tokenstream:  tokenstream,
		ErrorManager: report.NewErrorManager(),
	}
	p.currentToken = p.tokenstream.Next()
	p.lookahead = p.tokenstream.Next()
	if p.currentToken != nil {
		p.loc = &p.currentToken.Loc
	}
	return p
}

// -----------------------------------------------------------------------------
// Public API
// -----------------------------------------------------------------------------

// ParseTokenStream performs syntactic analysis on the given token stream and returns
// the AST file block along with any diagnostic items encountered.
func ParseTokenStream(tokenStream *tokens.TokenStream) (*ast.FileBlock, []*report.DiagnosticItem) {
	p := NewParser(tokenStream)
	fileBlock := p.fileBlock()
	return fileBlock, p.Errors()
}

// ParseParadoxFile is the high-level entry point that reads, tokenizes, and parses a
// Paradox file into an AST.
func ParseParadoxFile(file files.ParadoxFile) (*ast.AST, error) {
	// Read file content from disk.
	content, err := utils.ReadFileWithUTF8BOM(file.FullPath())
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	diagnostics := []*report.DiagnosticItem{}

	// Lexical analysis: scan the file content into tokens.
	tokenStream, lexerErrors := lexer.Scan(file, content)
	diagnostics = append(diagnostics, lexerErrors...)

	// Syntactic analysis: parse tokens to create the file block.
	fileBlock, parserErrors := ParseTokenStream(tokenStream)
	diagnostics = append(diagnostics, parserErrors...)

	// Build the AST using the parsed file block.
	astTree := &ast.AST{
		Filename: file.FileName(),
		Fullpath: file.FullPath(),
		Block:    fileBlock,
	}

	// Report diagnostic information.
	reportDiagnostics(diagnostics)

	return astTree, nil
}

// -----------------------------------------------------------------------------
// Parser Internal Methods
// -----------------------------------------------------------------------------

// nextToken advances the current token and the lookahead pointer.
func (p *Parser) nextToken() {
	p.currentToken = p.lookahead
	p.lookahead = p.tokenstream.Next()
	if p.currentToken != nil {
		p.loc = &p.currentToken.Loc
	} else {
		p.loc = nil
	}
}

// -----------------------------------------------------------------------------
// Diagnostic Reporting
// -----------------------------------------------------------------------------

// reportDiagnostics reports all collected diagnostic items.
func reportDiagnostics(diagnostics []*report.DiagnosticItem) {
	fileCache := cache.NewFileCache()
	for _, diag := range diagnostics {
		printDiagnostic(diag, fileCache)
	}
}

// printDiagnostic formats and prints a single diagnostic message.
func printDiagnostic(diag *report.DiagnosticItem, fileCache *cache.FileCache) {
	color := diag.Severity.Color()
	filename, _ := diag.Pointer.Loc.Filename()
	column := diag.Pointer.Loc.Column
	line := diag.Pointer.Loc.Line

	// Special-case: if the error is at the very beginning, output minimal information.
	if line == 1 && column == 1 {
		color.Printf("[%s:%d:%d]: %s\n", filename, line, column, diag.Msg)
		return
	}

	errLine := getErrorLine(fileCache, diag, column)
	color.Printf("[%s:%d:%d]: %s, got %s\n", filename, line, column, diag.Msg, strconv.Quote(errLine))
}

// getErrorLine returns the specific portion of a line where the error occurred.
func getErrorLine(fileCache *cache.FileCache, diag *report.DiagnosticItem, column uint16) string {
	// Retrieve the full line of text for the error's location.
	lineText := fileCache.GetLine(&diag.Pointer.Loc)
	// Normalize tabs into spaces for consistent column matching.
	normalizedLine := strings.ReplaceAll(lineText, "\t", "    ")
	errorEndIndex := int(column) + int(diag.Pointer.Length) - 1
	// Clamp errorEndIndex to the length of the line.
	if errorEndIndex > len(normalizedLine) {
		errorEndIndex = len(normalizedLine)
	}
	return normalizedLine[:errorEndIndex]
}
