package ast

import (
	"testing"

	"github.com/unLomTrois/gock3/pkg/lexer/tokens"
)

// validateField performs common assertions for a field in a FileBlock.
func validateField(t *testing.T, field *Field) {
	if field.Key == nil {
		t.Error("expected field.Key to be non-nil")
	} else if field.Key.Type != tokens.WORD {
		t.Errorf("expected field.Key.Type to be %v, got %v", tokens.WORD, field.Key.Type)
	}

	if field.Operator == nil {
		t.Error("expected field.Operator to be non-nil")
	} else if field.Operator.Type != tokens.EQUALS {
		t.Errorf("expected field.Operator.Type to be %v, got %v", tokens.EQUALS, field.Operator.Type)
	}

	switch v := field.Value.(type) {
	case *tokens.Token:
		if v.Type != tokens.WORD {
			t.Errorf("expected token Value.Type to be %v, got %v", tokens.WORD, v.Type)
		}
		if v.Value == "" {
			t.Error("expected token Value to be non-empty")
		}
	case *FieldBlock:
		if len(v.Values) == 0 {
			t.Error("expected FieldBlock.Values to contain at least one Field")
		}
		for _, fb := range v.Values {
			if fb.Key == nil || fb.Operator == nil || fb.Value == nil {
				t.Error("FieldBlock contains nil Field components")
			}
		}
	case *TokenBlock:
		if len(v.Values) == 0 {
			t.Error("expected TokenBlock.Values to contain at least one Token")
		}
		for _, tb := range v.Values {
			if tb.Type != tokens.WORD {
				t.Errorf("expected TokenBlock.Token.Type to be %v, got %v", tokens.WORD, tb.Type)
			}
			if tb.Value == "" {
				t.Error("expected TokenBlock.Token.Value to be non-empty")
			}
		}
	default:
		t.Errorf("unexpected type for field.Value: %T", v)
	}
}

func TestFileBlock(t *testing.T) {
	tests := []struct {
		name      string
		fileBlock FileBlock
	}{
		{
			name: "Namespace and Event Blocks",
			fileBlock: FileBlock{
				Values: []*Field{
					{
						Key:      &tokens.Token{Value: "namespace", Type: tokens.WORD},
						Operator: &tokens.Token{Value: "=", Type: tokens.EQUALS},
						Value:    &tokens.Token{Value: "test", Type: tokens.WORD},
					},
					{
						Key: &tokens.Token{Value: "event", Type: tokens.WORD},
						Operator: &tokens.Token{
							Value: "=", Type: tokens.EQUALS,
						},
						Value: &FieldBlock{
							Values: []*Field{
								{
									Key: &tokens.Token{Value: "type", Type: tokens.WORD},
									Operator: &tokens.Token{
										Value: "=", Type: tokens.EQUALS,
									},
									Value: &tokens.Token{Value: "character_event", Type: tokens.WORD},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Color Token Block",
			fileBlock: FileBlock{
				Values: []*Field{
					{
						Key:      &tokens.Token{Value: "color", Type: tokens.WORD},
						Operator: &tokens.Token{Value: "=", Type: tokens.EQUALS},
						Value: &TokenBlock{
							Values: []*tokens.Token{
								{Value: "255", Type: tokens.WORD},
								{Value: "255", Type: tokens.WORD},
								{Value: "255", Type: tokens.WORD},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.fileBlock.Values) == 0 {
				t.Error("expected at least one field in the file block")
			}

			for _, field := range tt.fileBlock.Values {
				validateField(t, field)
			}
		})
	}
}
