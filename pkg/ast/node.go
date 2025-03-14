package ast

import "github.com/unLomTrois/gock3/pkg/tokens"

// Node represents any node in the AST.
// All AST elements should implement this interface.
type Node interface {
	// GetLoc returns the location information for this node
	GetLoc() tokens.Loc
}
