package the_platinum_searcher

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"
)

type extendedGrep struct {
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

	f.Seek(0, 0)

	var reader io.Reader
	if r := newDecodeReader(f, encoding); r != nil {
		// decode file from shift-jis or euc-jp.
		reader = r
	} else {
		reader = f
	}

	match := match{path: f.Name()}
	scanner := bufio.NewScanner(reader)
	line := 1
	for scanner.Scan() {
		if g.pattern.regexp.Match(scanner.Bytes()) {
			match.add(line, scanner.Text())
		}
		line++
	}
	g.printer.print(match)
}
