package the_platinum_searcher

import (
	"fmt"
	"io"
	"sync"
)

type printer struct {
	mu  *sync.Mutex
	out io.Writer
}

func newPrinter(out io.Writer) printer {
	return printer{
		mu:  new(sync.Mutex),
		out: out,
	}
}

func (p printer) print(match match) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, line := range match.lines {
		fmt.Fprintf(p.out, "%s:%d:%s\n", match.path, line.num, line.text)
	}
}
