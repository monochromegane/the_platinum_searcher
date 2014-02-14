package pattern

import (
	"regexp"
)

type Pattern struct {
	Pattern    string
	Regexp     *regexp.Regexp
	IgnoreCase bool
}

func NewPattern(pattern string, smartCase, ignoreCase bool) *Pattern {

	if smartCase {
		if regexp.MustCompile(`[[:upper:]]`).MatchString(pattern) {
			ignoreCase = false
		} else {
			ignoreCase = true
		}
	}

	return &Pattern{
		Pattern:    pattern,
		Regexp:     regexp.MustCompile(`(?i)(` + pattern + `)`),
		IgnoreCase: ignoreCase,
	}

}
