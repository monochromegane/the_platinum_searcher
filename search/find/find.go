package find

import (
	"github.com/monochromegane/the_platinum_searcher/search/file"
	"github.com/monochromegane/the_platinum_searcher/search/grep"
	"github.com/monochromegane/the_platinum_searcher/search/ignore"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Finder struct {
	Out    chan *grep.Params
	Option *option.Option
}

func (self *Finder) Find(root string, pattern *pattern.Pattern) {
	if self.Option.SearchStream {
		self.findStream(pattern)
	} else {
		self.findFile(root, pattern)
	}
}

func (self *Finder) findStream(pattern *pattern.Pattern) {
	// TODO: File type is fixed in ASCII because it can not determine the character code.
	self.Out <- &grep.Params{"", file.ASCII, pattern}
	close(self.Out)
}

func (self *Finder) findFile(root string, pattern *pattern.Pattern) {
	if self.Option.NoPtIgnore == false {
		usr, err := user.Current()
		var homeDir string
		if err == nil {
			homeDir = usr.HomeDir
		} else {
			// Maybe it's cross compilation without cgo support. (darwin, unix)
			homeDir = os.Getenv("HOME")
		}
		if homeDir != "" {
			self.Option.Ignore = append(self.Option.Ignore, ignore.IgnorePatterns(homeDir, []string{".ptignore"})...)
		}
	}
	Walk(root, self.Option.Ignore, self.Option.Follow, func(path string, info *FileInfo, depth int, ig ignore.Ignore, err error) (error, ignore.Ignore) {
		if info.IsDir() {
			if depth > self.Option.Depth+1 {
				return filepath.SkipDir, ig
			}
			//Current Directory skipping should be checked first before loading ignores
			//within this directory
			if !isRoot(depth) && isHidden(info.Name()) {
				return filepath.SkipDir, ig
			} else {
				for _, p := range ig.Patterns {
					val, _:= filepath.Match(p, filepath.Base(path) + "/") 
					if val {
						return filepath.SkipDir, ig
					}
				}
			}
			ig.Patterns = append(ig.Patterns, ignore.IgnorePatterns(path, self.Option.VcsIgnores())...)
			return nil, ig
		}
		if !info.follow && info.IsSymlink() {
			return nil, ig
		}
		if !isRoot(depth) && isHidden(info.Name()) {
			return nil, ig
		}
		for _, p := range ig.Patterns {
			val, _:= filepath.Match(p, filepath.Base(path)) 
			if val {
				return nil, ig
			}
		}
		if pattern.FileRegexp != nil && !pattern.FileRegexp.MatchString(path) {
			return nil, ig
		}
		fileType := ""
		if self.Option.FilesWithRegexp == "" {
			fileType = file.IdentifyType(path)
			if fileType == file.ERROR || fileType == file.BINARY {
				return nil, ig
			}
		}
		self.Out <- &grep.Params{path, fileType, pattern}
		return nil, ig
	})
	close(self.Out)
}

type WalkFunc func(path string, info *FileInfo, depth int, ig ignore.Ignore, err error) (error, ignore.Ignore)

func Walk(root string, ignorePatterns []string, follow bool, walkFn WalkFunc) error {
	info, err := os.Lstat(root)
	fileInfo := newFileInfo(root, info, follow)
	if err != nil {
		walkError, _ := walkFn(root, fileInfo, 1, ignore.Ignore{}, err)
		return walkError
	}
	return walk(root, fileInfo, 1, ignore.Ignore{Patterns: ignorePatterns}, walkFn)
}

func walkOnGoRoutine(path string, info *FileInfo, notify chan int, depth int, parentIgnore ignore.Ignore, walkFn WalkFunc) {
	walk(path, info, depth, parentIgnore, walkFn)
	notify <- 0
}

func walk(path string, info *FileInfo, depth int, parentIgnore ignore.Ignore, walkFn WalkFunc) error {
	err, ig := walkFn(path, info, depth, parentIgnore, nil)
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
	notify := make(chan int, len(list))
	for _, l := range list {
		fileInfo := newFileInfo(path, l, info.follow)
		if isDirectRoot(depth) {
			go walkOnGoRoutine(filepath.Join(path, fileInfo.Name()), fileInfo, notify, depth, ig, walkFn)

		} else {
			walk(filepath.Join(path, fileInfo.Name()), fileInfo, depth, ig, walkFn)
		}
	}
	if isDirectRoot(depth) {
		for i := 0; i < cap(notify); i++ {
			<-notify
		}
	}
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
