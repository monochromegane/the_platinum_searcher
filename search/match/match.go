package match

import (
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"strings"
)

type Match struct {
	Matched bool
	*Line
	beforeNum int
	afterNum  int
	Befores   []*Line
	Afters    []*Line
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

func (self *Match) newMatch() *Match {
	return &Match{
		beforeNum: self.beforeNum,
		afterNum:  self.afterNum,
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
	self.Matched = true
}

func (self *Match) setBefore(num int, s string) {
	befores := self.Befores
	if len(self.Befores) >= self.beforeNum {
		befores = self.Befores[1:]
	}
	self.Befores = append(befores, &Line{num, s})
}

func (self *Match) setAfter(num int, s string) bool {
	if len(self.Afters) >= self.afterNum {
		return false
	}
	self.Afters = append(self.Afters, &Line{num, s})
	return true
}

func (self *Match) setUpNewMatch(num int, s string) (*Match, bool) {
	// already match
	if self.Matched {
		newMatch := self.newMatch()
		newMatch.setMatch(num, s)
		return newMatch, true
	}
	self.setMatch(num, s)
	if self.afterNum == 0 {
		return self.newMatch(), true
	} else {
		// set after line
		return nil, false
	}
}

func (self *Match) IsMatch(pattern *pattern.Pattern, num int, s string) (*Match, bool) {
	if pattern.IgnoreCase {
		if strings.Contains(strings.ToUpper(s), strings.ToUpper(pattern.Pattern)) {
			return self.setUpNewMatch(num, s)
		}
	} else if strings.Contains(s, pattern.Pattern) {
		return self.setUpNewMatch(num, s)
	}
	if !self.Matched && self.beforeNum > 0 {
		self.setBefore(num, s)
	}
	if self.Matched && self.afterNum > 0 {
		if !self.setAfter(num, s) {
			newMatch := self.newMatch()
			if self.beforeNum > 0 {
				newMatch.setBefore(num, s)
			}
			return newMatch, true
		}
	}
	return nil, false
}

func (self *Match) FirstLineNum() int {
	if len(self.Befores) == 0 {
		return self.Line.Num
	} else {
		return self.Befores[0].Num
	}
}

func (self *Match) LastLineNum() int {
	if len(self.Afters) == 0 {
		return self.Line.Num
	} else {
		return self.Afters[len(self.Afters)-1].Num
	}
}
