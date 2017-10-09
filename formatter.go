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
	printError(err error)
}

func newFormatPrinter(
	pattern pattern,
	w,
	outErr io.Writer,
	opts Option,
) formatPrinter {
	writer := newWriter(w, opts)
	errorWriter := newWriter(outErr, opts)
	decorator := newDecorator(pattern, opts)

	switch {
	case opts.SearchOption.SearchStream:
		return matchLine{decorator: decorator, w: writer, errorWriter: errorWriter}
	case opts.OutputOption.FilesWithMatches:
		return fileWithMatch{decorator: decorator, w: writer, errorWriter: errorWriter, useNull: opts.OutputOption.Null}
	case opts.OutputOption.Count:
		return count{decorator: decorator, w: writer}
	case opts.OutputOption.EnableGroup:
		return group{decorator: decorator, w: writer, errorWriter: errorWriter, useNull: opts.OutputOption.Null, enableLineNumber: opts.OutputOption.EnableLineNumber}
	default:
		return noGroup{decorator: decorator, w: writer, errorWriter: errorWriter, enableLineNumber: opts.OutputOption.EnableLineNumber}
	}
}

type matchLine struct {
	w           io.Writer
	errorWriter io.Writer
	decorator   decorator
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

func (f matchLine) printError(err error) {
	fmt.Fprintln(f.errorWriter, f.decorator.error(err))
}

type fileWithMatch struct {
	w           io.Writer
	errorWriter io.Writer
	decorator   decorator
	useNull     bool
}

func (f fileWithMatch) print(match match) {
	if f.useNull {
		fmt.Fprint(f.w, f.decorator.path(match.path))
		fmt.Fprint(f.w, "\x00")
	} else {
		fmt.Fprintln(f.w, f.decorator.path(match.path))
	}
}

func (f fileWithMatch) printError(err error) {
	fmt.Fprintln(f.errorWriter, f.decorator.error(err))
}

type count struct {
	w           io.Writer
	errorWriter io.Writer
	decorator   decorator
}

func (f count) print(match match) {
	count := len(match.lines)
	fmt.Fprintln(f.w,
		f.decorator.path(match.path)+
			SeparatorColon+
			f.decorator.lineNumber(count),
	)
}

func (f count) printError(err error) {
	fmt.Fprintln(f.errorWriter, f.decorator.error(err))
}

type group struct {
	w                io.Writer
	errorWriter      io.Writer
	decorator        decorator
	useNull          bool
	enableLineNumber bool
}

func (f group) print(match match) {
	if f.useNull {
		fmt.Fprint(f.w, f.decorator.path(match.path))
		fmt.Fprint(f.w, "\x00")
	} else {
		fmt.Fprintln(f.w, f.decorator.path(match.path))
	}

	for _, line := range match.lines {
		sep := SeparatorColon
		if !line.matched {
			sep = SeparatorHyphen
		}
		column := ""
		if line.matched && line.column > 0 {
			column = f.decorator.columnNumber(line.column) + SeparatorColon
		}
		lineNum := ""
		if f.enableLineNumber {
			lineNum = f.decorator.lineNumber(line.num) + sep
		}
		fmt.Fprintln(f.w,
			lineNum+
				column+
				f.decorator.match(line.text, line.matched),
		)
	}
	fmt.Fprintln(f.w)
}

func (f group) printError(err error) {
	fmt.Fprintln(f.errorWriter, f.decorator.error(err))
}

type noGroup struct {
	w                io.Writer
	errorWriter      io.Writer
	decorator        decorator
	enableLineNumber bool
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
		lineNum := ""
		if f.enableLineNumber {
			lineNum = f.decorator.lineNumber(line.num) + sep
		}
		fmt.Fprintln(f.w,
			path+
				lineNum+
				column+
				f.decorator.match(line.text, line.matched),
		)
	}
}

func (f noGroup) printError(err error) {
	fmt.Fprintln(f.errorWriter, f.decorator.error(err))
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
