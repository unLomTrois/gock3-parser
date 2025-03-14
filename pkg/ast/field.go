package ast

import "github.com/unLomTrois/gock3/pkg/tokens"

// Field represents a single key-operator-value triple in the AST.
type Field struct {
	Key      *tokens.Token `json:"key"`
	Operator *tokens.Token `json:"operator"`
	Value    BlockOrValue  `json:"value"`
}
