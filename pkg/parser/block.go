// block.go
package parser

import (
	"fmt"

	"github.com/unLomTrois/gock3/pkg/lexer/tokens"
	"github.com/unLomTrois/gock3/pkg/parser/ast"
	"github.com/unLomTrois/gock3/pkg/report"
	"github.com/unLomTrois/gock3/pkg/report/severity"
)

// Block parses a block and returns the corresponding AST node.
func (p *Parser) Block() ast.Block {
	// Expect the start of a block.
	p.Expect(tokens.START)
	loc := *p.loc

	// Handle an empty block.
	if p.currentToken.Type == tokens.END {
		p.Expect(tokens.END)
		return &ast.FieldBlock{Values: []*ast.Field{}, Loc: loc}
	}

	var block ast.Block

	for p.currentToken != nil && p.currentToken.Type != tokens.END {
		switch p.currentToken.Type {
		case tokens.NEXTLINE:
			p.skipTokens(tokens.NEXTLINE)
			continue
		case tokens.WORD, tokens.DATE, tokens.QUOTED_STRING, tokens.NUMBER:
			if p.isNextField() {
				block = p.FieldBlock(loc)
			} else {
				block = p.TokenBlock()
			}
		default:
			errorMsg := fmt.Sprintf(errBlockUnexpectedToken, p.currentToken.Value, p.currentToken.Type)
			err := report.FromToken(p.currentToken, severity.Error, errorMsg)
			p.AddError(err)
			p.synchronize(BlockRecovery)
			continue
		}
		// Once a block is parsed, exit the loop.
		break
	}

	// Expect the closing token for the block.
	p.Expect(tokens.END)
	return block
}

func (p *Parser) skipTokens(types ...tokens.TokenType) {
	for p.currentToken != nil {
		matched := false
		for _, t := range types {
			if p.currentToken.Type == t {
				p.Expect(t)
				matched = true
				break
			}
		}
		if !matched {
			break
		}
	}
}

// FieldBlock parses a block of fields and returns the corresponding AST node.
func (p *Parser) FieldBlock(loc tokens.Loc) *ast.FieldBlock {
	fields := p.FieldList(tokens.END)
	return &ast.FieldBlock{Values: fields, Loc: loc}
}

// TokenBlock parses a block of tokens and returns the corresponding AST node.
func (p *Parser) TokenBlock() *ast.TokenBlock {
	tokensList := p.TokenList(tokens.END)
	return &ast.TokenBlock{Values: tokensList}
}

// TokenList parses a list of tokens until a specified stop token is encountered.
func (p *Parser) TokenList(stopLookahead ...tokens.TokenType) []*tokens.Token {
	var tokensList []*tokens.Token
	for p.currentToken != nil {
		if len(stopLookahead) > 0 && p.currentToken.Type == stopLookahead[0] {
			break
		}
		switch p.currentToken.Type {
		case tokens.NEXTLINE:
			p.Expect(tokens.NEXTLINE)
			// continue
		case tokens.NUMBER, tokens.QUOTED_STRING, tokens.WORD:
			if token := p.Literal(); token != nil {
				tokensList = append(tokensList, token)
			}
		default:
			errMsg := fmt.Sprintf(errTokenListUnexpectedToken, p.currentToken.Value, p.currentToken.Type)
			err := report.FromToken(p.currentToken, severity.Error, errMsg)
			p.AddError(err)
			recoveryPoint := RecoveryPoint{
				TokenTypes: []tokens.TokenType{tokens.END, tokens.WORD, tokens.DATE},
				Context:    "TokenList",
			}
			if _, recovered := p.synchronize(recoveryPoint); !recovered {
				return tokensList
			}
		}
	}
	return tokensList
}
