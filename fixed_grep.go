package the_platinum_searcher

import (
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
}

func (g fixedGrep) grep(path string, sem chan struct{}, wg *sync.WaitGroup) {
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
			return bytes.Contains(b, g.pattern.pattern)
		}, func(b []byte) int {
			return bytes.Index(b, g.pattern.pattern) + 1
		})
		return
	}

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
				g.grepEachLines(f, encoding, func(b []byte) bool {
					return bytes.Contains(b, g.pattern.pattern)
				}, func(b []byte) int {
					return bytes.Index(b, g.pattern.pattern) + 1
				})
				break
			}
		}

		// grep from buffer.
		if bytes.Contains(buf[:c], pattern) {
			// grep each lines.
			g.grepEachLines(f, encoding, func(b []byte) bool {
				return bytes.Contains(b, g.pattern.pattern)
			}, func(b []byte) int {
				return bytes.Index(b, g.pattern.pattern) + 1
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
