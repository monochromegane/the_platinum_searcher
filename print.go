package the_platinum_searcher

import (
	"io"
)

type printer struct {
	in        chan match
	opts      Option
	formatter formatPrinter
	done      chan struct{}
}

func newPrinter(pattern pattern, out io.Writer, opts Option) printer {
	p := printer{
		in:        make(chan match, 200),
		opts:      opts,
		formatter: newFormatPrinter(pattern, out, opts),
		done:      make(chan struct{}),
	}

	go p.loop()
	return p
}

func (p printer) print(match match) {
	if match.size() == 0 {
		return
	}

	p.in <- match
}

func (p printer) loop() {
	defer func() {
		p.done <- struct{}{}
	}()

	for match := range p.in {
		p.formatter.print(match)
	}
}
