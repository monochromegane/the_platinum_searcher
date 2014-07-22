package the_platinum_searcher

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileInfo struct {
	path string
	os.FileInfo
	follow bool
}

func (f *FileInfo) IsDir() bool {
	if f.follow && f.IsSymlink() {
		if _, err := ioutil.ReadDir(filepath.Join(f.path, f.FileInfo.Name())); err == nil {
			return true
		} else {
			return false
		}
	} else {
		return f.FileInfo.IsDir()
	}
}

func (f *FileInfo) IsSymlink() bool {
	return f.FileInfo.Mode()&os.ModeSymlink == os.ModeSymlink
}

func newFileInfo(path string, info os.FileInfo, follow bool) *FileInfo {
	return &FileInfo{path, info, follow}
}
