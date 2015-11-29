package the_platinum_searcher

import "path/filepath"

type find struct {
	out  chan string
	opts Option
}

func (f find) start(root string) {
	f.findFile(root)
}

func (f find) findFile(root string) {
	var ignores ignoreMatchers
	concurrentWalk(root, ignores, func(path string, info fileInfo, depth int, ignores ignoreMatchers) (ignoreMatchers, error) {
		if info.isDir(f.opts.SearchOption.Follow) {
			if depth > f.opts.SearchOption.Depth+1 {
				return ignores, filepath.SkipDir
			}

			if !f.opts.SearchOption.Hidden && isHidden(info.Name()) {
				return ignores, filepath.SkipDir
			}

			if ignores.Match(path, true) {
				return ignores, filepath.SkipDir
			}

			if !f.opts.SearchOption.SkipVcsIgnore {
				ignores = append(ignores, newIgnoreMatchers(path, f.opts.SearchOption.VcsIgnore)...)
			}
			return ignores, nil
		}
		if !f.opts.SearchOption.Follow && info.isSymlink() {
			return ignores, nil
		}

		if info.isNamedPipe() {
			return ignores, nil
		}

		if !f.opts.SearchOption.Hidden && isHidden(info.Name()) {
			return ignores, filepath.SkipDir
		}

		if ignores.Match(path, false) {
			return ignores, nil
		}

		f.out <- path
		return ignores, nil
	})
	close(f.out)
}

func isHidden(name string) bool {
	return len(name) > 1 && name[0] == '.'
}
