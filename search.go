package the_platinum_searcher

import (
	"io"
)

type search struct {
	roots []string
	out   io.Writer
}

func (s search) start(pathPattern, contentPattern string) error {
	grepChan := make(chan string, 5000)
	done := make(chan struct{})

	p, err := newPattern(contentPattern, opts)
	if err != nil {
		return err
	}

	regFile, err := newPathPattern(pathPattern)
	if err != nil {
		return err
	}

	go find{
		out:  grepChan,
		opts: opts,
	}.start(s.roots, regFile)

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
