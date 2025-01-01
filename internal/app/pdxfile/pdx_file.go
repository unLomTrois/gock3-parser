package pdxfile

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/unLomTrois/gock3/internal/app/files"
	"github.com/unLomTrois/gock3/internal/app/lexer"
	"github.com/unLomTrois/gock3/internal/app/parser"
	"github.com/unLomTrois/gock3/internal/app/parser/ast"
	"github.com/unLomTrois/gock3/internal/app/utils"
	"github.com/unLomTrois/gock3/pkg/cache"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

func ParseFile(entry *files.FileEntry) (*ast.AST, error) {
	content, err := utils.ReadFileWithUTF8BOM(entry.FullPath())
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	diagnostics := []*report.DiagnosticItem{}

	tokenStream, lexerErrors := lexer.Scan(entry, content)
	diagnostics = append(diagnostics, lexerErrors...)

	fileBlock, parserErrors := parser.Parse(tokenStream)
	diagnostics = append(diagnostics, parserErrors...)

	ast := &ast.AST{
		Filename: entry.FileName(),
		Fullpath: entry.FullPath(),
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
	color := getSeverityColor(err.Severity)
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

func getSeverityColor(sev severity.Severity) *color.Color {
	switch sev {
	case severity.Error:
		return color.New(color.FgRed)
	case severity.Warning:
		return color.New(color.FgYellow)
	case severity.Info:
		return color.New(color.FgCyan)
	case severity.Critical:
		return color.New(color.FgHiMagenta)
	default:
		return color.New(color.Reset)
	}
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
