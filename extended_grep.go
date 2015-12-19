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
}

func (g extendedGrep) grep(path string, sem chan struct{}, wg *sync.WaitGroup) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}

	defer func() {
		f.Close()
		<-sem
		wg.Done()
	}()

	if f == os.Stdin {
		// TODO: File type is fixed in ASCII because it can not determine the character code.
		g.grepEachLines(f, ASCII, func(b []byte) bool {
			return g.pattern.regexp.Match(b)
		}, func(b []byte) int {
			return g.pattern.regexp.FindIndex(b)[0] + 1
		})
		return
	}

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
	g.grepEachLines(f, encoding, func(b []byte) bool {
		return g.pattern.regexp.Match(b)
	}, func(b []byte) int {
		return g.pattern.regexp.FindIndex(b)[0] + 1
	})
}
