package the_platinum_searcher

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type ignoreMatchers []ignoreMatcher

type ignoreMatcher interface {
	Match(path string, isDir bool, depth int) bool
}

type genericIgnore []string

func (gi genericIgnore) Match(path string, isDir bool, depth int) bool {
	for _, p := range gi {
		val, _ := filepath.Match(p, path)
		if val {
			return true
		}
	}
	return false
}

func newIgnoreMatchers(path string, ignores []string, depth int) ignoreMatchers {
	var matchers ignoreMatchers
	for _, i := range ignores {
		matchers = append(matchers, newIgnoreMatcher(path, i, depth))
	}
	return matchers
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
		return NewGitIgnore(depth, patterns)
	} else {
		return genericIgnore(patterns)
	}
}
