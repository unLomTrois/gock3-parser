// helpers.go
package parser

import (
	"fmt"
	"strings"

	"github.com/unLomTrois/gock3/pkg/tokens"
)

// isNextField determines if the upcoming tokens likely form a field.
func (p *Parser) isNextField() bool {
	return isKeyToken(p.currentToken.Type) && isOperatorToken(p.lookahead.Type)
}

// isKeyToken checks if a token type is a valid key.
func isKeyToken(tokenType tokens.TokenType) bool {
	return tokenType == tokens.WORD || tokenType == tokens.DATE || tokenType == tokens.NUMBER
}

// isOperatorToken checks if a token type is a valid operator.
func isOperatorToken(tokenType tokens.TokenType) bool {
	return isEqualOperatorToken(tokenType) || tokenType == tokens.COMPARISON
}

// isEqualOperatorToken checks if a token type is an equality operator.
func isEqualOperatorToken(tokenType tokens.TokenType) bool {
	return tokenType == tokens.EQUALS || tokenType == tokens.QUESTION_EQUALS
}

// isLiteralType checks if a token type represents a literal value.
func isLiteralType(tokenType tokens.TokenType) bool {
	switch tokenType {
	case tokens.WORD, tokens.NUMBER, tokens.BOOL, tokens.QUOTED_STRING:
		return true
	default:
		return false
	}
}

// formatTokenTypes formats a slice of TokenType into a human-readable string.
func formatTokenTypes(types []tokens.TokenType) string {
	if len(types) == 0 {
		return "no token types specified"
	}
	parts := make([]string, len(types))
	for i, t := range types {
		parts[i] = fmt.Sprintf("%q", t)
	}
	if len(parts) == 1 {
		return parts[0]
	} else if len(parts) == 2 {
		return fmt.Sprintf("%s or %s", parts[0], parts[1])
	}
	return fmt.Sprintf("%s, or %s", strings.Join(parts[:len(parts)-1], ", "), parts[len(parts)-1])
}
