package pt

import (
	"bufio"
	"fmt"
	"github.com/monochromegane/the_platinum_searcher/util"
	"os"
	"path/filepath"
	"strings"
)

type Searcher struct {
	Root, Pattern string
}

func (self *Searcher) Search() {
	grep := make(chan string, 2)
	match := make(chan string, 2)
	done := make(chan bool)
	go self.find(self.Root, grep)
	go self.grep(self.Pattern, grep, match)
	go self.print(match, done)
	<-done
}

func (self *Searcher) find(root string, grep chan string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
                fileType := pt.IdentifyFileType(path)
                if fileType == pt.BINARY {
                        return nil
                }
		grep <- path
		return nil
	})
	grep <- "end"
}

func (self *Searcher) grep(pattern string, grep chan string, match chan string) {
	for {
		path := <-grep
		if path == "end" {
			break
		}

		fh, err := os.Open(path)
		f := bufio.NewReader(fh)
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 1024)

		for {
			buf, _, err = f.ReadLine()
			if err != nil {
				break
			}

			s := string(buf)
			if strings.Contains(s, pattern) {
				match <- s
			}
		}
                fh.Close()

	}
	match <- "end"
}

func (self *Searcher) print(match chan string, done chan bool) {
	for {
		matched := <-match
		if matched == "end" {
			break
		}
		fmt.Printf("matched: %s\n", matched)
	}
	done <- true
}
