package pt

import (
	"bufio"
	"code.google.com/p/mahonia"
	"fmt"
	"github.com/monochromegane/the_platinum_searcher/util"
	"os"
	"path/filepath"
	"strings"
)

type Searcher struct {
	Root, Pattern string
}

type GrepArgument struct {
	Path, Pattern, Encode string
}

type Match struct {
	LineNum int
	Match   string
}

type PrintArgument struct {
	Pattern string
	Path    string
	Matches []*Match
}

func (self *Searcher) Search() {
	grep := make(chan *GrepArgument, 2)
	match := make(chan *PrintArgument, 2)
	done := make(chan bool)
	go self.find(grep)
	go self.grep(grep, match)
	go self.print(match, done)
	<-done
}

func (self *Searcher) find(grep chan *GrepArgument) {
	filepath.Walk(self.Root, func(path string, info os.FileInfo, err error) error {
		if len(info.Name()) > 1 && strings.Index(info.Name(), ".") == 0 {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		fileType := pt.IdentifyFileType(path)
		if fileType == pt.BINARY {
			return nil
		}
		grep <- &GrepArgument{path, self.Pattern, fileType}
		return nil
	})
	grep <- nil
}

func (self *Searcher) grep(grep chan *GrepArgument, match chan *PrintArgument) {
	for {
		arg := <-grep
		if arg == nil {
			break
		}

		fh, err := os.Open(arg.Path)
		f := bufio.NewReader(fh)
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 1024)
		decoder := mahonia.NewDecoder(arg.Encode)

		m := make([]*Match, 0)

		var lineNum = 1
		for {
			buf, _, err = f.ReadLine()
			if err != nil {
				break
			}

			s := string(buf)
			if decoder != nil && arg.Encode != pt.UTF8 && arg.Encode != pt.ASCII {
				s = decoder.ConvertString(s)
			}
			if strings.Contains(s, arg.Pattern) {
				m = append(m, &Match{lineNum, s})
			}
			lineNum++
		}
		match <- &PrintArgument{arg.Pattern, arg.Path, m}
		fh.Close()

	}
	match <- nil
}

func (self *Searcher) print(match chan *PrintArgument, done chan bool) {
	for {
		arg := <-match
		if arg == nil {
			break
		}
		if len(arg.Matches) == 0 {
			continue
		}
		pt.PrintPath(arg.Path)
		fmt.Printf("\n")
		for _, v := range arg.Matches {
			if v == nil {
				continue
			}
			pt.PrintLineNumber(v.LineNum)
			pt.PrintMatch(arg.Pattern, v.Match)
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}
	done <- true
}
