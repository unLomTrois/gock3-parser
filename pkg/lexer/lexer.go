// pkg/lexer/lexer.go
package lexer

import (
	"bytes"
	"fmt"

	"github.com/unLomTrois/gock3/pkg/files"
	"github.com/unLomTrois/gock3/pkg/lexer/tokens"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

type Lexer struct {
	file           files.ParadoxFile
	text           []byte
	cursor         int
	line           int
	column         int
	patternMatcher *TokenPatternMatcher
	*report.ErrorManager
}

// NewLexer creates a new Lexer instance.
func NewLexer(file files.ParadoxFile, text []byte) *Lexer {
	return &Lexer{
		file:           file,
		text:           NormalizeText(text),
		cursor:         0,
		line:           1,
		column:         1,
		patternMatcher: NewTokenPatternMatcher(),
		ErrorManager:   report.NewErrorManager(),
	}
}

// NormalizeText replaces CRLF with LF.
func NormalizeText(text []byte) []byte {
	// Optionally, you could trim spaces if needed:
	// text = bytes.TrimSpace(text)
	return bytes.ReplaceAll(text, []byte("\r\n"), []byte("\n"))
}

// hasMoreTokens checks if there are unprocessed tokens.
func (lex *Lexer) hasMoreTokens() bool {
	return lex.cursor < len(lex.text)
}

// Scan tokenizes the entire input text.
func Scan(file files.ParadoxFile, text []byte) (*tokens.TokenStream, []*report.DiagnosticItem) {
	lex := NewLexer(file, text)
	tokenStream := tokens.NewTokenStream()

	for lex.hasMoreTokens() {
		if token := lex.getNextToken(); token != nil {
			tokenStream.Push(token)
		}
	}

	return tokenStream, lex.Errors()
}

// remainder returns the unprocessed text.
func (lex *Lexer) remainder() []byte {
	return lex.text[lex.cursor:]
}

// getNextToken retrieves the next token from the input text.
func (lex *Lexer) getNextToken() *tokens.Token {
	if !lex.hasMoreTokens() {
		return nil
	}

	remaining := lex.remainder()
	startLine, startColumn := lex.line, lex.column

	// Attempt to match tokens in the specified order.
	for _, tokenType := range tokens.TokenCheckOrder {
		if match := lex.patternMatcher.MatchToken(tokenType, remaining); match != nil {
			return lex.processMatch(tokenType, match, startLine, startColumn)
		}
	}

	lex.reportUnexpectedToken()
	return nil
}

// processMatch handles a successful token match.
func (lex *Lexer) processMatch(tokenType tokens.TokenType, match []byte, startLine, startColumn int) *tokens.Token {
	tokenValue := string(match)
	lex.cursor += len(match)

	switch tokenType {
	case tokens.TAB:
		// Consider tab width as 4 spaces.
		lex.column += 4
		return nil
	case tokens.NEXTLINE:
		lex.line++
		lex.column = 1
		loc := tokens.LocFromParadoxFile(lex.file)
		loc.Line = uint32(lex.line)
		loc.Column = uint16(lex.column)
		return tokens.New(tokenValue, tokenType, *loc)
	case tokens.WHITESPACE:
		lex.column++
		return nil
	case tokens.COMMENT:
		return nil
	default:
		lex.column += len(match)
		loc := tokens.LocFromParadoxFile(lex.file)
		loc.Line = uint32(startLine)
		loc.Column = uint16(startColumn)
		return tokens.New(tokenValue, tokenType, *loc)
	}
}

// reportUnexpectedToken logs an error for an unexpected token and advances the cursor.
func (lex *Lexer) reportUnexpectedToken() {
	remaining := lex.remainder()
	unexpectedChar := remaining[0]

	loc := tokens.LocFromParadoxFile(lex.file)
	loc.Line = uint32(lex.line)
	loc.Column = uint16(lex.column)
	err := report.FromLoc(*loc, severity.Critical, fmt.Sprintf("unexpected token '%c'", unexpectedChar))
	lex.AddError(err)

	// Advance to prevent an infinite loop.
	lex.cursor++
	lex.column++
}
