package the_platinum_searcher

import (
	"path/filepath"
	"strings"
)

type gitIgnore struct {
	ignorePatterns patterns
	acceptPatterns patterns
	path           string
	depth          int
}

func newGitIgnore(path string, depth int, patterns []string) gitIgnore {
	g := gitIgnore{path: path, depth: depth}
	g.parse(patterns)
	return g
}

func (g *gitIgnore) parse(patterns []string) {
	for _, p := range patterns {
		p := strings.Trim(string(p), " ")
		if len(p) == 0 || strings.HasPrefix(p, "#") {
			continue
		}

		if strings.HasPrefix(p, "!") {
			g.acceptPatterns = append(g.acceptPatterns,
				pattern{strings.TrimPrefix(p, "!"), g.path, g.depth - 1})
		} else {
			g.ignorePatterns = append(g.ignorePatterns, pattern{p, g.path, g.depth - 1})
		}
	}
}

func (g gitIgnore) Match(path string, isDir bool) bool {
	if match := g.acceptPatterns.match(path, isDir); match {
		return false
	}
	return g.ignorePatterns.match(path, isDir)
}

type patterns []pattern

func (ps patterns) match(path string, isDir bool) bool {
	for _, p := range ps {
		match := p.match(path, isDir)
		if match {
			return true
		}
	}
	return false
}

type pattern struct {
	path  string
	base  string
	depth int
}

func (p pattern) match(path string, isDir bool) bool {

	if p.hasDirSuffix() && !isDir {
		return false
	}

	pattern := p.trimedPattern()

	var match bool
	if p.hasRootPrefix() {
		// absolute pattern
		match, _ = filepath.Match(filepath.Join(p.base, p.path), path)
	} else {
		// relative pattern
		match, _ = filepath.Match(pattern, p.equalizeDepth(path))
	}
	return match
}

func (p pattern) equalizeDepth(path string) string {
	trimedPath := strings.TrimPrefix(path, p.base)
	patternDepth := strings.Count(p.path, "/")
	pathDepth := strings.Count(trimedPath, string(filepath.Separator))
	start := 0
	if diff := pathDepth - patternDepth; diff > 0 {
		start = diff
	}
	return filepath.Join(strings.Split(trimedPath, string(filepath.Separator))[start:]...)
}

func (p pattern) prefix() string {
	return string(p.path[0])
}

func (p pattern) suffix() string {
	return string(p.path[len(p.path)-1])
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
	return strings.Trim(p.path, "/")
}
