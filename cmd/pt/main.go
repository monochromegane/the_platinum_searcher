package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		return nil
	})
}

type WalkFunc func(path string, info os.FileInfo, err error) error

func Walk(root string, walkFn WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	sem := make(chan struct{}, 256)
	return walk(root, info, walkFn, sem)
}

func walk(path string, info os.FileInfo, walkFn WalkFunc, sem chan struct{}) error {
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

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return walkFn(path, info, err)
	}

	wg := &sync.WaitGroup{}
	for _, file := range files {
		select {
		case sem <- struct{}{}:
			wg.Add(1)
			go func(path string, file os.FileInfo, wg *sync.WaitGroup) {
				defer wg.Done()
				defer func() { <-sem }()
				walk(path, file, walkFn, sem)
			}(filepath.Join(path, file.Name()), file, wg)
		default:
			walk(filepath.Join(path, file.Name()), file, walkFn, sem)
		}
	}
	wg.Wait()
	return nil
}
