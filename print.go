package the_platinum_searcher

import (
	"io"
	"sync"
)

type printer struct {
	mu        *sync.Mutex
	opts      Option
	formatter formatPrinter
}

func newPrinter(
	pattern pattern,
	out,
	errorWriter io.Writer,
	opts Option,
) printer {
	return printer{
		mu:        new(sync.Mutex),
		opts:      opts,
		formatter: newFormatPrinter(pattern, out, errorWriter, opts),
	}
}

func (p printer) print(match match) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if match.size() == 0 {
		return
	}

	p.formatter.print(match)
}

func (p printer) printError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.formatter.printError(err)
}
