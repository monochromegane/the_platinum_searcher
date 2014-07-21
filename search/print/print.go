package print

import (
	"fmt"
	"io"
	"os"
	"strings"

	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
	"github.com/monochromegane/the_platinum_searcher/search/match"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"github.com/shiena/ansicolor"
)

var FileMatchCount, MatchCount uint

const (
	ColorReset      = "\x1b[0m\x1b[K"
	ColorLineNumber = "\x1b[1;33m"  /* yellow with black background */
	ColorPath       = "\x1b[1;32m"  /* bold green */
	ColorMatch      = "\x1b[30;43m" /* black with yellow background */
)

type Params struct {
	Pattern *pattern.Pattern
	Path    string
	Matches []*match.Match
}

type Printer struct {
	In     chan *Params
	Done   chan bool
	Option *option.Option
	writer io.Writer
}

func NewPrinter(in chan *Params, done chan bool, option *option.Option) *Printer {
	return &Printer{in, done, option, createWriter(option)}
}

func (p *Printer) Print() {
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
	p.Done <- true
}

func (p *Printer) printPath(path string) {
	if p.Option.EnableColor {
		fmt.Fprintf(p.writer, "%s%s%s", ColorPath, path, ColorReset)
	} else {
		fmt.Fprintf(p.writer, "%s", path)
	}
	if !p.Option.FilesWithMatches && p.Option.FilesWithRegexp == "" {
		fmt.Fprintf(p.writer, ":")
	}
}
func (p *Printer) printLineNumber(lineNum int, sep string) {
	if p.Option.EnableColor {
		fmt.Fprintf(p.writer, "%s%d%s%s", ColorLineNumber, lineNum, ColorReset, sep)
	} else {
		fmt.Fprintf(p.writer, "%d%s", lineNum, sep)
	}
}
func (p *Printer) printMatch(pattern *pattern.Pattern, line *match.Line) {
	p.printLineNumber(line.Num, ":")
	if !p.Option.EnableColor {
		fmt.Fprintf(p.writer, "%s", line.Str)
	} else if pattern.UseRegexp || pattern.IgnoreCase {
		fmt.Fprintf(p.writer, "%s", pattern.Regexp.ReplaceAllString(line.Str, ColorMatch+"${1}"+ColorReset))
	} else {
		fmt.Fprintf(p.writer, "%s", strings.Replace(line.Str, pattern.Pattern, ColorMatch+pattern.Pattern+ColorReset, -1))
	}
}

func (p *Printer) printContext(lines []*match.Line) {
	for _, line := range lines {
		p.printLineNumber(line.Num, "-")
		fmt.Fprintf(p.writer, "%s", line.Str)
		fmt.Fprintln(p.writer)
	}
}

func createWriter(option *option.Option) io.Writer {
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
