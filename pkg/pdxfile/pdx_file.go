// TODO: rename package
package pdxfile

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/unLomTrois/gock3/internal/utils"
	"github.com/unLomTrois/gock3/pkg/cache"
	"github.com/unLomTrois/gock3/pkg/files"
	"github.com/unLomTrois/gock3/pkg/lexer"
	"github.com/unLomTrois/gock3/pkg/parser"
	"github.com/unLomTrois/gock3/pkg/parser/ast"
	"github.com/unLomTrois/gock3/pkg/report"
)

func ParseFile(file files.ParadoxFile) (*ast.AST, error) {
	content, err := utils.ReadFileWithUTF8BOM(file.FullPath())
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	diagnostics := []*report.DiagnosticItem{}

	tokenStream, lexerErrors := lexer.Scan(file, content)
	diagnostics = append(diagnostics, lexerErrors...)

	fileBlock, parserErrors := parser.Parse(tokenStream)
	diagnostics = append(diagnostics, parserErrors...)

	ast := &ast.AST{
		Filename: file.FileName(),
		Fullpath: file.FullPath(),
		Block:    fileBlock,
	}

	finalize(diagnostics)

	return ast, nil
}

func finalize(errs []*report.DiagnosticItem) {
	fileCache := cache.NewFileCache()

	for _, err := range errs {
		printDiagnostic(err, fileCache)
	}
}

func printDiagnostic(err *report.DiagnosticItem, fileCache *cache.FileCache) {
	color := err.Severity.Color()
	filename, _ := err.Pointer.Loc.Filename()
	column := err.Pointer.Loc.Column
	line := err.Pointer.Loc.Line

	if err.Pointer.Loc.Line == 1 && err.Pointer.Loc.Column == 1 {
		color.Printf("[%s:%d:%d]: %s\n", filename, line, column, err.Msg)
		return
	}

	errLine := getErrorLine(fileCache, err, column)
	color.Printf("[%s:%d:%d]: %s, got %s\n", filename, line, column, err.Msg, strconv.Quote(errLine))
}

func getErrorLine(fileCache *cache.FileCache, err *report.DiagnosticItem, column uint16) string {
	// Get the full line of text where the error occurred
	line := fileCache.GetLine(&err.Pointer.Loc)

	// Normalize tabs to spaces for consistent column counting
	normalizedLine := strings.ReplaceAll(line, "\t", "    ")

	// Calculate the end index of the error span
	errorEndIndex := column + uint16(err.Pointer.Length) - 1

	// Return the portion of the line up to the error end
	return normalizedLine[:errorEndIndex]
}
