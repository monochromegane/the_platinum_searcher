package the_platinum_searcher

import "github.com/jessevdk/go-flags"

// Top level options
type Option struct {
	Version      bool          `long:"version" description:"Show version"`
	OutputOption *OutputOption `group:"Output Options"`
	SearchOption *SearchOption `group:"Search Options"`
}

// Output options.
type OutputOption struct {
	Color       func() `long:"color" description:"Print color codes in results (default: true)"`
	NoColor     func() `long:"nocolor" description:"Don't print color codes in results (default: false)"`
	EnableColor bool   // Enable color. Not user option.
}

func newOutputOption() *OutputOption {
	opt := &OutputOption{}
	opt.Color = opt.SetEnableColor
	opt.NoColor = opt.SetDisableColor
	opt.EnableColor = true
	return opt
}

func (o *OutputOption) SetEnableColor() {
	o.EnableColor = true
}

func (o *OutputOption) SetDisableColor() {
	o.EnableColor = false
}

// Search options.
type SearchOption struct {
	Depth int `long:"depth" default:"25" description:"Search up to NUM directories deep"`
}

func newOptionParser(opts *Option) *flags.Parser {
	output := flags.NewNamedParser("pt", flags.Default)
	output.AddGroup("Output Options", "", &OutputOption{})

	search := flags.NewNamedParser("pt", flags.Default)
	search.AddGroup("Search Options", "", &SearchOption{})

	opts.OutputOption = newOutputOption()

	parser := flags.NewParser(opts, flags.Default)
	parser.Name = "pt"
	parser.Usage = "[OPTIONS] PATTERN [PATH]"
	return parser
}
