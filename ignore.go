package the_platinum_searcher

import (
	"bufio"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type ignoreMatchers []ignoreMatcher

func (im ignoreMatchers) Match(path string, isDir bool) bool {
	for _, ig := range im {
		if ig.Match(path, isDir) {
			return true
		}
	}
	return false
}

func newIgnoreMatchers(path string, ignores []string, depth int) ignoreMatchers {
	var matchers ignoreMatchers
	for _, i := range ignores {
		if matcher := newIgnoreMatcher(path, i, depth); matcher != nil {
			matchers = append(matchers, matcher)
		}
	}
	return matchers
}

type ignoreMatcher interface {
	Match(path string, isDir bool) bool
}

func newIgnoreMatcher(path string, ignore string, depth int) ignoreMatcher {

	file, err := os.Open(filepath.Join(path, ignore))
	if err != nil {
		return nil
	}

	defer file.Close()
	reader := bufio.NewReader(file)

	var patterns []string
	for {
		buf, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		line := strings.Trim(string(buf), " ")
		if len(line) == 0 {
			continue
		}
		patterns = append(patterns, line)
	}

	if ignore == ".ptignore" || ignore == ".gitignore" {
		return newGitIgnore(path, depth, patterns)
	} else {
		return genericIgnore(patterns)
	}
}

type genericIgnore []string

func (gi genericIgnore) Match(path string, isDir bool) bool {
	for _, p := range gi {
		val, _ := filepath.Match(p, filepath.Base(path))
		if val {
			return true
		}
	}
	return false
}

func homePtIgnore() ignoreMatcher {
	homeDir := getHomeDir()
	if homeDir != "" {
		return newIgnoreMatcher(homeDir, ".ptignore", -1)
	}
	return nil
}

func globalGitIgnore() ignoreMatcher {
	homeDir := getHomeDir()
	if homeDir != "" {
		globalIgnore := globalGitIgnoreName()
		if globalIgnore != "" {
			return newIgnoreMatcher(homeDir, globalIgnore, -1)
		}
	}
	return nil
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
