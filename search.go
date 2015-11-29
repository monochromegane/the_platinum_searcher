package the_platinum_searcher

import (
	"io"
	"regexp"
)

type search struct {
	root string
	out  io.Writer
}

func (s search) start(pattern string) error {
	grepChan := make(chan string, 5000)
	done := make(chan struct{})

	if opts.SearchOption.WordRegexp {
		opts.SearchOption.Regexp = true
		pattern = "\\b" + pattern + "\\b"
	}

	if opts.SearchOption.SmartCase {
		if !regexp.MustCompile(`[[:upper:]]`).MatchString(pattern) {
			opts.SearchOption.IgnoreCase = true
		}
	}

	if opts.SearchOption.IgnoreCase {
		opts.SearchOption.Regexp = true
	}

	p, err := newPattern(pattern, opts)
	if err != nil {
		return err
	}

	if opts.OutputOption.Context > 0 {
		opts.OutputOption.Before = opts.OutputOption.Context
		opts.OutputOption.After = opts.OutputOption.Context
	}

	go find{
		out:  grepChan,
		opts: opts,
	}.start(s.root)

	go newGrep(
		p,
		grepChan,
		done,
		opts,
		newPrinter(p, s.out, opts),
	).start()

	<-done

	return nil
}
