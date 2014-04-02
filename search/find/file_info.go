package find

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

func (self *FileInfo) IsDir() bool {
	if self.follow && self.IsSymlink() {
		if _, err := ioutil.ReadDir(filepath.Join(self.path, self.FileInfo.Name())); err == nil {
			return true
		} else {
			return false
		}
	} else {
		return self.FileInfo.IsDir()
	}
}

func (self *FileInfo) IsSymlink() bool {
	return self.FileInfo.Mode()&os.ModeSymlink == os.ModeSymlink
}

func newFileInfo(path string, info os.FileInfo, follow bool) *FileInfo {
	return &FileInfo{path, info, follow}
}
