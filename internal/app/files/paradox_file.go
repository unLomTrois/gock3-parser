package files

import (
	"os"
	"path/filepath"
)

// FileKind represents the type of Paradox file (vanilla or mod).
type FileKind uint8

const (
	// Vanilla file
	Vanilla FileKind = iota
	// Mod file
	Mod
)

type ParadoxFile interface {
	FullPath() string
	FileName() string
	Kind() FileKind
	PathIdx() *PathTableIndex
	StoreInPathTable() *PathTableIndex
}

type ParadoxTxtFile struct {
	// The full filesystem path of this file
	fullpath string
	// Whether it's a vanilla or mod file
	kind FileKind
	// Index into the PathTable (optional, using *PathTableIndex to allow nil)
	idx *PathTableIndex
}

// NewParadoxTxtFile is the constructor for ParadoxFile.
// Ensures the path is valid and not empty.
func NewParadoxTxtFile(fullpath string, kind FileKind) *ParadoxTxtFile {
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		panic("Invalid path: path does not exist")
	}

	return &ParadoxTxtFile{
		fullpath: fullpath,
		kind:     kind,
		idx:      nil,
	}
}

// Kind returns the file kind (vanilla or mod).
func (file *ParadoxTxtFile) Kind() FileKind {
	return file.kind
}

// FullPath returns the full filesystem path.
func (file *ParadoxTxtFile) FullPath() string {
	return file.fullpath
}

// FileName returns the file name, ensuring it's not empty.
func (file *ParadoxTxtFile) FileName() string {
	return filepath.Base(file.fullpath)
}

// StoreInPathTable stores the file in the PathTable and returns the index.
func (file *ParadoxTxtFile) StoreInPathTable() *PathTableIndex {
	if file.idx != nil {
		return file.idx
	}
	file.idx = PATHTABLE.Store(file.fullpath)
	return file.idx
}

// PathIdx returns the index into the PathTable if it exists, otherwise nil.
func (file *ParadoxTxtFile) PathIdx() *PathTableIndex {
	return file.idx
}
