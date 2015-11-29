package the_platinum_searcher

import (
	"io"
	"log"
	"os"
	"sync"
)

type extendedGrep struct {
	lineGrep
	pattern pattern
	printer printer
}

func (g extendedGrep) grep(path string, sem chan struct{}, wg *sync.WaitGroup) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}

	defer func() {
		f.Close()
		<-sem
		wg.Done()
	}()

	buf := make([]byte, 512)

	c, err := f.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatalf("read: %s\n", err)
	}

	if err == io.EOF {
		return
	}

	// detect encoding.
	limit := c
	if limit > 512 {
		limit = 512
	}

	encoding := detectEncoding(buf[:limit])
	if encoding == ERROR || encoding == BINARY {
		return
	}

	// grep each lines.
	g.grepEachLines(f, encoding, g.printer, func(b []byte) bool {
		return g.pattern.regexp.Match(b)
	}, func(b []byte) int {
		return g.pattern.regexp.FindIndex(b)[0] + 1
	})
}
