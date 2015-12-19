package the_platinum_searcher

import (
	"fmt"
	"io"

	"github.com/shiena/ansicolor"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type formatPrinter interface {
	print(match match)
}

func newFormatPrinter(pattern pattern, w io.Writer, opts Option) formatPrinter {
	writer := newWriter(w, opts)
	decorator := newDecorator(pattern, opts)

	switch {
	case opts.SearchOption.SearchStream:
		return matchLine{decorator: decorator, w: writer}
	case opts.OutputOption.FilesWithMatches:
		return fileWithMatch{decorator: decorator, w: writer}
	case opts.OutputOption.Count:
		return count{decorator: decorator, w: writer}
	case opts.OutputOption.EnableGroup:
		return group{decorator: decorator, w: writer}
	default:
		return noGroup{decorator: decorator, w: writer}
	}
}

type matchLine struct {
	w         io.Writer
	decorator decorator
}

func (f matchLine) print(match match) {
	for _, line := range match.lines {
		column := ""
		if line.matched && line.column > 0 {
			column = f.decorator.columnNumber(line.column) + SeparatorColon
		}
		fmt.Fprintln(f.w,
			column+
				f.decorator.match(line.text, line.matched),
		)
	}
}

type fileWithMatch struct {
	w         io.Writer
	decorator decorator
}

func (f fileWithMatch) print(match match) {
	fmt.Fprintln(f.w, f.decorator.path(match.path))
}

type count struct {
	w         io.Writer
	decorator decorator
}

func (f count) print(match match) {
	count := len(match.lines)
	fmt.Fprintln(f.w,
		f.decorator.path(match.path)+
			SeparatorColon+
			f.decorator.lineNumber(count),
	)
}

type group struct {
	w         io.Writer
	decorator decorator
}

func (f group) print(match match) {
	fmt.Fprintln(f.w, f.decorator.path(match.path))
	for _, line := range match.lines {
		sep := SeparatorColon
		if !line.matched {
			sep = SeparatorHyphen
		}
		column := ""
		if line.matched && line.column > 0 {
			column = f.decorator.columnNumber(line.column) + SeparatorColon
		}
		fmt.Fprintln(f.w,
			f.decorator.lineNumber(line.num)+
				sep+
				column+
				f.decorator.match(line.text, line.matched),
		)
	}
	fmt.Fprintln(f.w)
}

type noGroup struct {
	w         io.Writer
	decorator decorator
}

func (f noGroup) print(match match) {
	path := f.decorator.path(match.path) + SeparatorColon
	for _, line := range match.lines {
		sep := SeparatorColon
		if !line.matched {
			sep = SeparatorHyphen
		}
		column := ""
		if line.matched && line.column > 0 {
			column = f.decorator.columnNumber(line.column) + SeparatorColon
		}
		fmt.Fprintln(f.w,
			path+
				f.decorator.lineNumber(line.num)+
				sep+
				column+
				f.decorator.match(line.text, line.matched),
		)
	}
}

func newWriter(out io.Writer, opts Option) io.Writer {
	encoder := func() io.Writer {
		switch opts.OutputOption.OutputEncode {
		case "sjis":
			return transform.NewWriter(out, japanese.ShiftJIS.NewEncoder())
		case "euc":
			return transform.NewWriter(out, japanese.EUCJP.NewEncoder())
		case "jis":
			return transform.NewWriter(out, japanese.ISO2022JP.NewEncoder())
		default:
			return out
		}
	}()
	if opts.OutputOption.EnableColor {
		return ansicolor.NewAnsiColorWriter(encoder)
	}
	return encoder
}
