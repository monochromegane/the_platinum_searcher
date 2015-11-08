package the_platinum_searcher

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var newLine = []byte("\n")

type grep struct {
	in   chan string
	done chan struct{}
}

func (g grep) start(pattern string) {
	sem := make(chan struct{}, 256)
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

	for {
		c, err := f.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("read: %s\n", err)
		}

		if err == io.EOF {
			break
		}

		if !identified {
			limit := c
			if limit > 512 {
				limit = 512
			}
			detectEncoding(buf[:limit])
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
				g.grepAsLines(f, pattern)
				break
			}
		}

		// grep from buffer.
		if bytes.Contains(buf[:c], pattern) {
			g.grepAsLines(f, pattern)
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

func (g grep) grepAsLines(f *os.File, pattern []byte) {
	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if bytes.Contains(scanner.Bytes(), pattern) {
			fmt.Printf("%s\n", scanner.Text())
		}
	}
}
