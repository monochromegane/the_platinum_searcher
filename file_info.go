package the_platinum_searcher

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type fileInfo struct {
	path, name string
	typ        os.FileMode
}

func (f *fileInfo) isDir(follow bool) bool {
	if follow && f.isSymlink() {
		if _, err := ioutil.ReadDir(filepath.Join(f.path, f.name)); err == nil {
			return true
		} else {
			return false
		}
	} else {
		return f.typ&os.ModeDir == os.ModeDir
	}
}

func (f *fileInfo) isSymlink() bool {
	return f.typ&os.ModeSymlink == os.ModeSymlink
}

func (f *fileInfo) isNamedPipe() bool {
	return f.typ&os.ModeNamedPipe == os.ModeNamedPipe
}

func newFileInfo(path, name string, typ os.FileMode) *fileInfo {
	return &fileInfo{
		path: path,
		name: name,
		typ:  typ,
	}
}
