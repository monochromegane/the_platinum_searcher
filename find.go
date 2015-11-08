package the_platinum_searcher

import (
	"os"
	"path/filepath"
)

type find struct {
	out chan string
}

func (f find) start(root string) {
	f.findFile(root)
}

func (f find) findFile(root string) {
	concurrentWalk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}
		f.out <- path
		return nil
	})
	close(f.out)
}
