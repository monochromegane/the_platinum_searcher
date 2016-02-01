package the_platinum_searcher

import (
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/monochromegane/go-gitignore"
)

type find struct {
	out  chan string
	opts Option
}

func (f find) start(roots []string, regexp *regexp.Regexp) {
	defer close(f.out)

	if f.opts.SearchOption.SearchStream {
		f.out <- ""
	} else if len(roots) == 1 {
		f.findFile(roots[0], regexp)
	} else {
		f.findFiles(roots, regexp)
	}
}

func (f find) findFiles(roots []string, reg *regexp.Regexp) {
	wg := &sync.WaitGroup{}
	wg.Add(len(roots))
	for _, r := range roots {
		go func(root string, reg *regexp.Regexp, wg *sync.WaitGroup) {
			defer wg.Done()
			f.findFile(root, reg)
		}(r, reg, wg)
	}
	wg.Wait()
}

func (f find) findFile(root string, regexp *regexp.Regexp) {
	var ignores ignoreMatchers

	// add ignores from ignore option.
	if len(f.opts.SearchOption.Ignore) > 0 {
		ignores = append(ignores, gitignore.NewGitIgnoreFromReader(
			root,
			strings.NewReader(strings.Join(f.opts.SearchOption.Ignore, "\n")),
		))
	}

	// add global gitignore.
	if f.opts.SearchOption.GlobalGitIgnore {
		if ignore := globalGitIgnore(root); ignore != nil {
			ignores = append(ignores, ignore)
		}
	}

	// add home ptignore.
	if f.opts.SearchOption.HomePtIgnore {
		if ignore := homePtIgnore(root); ignore != nil {
			ignores = append(ignores, ignore)
		}
	}

	followed := f.opts.SearchOption.Follow
	concurrentWalk(root, ignores, followed, func(path string, info fileInfo, depth int, ignores ignoreMatchers) (ignoreMatchers, error) {
		if info.isDir(followed) {
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

		if regexp != nil && !regexp.MatchString(path) {
			return ignores, nil
		}

		f.out <- path
		return ignores, nil
	})
}

func isHidden(name string) bool {
	if name == "." || name == ".." {
		return false
	}
	return len(name) > 1 && name[0] == '.'
}
