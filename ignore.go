package the_platinum_searcher

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/monochromegane/go-home"
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

func newIgnoreMatchers(path string, ignores []string) ignoreMatchers {
	var matchers ignoreMatchers
	for _, i := range ignores {
		if matcher := newIgnoreMatcher(path, i); matcher != nil {
			matchers = append(matchers, matcher)
		}
	}
	return matchers
}

type ignoreMatcher interface {
	Match(path string, isDir bool) bool
}

func newIgnoreMatcher(path string, ignore string) ignoreMatcher {

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
		return newGitIgnore(path, patterns)
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
	homeDir := home.Dir()
	if homeDir != "" {
		return newIgnoreMatcher(homeDir, ".ptignore")
	}
	return nil
}

func globalGitIgnore() ignoreMatcher {
	homeDir := home.Dir()
	if homeDir != "" {
		path, globalIgnore := globalGitIgnorePath()
		if globalIgnore != "" {
			return newIgnoreMatcher(path, globalIgnore)
		}
	}
	return nil
}

func globalGitIgnorePath() (string, string) {
	gitCmd, err := exec.LookPath("git")
	if err != nil {
		return "", ""
	}

	file, err := exec.Command(gitCmd, "config", "--get", "core.excludesfile").Output()
	var filename string
	var path string
	if err != nil {
		path, filename = defaultGitIgnorePath()
	} else {
		path = home.Dir()
		filename = strings.TrimSpace(filepath.Base(string(file)))
	}
	return path, filename
}

func defaultGitIgnorePath() (string, string) {
	file_exists := func(path string) bool {
		_, err := os.Stat(path)
		return err == nil
	}

	if path := os.Getenv("XDG_CONFIG_HOME"); path != "" && file_exists(fmt.Sprintf("%s/git/ignore", path)) {
		return fmt.Sprintf("%s/git", path), "ignore"
	} else if path := fmt.Sprintf("%s/.config/git", home.Dir()); file_exists(fmt.Sprintf("%s/ignore", path)) {
		return path, "ignore"
	} else if path := home.Dir(); file_exists(fmt.Sprintf("%s/.gitignore", path)) {
		return path, ".gitignore"
	}
	return "", ""
}
