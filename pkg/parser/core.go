// core.go
package parser

import (
	"fmt"
	"strconv"

	"github.com/unLomTrois/gock3/pkg/lexer/tokens"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

// Expect verifies that the current token matches one of the expected types.
// If it does, it consumes and returns the token.
// Otherwise, it reports an error and attempts to recover.
func (p *Parser) Expect(expectedTypes ...tokens.TokenType) *tokens.Token {
	token := p.currentToken
	if token == nil {
		errMsg := fmt.Sprintf(errUnexpectedEOF, formatTokenTypes(expectedTypes))
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	for _, et := range expectedTypes {
		if token.Type == et {
			p.nextToken()
			return token
		}
	}

	errMsg := fmt.Sprintf(errUnexpectedToken, token.Value, token.Type, formatTokenTypes(expectedTypes))
	err := report.FromToken(token, severity.Error, errMsg)
	p.AddError(err)

	recoveryPoint := RecoveryPoint{
		TokenTypes: expectedTypes,
		Context:    "expected " + formatTokenTypes(expectedTypes),
	}

	p.synchronize(recoveryPoint) // attempt recovery regardless
	return nil
}

// unquoteExpect parses and unquotes a quoted string token.
func (p *Parser) unquoteExpect(expectedType tokens.TokenType) *tokens.Token {
	token := p.Expect(expectedType)
	if token == nil {
		return nil
	}

	unquotedValue, err := strconv.Unquote(token.Value)
	if err != nil {
		errMsg := fmt.Sprintf(errFailedUnquoteString, token.Value)
		diag := report.FromToken(token, severity.Error, errMsg)
		p.AddError(diag)
		// Keep the original value if unquoting fails
		return token
	}

	token.Value = unquotedValue
	return token
}
