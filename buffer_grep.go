package the_platinum_searcher

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type bufferGrep struct {
	printer
	pattern pattern
	column  bool
}

func (g bufferGrep) grep(path string, buf []byte) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}
	defer f.Close()

	identified := false
	var encoding int
	pattern := g.pattern.pattern
	match := match{path: path}
	offset, read := 0, 0

loop:
	for {
		n, err := f.Read(buf[offset:])
		if err == io.EOF {
			// Scan remain (For last line without new line.)
			scan(&match, buf[:offset], pattern, read, encoding, g.column)
			break
		}
		if err != nil {
			panic(err)
		}
		cbuf := buf[0 : offset+n]

		// detect encoding.
		if !identified {
			limit := n
			if limit > 512 {
				limit = 512
			}

			if f == os.Stdin {
				// TODO: File type is fixed in ASCII because it can not determine the character code.
				encoding = ASCII
			} else {
				encoding = detectEncoding(cbuf[:limit])
			}
			if encoding == ERROR || encoding == BINARY {
				break
			}

			if r := newEncodeReader(bytes.NewReader(pattern), encoding); r != nil {
				// encode pattern to shift-jis or euc-jp.
				pattern, _ = ioutil.ReadAll(r)
			}
			identified = true
		}

		newLine := bytes.LastIndexByte(cbuf, '\n')
		// fmt.Printf("offset: %d, newLine: %d\n", offset, newLine)
		if newLine >= 0 {
			c := scan(&match, cbuf[0:newLine], pattern, read, encoding, g.column)
			// matchLines = append(matchLines, m...)
			offset = len(cbuf[newLine+1:])
			for i, _ := range cbuf[newLine+1:] {
				buf[0+i] = cbuf[newLine+1+i]
			}
			read += c
		} else {
			grow := make([]byte, len(cbuf)*2)
			copy(grow, buf)
			buf = grow
			offset = len(cbuf)
			continue loop
		}
	}
	g.printer.print(match)
}

var NewLineBytes = []byte{10}

func scanNewLine(buf []byte) int {
	return bytes.Count(buf, NewLineBytes)
}

func scan(match *match, buf, pattern []byte, base, encoding int, column bool) int {
	offset, newLineCount := 0, 0
	for {
		if offset > len(buf) {
			break
		}
		cbuf := buf[offset:]
		idx := bytes.Index(cbuf, pattern)
		if idx == -1 {
			newLineCount += scanNewLineCount(cbuf)
			break
		}
		beforeNewLine := bytes.LastIndexByte(cbuf[:idx], '\n')
		if beforeNewLine != -1 {
			newLineCount += (scanNewLineCount(cbuf[:beforeNewLine]) + 1)
		}
		num := base + newLineCount + 1
		afterNewLine := bytes.IndexByte(cbuf[idx+len(pattern):], '\n')
		if afterNewLine == -1 {
			afterNewLine = len(cbuf) - (idx + len(pattern))
		} else {
			newLineCount++
		}
		mbuf := cbuf[beforeNewLine+1 : idx+len(pattern)+afterNewLine]
		line := make([]byte, len(mbuf))
		copy(line, mbuf)

		// decode bytes from shift-jis or euc-jp.
		if r := newDecodeReader(bytes.NewReader(line), encoding); r != nil {
			line, _ = ioutil.ReadAll(r)
		}
		c := 0
		if column {
			if beforeNewLine == -1 {
				c = idx + 1
			} else {
				c = idx - beforeNewLine
			}
		}
		match.add(num, c, string(line), true)
		offset += idx + len(pattern) + afterNewLine + 1
	}
	return newLineCount + 1
}

func scanNewLineCount(buf []byte) int {
	return bytes.Count(buf, NewLineBytes)
}
