package match

import (
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"strings"
)

type Match struct {
	matched bool
	*Line
	beforeNum int
	afterNum  int
	Befores   []*Line
}

type Line struct {
	Num int
	Str string
}

func NewMatch(before, after int) *Match {
	return &Match{
		beforeNum: before,
		afterNum:  after,
		Line:      &Line{},
	}
}

func (self *Match) LineNum() int {
	return self.Line.Num
}

func (self *Match) Match() string {
	return self.Line.Str
}

func (self *Match) setMatch(num int, s string) {
	self.Line.Num = num
	self.Line.Str = s
}

func (self *Match) setBefore(num int, s string) {
	befores := self.Befores
	if len(self.Befores) >= self.beforeNum {
		befores = self.Befores[1:]
	}
	self.Befores = append(befores, &Line{num, s})
}

func (self *Match) IsMatch(pattern *pattern.Pattern, num int, s string) bool {
	if pattern.IgnoreCase {
		if pattern.Regexp.MatchString(s) {
			self.setMatch(num, s)
			return true
		}
	} else if strings.Contains(s, pattern.Pattern) {
		self.setMatch(num, s)
		return true
	}
	if self.beforeNum > 0 {
		self.setBefore(num, s)
	}
	return false
}
