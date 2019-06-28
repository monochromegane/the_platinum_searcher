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
	offset := byteCount(0)
	read := lineCount(0)
	newLineBytes := []byte{'\n'}

loop:
	for {
		n, err := readFile(f, buf[offset:])
		if err == io.EOF {
			// Scan remain (For last line without new line.)
			scan(&match, buf[:offset], pattern, read, encoding, newLineBytes, g.column)
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

			if r := newEncodeReader(bytes.NewReader(newLineBytes), encoding); r != nil {
				newLineBytes, _ = ioutil.ReadAll(r)
			}

			identified = true
		}

		newLine := byteCount(bytes.LastIndex(cbuf, newLineBytes))
		// fmt.Printf("offset: %d, newLine: %d\n", offset, newLine)
		if newLine >= 0 {
			c := scan(&match, cbuf[0:newLine], pattern, read, encoding, newLineBytes, g.column)
			// matchLines = append(matchLines, m...)
			offset = lenb(cbuf[newLine+lenb(newLineBytes):])
			for i := range cbuf[newLine+lenb(newLineBytes):] {
				buf[0+i] = cbuf[newLine+lenb(newLineBytes)+byteCount(i)]
			}
			read += c
		} else {
			grow := make([]byte, len(cbuf)*2)
			copy(grow, buf)
			buf = grow
			offset = lenb(cbuf)
			continue loop
		}
	}
	g.printer.print(match)
}

type byteCount int
type lineCount int

func scan(match *match, buf, pattern []byte, base lineCount, encoding int, newLineBytes []byte, column bool) lineCount {
	offset := byteCount(0)
	newLineCount := lineCount(0)

	for {
		if offset > lenb(buf) {
			break
		}
		cbuf := buf[offset:]
		idx := byteCount(bytes.Index(cbuf, pattern))
		if idx == -1 {
			newLineCount += scanNewLineCount(cbuf, newLineBytes)
			break
		}
		beforeNewLine := byteCount(bytes.LastIndex(cbuf[:idx], newLineBytes))
		if beforeNewLine != -1 {
			newLineCount += (scanNewLineCount(cbuf[:beforeNewLine], newLineBytes) + 1)
		} else {
			beforeNewLine = -lenb(newLineBytes)
		}
		num := base + newLineCount + 1
		afterNewLine := byteCount(bytes.Index(cbuf[idx+lenb(pattern):], newLineBytes))
		if afterNewLine == -1 {
			afterNewLine = lenb(cbuf) - (idx + lenb(pattern))
		} else {
			newLineCount++
		}
		mbuf := cbuf[beforeNewLine+lenb(newLineBytes) : idx+lenb(pattern)+afterNewLine]
		line := make([]byte, lenb(mbuf))
		copy(line, mbuf)

		// decode bytes from shift-jis or euc-jp.
		if r := newDecodeReader(bytes.NewReader(line), encoding); r != nil {
			line, _ = ioutil.ReadAll(r)
		}
		c := byteCount(0)
		if column {
			if beforeNewLine < 0 {
				c = idx + lenb(newLineBytes)
			} else {
				c = idx - beforeNewLine
			}
		}
		match.add(int(num), int(c), string(line), true)
		offset += idx + lenb(pattern) + afterNewLine + lenb(newLineBytes)
	}
	return newLineCount + 1
}

func scanNewLineCount(buf, newLineBytes []byte) lineCount {
	return lineCount(bytes.Count(buf, newLineBytes))
}

func lenb(buf []byte) byteCount {
	return byteCount(len(buf))
}

func readFile(f *os.File, buffer []byte) (byteCount, error) {
	l, err := f.Read(buffer)
	return byteCount(l), err
}
