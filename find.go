package the_platinum_searcher

import (
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
)

type Finder struct {
	Out    chan *GrepParams
	Option *option.Option
}

func (f *Finder) Find(root string, pattern *pattern.Pattern) {
	if f.Option.SearchStream {
		f.findStream(pattern)
	} else {
		f.findFile(root, pattern)
	}
}

func (f *Finder) findStream(pattern *pattern.Pattern) {
	// TODO: File type is fixed in ASCII because it can not determine the character code.
	f.Out <- &GrepParams{"", ASCII, pattern}
	close(f.Out)
}

func (f *Finder) findFile(root string, pattern *pattern.Pattern) {
	if f.Option.NoPtIgnore == false {
		f.addHomePtIgnore()
	}

	if f.Option.NoGlobalGitIgnore == false {
		f.addGlobalGitIgnore()
	}

	Walk(root, f.Option.Ignore, f.Option.Follow, func(path string, info *FileInfo, depth int, ig Ignore, err error) (error, Ignore) {
		if info.IsDir() {
			if depth > f.Option.Depth+1 {
				return filepath.SkipDir, ig
			}
			//Current Directory skipping should be checked first before loading ignores
			//within this directory
			if !isRoot(depth) && isHidden(info.Name()) {
				return filepath.SkipDir, ig
			} else {
				for _, p := range ig.Patterns {
					val, _ := filepath.Match(p, filepath.Base(path)+"/")
					if val {
						return filepath.SkipDir, ig
					}
				}
			}
			ig.Patterns = append(ig.Patterns, IgnorePatterns(path, f.Option.VcsIgnores())...)
			return nil, ig
		}
		if !info.follow && info.IsSymlink() {
			return nil, ig
		}
		if !isRoot(depth) && isHidden(info.Name()) {
			return nil, ig
		}
		for _, p := range ig.Patterns {
			val, _ := filepath.Match(p, filepath.Base(path))
			if val {
				return nil, ig
			}
		}
		if pattern.FileRegexp != nil && !pattern.FileRegexp.MatchString(path) {
			return nil, ig
		}
		fileType := ""
		if f.Option.FilesWithRegexp == "" {
			fileType = IdentifyType(path)
			if fileType == ERROR || fileType == BINARY {
				return nil, ig
			}
		}
		f.Out <- &GrepParams{path, fileType, pattern}
		return nil, ig
	})
	close(f.Out)
}

type WalkFunc func(path string, info *FileInfo, depth int, ig Ignore, err error) (error, Ignore)

func Walk(root string, ignorePatterns []string, follow bool, walkFn WalkFunc) error {
	info, err := os.Lstat(root)
	fileInfo := newFileInfo(root, info, follow)
	if err != nil {
		walkError, _ := walkFn(root, fileInfo, 1, Ignore{}, err)
		return walkError
	}
	return walk(root, fileInfo, 1, Ignore{Patterns: ignorePatterns}, walkFn)
}

func walkOnGoRoutine(path string, info *FileInfo, notify chan int, depth int, parentIgnore Ignore, walkFn WalkFunc) {
	walk(path, info, depth, parentIgnore, walkFn)
	notify <- 0
}

func walk(path string, info *FileInfo, depth int, parentIgnore Ignore, walkFn WalkFunc) error {
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

func (f *Finder) addHomePtIgnore() {
	homeDir := setHomeDir()
	if homeDir != "" {
		f.Option.Ignore = append(f.Option.Ignore, IgnorePatterns(homeDir, []string{".ptignore"})...)
	}
}

func setHomeDir() string {
	usr, err := user.Current()
	var homeDir string
	if err == nil {
		homeDir = usr.HomeDir
	} else {
		// Maybe it's cross compilation without cgo support. (darwin, unix)
		homeDir = os.Getenv("HOME")
	}
	return homeDir
}

func (f *Finder) addGlobalGitIgnore() {
	homeDir := setHomeDir()
	if homeDir != "" {
		globalIgnore := globalGitIgnore()
		if globalIgnore != "" {
			f.Option.Ignore = append(f.Option.Ignore, IgnorePatterns(homeDir, []string{globalIgnore})...)
		}
	}
}

func globalGitIgnore() string {
	gitCmd, err := exec.LookPath("git")
	if err != nil {
		return ""
	}

	file, err := exec.Command(gitCmd, "config", "--get", "core.excludesfile").Output()
	var filename string
	if err != nil {
		filename = ""
	} else {
		filename = strings.TrimSpace(filepath.Base(string(file)))
	}
	return filename
}
