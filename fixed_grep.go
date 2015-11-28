package the_platinum_searcher

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type fixedGrep struct {
	lineGrep
	pattern pattern
	printer printer
}

func (g fixedGrep) grep(path string, sem chan struct{}, wg *sync.WaitGroup) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}

	defer func() {
		f.Close()
		<-sem
		wg.Done()
	}()

	buf := make([]byte, 8196)
	var stash []byte
	identified := false
	var encoding int
	pattern := g.pattern.pattern

	for {
		c, err := f.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("read: %s\n", err)
		}

		if err == io.EOF {
			break
		}

		// detect encoding.
		if !identified {
			limit := c
			if limit > 512 {
				limit = 512
			}

			encoding = detectEncoding(buf[:limit])
			if encoding == ERROR || encoding == BINARY {
				break
			}

			if r := newEncodeReader(bytes.NewReader(pattern), encoding); r != nil {
				// encode pattern to shift-jis or euc-jp.
				pattern, _ = ioutil.ReadAll(r)
			}
			identified = true
		}

		// repair first line from previous last line.
		if len(stash) > 0 {
			var repaired []byte
			index := bytes.Index(buf[:c], newLine)
			if index == -1 {
				repaired = append(stash, buf[:c]...)
			} else {
				repaired = append(stash, buf[:index]...)
			}
			// grep from repaied line.
			if bytes.Contains(repaired, pattern) {
				// grep each lines.
				g.grepAsLines(f, encoding, g.printer, func(b []byte) bool {
					return bytes.Contains(b, g.pattern.pattern)
				})
				break
			}
		}

		// grep from buffer.
		if bytes.Contains(buf[:c], pattern) {
			// grep each lines.
			g.grepAsLines(f, encoding, g.printer, func(b []byte) bool {
				return bytes.Contains(b, g.pattern.pattern)
			})
			break
		}

		// stash last line.
		index := bytes.LastIndex(buf[:c], newLine)
		if index == -1 {
			stash = append(stash, buf[:c]...)
		} else {
			stash = make([]byte, c-index)
			copy(stash, buf[index:c])
		}
	}
}

type lineGrep struct {
}

type matchFunc func(b []byte) bool

func (g lineGrep) grepAsLines(f *os.File, encoding int, printer printer, matchFn matchFunc) {
	f.Seek(0, 0)
	match := match{path: f.Name()}

	var reader io.Reader
	if r := newDecodeReader(f, encoding); r != nil {
		// decode file from shift-jis or euc-jp.
		reader = r
	} else {
		reader = f
	}

	scanner := bufio.NewScanner(reader)
	line := 1
	for scanner.Scan() {
		if matchFn(scanner.Bytes()) {
			match.add(line, scanner.Text())
		}
		line++
	}
	printer.print(match)
}
