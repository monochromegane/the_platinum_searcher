package the_platinum_searcher

type PlatinumSearcher struct {
	Root, Pattern string
	Option        *Option
}

func (p *PlatinumSearcher) Search() error {
	pattern, err := p.pattern()
	if err != nil {
		return err
	}
	grep := make(chan *GrepParams, p.Option.Proc)
	match := make(chan *PrintParams, p.Option.Proc)
	done := make(chan bool)
	go p.find(grep, pattern)
	go p.grep(grep, match)
	go p.print(match, done)
	<-done
	return nil
}

func (p *PlatinumSearcher) pattern() (*Pattern, error) {
	fileRegexp := p.Option.FileSearchRegexp
	if p.Option.FilesWithRegexp != "" {
		fileRegexp = p.Option.FilesWithRegexp
	}
	return NewPattern(
		p.Pattern,
		fileRegexp,
		p.Option.SmartCase,
		p.Option.IgnoreCase,
		p.Option.Regexp,
	)
}

func (p *PlatinumSearcher) find(out chan *GrepParams, pattern *Pattern) {
	Find(p.Root, pattern, out, p.Option)
}

func (p *PlatinumSearcher) grep(in chan *GrepParams, out chan *PrintParams) {
	Grep(in, out, p.Option)
}

func (p *PlatinumSearcher) print(in chan *PrintParams, done chan bool) {
	printer := NewPrinter(in, done, p.Option)
	printer.Print()
}
