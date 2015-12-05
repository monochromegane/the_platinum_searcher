package the_platinum_searcher

import (
	"bufio"
	"io"
	"os"
)

type lineGrep struct {
	printer printer
	before  int
	after   int
	column  bool
}

func newLineGrep(printer printer, opts Option) lineGrep {
	return lineGrep{
		printer: printer,
		before:  opts.OutputOption.Before,
		after:   opts.OutputOption.After,
		column:  opts.OutputOption.Column,
	}
}

func (g lineGrep) enableContext() bool {
	return g.before > 0 || g.after > 0
}

type matchFunc func(b []byte) bool
type countFunc func(b []byte) int

func (g lineGrep) grepEachLines(f *os.File, encoding int, matchFn matchFunc, countFn countFunc) {
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
						match.add(lineNum, 0, scanner.Text(), matched)
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
					column := 0
					if g.column {
						column = countFn(scanner.Bytes())
					}
					match.add(lineNum, column, scanner.Text(), matched)
					// reset after count.
					afterCount = 0
				}
			} else if matched {
				// print only match line.
				column := 0
				if g.column {
					column = countFn(scanner.Bytes())
				}
				match.add(lineNum, column, scanner.Text(), matched)
			}
		}
		lineNum++
	}
	g.printer.print(match)
}

func (g lineGrep) storeBeforeMatch(beforeMatches []line, lineNum int, text string, matched bool) []line {
	if g.before == 0 {
		return beforeMatches
	}
	if len(beforeMatches) >= g.before {
		beforeMatches = beforeMatches[1:]
	}
	return append(beforeMatches, line{
		num:     lineNum,
		column:  0,
		text:    text,
		matched: matched,
	})
}
