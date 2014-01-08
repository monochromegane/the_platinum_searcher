package find

import (
	"github.com/monochromegane/the_platinum_searcher/search/file"
	"github.com/monochromegane/the_platinum_searcher/search/grep"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Finder struct {
	Out    chan *grep.Params
	Option *option.Option
}

func (self *Finder) Find(root, pattern string) {
	Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if len(info.Name()) > 1 && strings.Index(info.Name(), ".") == 0 {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
		fileType := file.IdentifyType(path)
		if fileType == file.BINARY {
			return nil
		}
		self.Out <- &grep.Params{path, pattern, fileType}
		return nil
	})
	close(self.Out)
}

type WalkFunc func(path string, info os.FileInfo, err error) error

func Walk(root string, walkFn WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	return walk(root, info, walkFn)
}

func walk(path string, info os.FileInfo, walkFn WalkFunc) error {
	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	list, err := ioutil.ReadDir(path)
	if err != nil {
		return walkFn(path, info, err)
	}

	for _, fileInfo := range list {
		err = walk(filepath.Join(path, fileInfo.Name()), fileInfo, walkFn)
		if err != nil {
			if !fileInfo.IsDir() || err != filepath.SkipDir {
				return err
			}
		}
	}
	return nil
}
