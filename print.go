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

func newPrinter(out io.Writer, opts Option) printer {
	return printer{
		mu:        new(sync.Mutex),
		opts:      opts,
		formatter: newFormatPrinter(out, opts),
	}
}

func (p printer) print(match match) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.formatter.print(match)
}
