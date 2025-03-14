// expression.go
package parser

import (
	"fmt"

	"github.com/unLomTrois/gock3/pkg/ast"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
	"github.com/unLomTrois/gock3/pkg/tokens"
)

// fileBlock parses the entire file and constructs the AST's FileBlock.
func (p *Parser) fileBlock() *ast.FileBlock {
	if p.currentToken == nil {
		// Empty file.
		return &ast.FileBlock{Values: []*ast.Field{}, Loc: tokens.Loc{}}
	}
	loc := p.loc
	fields := p.FieldList()
	return &ast.FileBlock{Values: fields, Loc: *loc}
}

// FieldList parses a list of fields until a stop token is encountered.
func (p *Parser) FieldList(stopLookahead ...tokens.TokenType) []*ast.Field {
	var fields []*ast.Field

	for p.currentToken != nil {
		// Check for stop tokens to end the field list
		if len(stopLookahead) > 0 && p.currentToken.Type == stopLookahead[0] {
			break
		}

		switch p.currentToken.Type {
		case tokens.NEXTLINE:
			p.skipTokens(tokens.NEXTLINE)
			// continue
		case tokens.WORD, tokens.DATE, tokens.NUMBER:
			if field := p.Field(); field != nil {
				fields = append(fields, field)
			}
		default:
			// Handle unexpected token
			errMsg := fmt.Sprintf(errFieldListUnexpectedToken, p.currentToken.Value, p.currentToken.Type)
			err := report.FromToken(p.currentToken, severity.Error, errMsg)
			p.AddError(err)

			if fullPath, errPath := p.currentToken.Loc.Fullpath(); errPath == nil {
				_ = fullPath // Optionally, log or use the full path.
			} else {
				p.AddError(report.FromLoc(*p.loc, severity.Error, errPath.Error()))
			}

			if _, recovered := p.synchronize(FieldListRecovery); !recovered {
				return fields
			}
		}
	}
	return fields
}

// Field parses a single field and returns the corresponding AST node.
func (p *Parser) Field() *ast.Field {
	switch p.currentToken.Type {
	case tokens.WORD, tokens.DATE, tokens.NUMBER:
		return p.ExpressionNode()
	default:
		errMsg := fmt.Sprintf(errUnexpectedFieldToken, p.currentToken.Value, p.currentToken.Type)
		err := report.FromToken(p.currentToken, severity.Error, errMsg)
		p.AddError(err)
		p.synchronize(FieldRecovery)
		return nil
	}
}

// ExpressionNode parses an expression and returns an AST Field node.
func (p *Parser) ExpressionNode() *ast.Field {
	key := p.Key()
	if key == nil {
		return nil
	}

	operator := p.Operator()
	if operator == nil {
		return nil
	}

	value := p.Value()
	if value == nil {
		return nil
	}

	return &ast.Field{
		Key:      key,
		Operator: operator,
		Value:    value,
	}
}

// Key parses the key of a field and returns the corresponding token.
func (p *Parser) Key() *tokens.Token {
	if p.currentToken == nil {
		errMsg := "Expected a key, but reached end of input"
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	switch p.currentToken.Type {
	case tokens.WORD, tokens.DATE, tokens.NUMBER:
		return p.Expect(tokens.WORD, tokens.DATE, tokens.NUMBER)
	default:
		errMsg := fmt.Sprintf("Expected a key (WORD, DATE, or NUMBER), but found %q of type %q", p.currentToken.Value, p.currentToken.Type)
		err := report.FromToken(p.currentToken, severity.Error, errMsg)
		p.AddError(err)
		p.synchronize(KeyRecovery)
		return nil
	}
}

// Operator parses the operator of a field and returns the corresponding token.
func (p *Parser) Operator() *tokens.Token {
	if p.currentToken == nil {
		errMsg := errOperatorExpectedEOF
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	switch p.currentToken.Type {
	case tokens.QUESTION_EQUALS, tokens.EQUALS, tokens.COMPARISON:
		return p.Expect(p.currentToken.Type)
	default:
		errMsg := fmt.Sprintf(errOperatorUnexpectedToken, p.currentToken.Value, p.currentToken.Type)
		err := report.FromToken(p.currentToken, severity.Error, errMsg)
		p.AddError(err)
		p.synchronize(ValueRecovery)
		return nil
	}
}

// Value parses the value of a field and returns the corresponding AST node.
func (p *Parser) Value() ast.BlockOrValue {
	if p.currentToken == nil {
		errMsg := errValueExpectedEOF
		err := report.FromLoc(*p.loc, severity.Error, errMsg)
		p.AddError(err)
		return nil
	}

	switch p.currentToken.Type {
	case tokens.NEXTLINE:
		p.Expect(tokens.NEXTLINE)
		return p.EmptyValue()
	case tokens.WORD, tokens.NUMBER, tokens.QUOTED_STRING, tokens.BOOL, tokens.DATE:
		return p.Literal()
	case tokens.START:
		return p.Block()
	default:
		errMsg := fmt.Sprintf(errValueUnexpectedToken, p.currentToken.Value, p.currentToken.Type)
		err := report.FromToken(p.currentToken, severity.Error, errMsg)
		p.AddError(err)
		p.synchronize(ValueRecovery)
		return nil
	}
}

// EmptyValue returns an empty value AST node.
func (p *Parser) EmptyValue() ast.BlockOrValue {
	return ast.EmptyValue{
		Loc: *p.loc,
	}
}

// Literal parses a literal token and returns it.
func (p *Parser) Literal() *tokens.Token {
	if p.currentToken == nil {
		err := report.FromLoc(*p.loc, severity.Error, errLiteralExpectedEOF)
		p.AddError(err)
		return nil
	}

	switch p.currentToken.Type {
	case tokens.WORD, tokens.NUMBER, tokens.BOOL, tokens.DATE:
		return p.Expect(p.currentToken.Type)
	case tokens.QUOTED_STRING:
		return p.unquoteExpect(tokens.QUOTED_STRING)
	default:
		errMsg := fmt.Sprintf(errLiteralUnexpectedToken, p.currentToken.Value, p.currentToken.Type)
		err := report.FromToken(p.currentToken, severity.Error, errMsg)
		p.AddError(err)
	}

	if token, recovered := p.synchronize(LiteralRecovery); recovered {
		// Try parsing literal again from recovery point
		if isLiteralType(token.Type) {
			return p.Literal()
		}
		// If recovered to a non-literal token, give up
		errMsg := fmt.Sprintf(errRecoveredNonLiteralToken, token.Value, token.Type)
		err := report.FromToken(token, severity.Error, errMsg)
		p.AddError(err)
	}

	return nil
}
