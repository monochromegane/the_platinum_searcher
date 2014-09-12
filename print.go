package the_platinum_searcher

import (
	"fmt"
	"io"
	"os"
	"strings"

	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
	"github.com/homburg/tree"
	"github.com/shiena/ansicolor"
)

var FileMatchCount, MatchCount uint

const (
	ColorReset      = "\x1b[0m\x1b[K"
	ColorLineNumber = "\x1b[1;33m"  /* yellow with black background */
	ColorPath       = "\x1b[1;32m"  /* bold green */
	ColorMatch      = "\x1b[30;43m" /* black with yellow background */
)

type PrintParams struct {
	Pattern *Pattern
	Path    string
	Matches []*Match
}

type print struct {
	In     chan *PrintParams
	Done   chan bool
	Option *Option
	writer io.Writer
}

func Print(in chan *PrintParams, done chan bool, option *Option) {
	print := &print{
		In:     in,
		Done:   done,
		Option: option,
		writer: createWriter(option),
	}
	print.Start()
}

func (p *print) Start() {
	FileMatchCount = 0
	MatchCount = 0

	var files []string

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

		if p.Option.Tree {
			files = append(files, arg.Path)
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

	if p.Option.Tree {
		t := tree.New("/")
		t.EatLines(files)
		fmt.Fprint(p.writer, t.Format())
	}

	if p.Option.Stats {
		fmt.Printf("%d Files Matched\n", FileMatchCount)
		fmt.Printf("%d Total Text Matches\n", MatchCount)
	}
	p.Done <- true
}

func (p *print) printPath(path string) {
	if p.Option.EnableColor {
		fmt.Fprintf(p.writer, "%s%s%s", ColorPath, path, ColorReset)
	} else {
		fmt.Fprintf(p.writer, "%s", path)
	}
	if !p.Option.FilesWithMatches && p.Option.FilesWithRegexp == "" {
		fmt.Fprintf(p.writer, ":")
	}
}

func (p *print) printLineNumber(lineNum int, sep string) {
	if p.Option.EnableColor {
		fmt.Fprintf(p.writer, "%s%d%s%s", ColorLineNumber, lineNum, ColorReset, sep)
	} else {
		fmt.Fprintf(p.writer, "%d%s", lineNum, sep)
	}
}

func (p *print) printMatch(pattern *Pattern, line *Line) {
	p.printLineNumber(line.Num, ":")
	if !p.Option.EnableColor {
		fmt.Fprintf(p.writer, "%s", line.Str)
	} else if pattern.UseRegexp || pattern.IgnoreCase {
		fmt.Fprintf(p.writer, "%s", pattern.Regexp.ReplaceAllString(line.Str, ColorMatch+"${1}"+ColorReset))
	} else {
		fmt.Fprintf(p.writer, "%s", strings.Replace(line.Str, pattern.Pattern, ColorMatch+pattern.Pattern+ColorReset, -1))
	}
}

func (p *print) printContext(lines []*Line) {
	for _, line := range lines {
		p.printLineNumber(line.Num, "-")
		fmt.Fprintf(p.writer, "%s", line.Str)
		fmt.Fprintln(p.writer)
	}
}

func createWriter(option *Option) io.Writer {
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
