package the_platinum_searcher

import (
	"fmt"
	"sync"
)

type printer struct {
	mu *sync.Mutex
}

func newPrinter() printer {
	return printer{
		mu: new(sync.Mutex),
	}
}

func (p printer) print(match match) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, line := range match.lines {
		fmt.Printf("%s:%d:%s\n", match.path, line.num, line.text)
	}
}
