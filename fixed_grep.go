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
				g.grepAsLines(f, pattern, encoding)
				break
			}
		}

		// grep from buffer.
		if bytes.Contains(buf[:c], pattern) {
			g.grepAsLines(f, pattern, encoding)
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

func (g fixedGrep) grepAsLines(f *os.File, pattern []byte, encoding int) {
	f.Seek(0, 0)
	match := match{path: f.Name()}
	scanner := bufio.NewScanner(f)
	line := 1
	for scanner.Scan() {
		if bytes.Contains(scanner.Bytes(), pattern) {
			var matched []byte
			if r := newDecodeReader(bytes.NewReader(scanner.Bytes()), encoding); r != nil {
				// decode matched line from shift-jis or euc-jp.
				matched, _ = ioutil.ReadAll(r)
			} else {
				matched = scanner.Bytes()
			}
			match.add(line, string(matched))
		}
		line++
	}
	g.printer.print(match)
}
