package match

type Match struct {
	matched bool
	*Line
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

type Line struct {
	Num int
	Str string
}
