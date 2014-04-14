package print

import (
	"fmt"
	"github.com/monochromegane/the_platinum_searcher/search/match"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"strings"
)

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
}

func (self *Printer) Print() {
	for arg := range self.In {

		if self.Option.FilesWithRegexp != "" {
			self.printPath(arg.Path)
			fmt.Println()
			continue
		}

		if len(arg.Matches) == 0 {
			continue
		}

		if self.Option.FilesWithMatches {
			self.printPath(arg.Path)
			fmt.Println()
			continue
		}
		if !self.Option.NoGroup {
			self.printPath(arg.Path)
			fmt.Println()
		}
		lastLineNum := 0
		enableContext := self.Option.Before > 0 || self.Option.After > 0
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
			if self.Option.NoGroup {
				self.printPath(arg.Path)
			}
			self.printContext(v.Befores)
			self.printMatch(arg.Pattern, v.Line)
			fmt.Println()
			self.printContext(v.Afters)
		}
		if !self.Option.NoGroup {
			fmt.Println()
		}
	}
	self.Done <- true
}

func (self *Printer) printPath(path string) {
	if self.Option.NoColor {
		fmt.Printf("%s", path)
	} else {
		fmt.Printf("%s%s%s", ColorPath, path, ColorReset)
	}
	if !self.Option.FilesWithMatches && self.Option.FilesWithRegexp == "" {
		fmt.Printf(":")
	}
}
func (self *Printer) printLineNumber(lineNum int, sep string) {
	if self.Option.NoColor {
		fmt.Printf("%d%s", lineNum, sep)
	} else {
		fmt.Printf("%s%d%s%s", ColorLineNumber, lineNum, ColorReset, sep)
	}
}
func (self *Printer) printMatch(pattern *pattern.Pattern, line *match.Line) {
	self.printLineNumber(line.Num, ":")
	if self.Option.NoColor {
		fmt.Printf("%s", line.Str)
	} else if pattern.IgnoreCase {
		fmt.Printf("%s", pattern.Regexp.ReplaceAllString(line.Str, ColorMatch+"${1}"+ColorReset))
	} else {
		fmt.Printf("%s", strings.Replace(line.Str, pattern.Pattern, ColorMatch+pattern.Pattern+ColorReset, -1))
	}
}

func (self *Printer) printContext(lines []*match.Line) {
	for _, line := range lines {
		self.printLineNumber(line.Num, "-")
		fmt.Printf("%s", line.Str)
		fmt.Println()
	}
}
