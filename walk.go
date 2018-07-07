package the_platinum_searcher

import (
	"os"
	"path/filepath"
	"sync"
)

type walkFunc func(path string, info *fileInfo, depth int, ignores ignoreMatchers) (ignoreMatchers, error)

func concurrentWalk(root string, ignores ignoreMatchers, followed bool, walkFn walkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return err
	}
	sem := make(chan struct{}, 16)
	return walk(root, newFileInfo(root, info.Name(), info.Mode()), 1, ignores, followed, walkFn, sem)
}

func walk(path string, info *fileInfo, depth int, parentIgnores ignoreMatchers, followed bool, walkFn walkFunc, sem chan struct{}) error {
	ignores, walkError := walkFn(path, info, depth, parentIgnores)
	if walkError != nil {
		if info.isDir(false) && walkError == filepath.SkipDir {
			return nil
		}
		return walkError
	}

	if !info.isDir(followed) {
		return nil
	}

	files, err := readDir(path)
	if err != nil {
		return err
	}

	depth++
	wg := &sync.WaitGroup{}
	for _, file := range files {
		select {
		case sem <- struct{}{}:
			wg.Add(1)
			go func(path string, file *fileInfo, depth int, ignores ignoreMatchers, wg *sync.WaitGroup) {
				defer func() {
					wg.Done()
					<-sem
				}()
				walk(path, file, depth, ignores, followed, walkFn, sem)
			}(path+string(os.PathSeparator)+file.name, file, depth, ignores, wg)
		default:
			walk(path+string(os.PathSeparator)+file.name, file, depth, ignores, followed, walkFn, sem)
		}
	}
	wg.Wait()
	return nil
}
