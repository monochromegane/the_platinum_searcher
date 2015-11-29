package the_platinum_searcher

import "regexp"

// pattern should be interface, but match(large []byte) method called through interface is too slow.
type pattern struct {
	pattern []byte
	regexp  *regexp.Regexp
	opts    Option
}

func newPattern(p string, opts Option) (pattern, error) {
	pattern := pattern{pattern: []byte(p), opts: opts}

	if opts.SearchOption.Regexp {
		var reg *regexp.Regexp
		var err error
		if opts.SearchOption.IgnoreCase {
			reg, err = regexp.Compile(`(?i)(` + p + `)`)
		} else {
			reg, err = regexp.Compile(`(` + p + `)`)
		}
		if err != nil {
			return pattern, err
		}
		pattern.regexp = reg
	}
	return pattern, nil
}
