package the_platinum_searcher

import "github.com/jessevdk/go-flags"

// Top level options
type Option struct {
	Version        bool            `long:"version" description:"Show version"`
	OutputOption   *OutputOption   `group:"Output Options"`
	SearchOption   *SearchOption   `group:"Search Options"`
	FileTypeOption *FileTypeOption `group:"File Type Options"`
}

// Output options.
type OutputOption struct {
	Color               func()       `long:"color" description:"Print color codes in results (default: true)"`
	NoColor             func()       `long:"nocolor" description:"Don't print color codes in results (default: false)"`
	ForceColor          bool         // Force color. Not user option.
	EnableColor         bool         // Enable color. Not user option.
	ColorLineNumber     func(string) `long:"color-line-number" description:"Color codes for line numbers (default: 1;33)"`
	ColorPath           func(string) `long:"color-path" description:"Color codes for path names (default: 1;32)"`
	ColorMatch          func(string) `long:"color-match" description:"Color codes for result matches (default: 30;43)"`
	ColorCodeLineNumber string       // Color line numbers. Not user option.
	ColorCodePath       string       // Color path names. Not user option.
	ColorCodeMatch      string       // Color result matches. Not user option.
	Group               func()       `long:"group" description:"Print file name at header (default: true)"`
	NoGroup             func()       `long:"nogroup" description:"Don't print file name at header (default: false)"`
	ForceGroup          bool         // Force group. Not user option.
	EnableGroup         bool         // Enable group. Not user option.
	Null                bool         `short:"0" long:"null" description:"Separate filenames with null (for 'xargs -0') (default: false)"`
	Column              bool         `long:"column" description:"Print column (default: false)"`
	LineNumber          func()       `long:"numbers" description:"Print Line number. (default: true)"`
	NoLineNumber        func()       `short:"N" long:"nonumbers" description:"Omit Line number. (default: false)"`
	ForceLineNumber     bool         // Force line number. Not user option.
	EnableLineNumber    bool         // Enable line number. Not user option.
	After               int          `short:"A" long:"after" description:"Print lines after match"`
	Before              int          `short:"B" long:"before" description:"Print lines before match"`
	Context             int          `short:"C" long:"context" description:"Print lines before and after match"`
	FilesWithMatches    bool         `short:"l" long:"files-with-matches" description:"Only print filenames that contain matches"`
	Count               bool         `short:"c" long:"count" description:"Only print the number of matching lines for each input file."`
	OutputEncode        string       `short:"o" long:"output-encode" description:"Specify output encoding (none, jis, sjis, euc)"`
}

func newOutputOption() *OutputOption {
	opt := &OutputOption{}

	opt.Color = opt.SetEnableColor
	opt.NoColor = opt.SetDisableColor
	opt.EnableColor = true

	opt.Group = opt.SetEnableGroup
	opt.NoGroup = opt.SetDisableGroup
	opt.EnableGroup = true

	opt.LineNumber = opt.SetEnableLineNumber
	opt.NoLineNumber = opt.SetDisableLineNumber
	opt.EnableLineNumber = true

	opt.ColorLineNumber = opt.SetColorLineNumber
	opt.ColorPath = opt.SetColorPath
	opt.ColorMatch = opt.SetColorMatch
	opt.ColorCodeLineNumber = "1;33" // yellow with black background
	opt.ColorCodePath = "1;32"       // bold green
	opt.ColorCodeMatch = "30;43"     // black with yellow background

	return opt
}

func (o *OutputOption) SetEnableColor() {
	o.ForceColor = true
	o.EnableColor = true
}

func (o *OutputOption) SetDisableColor() {
	o.EnableColor = false
}

func (o *OutputOption) SetEnableLineNumber() {
	o.ForceLineNumber = true
	o.EnableLineNumber = true
}

func (o *OutputOption) SetDisableLineNumber() {
	o.EnableLineNumber = false
}

func (o *OutputOption) SetEnableGroup() {
	o.ForceGroup = true
	o.EnableGroup = true
}

func (o *OutputOption) SetDisableGroup() {
	o.EnableGroup = false
}

func (o *OutputOption) SetColorLineNumber(code string) {
	o.ColorCodeLineNumber = code
}

func (o *OutputOption) SetColorPath(code string) {
	o.ColorCodePath = code
}

func (o *OutputOption) SetColorMatch(code string) {
	o.ColorCodeMatch = code
}

// Search options.
type SearchOption struct {
	Regexp                 bool         `short:"e" description:"Parse PATTERN as a regular expression (default: false). Accepted syntax is the same as https://github.com/google/re2/wiki/Syntax except from \\C"`
	IgnoreCase             bool         `short:"i" long:"ignore-case" description:"Match case insensitively"`
	SmartCase              bool         `short:"S" long:"smart-case" description:"Match case insensitively unless PATTERN contains uppercase characters"`
	WordRegexp             bool         `short:"w" long:"word-regexp" description:"Only match whole words"`
	Ignore                 []string     `long:"ignore" description:"Ignore files/directories matching pattern"`
	VcsIgnore              []string     `long:"vcs-ignore" description:"VCS ignore files" default:".gitignore"`
	GlobalGitIgnore        bool         `long:"global-gitignore" description:"Use git's global gitignore file for ignore patterns"`
	HomePtIgnore           bool         `long:"home-ptignore" description:"Use $Home/.ptignore file for ignore patterns"`
	SkipVcsIgnore          bool         `short:"U" long:"skip-vcs-ignores" description:"Don't use VCS ignore file for ignore patterns"`
	FilesWithRegexp        func(string) `short:"g" description:"Print filenames matching PATTERN"`
	EnableFilesWithRegexp  bool         // Enable files with regexp. Not user option.
	PatternFilesWithRegexp string       // Pattern files with regexp. Not user option.
	FileSearchRegexp       string       `short:"G" long:"file-search-regexp" description:"PATTERN Limit search to filenames matching PATTERN"`
	Depth                  int          `long:"depth" default:"25" description:"Search up to NUM directories deep"`
	Follow                 bool         `short:"f" long:"follow" description:"Follow symlinks"`
	Hidden                 bool         `long:"hidden" description:"Search hidden files and directories"`
	SearchStream           bool         // Input from pipe. Not user option.
}

func (o *SearchOption) SetFilesWithRegexp(p string) {
	o.EnableFilesWithRegexp = true
	o.PatternFilesWithRegexp = p
}

func newSearchOption() *SearchOption {
	opt := &SearchOption{}
	opt.FilesWithRegexp = opt.SetFilesWithRegexp
	return opt
}

// File Type options.
type FileTypeOption struct {
	ListFileTypes bool `long:"list-file-types" description:"List available file type options"`

	ActionScript     bool `hidden:"true" long:"actionscript" description:".as .mxml"`
	Ada              bool `hidden:"true" long:"ada" description:".ada .abd .ads"`
	Asm              bool `hidden:"true" long:"asm" description:".asm .s"`
	Batch            bool `hidden:"true" long:"batch" description:".bat .cmd"`
	Bitbake          bool `hidden:"true" long:"bitbake" description:".bb .bbappend .bbclass .inc"`
	Bro              bool `hidden:"true" long:"bro" description:".bro .bif"`
	CC               bool `hidden:"true" long:"cc" description:".c .h .xs"`
	Cfmx             bool `hidden:"true" long:"cfmx" description:".cfc .cfm .cfml"`
	Chpl             bool `hidden:"true" long:"chpl" description:".chpl"`
	Clojure          bool `hidden:"true" long:"clojure" description:".clj .cljs .cljc .cljx"`
	Coffee           bool `hidden:"true" long:"coffee" description:".coffee .cjsx"`
	Cpp              bool `hidden:"true" long:"cpp" description:".cpp .cc .C .cxx .m .hpp .hh .h .H .hxx .tpp"`
	Crystal          bool `hidden:"true" long:"crystal" description:".cr .ecr"`
	Csharp           bool `hidden:"true" long:"csharp" description:".cs"`
	CSS              bool `hidden:"true" long:"css" description:".css"`
	Cython           bool `hidden:"true" long:"cython" description:".pyx .pxd .pxi"`
	Delphi           bool `hidden:"true" long:"delphi" description:".pas .int .dfm .nfm .dof .dpk .dpr .dproj .groupproj .bdsgroup .bdsproj"`
	Ebuild           bool `hidden:"true" long:"ebuild" description:".ebuild .eclass"`
	Elisp            bool `hidden:"true" long:"elisp" description:".el"`
	Elixer           bool `hidden:"true" long:"elixer" description:".ex .eex .exs"`
	Erlang           bool `hidden:"true" long:"erlang" description:".erl .hrl"`
	Factor           bool `hidden:"true" long:"factor" description:".factor"`
	Fortran          bool `hidden:"true" long:"fortran" description:".f .f77 .f90 .f95 .f03 .for .ftn .fpp"`
	Fsharp           bool `hidden:"true" long:"fsharp" description:".fs .fsi .fsx"`
	GetText          bool `hidden:"true" long:"gettext" description:".po .pot .mo"`
	Glsl             bool `hidden:"true" long:"glsl" description:".vert .tesc .tese .geom .frag .comp"`
	Go               bool `hidden:"true" long:"go" description:".go"`
	Groovy           bool `hidden:"true" long:"groovy" description:".groovy .gtmpl .gpp .grunit"`
	Haml             bool `hidden:"true" long:"haml" description:".haml"`
	Haskell          bool `hidden:"true" long:"haskell" description:".hs .lhs"`
	HH               bool `hidden:"true" long:"hh" description:".h"`
	HTML             bool `hidden:"true" long:"html" description:".htm .html .shtml .xhtml"`
	INI              bool `hidden:"true" long:"ini" description:".ini"`
	Jade             bool `hidden:"true" long:"jade" description:".jade"`
	Java             bool `hidden:"true" long:"java" description:".java .properties"`
	JS               bool `hidden:"true" long:"js" description:".js .jsx .vue"`
	JSON             bool `hidden:"true" long:"json" description:".json"`
	Jsp              bool `hidden:"true" long:"jsp" description:".jsp .jspx .jhtm .jhtml"`
	Julia            bool `hidden:"true" long:"julia" description:".jl"`
	Kotlin           bool `hidden:"true" long:"kotlin" description:".kt"`
	Less             bool `hidden:"true" long:"less" description:".less"`
	Liquid           bool `hidden:"true" long:"liquid" description:".liquid"`
	Lisp             bool `hidden:"true" long:"lisp" description:".lisp .lsp"`
	Log              bool `hidden:"true" long:"log" description:".log"`
	Lua              bool `hidden:"true" long:"lua" description:".lua"`
	M4               bool `hidden:"true" long:"m4" description:".m4"`
	Make             bool `hidden:"true" long:"make" description:".Makefiles .mk .mak"`
	Mako             bool `hidden:"true" long:"mako" description:".mako"`
	Markdown         bool `hidden:"true" long:"markdown" description:".markdown .mdown .mdwn .mkdn .mkd .md"`
	Mason            bool `hidden:"true" long:"mason" description:".mas .mhtml .mpl .mtxt"`
	Matlab           bool `hidden:"true" long:"matlab" description:".m"`
	Mathematica      bool `hidden:"true" long:"mathematica" description:".m .wl"`
	Mercury          bool `hidden:"true" long:"mercury" description:".m .moo"`
	Nim              bool `hidden:"true" long:"nim" description:".nim"`
	ObjC             bool `hidden:"true" long:"objc" description:".m .h"`
	ObjCpp           bool `hidden:"true" long:"objcpp" description:".mm .h"`
	OCaml            bool `hidden:"true" long:"ocaml" description:".ml .mli .mll .mly"`
	Octave           bool `hidden:"true" long:"octave" description:".m"`
	Parrot           bool `hidden:"true" long:"parrot" description:".pir .pasm .pmc .ops .pod .pg .tg"`
	Perl             bool `hidden:"true" long:"perl" description:".pl .pm .pm6 .pod .t"`
	PHP              bool `hidden:"true" long:"php" description:".php .phpt .php3 .php4 .php5 .phtml"`
	Pike             bool `hidden:"true" long:"pike" description:".pike .pmod"`
	Plone            bool `hidden:"true" long:"plone" description:".pt .cpt .metadata .cpy .py .xml .zcml"`
	Proto            bool `hidden:"true" long:"proto" description:".proto"`
	Puppet           bool `hidden:"true" long:"puppet" description:".pp"`
	Python           bool `hidden:"true" long:"python" description:".py"`
	QML              bool `hidden:"true" long:"qml" description:".qml"`
	Racket           bool `hidden:"true" long:"racket" description:".rkt .ss .scm"`
	Rake             bool `hidden:"true" long:"rake" description:".Rakefile"`
	RestructuredText bool `hidden:"true" long:"restructuredtext" description:".rst"`
	RS               bool `hidden:"true" long:"rs" description:".rs"`
	R                bool `hidden:"true" long:"r" description:".R .Rmd .Rnw .Rtex .Rrst"`
	Rdoc             bool `hidden:"true" long:"rdoc" description:".rdoc"`
	Ruby             bool `hidden:"true" long:"ruby" description:".rb .rhtml .rjs .rxml .erb .rake .spec"`
	Rust             bool `hidden:"true" long:"rust" description:".rs"`
	Salt             bool `hidden:"true" long:"salt" description:".sls"`
	Sass             bool `hidden:"true" long:"sass" description:".sass .scss"`
	Scala            bool `hidden:"true" long:"scala" description:".scala"`
	Scheme           bool `hidden:"true" long:"scheme" description:".scm .ss"`
	Shell            bool `hidden:"true" long:"shell" description:".sh .bash .csh .tcsh .ksh .zsh .fish"`
	Smalltalk        bool `hidden:"true" long:"smalltalk" description:".st"`
	SML              bool `hidden:"true" long:"sml" description:".sml .fun .mlb .sig"`
	SQL              bool `hidden:"true" long:"sql" description:".sql .ctl"`
	Stylus           bool `hidden:"true" long:"stylus" description:".styl"`
	Swift            bool `hidden:"true" long:"swift" description:".swift"`
	TCL              bool `hidden:"true" long:"tcl" description:".tcl .itcl .itk"`
	Tex              bool `hidden:"true" long:"tex" description:".tex .cls .sty"`
	TT               bool `hidden:"true" long:"tt" description:".tt .tt2 .ttml"`
	TOML             bool `hidden:"true" long:"toml" description:".toml"`
	TS               bool `hidden:"true" long:"ts" description:".ts .tsx"`
	Vala             bool `hidden:"true" long:"vala" description:".vala .vapi"`
	VB               bool `hidden:"true" long:"vb" description:".bas .cls .frm .ctl .vb .resx"`
	Velocity         bool `hidden:"true" long:"velocity" description:".vm .vtl .vsl"`
	Verilog          bool `hidden:"true" long:"verilog" description:".v .vh .sv"`
	VHDL             bool `hidden:"true" long:"vhdl" description:".vhd .vhdl"`
	Vim              bool `hidden:"true" long:"vim" description:".vim"`
	Wix              bool `hidden:"true" long:"wix" description:".wxi .wxs"`
	WSDL             bool `hidden:"true" long:"wsdl" description:".wsdl"`
	WADL             bool `hidden:"true" long:"wadl" description:".wadl"`
	XML              bool `hidden:"true" long:"xml" description:".xml .dtd .xsl .xslt .ent .tld"`
	YAML             bool `hidden:"true" long:"yaml" description:".yaml .yml"`
}

func newFileTypeOption() *FileTypeOption {
	opt := &FileTypeOption{}
	return opt
}

func newOptionParser(opts *Option) *flags.Parser {
	output := flags.NewNamedParser("pt", flags.Default)
	output.AddGroup("Output Options", "", &OutputOption{})

	search := flags.NewNamedParser("pt", flags.Default)
	search.AddGroup("Search Options", "", &SearchOption{})

	filetypes := flags.NewNamedParser("pt", flags.Default)
	filetypes.AddGroup("File Type Options", "", &FileTypeOption{})

	opts.OutputOption = newOutputOption()
	opts.SearchOption = newSearchOption()
	opts.FileTypeOption = newFileTypeOption()

	parser := flags.NewParser(opts, flags.Default)
	parser.Name = "pt"
	parser.Usage = "[OPTIONS] PATTERN [PATH]"
	return parser
}
