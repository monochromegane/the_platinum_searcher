package pattern

import (
	"regexp"
)

type Pattern struct {
	Pattern    string
	Error      error
	Regexp     *regexp.Regexp
	FileRegexp *regexp.Regexp
	IgnoreCase bool
}

func NewPattern(pattern, filePattern string, smartCase, ignoreCase bool) *Pattern {

	if smartCase {
		if regexp.MustCompile(`[[:upper:]]`).MatchString(pattern) {
			ignoreCase = false
		} else {
			ignoreCase = true
		}
	}

	var regIgnoreCase *regexp.Regexp
	var ignoreErr error
	if ignoreCase {
		regIgnoreCase, ignoreErr = regexp.Compile(`(?i)(` + pattern + `)`)
	}

	var regFile *regexp.Regexp
	var fileErr error
	if filePattern != "" {
		regFile, fileErr = regexp.Compile(filePattern)
	}

	var err error
	switch {
	case ignoreErr != nil:
		err = ignoreErr
	case fileErr != nil:
		err = fileErr
	default:
		err = nil
	}

	return &Pattern{
		Pattern:    pattern,
		Error:      err,
		Regexp:     regIgnoreCase,
		FileRegexp: regFile,
		IgnoreCase: ignoreCase,
	}

}
