package the_platinum_searcher

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"github.com/shiena/ansicolor"
)

var FileMatchCount, MatchCount uint

type PrintParams struct {
	Pattern *Pattern
	Path    string
	Matches []*Match
}

type print struct {
	In        chan *PrintParams
	Done      chan struct{}
	Option    *Option
	writer    io.Writer
	decorator Decorator
}

func Print(in chan *PrintParams, done chan struct{}, option *Option) {
	print := &print{
		In:        in,
		Done:      done,
		Option:    option,
		writer:    newWriter(option),
		decorator: newDecorator(option),
	}
	print.Start()
}

func (p *print) Start() {
	FileMatchCount = 0
	MatchCount = 0
	for arg := range p.In {

		if p.Option.FilesWithRegexp != "" {
			p.printPath(arg.Path)
			fmt.Println()
			FileMatchCount++
			continue
		}

		if len(arg.Matches) == 0 {
			continue
		}

		if p.Option.FilesWithMatches {
			p.printPath(arg.Path)
			fmt.Println()
			FileMatchCount++
			continue
		}
		if !p.Option.NoGroup {
			p.printPath(arg.Path)
			fmt.Println()
			FileMatchCount++
		}
		lastLineNum := 0
		enableContext := p.Option.Before > 0 || p.Option.After > 0
		for _, v := range arg.Matches {
			if v == nil {
				continue
			}
			if enableContext {
				if lastLineNum > 0 && lastLineNum+1 != v.FirstLineNum() {
					fmt.Println("--")
				}
				lastLineNum = v.LastLineNum()
			}
			if p.Option.NoGroup {
				p.printPath(arg.Path)
			}
			p.printContext(v.Befores)
			p.printMatch(arg.Pattern, v.Line)
			MatchCount++
			fmt.Println()
			p.printContext(v.Afters)
		}
		if !p.Option.NoGroup {
			fmt.Println()
		}
	}
	if p.Option.Stats {
		fmt.Printf("%d Files Matched\n", FileMatchCount)
		fmt.Printf("%d Total Text Matches\n", MatchCount)
	}
	p.Done <- struct{}{}
}

func (p *print) printPath(path string) {
	fmt.Fprint(p.writer, p.decorator.path(path))
	if !p.Option.FilesWithMatches && p.Option.FilesWithRegexp == "" {
		fmt.Fprintf(p.writer, ":")
	}
}

func (p *print) printLineNumber(lineNum int, sep string) {
	fmt.Fprint(p.writer, p.decorator.lineNumber(lineNum, sep))
}

func (p *print) printMatch(pattern *Pattern, line *Line) {
	p.printLineNumber(line.Num, ":")
	fmt.Fprint(p.writer, p.decorator.match(pattern, line))
}

func (p *print) printContext(lines []*Line) {
	for _, line := range lines {
		p.printLineNumber(line.Num, "-")
		fmt.Fprint(p.writer, line.Str)
		fmt.Fprintln(p.writer)
	}
}

func newWriter(option *Option) io.Writer {
	encoder := func() io.Writer {
		switch option.OutputEncode {
		case "sjis":
			return transform.NewWriter(os.Stdout, japanese.ShiftJIS.NewEncoder())
		case "euc":
			return transform.NewWriter(os.Stdout, japanese.EUCJP.NewEncoder())
		case "jis":
			return transform.NewWriter(os.Stdout, japanese.ISO2022JP.NewEncoder())
		default:
			return os.Stdout
		}
	}()
	if option.EnableColor {
		return ansicolor.NewAnsiColorWriter(encoder)
	}
	return encoder
}
