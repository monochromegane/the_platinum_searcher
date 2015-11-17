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
	var ignores ignoreMatchers
	concurrentWalk(root, ignores, func(path string, info os.FileInfo, ignores ignoreMatchers) (ignoreMatchers, error) {
		if info.IsDir() {
			if info.Name() == ".git" {
				return ignores, filepath.SkipDir
			}

			if ignores.Match(path, true) {
				return ignores, filepath.SkipDir
			}

			ignores = append(ignores, newIgnoreMatchers(path, []string{".gitignore"})...)
			return ignores, nil
		}
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			return ignores, nil
		}

		if ignores.Match(path, false) {
			return ignores, nil
		}

		f.out <- path
		return ignores, nil
	})
	close(f.out)
}
