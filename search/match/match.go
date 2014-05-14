package match

import (
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"strings"
	"bytes"
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

func (self *Match) FindMatches(pattern *pattern.Pattern, buf []byte, matches *[]*Match) {
	var matchIndexes [][]int
	if pattern.IgnoreCase {
		matchIndexes = pattern.Regexp.FindAllIndex(buf, -1)
	} else {
		for i := 0; i < len(buf); {
			matchIndex := make([]int, 2)
			matchIndex[0] = bytes.Index(buf[i:], []byte (pattern.Pattern))
			if (-1 == matchIndex[0]) { break }
			matchIndex[0] = matchIndex[0] + i
			matchIndex[1] = matchIndex[0] + len(pattern.Pattern)
			matchIndexes = append(matchIndexes, matchIndex)
			i = matchIndex[1] + 1
		}
	}
	// Return right away if there were no matches found
	if (nil == matchIndexes) {
		return
	}
	// Found a match so find newlines
	tempIndex := 0
	line := 0
	lineStartIndex := 0
	var prevMatch *Match
	for i := 0; i < len (matchIndexes); i++ {
		currentMatch := NewMatch(self.beforeNum, self.afterNum)
		for {
			//If no more new lines before the current match, go to next match loop
			lineEndIndex := bytes.Index(buf[tempIndex:matchIndexes[i][0]], []byte("\n"))
			if ( -1 == lineEndIndex) { break }
			lineEndIndex = lineStartIndex + lineEndIndex
			line++
			if self.beforeNum > 0 {
				currentMatch.setBefore(line, string (buf[lineStartIndex:lineEndIndex]))
			}
			if (self.afterNum > 0) && (prevMatch != nil) {
				prevMatch.setAfter(line, string (buf[lineStartIndex:lineEndIndex]))
			}
			lineStartIndex = lineEndIndex + 1
			tempIndex = lineStartIndex
		}
		line++
		//Setup new match and append here
		currentMatch.Matched = true
		currentMatch.setMatch(line, string (buf[lineStartIndex:matchIndexes[i][1]]))
		*matches = append(*matches, currentMatch)
		prevMatch = currentMatch
		//Search for newlines within the match (multiline regex case)
		for j := lineStartIndex; j < matchIndexes[i][1]; {
			tempIndex2 := bytes.Index(buf[j:matchIndexes[i][1]], []byte("\n"))
			if( -1 == tempIndex2 ) { break }
			line++
			j = j + tempIndex2 + 1
		}
		//Next newline search should start from current match's end
		tempIndex = matchIndexes[i][1] + 1
		lineStartIndex = tempIndex
	}
	//After context for last match should be done here as a special case
	lastAfterCtr := 0
	for i := lineStartIndex; (i < len(buf)) && (lastAfterCtr < self.afterNum) ; {
		idx := bytes.Index(buf[i:], []byte("\n"))
		if (-1 == idx) { break }
		line++
		lastAfterCtr++
		prevMatch.setAfter(line, string (buf[i:i + idx]))
		i = i + idx + 1
	}
}
