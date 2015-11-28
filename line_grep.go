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

	lineNum := 1
	matchState := newMatchState()
	afterCount := 0
	beforeMatches := make([]line, 0, g.before)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if matched := matchFn(scanner.Bytes()); matched || g.enableContext() {
			if g.enableContext() {
				// print match and context lines.
				matchState = matchState.transition(matched)

				if matchState.isBefore() {
					// store before line.
					beforeMatches = g.storeBeforeMatch(beforeMatches, lineNum, scanner.Text(), matched)
				} else if matchState.isAfter() {
					if g.after > 0 {
						// append after line.
						match.add(lineNum, scanner.Text(), matched)
						afterCount++
					}
					if afterCount >= g.after {
						// reset to before match
						matchState = matchState.reset()
						afterCount = 0
					}
				} else if matchState.isMatching() {
					// append and reset before lines.
					match.lines = append(match.lines, beforeMatches...)
					beforeMatches = make([]line, 0, g.before)
					// append match line.
					match.add(lineNum, scanner.Text(), matched)
					// reset after count.
					afterCount = 0
				}
			} else if matched {
				// print only match line.
				match.add(lineNum, scanner.Text(), matched)
			}
		}
		lineNum++
	}
	printer.print(match)
}

func (g lineGrep) storeBeforeMatch(beforeMatches []line, lineNum int, text string, matched bool) []line {
	if g.before == 0 {
		return beforeMatches
	}
	if len(beforeMatches) >= g.before {
		beforeMatches = beforeMatches[1:]
	}
	return append(beforeMatches, line{lineNum, text, matched})
}
