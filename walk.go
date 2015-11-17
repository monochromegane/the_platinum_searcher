package the_platinum_searcher

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type walkFunc func(path string, info os.FileInfo, ignores ignoreMatchers) (ignoreMatchers, error)

func concurrentWalk(root string, ignores ignoreMatchers, walkFn walkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return err
	}
	sem := make(chan struct{}, 16)
	return walk(root, info, ignores, walkFn, sem)
}

func walk(path string, info os.FileInfo, parentIgnores ignoreMatchers, walkFn walkFunc, sem chan struct{}) error {
	ignores, walkError := walkFn(path, info, parentIgnores)
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

	wg := &sync.WaitGroup{}
	for _, file := range files {
		select {
		case sem <- struct{}{}:
			wg.Add(1)
			go func(path string, file os.FileInfo, ignores ignoreMatchers, wg *sync.WaitGroup) {
				defer wg.Done()
				defer func() { <-sem }()
				walk(path, file, ignores, walkFn, sem)
			}(filepath.Join(path, file.Name()), file, ignores, wg)
		default:
			walk(filepath.Join(path, file.Name()), file, ignores, walkFn, sem)
		}
	}
	wg.Wait()
	return nil
}
