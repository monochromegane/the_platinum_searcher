package the_platinum_searcher

import (
	"fmt"
	"io"
	"sync"

	"github.com/shiena/ansicolor"
)

type printer struct {
	mu        *sync.Mutex
	opts      Option
	writer    io.Writer
	decorator decorator
}

func newPrinter(out io.Writer, opts Option) printer {
	return printer{
		mu:        new(sync.Mutex),
		opts:      opts,
		writer:    newWriter(out, opts),
		decorator: newDecorator(opts),
	}
}

func (p printer) print(match match) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, line := range match.lines {
		fmt.Fprintln(
			p.writer,
			p.decorator.path(match.path)+
				p.decorator.lineNumber(line.num)+
				p.decorator.match(match.pattern, line.text),
		)
	}
}

func newWriter(out io.Writer, opts Option) io.Writer {
	if opts.OutputOption.EnableColor {
		return ansicolor.NewAnsiColorWriter(out)
	}
	return out
}
