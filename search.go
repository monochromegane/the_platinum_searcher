package the_platinum_searcher

import "io"

type search struct {
	root string
	out  io.Writer
	opts Option
}

func (s search) start(pattern string) error {
	grepChan := make(chan string, 5000)
	done := make(chan struct{})

	if opts.SearchOption.WordRegexp {
		opts.SearchOption.Regexp = true
		pattern = "\\b" + pattern + "\\b"
	}

	p, err := newPattern(pattern, s.opts)
	if err != nil {
		return err
	}

	go find{
		out:  grepChan,
		opts: s.opts,
	}.start(s.root)

	go newGrep(
		p,
		grepChan,
		done,
		opts,
		newPrinter(p, s.out, s.opts),
	).start()

	<-done

	return nil
}
