package pattern

import (
	"regexp"
)

type Pattern struct {
	Pattern    string
	Regexp     *regexp.Regexp
	FileRegexp *regexp.Regexp
	IgnoreCase bool
	UseRegexp  bool
}

func NewPattern(pattern, filePattern string, smartCase, ignoreCase, useRegexp bool) (*Pattern, error) {

	if smartCase {
		if regexp.MustCompile(`[[:upper:]]`).MatchString(pattern) {
			ignoreCase = false
		} else {
			ignoreCase = true
		}
	}

	var regPattern *regexp.Regexp
	var patternErr error
	if ignoreCase {
		regPattern, patternErr = regexp.Compile(`(?i)(` + pattern + `)`)
	} else if useRegexp {
		regPattern, patternErr = regexp.Compile(`(` + pattern + `)`)
	}

	var regFile *regexp.Regexp
	var fileErr error
	if filePattern != "" {
		regFile, fileErr = regexp.Compile(filePattern)
	}

	var err error
	switch {
	case patternErr != nil:
		err = patternErr
	case fileErr != nil:
		err = fileErr
	default:
		err = nil
	}

	return &Pattern{
		Pattern:    pattern,
		Regexp:     regPattern,
		FileRegexp: regFile,
		IgnoreCase: ignoreCase,
		UseRegexp:  useRegexp,
	}, err

}
