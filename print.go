package the_platinum_searcher

import (
	"io"
	"sync"

	"github.com/shiena/ansicolor"
)

type printer struct {
	mu        *sync.Mutex
	opts      Option
	formatter formatPrinter
}

func newPrinter(out io.Writer, opts Option) printer {
	return printer{
		mu:   new(sync.Mutex),
		opts: opts,
		formatter: newFormatPrinter(
			newWriter(out, opts),
			newDecorator(opts),
			opts,
		),
	}
}

func (p printer) print(match match) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.formatter.print(match)
}

func newWriter(out io.Writer, opts Option) io.Writer {
	if opts.OutputOption.EnableColor {
		return ansicolor.NewAnsiColorWriter(out)
	}
	return out
}
