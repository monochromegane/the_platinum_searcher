package the_platinum_searcher

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type StringMatcher interface {
	// Match matcher with pattern and depth
	Match(string, int) bool
}

type Ignore struct {
	Patterns []StringMatcher
}

type gitIgnoreMatcher struct {
	patterns []string
	depth    int
}

func (ps gitIgnoreMatcher) Match(file string, depth int) bool {
	negatedIgnoreMatch := false
	ignoreMatch := false

	for _, p := range ps.patterns {
		if len(p) == 0 {
			continue
		}

		if p[0] == '!' {
			negatedIgnoreMatch, _ = filepath.Match(p[1:], file)
		} else if !ignoreMatch {
			if p[0] == '/' {
				// Only match with "/" prefix on current depth
				if ps.depth == depth || ps.depth == -1 {
					ignoreMatch, _ = filepath.Match(p[1:], file)
				}
			} else {
				ignoreMatch, _ = filepath.Match(p, file)
			}
		}
	}

	return ignoreMatch && !negatedIgnoreMatch
}

type genericIgnoreMatcher []string

func (im genericIgnoreMatcher) Match(file string, depth int) bool {
	for _, p := range im {
		val, _ := filepath.Match(p, file)
		if val {
			return true
		}
	}
	return false
}

func IgnorePatterns(path string, ignores []string, depth int) []StringMatcher {
	var patterns []StringMatcher
	for _, ignore := range ignores {
		file, err := os.Open(filepath.Join(path, ignore))
		if err != nil {
			continue
		}
		reader := bufio.NewReader(file)
		buf := make([]byte, 1024)

		var thesePatterns []string
		for {
			buf, _, err = reader.ReadLine()
			if err != nil {
				break
			}
			s := strings.Trim(string(buf), " ")

			if len(s) == 0 || strings.HasPrefix(s, "#") {
				continue
			}
			thesePatterns = append(thesePatterns, s)
		}

		if len(thesePatterns) > 0 {
			if ignore == ".gitignore" {
				patterns = append(patterns, gitIgnoreMatcher{thesePatterns, depth})
			} else {
				patterns = append(patterns, genericIgnoreMatcher(thesePatterns))
			}
		}
		file.Close()
	}
	return patterns
}
