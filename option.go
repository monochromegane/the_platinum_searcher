package the_platinum_searcher

type Option struct {
	Color             func()   `long:"color" description:"Print color codes in results (Enabled by default)"`
	NoColor           func()   `long:"nocolor" description:"Don't print color codes in results (Disabled by default)"`
	ForceColor        bool     // Force color.  Not user option.
	EnableColor       bool     // Enable color. Not user option.
	NoGroup           bool     `long:"nogroup" description:"Don't print file name at header (Disabled by default)"`
	Column            bool     `long:"column" description:"Print column (Disabled by default)"`
	FilesWithMatches  bool     `short:"l" long:"files-with-matches" description:"Only print filenames that contain matches"`
	VcsIgnore         []string `long:"vcs-ignore" description:"VCS ignore files" default:".gitignore" default:".hgignore" default:".ptignore"`
	NoPtIgnore        bool     `long:"noptignore" description:"Don't use default ($Home/.ptignore) file for ignore patterns"`
	NoGlobalGitIgnore bool     `long:"noglobal-gitignore" description:"Don't use git's global gitignore file for ignore patterns"`
	SkipVcsIgnore     func()   `short:"U" long:"skip-vsc-ignores" description:"Don't use VCS ignore file for ignore patterns. Still obey .ptignore"`
	skipVcsIgnore     bool     // Skip VCS ignore file. Not user option.
	Hidden            bool     `short:"H" long:"hidden" description:"Search hidden files and directories"`
	Ignore            []string `long:"ignore" description:"Ignore files/directories matching pattern"`
	IgnoreCase        bool     `short:"i" long:"ignore-case" description:"Match case insensitively"`
	SmartCase         bool     `short:"S" long:"smart-case" description:"Match case insensitively unless PATTERN contains uppercase characters"`
	FilesWithRegexp   string   `short:"g" description:"Print filenames matching PATTERN"`
	FileSearchRegexp  string   `short:"G" long:"file-search-regexp" description:"PATTERN Limit search to filenames matching PATTERN"`
	Depth             int      `long:"depth" default:"25" default-mask:"-" description:"Search up to NUM directories deep (Default: 25)"`
	Follow            bool     `short:"f" long:"follow" description:"Follow symlinks"`
	After             int      `short:"A" long:"after" description:"Print lines after match"`
	Before            int      `short:"B" long:"before" description:"Print lines before match"`
	Context           int      `short:"C" long:"context" description:"Print lines before and after match"`
	OutputEncode      string   `short:"o" long:"output-encode" description:"Specify output encoding (none, jis, sjis, euc)"`
	SearchStream      bool     // Input from pipe. Not user option.
	Regexp            bool     `short:"e" description:"Parse PATTERN as a regular expression (Disabled by default)"`
	WordRegexp        bool     `short:"w" long:"word-regexp" description:"Only match whole words"`
	Proc              int      // Number of goroutine. Not user option.
	Stats             bool     `long:"stats" description:"Print stats about files scanned, time taken, etc"`
	Parallel          bool     `long:"parallel" description:"Use as many concurrent finders as possible, this will lead the result disorder"`
	Version           bool     `long:"version" description:"Show version"`
}

func (o *Option) VcsIgnores() []string {
	if o.skipVcsIgnore {
		return []string{}
	}
	return o.VcsIgnore
}

func (o *Option) SkipVcsIgnores() {
	o.skipVcsIgnore = true
	o.NoGlobalGitIgnore = true
}

func (o *Option) SetEnableColor() {
	o.ForceColor = true
	o.EnableColor = true
}

func (o *Option) SetDisableColor() {
	o.EnableColor = false
}
