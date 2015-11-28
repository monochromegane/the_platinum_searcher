package the_platinum_searcher

import (
	"bufio"
	"io"
	"os"
)

type lineGrep struct {
	before int
	after  int
}

func newLineGrep(opts Option) lineGrep {
	return lineGrep{
		before: opts.OutputOption.Before,
		after:  opts.OutputOption.After,
	}
}

func (g lineGrep) enableContext() bool {
	return g.before > 0 || g.after > 0
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

	line := 1
	matchState := newMatchState()
	afterCount := 0

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if matched := matchFn(scanner.Bytes()); matched || g.enableContext() {
			if g.enableContext() {
				// print match and context lines.
				matchState = matchState.transition(matched)

				if matchState.isBefore() {
				} else if matchState.isMatching() {
					match.add(line, scanner.Text(), matched)
					afterCount = 0
				} else if matchState.isAfter() {
					match.add(line, scanner.Text(), matched)
					afterCount++
					if afterCount >= g.after {
						matchState = matchState.reset()
						afterCount = 0
					}
				}
			} else if matched {
				// print only match line.
				match.add(line, scanner.Text(), matched)
			}
		}
		line++
	}
	printer.print(match)
}
