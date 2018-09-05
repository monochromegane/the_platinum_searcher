package the_platinum_searcher

import (
	"os"
	"path/filepath"
)

type walkFunc func(path string, info *fileInfo, depth int, ignores ignoreMatchers) (ignoreMatchers, error)

func walkIn(root string, ignores ignoreMatchers, followed bool, walkFn walkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return err
	}
	return walk(root, newFileInfo(root, info.Name(), info.Mode()), 1, ignores, followed, walkFn)
}

func walk(path string, info *fileInfo, depth int, parentIgnores ignoreMatchers, followed bool, walkFn walkFunc) error {
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
	for _, file := range files {
		walk(path+string(os.PathSeparator)+file.name, file, depth, ignores, followed, walkFn)
	}
	return nil
}
