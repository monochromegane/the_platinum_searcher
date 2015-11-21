package the_platinum_searcher

import (
	"os"
	"path/filepath"
)

type find struct {
	out  chan string
	opts Option
}

func (f find) start(root string) {
	f.findFile(root)
}

func (f find) findFile(root string) {
	var ignores ignoreMatchers
	concurrentWalk(root, ignores, func(path string, info os.FileInfo, depth int, ignores ignoreMatchers) (ignoreMatchers, error) {
		if info.IsDir() {
			if depth > f.opts.SearchOption.Depth+1 {
				return ignores, filepath.SkipDir
			}
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
