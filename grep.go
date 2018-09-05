package the_platinum_searcher

import (
	"math"
	"runtime"
	"sync"
)

var newLine = []byte("\n")

type grep struct {
	in      chan string
	done    chan struct{}
	grepper grepper
	printer printer
	opts    Option
}

func newGrep(pattern pattern, in chan string, done chan struct{}, opts Option, printer printer) grep {
	return grep{
		in:   in,
		done: done,
		grepper: newGrepper(
			pattern,
			printer,
			opts,
		),
		printer: printer,
		opts:    opts,
	}
}

func (g grep) start() {
	wg := &sync.WaitGroup{}
	worker := func() {
		defer wg.Done()
		buf := make([]byte, 16384)
		for path := range g.in {
			g.grepper.grep(path, buf)
		}
	}
	num := int(math.Max(float64(runtime.NumCPU()), 2.0))
	for i := 0; i < num; i++ {
		wg.Add(1)
		go worker()
	}

	wg.Wait()
	close(g.printer.in)
	g.done <- <-g.printer.done
}

type grepper interface {
	grep(path string, buf []byte)
}

func newGrepper(pattern pattern, printer printer, opts Option) grepper {
	if opts.SearchOption.EnableFilesWithRegexp {
		return passthroughGrep{
			printer: printer,
		}
	} else if opts.SearchOption.Regexp {
		return extendedGrep{
			pattern:  pattern,
			lineGrep: newLineGrep(printer, opts),
		}
	} else {
		if opts.OutputOption.Before > 0 || opts.OutputOption.After > 0 {
			return fixedGrep{
				pattern:  pattern,
				lineGrep: newLineGrep(printer, opts),
			}
		}
		return bufferGrep{
			pattern: pattern,
			printer: printer,
			column:  opts.OutputOption.Column,
		}
	}
}
