package the_platinum_searcher

import (
	"fmt"
	"path/filepath"
	"strings"
)

type patterns []pattern

func (ps patterns) IsMatch(path string, isDir, isRoot bool) bool {
	for _, p := range ps {
		match := p.IsMatch(path, isDir, isRoot)
		if match {
			return true
		}
	}
	return false
}

type pattern string

func (p pattern) IsMatch(path string, isDir, isRoot bool) bool {

	if p.hasRootPrefix() && !isRoot {
		return false
	}

	if p.hasDirSuffix() && !isDir {
		return false
	}

	pattern := p.trimedPattern()

	match, _ := filepath.Match(pattern, p.equalizeDepth(path))
	fmt.Printf("ptn:%s path:%s(%s) => %t\n", pattern, p.equalizeDepth(path), path, match)
	return match
}

func (p pattern) equalizeDepth(path string) string {
	patternDepth := strings.Count(string(p), "/")
	pathDepth := strings.Count(path, string(filepath.Separator))
	start := 0
	if diff := pathDepth - patternDepth; diff >= 0 {
		start = diff
	}
	return filepath.Join(strings.Split(path, "/")[start:]...)
}

func (p pattern) prefix() string {
	return string(p[0])
}

func (p pattern) suffix() string {
	return string(p[len(p)-1])
}

func (p pattern) hasRootPrefix() bool {
	return p.prefix() == "/"
}

func (p pattern) hasNegativePrefix() bool {
	return p.prefix() == "!"
}

func (p pattern) hasDirSuffix() bool {
	return p.suffix() == "/"
}

func (p pattern) trimedPattern() string {
	return strings.Trim(string(p), "/")
}
