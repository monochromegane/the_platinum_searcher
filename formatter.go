package the_platinum_searcher

import (
	"fmt"
	"io"

	"github.com/shiena/ansicolor"
)

type formatPrinter interface {
	print(match match)
	writer(encoding int) io.Writer
}

func newFormatPrinter(w io.Writer, opts Option) formatPrinter {
	decoder := newDecoder(newColorWriter(w, opts), opts)
	decorator := newDecorator(opts)

	switch {
	case opts.OutputOption.FilesWithMatches:
		return fileWithMatch{decorator: decorator, decoder: decoder}
	case opts.OutputOption.Count:
		return count{decorator: decorator, decoder: decoder}
	case opts.OutputOption.EnableGroup:
		return group{decorator: decorator, decoder: decoder}
	default:
		return noGroup{decorator: decorator, decoder: decoder}
	}
}

type fileWithMatch struct {
	decoder
	decorator decorator
}

func (f fileWithMatch) print(match match) {
	fmt.Fprintln(f.writer(match.encoding), f.decorator.path(match.path))
}

type count struct {
	decoder
	decorator decorator
}

func (f count) print(match match) {
	count := len(match.lines)
	fmt.Fprintln(f.writer(match.encoding),
		f.decorator.path(match.path)+
			SeparatorColon+
			f.decorator.lineNumber(count),
	)
}

type group struct {
	decoder
	decorator decorator
}

func (f group) print(match match) {
	w := f.writer(match.encoding)
	fmt.Fprintln(w, f.decorator.path(match.path))
	for _, line := range match.lines {
		fmt.Fprintln(w,
			f.decorator.lineNumber(line.num)+
				SeparatorColon+
				f.decorator.match(match.pattern, line.text),
		)
	}
	fmt.Fprintln(w)
}

type noGroup struct {
	decoder
	decorator decorator
}

func (f noGroup) print(match match) {
	w := f.writer(match.encoding)
	path := f.decorator.path(match.path) + SeparatorColon
	for _, line := range match.lines {
		fmt.Fprintln(w,
			path+
				f.decorator.lineNumber(line.num)+
				SeparatorColon+
				f.decorator.match(match.pattern, line.text),
		)
	}
}

func newColorWriter(out io.Writer, opts Option) io.Writer {
	if opts.OutputOption.EnableColor {
		return ansicolor.NewAnsiColorWriter(out)
	}
	return out
}
