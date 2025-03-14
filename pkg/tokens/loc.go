package tokens

import (
	"fmt"
	"path/filepath"

	"github.com/unLomTrois/gock3/pkg/files"
)

// Loc представляет позицию токена в файле
type Loc struct {
	idx    files.PathTableIndex `json:"-"`
	Line   uint32               `json:"line"`
	Column uint16               `json:"column"`
	kind   files.FileKind       `json:"-"`
}

// Filename возвращает имя файла из Loc
func (loc *Loc) Filename() (string, error) {
	path, err := files.PATHTABLE.LookupFullpath(loc.idx)
	if err != nil {
		return "", err
	}
	filename := filepath.Base(path)
	return filename, nil
}

// Pathname возвращает относительный путь из Loc
func (loc *Loc) Pathname() (string, error) {
	path, err := files.PATHTABLE.LookupFullpath(loc.idx)
	if err != nil {
		return "", err
	}
	return path, nil
}

// Fullpath возвращает полный путь из Loc
func (loc *Loc) Fullpath() (string, error) {
	fullpath, err := files.PATHTABLE.LookupFullpath(loc.idx)
	if err != nil {
		return "", err
	}

	fullpathWithLoc := fmt.Sprintf("%s:%d:%d", fullpath, loc.Line, loc.Column)

	return fullpathWithLoc, nil
}

// SameFile проверяет, ссылается ли Loc на тот же файл, что и другой Loc
func (loc *Loc) SameFile(other Loc) bool {
	return loc.idx == other.idx
}

// LocFromParadoxFile создает Loc из ParadoxFile
func LocFromParadoxFile(file files.ParadoxFile) *Loc {
	idx := file.StoreInPathTable()
	return &Loc{
		idx:    *idx,
		kind:   file.Kind(),
		Line:   1,
		Column: 1,
	}
}

func (loc *Loc) GetIdx() files.PathTableIndex {
	return loc.idx
}
