package the_platinum_searcher

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var newLine = []byte("\n")

type grep struct {
	in      chan string
	done    chan struct{}
	printer printer
}

func (g grep) start(pattern string) {
	sem := make(chan struct{}, 208)
	wg := &sync.WaitGroup{}

	p := []byte(pattern)
	for path := range g.in {
		sem <- struct{}{}
		wg.Add(1)
		go g.grep(path, p, sem, wg)
	}
	wg.Wait()
	g.done <- struct{}{}
}

func (g grep) grep(path string, pattern []byte, sem chan struct{}, wg *sync.WaitGroup) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}

	buf := make([]byte, 8196)
	var stash []byte
	identified := false
	var encoding int

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

			if encoder := getEncoder(encoding); encoder != nil {
				// encode pattern to shift-jis or euc-jp.
				pattern, _ = ioutil.ReadAll(transform.NewReader(bytes.NewReader(pattern), encoder))
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
	f.Close()
	<-sem
	wg.Done()
}

func (g grep) grepAsLines(f *os.File, pattern []byte, encoding int) {
	f.Seek(0, 0)
	match := match{path: f.Name(), pattern: pattern, encoding: encoding}
	scanner := bufio.NewScanner(f)
	line := 1
	for scanner.Scan() {
		if bytes.Contains(scanner.Bytes(), pattern) {
			match.add(line, scanner.Text())
		}
		line++
	}
	g.printer.print(match)
}

func getEncoder(encoding int) transform.Transformer {
	switch encoding {
	case EUCJP:
		return japanese.EUCJP.NewEncoder()
	case SHIFTJIS:
		return japanese.ShiftJIS.NewEncoder()
	}
	return nil
}
