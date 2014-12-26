package the_platinum_searcher

import (
	"io/ioutil"
	"os"
	"runtime"
	"sync"

	"path/filepath"
	"strings"
)

type find struct {
	Out    chan *GrepParams
	Option *Option
}

func Find(root string, pattern *Pattern, out chan *GrepParams, option *Option) {
	find := find{
		Out:    out,
		Option: option,
	}
	find.Start(root, pattern)
}

func (f *find) Start(root string, pattern *Pattern) {
	if f.Option.SearchStream {
		f.findStream(pattern)
	} else {
		f.findFile(root, pattern)
	}
}

func (f *find) findStream(pattern *Pattern) {
	// TODO: File type is fixed in ASCII because it can not determine the character code.
	f.Out <- &GrepParams{"", ASCII, pattern}
	close(f.Out)
}

func (f *find) findFile(root string, pattern *Pattern) {

	var ignores ignoreMatchers
	if f.Option.NoPtIgnore == false {
		if ignore := homePtIgnore(); ignore != nil {
			ignores = append(ignores, ignore)
		}
	}

	if f.Option.NoGlobalGitIgnore == false {
		if ignore := globalGitIgnore(); ignore != nil {
			ignores = append(ignores, ignore)
		}
	}

	ignores = append(ignores, genericIgnore(f.Option.Ignore))
	Walk(root, ignores, f.Option.Follow, func(path string, info *FileInfo, depth int, ignores ignoreMatchers, err error) (error, ignoreMatchers) {
		if info.IsDir() {
			if depth > f.Option.Depth+1 {
				return filepath.SkipDir, ignores
			}
			//Current Directory skipping should be checked first before loading ignores
			//within this directory
			if !isRoot(depth) && isHidden(info.Name()) {
				return filepath.SkipDir, ignores
			} else {
				if ignores.Match(path, info.IsDir()) {
					return filepath.SkipDir, ignores
				}
			}
			ignores = append(ignores, newIgnoreMatchers(path, f.Option.VcsIgnores())...)
			return nil, ignores
		}
		if !info.follow && info.IsSymlink() {
			return nil, ignores
		}
		if !isRoot(depth) && isHidden(info.Name()) {
			return nil, ignores
		}

		if ignores.Match(path, info.IsDir()) {
			return nil, ignores
		}

		if pattern.FileRegexp != nil && !pattern.FileRegexp.MatchString(path) {
			return nil, ignores
		}
		fileType := UNKNOWN
		if f.Option.FilesWithRegexp == "" {
			fileType = IdentifyType(path)
			if fileType == ERROR || fileType == BINARY {
				return nil, ignores
			}
		}
		f.Out <- &GrepParams{path, fileType, pattern}
		return nil, ignores
	})
	close(f.Out)
}

type WalkFunc func(path string, info *FileInfo, depth int, ignores ignoreMatchers, err error) (error, ignoreMatchers)

func Walk(root string, ignores ignoreMatchers, follow bool, walkFn WalkFunc) error {
	info, err := os.Lstat(root)
	fileInfo := newFileInfo(root, info, follow)
	if err != nil {
		walkError, _ := walkFn(root, fileInfo, 1, nil, err)
		return walkError
	}
	workerNum := runtime.NumCPU()
	pool := make(chan struct{}, workerNum)
	for i := 0; i < workerNum; i++ {
		pool <- struct{}{}
	}
	return walk(root, fileInfo, 1, ignores, walkFn, pool)
}

func walkOnGoRoutine(path string, info *FileInfo, depth int, parentIgnore ignoreMatchers, walkFn WalkFunc, pool chan struct{}) {
	walk(path, info, depth, parentIgnore, walkFn, pool)
	pool <- struct{}{}
}

func walk(path string, info *FileInfo, depth int, parentIgnores ignoreMatchers, walkFn WalkFunc, pool chan struct{}) error {
	err, ig := walkFn(path, info, depth, parentIgnores, nil)
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
		walkError, _ := walkFn(path, info, depth, ig, err)
		return walkError
	}

	depth++
	waiter := &sync.WaitGroup{}
	for _, l := range list {
		fileInfo := newFileInfo(path, l, info.follow)
		select {
		case <-pool:
			waiter.Add(1)
			go func() {
				walkOnGoRoutine(filepath.Join(path, fileInfo.Name()), fileInfo, depth, ig, walkFn, pool)
				waiter.Done()
			}()
		default:
			walk(filepath.Join(path, fileInfo.Name()), fileInfo, depth, ig, walkFn, pool)
		}
	}

	waiter.Wait()
	return nil
}

func isRoot(depth int) bool {
	return depth == 1
}

func isDirectRoot(depth int) bool {
	return depth == 2
}

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".") && len(name) > 1
}

func contains(path string, patterns *[]string) bool {
	for _, p := range *patterns {
		if p == path {
			return true
		}
	}
	return false
}
