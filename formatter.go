package the_platinum_searcher

import (
	"fmt"
	"io"
)

type formatPrinter interface {
	print(match match)
}

func newFormatPrinter(w io.Writer, d decorator, opts Option) formatPrinter {
	switch {
	case opts.OutputOption.FilesWithMatches:
		return fileWithMatch{writer: w, decorator: d}
	case opts.OutputOption.Count:
		return count{writer: w, decorator: d}
	case opts.OutputOption.EnableGroup:
		return group{writer: w, decorator: d}
	default:
		return noGroup{writer: w, decorator: d}
	}
}

type fileWithMatch struct {
	writer    io.Writer
	decorator decorator
}

func (f fileWithMatch) print(match match) {
	fmt.Fprintln(f.writer, f.decorator.path(match.path))
}

type count struct {
	writer    io.Writer
	decorator decorator
}

func (f count) print(match match) {
	count := len(match.lines)
	fmt.Fprintln(f.writer,
		f.decorator.path(match.path)+
			SeparatorColon+
			f.decorator.lineNumber(count),
	)
}

type group struct {
	writer    io.Writer
	decorator decorator
}

func (f group) print(match match) {
	fmt.Fprintln(f.writer, f.decorator.path(match.path))
	for _, line := range match.lines {
		fmt.Fprintln(
			f.writer,
			f.decorator.lineNumber(line.num)+
				SeparatorColon+
				f.decorator.match(match.pattern, line.text),
		)
	}
	fmt.Fprintln(f.writer)
}

type noGroup struct {
	writer    io.Writer
	decorator decorator
}

func (f noGroup) print(match match) {
	path := f.decorator.path(match.path) + SeparatorColon
	for _, line := range match.lines {
		fmt.Fprintln(
			f.writer,
			path+
				f.decorator.lineNumber(line.num)+
				SeparatorColon+
				f.decorator.match(match.pattern, line.text),
		)
	}
}
