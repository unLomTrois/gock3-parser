package ast

// AST represents the abstract syntax tree for a parsed Paradox file.
type AST struct {
	Filename string     `json:"filename"`
	Fullpath string     `json:"fullpath"`
	Block    *FileBlock `json:"data"`
}
