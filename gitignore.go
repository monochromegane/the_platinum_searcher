package the_platinum_searcher

import (
	"path/filepath"
	"strings"
)

type GitIgnore struct {
	ignorePatterns patterns
	acceptPatterns patterns
	depth          int
}

func NewGitIgnore(depth int, patterns []string) GitIgnore {
	g := GitIgnore{depth: depth}
	g.Parse(patterns)
	return g
}

func (g *GitIgnore) Parse(patterns []string) {
	for _, p := range patterns {
		p := strings.Trim(string(p), " ")
		if len(p) == 0 || strings.HasPrefix(p, "#") {
			continue
		}

		if strings.HasPrefix(p, "!") {
			g.acceptPatterns = append(g.acceptPatterns,
				pattern(strings.TrimPrefix(p, "!")))
		} else {
			g.ignorePatterns = append(g.ignorePatterns, pattern(p))
		}
	}
}

func (g GitIgnore) Match(path string, isDir bool, depth int) bool {
	if match := g.acceptPatterns.match(path, isDir, depth == g.depth); match {
		return false
	}
	return g.ignorePatterns.match(path, isDir, depth == g.depth)
}

type patterns []pattern

func (ps patterns) match(path string, isDir, isRoot bool) bool {
	for _, p := range ps {
		match := p.match(path, isDir, isRoot)
		if match {
			return true
		}
	}
	return false
}

type pattern string

func (p pattern) match(path string, isDir, isRoot bool) bool {

	if p.hasRootPrefix() && !isRoot {
		return false
	}

	if p.hasDirSuffix() && !isDir {
		return false
	}

	pattern := p.trimedPattern()

	match, _ := filepath.Match(pattern, p.equalizeDepth(path))
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
