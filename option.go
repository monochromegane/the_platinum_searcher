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
	Color            func() `long:"color" description:"Print color codes in results (default: true)"`
	NoColor          func() `long:"nocolor" description:"Don't print color codes in results (default: false)"`
	EnableColor      bool   // Enable color. Not user option.
	Group            func() `long:"group" description:"Print file name at header (default: true)"`
	NoGroup          func() `long:"nogroup" description:"Don't print file name at header (default: false)"`
	EnableGroup      bool   // Enable group. Not user option.
	After            int    `short:"A" long:"after" description:"Print lines after match"`
	Before           int    `short:"B" long:"before" description:"Print lines before match"`
	Context          int    `short:"C" long:"context" description:"Print lines before and after match"`
	FilesWithMatches bool   `short:"l" long:"files-with-matches" description:"Only print filenames that contain matches"`
	Count            bool   `short:"c" long:"count" description:"Only print the number of matches in each file."`
}

func newOutputOption() *OutputOption {
	opt := &OutputOption{}

	opt.Color = opt.SetEnableColor
	opt.NoColor = opt.SetDisableColor
	opt.EnableColor = true

	opt.Group = opt.SetEnableGroup
	opt.NoGroup = opt.SetDisableGroup
	opt.EnableGroup = true

	return opt
}

func (o *OutputOption) SetEnableColor() {
	o.EnableColor = true
}

func (o *OutputOption) SetDisableColor() {
	o.EnableColor = false
}

func (o *OutputOption) SetEnableGroup() {
	o.EnableGroup = true
}

func (o *OutputOption) SetDisableGroup() {
	o.EnableGroup = false
}

// Search options.
type SearchOption struct {
	Regexp        bool `short:"e" description:"Parse PATTERN as a regular expression (default: false). Accepted syntax is the same as https://github.com/google/re2/wiki/Syntax except from \\C"`
	IgnoreCase    bool `short:"i" long:"ignore-case" description:"Match case insensitively"`
	SmartCase     bool `short:"S" long:"smart-case" description:"Match case insensitively unless PATTERN contains uppercase characters"`
	WordRegexp    bool `short:"w" long:"word-regexp" description:"Only match whole words"`
	SkipVcsIgnore bool `short:"U" long:"skip-vcs-ignores" description:"Don't use VCS ignore file for ignore patterns"`
	Depth         int  `long:"depth" default:"25" description:"Search up to NUM directories deep"`
	Follow        bool `short:"f" long:"follow" description:"Follow symlinks"`
	Hidden        bool `long:"hidden" description:"Search hidden files and directories"`
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
