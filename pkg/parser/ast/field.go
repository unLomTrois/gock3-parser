package ast

import "github.com/unLomTrois/gock3/pkg/lexer/tokens"

type Field struct {
	Key      *tokens.Token `json:"key"`
	Operator *tokens.Token `json:"operator"`
	Value    BlockOrValue  `json:"value"`
}
