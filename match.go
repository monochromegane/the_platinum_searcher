package the_platinum_searcher

import (
	"strings"

	"github.com/monochromegane/the_platinum_searcher/search/pattern"
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

func (m *Match) newMatch() *Match {
	return &Match{
		beforeNum: m.beforeNum,
		afterNum:  m.afterNum,
		Line:      &Line{},
	}
}

func (m *Match) LineNum() int {
	return m.Line.Num
}

func (m *Match) Match() string {
	return m.Line.Str
}

func (m *Match) setMatch(num int, s string) {
	m.Line.Num = num
	m.Line.Str = s
	m.Matched = true
}

func (m *Match) setBefore(num int, s string) {
	befores := m.Befores
	if len(m.Befores) >= m.beforeNum {
		befores = m.Befores[1:]
	}
	m.Befores = append(befores, &Line{num, s})
}

func (m *Match) setAfter(num int, s string) bool {
	if len(m.Afters) >= m.afterNum {
		return false
	}
	m.Afters = append(m.Afters, &Line{num, s})
	return true
}

func (m *Match) setUpNewMatch(num int, s string) (*Match, bool) {
	// already match
	if m.Matched {
		newMatch := m.newMatch()
		newMatch.setMatch(num, s)
		return newMatch, true
	}
	m.setMatch(num, s)
	if m.afterNum == 0 {
		return m.newMatch(), true
	} else {
		// set after line
		return nil, false
	}
}

func (m *Match) IsMatch(pattern *pattern.Pattern, num int, s string) (*Match, bool) {
	if pattern.UseRegexp {
		if pattern.Regexp.MatchString(s) {
			return m.setUpNewMatch(num, s)
		}
	} else if pattern.IgnoreCase {
		if strings.Contains(strings.ToUpper(s), strings.ToUpper(pattern.Pattern)) {
			return m.setUpNewMatch(num, s)
		}
	} else if strings.Contains(s, pattern.Pattern) {
		return m.setUpNewMatch(num, s)
	}
	if !m.Matched && m.beforeNum > 0 {
		m.setBefore(num, s)
	}
	if m.Matched && m.afterNum > 0 {
		if !m.setAfter(num, s) {
			newMatch := m.newMatch()
			if m.beforeNum > 0 {
				newMatch.setBefore(num, s)
			}
			return newMatch, true
		}
	}
	return nil, false
}

func (m *Match) FirstLineNum() int {
	if len(m.Befores) == 0 {
		return m.Line.Num
	} else {
		return m.Befores[0].Num
	}
}

func (m *Match) LastLineNum() int {
	if len(m.Afters) == 0 {
		return m.Line.Num
	} else {
		return m.Afters[len(m.Afters)-1].Num
	}
}
