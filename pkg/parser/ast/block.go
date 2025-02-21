package ast

import (
	"github.com/unLomTrois/gock3/pkg/lexer/tokens"
)

// BlockOrValue represents an element in the AST that can be either a block or a literal value.
type BlockOrValue interface {
	IsBlockOrValue()
}

// Block represents a block of fields or tokens in the AST.
type Block interface {
	BlockOrValue
	IsBlock()
}

// FileBlock is an alias for FieldBlock representing the top-level block.
type FileBlock = FieldBlock

// FieldBlock represents a block that contains a list of fields.
type FieldBlock struct {
	Values []*Field   `json:"fields"`
	Loc    tokens.Loc `json:"-"`
}

func (fb *FieldBlock) IsBlock()        {}
func (fb *FieldBlock) IsBlockOrValue() {}

// GetValues returns the list of fields in the block.
func (fb *FieldBlock) GetValues() []*Field {
	return fb.Values
}

// GetField returns the first field that matches the given key.
func (fb *FieldBlock) GetField(key string) *Field {
	for _, field := range fb.Values {
		if field.Key.Value == key {
			return field
		}
	}
	return nil
}

// GetFieldValue returns the literal token value of the field with the given key, if present.
func (fb *FieldBlock) GetFieldValue(key string) *tokens.Token {
	field := fb.GetField(key)
	if field == nil {
		return nil
	}
	if token, ok := field.Value.(*tokens.Token); ok {
		return token
	}
	return nil
}

// GetFields returns all fields that match the given key.
func (fb *FieldBlock) GetFields(key string) []*Field {
	var res []*Field
	for _, field := range fb.Values {
		if field.Key.Value == key {
			res = append(res, field)
		}
	}
	return res
}

// GetFieldsValues returns the literal token values for all fields with the given key.
func (fb *FieldBlock) GetFieldsValues(key string) []*tokens.Token {
	fields := fb.GetFields(key)
	res := make([]*tokens.Token, 0, len(fields))
	for _, field := range fields {
		if token, ok := field.Value.(*tokens.Token); ok {
			res = append(res, token)
		}
	}
	return res
}

// GetFieldList returns a list of tokens if the field with the given key contains a TokenBlock.
func (fb *FieldBlock) GetFieldList(key string) []*tokens.Token {
	field := fb.GetField(key)
	if field == nil {
		return nil
	}

	if tb, ok := field.Value.(*TokenBlock); ok {
		return tb.Values
	}
	return nil
}

// GetFieldBlock returns the FieldBlock for the field with the given key, if it exists.
func (fb *FieldBlock) GetFieldBlock(key string) *FieldBlock {
	field := fb.GetField(key)
	if field == nil {
		return nil
	}
	if block, ok := field.Value.(*FieldBlock); ok {
		return block
	}
	return nil
}

// GetTokenBlock returns the TokenBlock for the field with the given key, if it exists.
func (fb *FieldBlock) GetTokenBlock(key string) *TokenBlock {
	field := fb.GetField(key)
	if field == nil {
		return nil
	}
	if block, ok := field.Value.(*TokenBlock); ok {
		return block
	}
	return nil
}

// TokenBlock represents a block that contains a list of literal tokens.
type TokenBlock struct {
	Values []*tokens.Token `json:"tokens"`
}

func (tb *TokenBlock) IsBlock()        {}
func (tb *TokenBlock) IsBlockOrValue() {}

// EmptyValue represents an empty value in the AST.
type EmptyValue struct {
	Loc tokens.Loc `json:"-"`
}

func (ev EmptyValue) IsBlockOrValue() {}
