package pattern

import (
	"regexp"
)

type Pattern struct {
	Pattern    string
	Error      error
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

	var regIgnoreCase *regexp.Regexp
	var err error
	if ignoreCase {
		regIgnoreCase, err = regexp.Compile(`(?i)(` + pattern + `)`)
	}

	return &Pattern{
		Pattern:    pattern,
		Error:      err,
		Regexp:     regIgnoreCase,
		IgnoreCase: ignoreCase,
	}

}
