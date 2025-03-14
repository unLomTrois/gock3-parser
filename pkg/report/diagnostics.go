package report

import (
	"fmt"

	"github.com/unLomTrois/gock3/pkg/ast"
	"github.com/unLomTrois/gock3/pkg/files"
	"github.com/unLomTrois/gock3/pkg/report/severity"
	"github.com/unLomTrois/gock3/pkg/tokens"
)

type DiagnosticItem struct {
	Severity severity.Severity
	Pointer  *DiagnosticPointer
	Msg      string
}

type DiagnosticPointer struct {
	Loc    tokens.Loc
	Length int
}

func (d *DiagnosticItem) Error() string {
	return fmt.Sprintf("%s: %s", d.Severity, d.Msg)
}

func NewDiagnosticItem(severity severity.Severity, msg string, pointer *DiagnosticPointer) *DiagnosticItem {
	return &DiagnosticItem{
		Severity: severity,
		Msg:      msg,
		Pointer:  pointer,
	}
}

func FromToken(token *tokens.Token, severity severity.Severity, msg string) *DiagnosticItem {
	return &DiagnosticItem{
		Severity: severity,
		Msg:      msg,
		Pointer: &DiagnosticPointer{
			Loc:    token.Loc,
			Length: len(token.Value),
		},
	}
}

func FromFile(file files.ParadoxFile, severity severity.Severity, msg string) *DiagnosticItem {
	loc := tokens.LocFromParadoxFile(file)

	return &DiagnosticItem{
		Severity: severity,
		Msg:      msg,
		Pointer: &DiagnosticPointer{
			Loc:    *loc,
			Length: 0,
		},
	}
}

func FromBlock(file_block *ast.FileBlock, severity severity.Severity, msg string) *DiagnosticItem {
	loc := file_block.Loc

	return &DiagnosticItem{
		Severity: severity,
		Msg:      msg,
		Pointer: &DiagnosticPointer{
			Loc:    loc,
			Length: 0,
		},
	}
}

// FromLoc creates a new DiagnosticItem from a loc
// Primary used in cases when you know the loc but you don't know the token
// Happens in Lexer
func FromLoc(loc tokens.Loc, severity severity.Severity, msg string) *DiagnosticItem {
	return &DiagnosticItem{
		Severity: severity,
		Msg:      msg,
		Pointer: &DiagnosticPointer{
			Loc:    loc,
			Length: 0,
		},
	}
}
