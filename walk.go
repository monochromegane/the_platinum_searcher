package the_platinum_searcher

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type walkFunc func(path string, info fileInfo, depth int, ignores ignoreMatchers) (ignoreMatchers, error)

func concurrentWalk(root string, ignores ignoreMatchers, walkFn walkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return err
	}
	sem := make(chan struct{}, 16)
	return walk(root, newFileInfo(root, info), 1, ignores, walkFn, sem)
}

func walk(path string, info fileInfo, depth int, parentIgnores ignoreMatchers, walkFn walkFunc, sem chan struct{}) error {
	ignores, walkError := walkFn(path, info, depth, parentIgnores)
	if walkError != nil {
		if info.IsDir() && walkError == filepath.SkipDir {
			return nil
		}
		return walkError
	}

	if !info.IsDir() {
		return nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	depth++
	wg := &sync.WaitGroup{}
	for _, file := range files {
		f := newFileInfo(path, file)
		select {
		case sem <- struct{}{}:
			wg.Add(1)
			go func(path string, file fileInfo, depth int, ignores ignoreMatchers, wg *sync.WaitGroup) {
				defer wg.Done()
				defer func() { <-sem }()
				walk(path, file, depth, ignores, walkFn, sem)
			}(filepath.Join(path, file.Name()), f, depth, ignores, wg)
		default:
			walk(filepath.Join(path, file.Name()), f, depth, ignores, walkFn, sem)
		}
	}
	wg.Wait()
	return nil
}
