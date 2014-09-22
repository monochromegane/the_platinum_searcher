package the_platinum_searcher

import (
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
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

	var matchers ignoreMatchers
	if f.Option.NoPtIgnore == false {
		matchers = append(matchers, homePtIgnore())
	}

	if f.Option.NoGlobalGitIgnore == false {
		matchers = append(matchers, globalGitIgnore())
	}

	matchers = append(matchers, genericIgnore(f.Option.Ignore))
	optionIgnoreMatchers := []StringMatcher{genericIgnoreMatcher(f.Option.Ignore)}
	Walk(root, optionIgnoreMatchers, f.Option.Follow, func(path string, info *FileInfo, depth int, ig Ignore, err error) (error, Ignore) {
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
					if p.Match(filepath.Base(path)+"/", depth) {
						return filepath.SkipDir, ig
					}
				}
			}
			ig.Patterns = append(ig.Patterns, IgnorePatterns(path, f.Option.VcsIgnores(), depth+1)...)
			return nil, ig
		}
		if !info.follow && info.IsSymlink() {
			return nil, ig
		}
		if !isRoot(depth) && isHidden(info.Name()) {
			return nil, ig
		}

		for _, p := range ig.Patterns {
			if p.Match(filepath.Base(path), depth) {
				return nil, ig
			}
		}

		if pattern.FileRegexp != nil && !pattern.FileRegexp.MatchString(path) {
			return nil, ig
		}
		fileType := UNKNOWN
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

func Walk(root string, ignorePatterns []StringMatcher, follow bool, walkFn WalkFunc) error {
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

func (f *find) addHomePtIgnore() {
	homeDir := getHomeDir()
	if homeDir != "" {
	}
}

func getHomeDir() string {
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

func (f *find) addGlobalGitIgnore() {
	homeDir := getHomeDir()
	if homeDir != "" {
		globalIgnore := globalGitIgnoreName()
		if globalIgnore != "" {
		}
	}
}

func globalGitIgnoreName() string {
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
