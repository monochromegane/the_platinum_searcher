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

func (self *Match) setBeforePrepend(num int, s string) {
	self.Befores = append(self.Befores, nil)
	copy(self.Befores[1:], self.Befores[:])
	self.Befores[0] = &Line{num, s}
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
		if pattern.Regexp.MatchString(s) {
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

func (self *Match) FindMatches(pattern *pattern.Pattern, buf string, matches *[]*Match) {
	var splitStrings []string
	var line = 1
	var tempBuf = buf
	var lastOffset = 0
	var curPtr = 0
	var numNewLines = 0
	m := NewMatch(self.beforeNum, self.afterNum)
	for {
		if pattern.IgnoreCase {
      tempIndex := pattern.Regexp.FindStringIndex(tempBuf)
      if (nil != tempIndex) {
				splitStrings = append(splitStrings, tempBuf[:tempIndex[1]])
				splitStrings = append(splitStrings, tempBuf[tempIndex[1]:len(tempBuf)])
      }
		} else {
			splitStrings = strings.SplitAfterN(tempBuf, pattern.Pattern, 2)
		}
		if(2 > len(splitStrings)) {
			break
		}
    sOffset := strings.LastIndex(splitStrings[0], "\n")
    s1Index := strings.Index(splitStrings[1], "\n")
    if (-1 == s1Index) { s1Index = 0 }
		lastOffset = len(splitStrings[0]) + s1Index
		curPtr = curPtr + len(splitStrings[0])
		s := tempBuf[sOffset + 1:lastOffset]
		newLines := strings.SplitAfter(buf[:curPtr], "\n")
		if (nil == newLines) { numNewLines = 0 } else { numNewLines = len(newLines) }
		line = numNewLines
		m.Matched = true
		newMatch, _ := m.setUpNewMatch(line, s)
		tempIdx := curPtr
		
		for i := 1; i <= self.beforeNum; i++ {
			tempIdx = strings.LastIndex(buf[:tempIdx], "\n")
			if (-1 == tempIdx) { break }
			tempIdx2 := strings.LastIndex(buf[:tempIdx], "\n")
			if (-1 == tempIdx2) { tempIdx2 = 0 } else { tempIdx2++ }
			m.setBeforePrepend(line - i, buf[tempIdx2:tempIdx])
		}

		curPtr = curPtr + s1Index + 1
		newLines = strings.SplitAfter(buf[curPtr:], "\n")
		if(nil != newLines) {
			for i:= 0; (i < self.afterNum) && (i < len(newLines)); i++ {
				if false == m.setAfter(line + i, newLines[i]) { break }
			}
		}
		*matches = append(*matches, m)
		m = newMatch
    if (len(splitStrings[1]) > s1Index) {
      tempBuf = splitStrings[1][s1Index + 1:len(splitStrings[1])]
    } else {
      break
    }

    splitStrings = nil
	}
}
