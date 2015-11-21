package the_platinum_searcher

import "github.com/jessevdk/go-flags"

type Option struct {
	Version      bool         `long:"version" description:"Show version"`
	OutputOption OutputOption `group:"Output Options"`
	SearchOption SearchOption `group:"Search Options"`
}

type OutputOption struct {
}

type SearchOption struct {
}

func newOptionParser(opts *Option) *flags.Parser {
	output := flags.NewNamedParser("pt", flags.Default)
	output.AddGroup("Output Options", "", &OutputOption{})

	search := flags.NewNamedParser("pt", flags.Default)
	search.AddGroup("Search Options", "", &SearchOption{})

	parser := flags.NewParser(opts, flags.Default)
	parser.Name = "pt"
	parser.Usage = "[OPTIONS] PATTERN [PATH]"
	return parser
}
