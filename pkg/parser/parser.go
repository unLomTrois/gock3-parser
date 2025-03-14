package parser

import (
	"fmt"
	"strconv"

	"github.com/unLomTrois/gock3/internal/utils"
	"github.com/unLomTrois/gock3/pkg/ast"
	"github.com/unLomTrois/gock3/pkg/cache"
	"github.com/unLomTrois/gock3/pkg/files"
	"github.com/unLomTrois/gock3/pkg/lexer"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/tokens"
)

// Parser represents the stateful parser for Paradox files.
type Parser struct {
	tokenstream  *tokens.TokenStream
	currentToken *tokens.Token
	lookahead    *tokens.Token
	loc          *tokens.Loc
	*report.ErrorManager
}

// NewParser creates and initializes a new Parser instance.
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
	content, err := utils.ReadFileWithUTF8BOM(file.FullPath())
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	diagnostics := []*report.DiagnosticItem{}
	tokenStream, lexerErrors := lexer.Scan(file, content)
	diagnostics = append(diagnostics, lexerErrors...)

	fileBlock, parserErrors := ParseTokenStream(tokenStream)
	diagnostics = append(diagnostics, parserErrors...)

	astTree := &ast.AST{
		Filename: file.FileName(),
		Fullpath: file.FullPath(),
		Block:    fileBlock,
	}

	reportDiagnostics(diagnostics)

	return astTree, nil
}

// nextToken advances the token stream.
func (p *Parser) nextToken() {
	p.currentToken = p.lookahead
	p.lookahead = p.tokenstream.Next()
	if p.currentToken != nil {
		p.loc = &p.currentToken.Loc
	} else {
		p.loc = nil
	}
}

// reportDiagnostics outputs all diagnostic items.
func reportDiagnostics(diagnostics []*report.DiagnosticItem) {
	fileCache := cache.NewFileCache()
	for _, diag := range diagnostics {
		printDiagnostic(diag, fileCache)
	}
}

func printDiagnostic(diag *report.DiagnosticItem, fileCache *cache.FileCache) {
	color := diag.Severity.Color()
	filename, _ := diag.Pointer.Loc.Filename()
	line := diag.Pointer.Loc.Line
	column := diag.Pointer.Loc.Column

	// Special-case: if the error is at the very beginning, output minimal information.
	if line == 1 && column == 1 {
		color.Printf("[%s:%d:%d]: %s\n", filename, line, column, diag.Msg)
		return
	}

	errLine := getErrorLine(fileCache, diag, column)
	color.Printf("[%s:%d:%d]: %s, got %s\n", filename, line, column, diag.Msg, strconv.Quote(errLine))
}

func getErrorLine(fileCache *cache.FileCache, diag *report.DiagnosticItem, column uint16) string {
	lineText := fileCache.GetLine(&diag.Pointer.Loc)
	normalizedLine := lineText // Optionally, normalize tabs/spaces if needed.
	errorEndIndex := int(column) + int(diag.Pointer.Length) - 1
	if errorEndIndex > len(normalizedLine) {
		errorEndIndex = len(normalizedLine)
	}
	return normalizedLine[:errorEndIndex]
}
