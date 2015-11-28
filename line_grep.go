package the_platinum_searcher

import (
	"bufio"
	"io"
	"os"
)

type lineGrep struct {
}

type matchFunc func(b []byte) bool

func (g lineGrep) grepEachLines(f *os.File, encoding int, printer printer, matchFn matchFunc) {
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
			match.add(line, scanner.Text(), true)
		}
		line++
	}
	printer.print(match)
}
