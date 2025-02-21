package files

import (
	"os"
	"path/filepath"
)

type FileKind uint8

const (
	Vanilla FileKind = iota
	Mod
)

type ParadoxFile struct {
	// The full filesystem path of this file
	fullpath string
	// Index into the PathTable (optional, using *PathTableIndex to allow nil)
	idx *PathTableIndex
	// Whether it's a vanilla or mod file
	kind FileKind
}

// NewParadoxFile is the constructor for ParadoxFile.
// Ensures the path is valid and not empty.
func NewParadoxFile(fullpath string, kind FileKind) *ParadoxFile {
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		panic("Invalid path: path does not exist")
	}

	return &ParadoxFile{
		fullpath: fullpath,
		kind:     kind,
		idx:      nil,
	}
}

// Kind returns the file kind (vanilla or mod).
func (fe *ParadoxFile) Kind() FileKind {
	return fe.kind
}

// FullPath returns the full filesystem path.
func (fe *ParadoxFile) FullPath() string {
	return fe.fullpath
}

// FileName returns the file name, ensuring it's not empty.
func (fe *ParadoxFile) FileName() string {
	return filepath.Base(fe.fullpath)
}

func (fe *ParadoxFile) StoreInPathTable() *PathTableIndex {
	if fe.idx != nil {
		return fe.idx
	}
	fe.idx = PATHTABLE.Store(fe.fullpath)
	return fe.idx
}

// PathIdx returns the index into the PathTable if it exists, otherwise nil.
func (fe *ParadoxFile) PathIdx() *PathTableIndex {
	return fe.idx
}
