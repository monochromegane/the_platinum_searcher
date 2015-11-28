package the_platinum_searcher

import "regexp"

type pattern struct {
	pattern string
	regexp  *regexp.Regexp
}

func newPattern(p string, useRegexp bool) pattern {
	pattern := pattern{pattern: p}
	if useRegexp {
		reg, err := regexp.Compile(`(` + p + `)`)
		if err == nil {
			pattern.regexp = reg
		}
	}
	return pattern
}
