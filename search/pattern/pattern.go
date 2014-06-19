package pattern

import (
	"regexp"
)

type Pattern struct {
	Pattern    string
	Regexp     *regexp.Regexp
	FileRegexp *regexp.Regexp
	IgnoreCase bool
	Literal    bool
}

func NewPattern(pattern, filePattern string, smartCase bool, ignoreCase bool, literal bool) (*Pattern, error) {

	if smartCase {
		if regexp.MustCompile(`[[:upper:]]`).MatchString(pattern) {
			ignoreCase = false
		} else {
			ignoreCase = true
		}
	}

	var regExp *regexp.Regexp
	var regExpErr error
	if ignoreCase {
		regExp, regExpErr = regexp.Compile(`(?i)(` + pattern + `)`)
	} else {
		regExp, regExpErr = regexp.Compile(`(` + pattern + `)`)
	}

	var regFile *regexp.Regexp
	var fileErr error
	if filePattern != "" {
		regFile, fileErr = regexp.Compile(filePattern)
	}

	var err error
	switch {
	case regExpErr != nil:
		err = regExpErr
	case fileErr != nil:
		err = fileErr
	default:
		err = nil
	}

	return &Pattern{
		Pattern:    pattern,
		Regexp:     regExp,
		FileRegexp: regFile,
		IgnoreCase: ignoreCase,
		Literal:    literal,
	}, err

}
