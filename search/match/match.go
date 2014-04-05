package match

import (
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"strings"
)

type Match struct {
	matched bool
	*Line
}

type Line struct {
	Num int
	Str string
}

func NewMatch(num int, str string) *Match {
	return &Match{
		Line: &Line{num, str},
	}
}

func (self *Match) LineNum() int {
	return self.Line.Num
}

func (self *Match) Match() string {
	return self.Line.Str
}

func (self *Match) IsMatch(pattern *pattern.Pattern, num int, s string) bool {
	if pattern.IgnoreCase {
		if pattern.Regexp.MatchString(s) {
			self.Line.Num = num
			self.Line.Str = s
			return true
		}
	} else if strings.Contains(s, pattern.Pattern) {
		self.Line.Num = num
		self.Line.Str = s
		return true
	}
	return false
}

