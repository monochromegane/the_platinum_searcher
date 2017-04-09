package the_platinum_searcher

import (
	"fmt"
	"strings"
)

// fileTypeSearchPattern builds a regexp pattern for any requested file types.
// If no file types were requested, an empty string is returned.
func fileTypeSearchPattern() string {
	var exts []string

	if opts.FileTypeOption.ActionScript {
		exts = append(exts, []string{"as", "mxml"}...)
	}
	if opts.FileTypeOption.Ada {
		exts = append(exts, []string{"ada", "abd", "ads"}...)
	}
	if opts.FileTypeOption.Asm {
		exts = append(exts, []string{"asm", "s"}...)
	}
	if opts.FileTypeOption.Batch {
		exts = append(exts, []string{"bat", "cmd"}...)
	}
	if opts.FileTypeOption.Bitbake {
		exts = append(exts, []string{"bb", "bbappend", "bbclass", "inc"}...)
	}
	if opts.FileTypeOption.Bro {
		exts = append(exts, []string{"bro", "bif"}...)
	}
	if opts.FileTypeOption.CC {
		exts = append(exts, []string{"c", "h", "xs"}...)
	}
	if opts.FileTypeOption.Cfmx {
		exts = append(exts, []string{"cfc", "cfm", "cfml"}...)
	}
	if opts.FileTypeOption.Chpl {
		exts = append(exts, "chpl")
	}
	if opts.FileTypeOption.Clojure {
		exts = append(exts, []string{"clj", "cljs", "cljc", "cljx"}...)
	}
	if opts.FileTypeOption.Coffee {
		exts = append(exts, []string{"coffee", "cjsx"}...)
	}
	if opts.FileTypeOption.Cpp {
		exts = append(exts, []string{"cpp", "cc", "C", "cxx", "m", "hpp", "hh", "h", "H", "hxx", "tpp"}...)
	}
	if opts.FileTypeOption.Crystal {
		exts = append(exts, []string{"cr", "ecr"}...)
	}
	if opts.FileTypeOption.Csharp {
		exts = append(exts, "cs")
	}
	if opts.FileTypeOption.CSS {
		exts = append(exts, "css")
	}
	if opts.FileTypeOption.Cython {
		exts = append(exts, []string{"pyx", "pxd", "pxi"}...)
	}
	if opts.FileTypeOption.Delphi {
		exts = append(exts, []string{"pas", "int", "dfm", "nfm", "dof", "dpk", "dpr", "dproj", "groupproj", "bdsgroup", "bdsproj"}...)
	}
	if opts.FileTypeOption.Ebuild {
		exts = append(exts, []string{"ebuild", "eclass"}...)
	}
	if opts.FileTypeOption.Elisp {
		exts = append(exts, "el")
	}
	if opts.FileTypeOption.Elixer {
		exts = append(exts, []string{"ex", "eex", "exs"}...)
	}
	if opts.FileTypeOption.Erlang {
		exts = append(exts, []string{"erl", "hrl"}...)
	}
	if opts.FileTypeOption.Factor {
		exts = append(exts, "factor")
	}
	if opts.FileTypeOption.Fortran {
		exts = append(exts, []string{"f", "f77", "f90", "f95", "f03", "for", "ftn", "fpp"}...)
	}
	if opts.FileTypeOption.Fsharp {
		exts = append(exts, []string{"fs", "fsi", "fsx"}...)
	}
	if opts.FileTypeOption.GetText {
		exts = append(exts, []string{"po", "pot", "mo"}...)
	}
	if opts.FileTypeOption.Glsl {
		exts = append(exts, []string{"vert", "tesc", "tese", "geom", "frag", "comp"}...)
	}
	if opts.FileTypeOption.Go {
		exts = append(exts, "go")
	}
	if opts.FileTypeOption.Groovy {
		exts = append(exts, []string{"groovy", "gtmpl", "gpp", "grunit"}...)
	}
	if opts.FileTypeOption.Haml {
		exts = append(exts, "haml")
	}
	if opts.FileTypeOption.Haskell {
		exts = append(exts, []string{"hs", "lhs"}...)
	}
	if opts.FileTypeOption.HH {
		exts = append(exts, "h")
	}
	if opts.FileTypeOption.HTML {
		exts = append(exts, []string{"htm", "html", "shtml", "xhtml"}...)
	}
	if opts.FileTypeOption.INI {
		exts = append(exts, "ini")
	}
	if opts.FileTypeOption.Jade {
		exts = append(exts, "jade")
	}
	if opts.FileTypeOption.Java {
		exts = append(exts, []string{"java", "properties"}...)
	}
	if opts.FileTypeOption.JS {
		exts = append(exts, []string{"js", "jsx", "vue"}...)
	}
	if opts.FileTypeOption.JSON {
		exts = append(exts, "json")
	}
	if opts.FileTypeOption.Jsp {
		exts = append(exts, []string{"jsp", "jspx", "jhtm", "jhtml"}...)
	}
	if opts.FileTypeOption.Julia {
		exts = append(exts, "jl")
	}
	if opts.FileTypeOption.Kotlin {
		exts = append(exts, "kt")
	}
	if opts.FileTypeOption.Less {
		exts = append(exts, "less")
	}
	if opts.FileTypeOption.Liquid {
		exts = append(exts, "liquid")
	}
	if opts.FileTypeOption.Lisp {
		exts = append(exts, []string{"lisp", "lsp"}...)
	}
	if opts.FileTypeOption.Log {
		exts = append(exts, "log")
	}
	if opts.FileTypeOption.Lua {
		exts = append(exts, "lua")
	}
	if opts.FileTypeOption.M4 {
		exts = append(exts, "m4")
	}
	if opts.FileTypeOption.Make {
		exts = append(exts, []string{"Makefiles", "mk", "mak"}...)
	}
	if opts.FileTypeOption.Mako {
		exts = append(exts, "mako")
	}
	if opts.FileTypeOption.Markdown {
		exts = append(exts, []string{"markdown", "mdown", "mdwn", "mkdn", "mkd", "md"}...)
	}
	if opts.FileTypeOption.Mason {
		exts = append(exts, []string{"mas", "mhtml", "mpl", "mtxt"}...)
	}
	if opts.FileTypeOption.Matlab {
		exts = append(exts, "m")
	}
	if opts.FileTypeOption.Mathematica {
		exts = append(exts, []string{"m", "wl"}...)
	}
	if opts.FileTypeOption.Mercury {
		exts = append(exts, []string{"m", "moo"}...)
	}
	if opts.FileTypeOption.Nim {
		exts = append(exts, "nim")
	}
	if opts.FileTypeOption.ObjC {
		exts = append(exts, []string{"m", "h"}...)
	}
	if opts.FileTypeOption.ObjCpp {
		exts = append(exts, []string{"mm", "h"}...)
	}
	if opts.FileTypeOption.OCaml {
		exts = append(exts, []string{"ml", "mli", "mll", "mly"}...)
	}
	if opts.FileTypeOption.Octave {
		exts = append(exts, "m")
	}
	if opts.FileTypeOption.Parrot {
		exts = append(exts, []string{"pir", "pasm", "pmc", "ops", "pod", "pg", "tg"}...)
	}
	if opts.FileTypeOption.Perl {
		exts = append(exts, []string{"pl", "pm", "pm6", "pod", "t"}...)
	}
	if opts.FileTypeOption.PHP {
		exts = append(exts, []string{"php", "phpt", "php3", "php4", "php5", "phtml"}...)
	}
	if opts.FileTypeOption.Pike {
		exts = append(exts, []string{"pike", "pmod"}...)
	}
	if opts.FileTypeOption.Plone {
		exts = append(exts, []string{"pt", "cpt", "metadata", "cpy", "py", "xml", "zcml"}...)
	}
	if opts.FileTypeOption.Proto {
		exts = append(exts, "proto")
	}
	if opts.FileTypeOption.Puppet {
		exts = append(exts, "pp")
	}
	if opts.FileTypeOption.Python {
		exts = append(exts, "py")
	}
	if opts.FileTypeOption.QML {
		exts = append(exts, "qml")
	}
	if opts.FileTypeOption.Racket {
		exts = append(exts, []string{"rkt", "ss", "scm"}...)
	}
	if opts.FileTypeOption.Rake {
		exts = append(exts, "Rakefile")
	}
	if opts.FileTypeOption.RestructuredText {
		exts = append(exts, "rst")
	}
	if opts.FileTypeOption.RS {
		exts = append(exts, "rs")
	}
	if opts.FileTypeOption.R {
		exts = append(exts, []string{"R", "Rmd", "Rnw", "Rtex", "Rrst"}...)
	}
	if opts.FileTypeOption.Rdoc {
		exts = append(exts, "rdoc")
	}
	if opts.FileTypeOption.Ruby {
		exts = append(exts, []string{"rb", "rhtml", "rjs", "rxml", "erb", "rake", "spec"}...)
	}
	if opts.FileTypeOption.Rust {
		exts = append(exts, "rs")
	}
	if opts.FileTypeOption.Salt {
		exts = append(exts, "sls")
	}
	if opts.FileTypeOption.Sass {
		exts = append(exts, []string{"sass", "scss"}...)
	}
	if opts.FileTypeOption.Scala {
		exts = append(exts, "scala")
	}
	if opts.FileTypeOption.Scheme {
		exts = append(exts, []string{"scm", "ss"}...)
	}
	if opts.FileTypeOption.Shell {
		exts = append(exts, []string{"sh", "bash", "csh", "tcsh", "ksh", "zsh", "fish"}...)
	}
	if opts.FileTypeOption.Smalltalk {
		exts = append(exts, "st")
	}
	if opts.FileTypeOption.SML {
		exts = append(exts, []string{"sml", "fun", "mlb", "sig"}...)
	}
	if opts.FileTypeOption.SQL {
		exts = append(exts, []string{"sql", "ctl"}...)
	}
	if opts.FileTypeOption.Stylus {
		exts = append(exts, "styl")
	}
	if opts.FileTypeOption.Swift {
		exts = append(exts, "swift")
	}
	if opts.FileTypeOption.TCL {
		exts = append(exts, []string{"tcl", "itcl", "itk"}...)
	}
	if opts.FileTypeOption.Tex {
		exts = append(exts, []string{"tex", "cls", "sty"}...)
	}
	if opts.FileTypeOption.TT {
		exts = append(exts, []string{"tt", "tt2", "ttml"}...)
	}
	if opts.FileTypeOption.TOML {
		exts = append(exts, "toml")
	}
	if opts.FileTypeOption.TS {
		exts = append(exts, []string{"ts", "tsx"}...)
	}
	if opts.FileTypeOption.Vala {
		exts = append(exts, []string{"vala", "vapi"}...)
	}
	if opts.FileTypeOption.VB {
		exts = append(exts, []string{"bas", "cls", "frm", "ctl", "vb", "resx"}...)
	}
	if opts.FileTypeOption.Velocity {
		exts = append(exts, []string{"vm", "vtl", "vsl"}...)
	}
	if opts.FileTypeOption.Verilog {
		exts = append(exts, []string{"v", "vh", "sv"}...)
	}
	if opts.FileTypeOption.VHDL {
		exts = append(exts, []string{"vhd", "vhdl"}...)
	}
	if opts.FileTypeOption.Vim {
		exts = append(exts, "vim")
	}
	if opts.FileTypeOption.Wix {
		exts = append(exts, []string{"wxi", "wxs"}...)
	}
	if opts.FileTypeOption.WSDL {
		exts = append(exts, "wsdl")
	}
	if opts.FileTypeOption.WADL {
		exts = append(exts, "wadl")
	}
	if opts.FileTypeOption.XML {
		exts = append(exts, []string{"xml", "dtd", "xsl", "xslt", "ent", "tld"}...)
	}
	if opts.FileTypeOption.YAML {
		exts = append(exts, []string{"yaml", "yml"}...)
	}

	if len(exts) == 0 {
		return ""
	}
	return `\.(` + strings.Join(exts, "|") + `)$`
}

func printFileTypeOptions() {
	fmt.Printf(`The following file type options are available:

      --actionscript        .as .mxml
      --ada                 .ada .abd .ads
      --asm                 .asm .s
      --batch               .bat .cmd
      --bitbake             .bb .bbappend .bbclass .inc
      --bro                 .bro .bif
      --cc                  .c .h .xs
      --cfmx                .cfc .cfm .cfml
      --chpl                .chpl
      --clojure             .clj .cljs .cljc .cljx
      --coffee              .coffee .cjsx
      --cpp                 .cpp .cc .C .cxx .m .hpp .hh .h .H .hxx .tpp
      --crystal             .cr .ecr
      --csharp              .cs
      --css                 .css
      --cython              .pyx .pxd .pxi
      --delphi              .pas .int .dfm .nfm .dof .dpk .dpr .dproj .groupproj .bdsgroup .bdsproj
      --ebuild              .ebuild .eclass
      --elisp               .el
      --elixer              .ex .eex .exs
      --erlang              .erl .hrl
      --factor              .factor
      --fortran             .f .f77 .f90 .f95 .f03 .for .ftn .fpp
      --fsharp              .fs .fsi .fsx
      --gettext             .po .pot .mo
      --glsl                .vert .tesc .tese .geom .frag .comp
      --go                  .go
      --groovy              .groovy .gtmpl .gpp .grunit
      --haml                .haml
      --haskell             .hs .lhs
      --hh                  .h
      --html                .htm .html .shtml .xhtml
      --ini                 .ini
      --jade                .jade
      --java                .java .properties
      --js                  .js .jsx .vue
      --json                .json
      --jsp                 .jsp .jspx .jhtm .jhtml
      --julia               .jl
      --kotlin              .kt
      --less                .less
      --liquid              .liquid
      --lisp                .lisp .lsp
      --log                 .log
      --lua                 .lua
      --m4                  .m4
      --make                .Makefiles .mk .mak
      --mako                .mako
      --markdown            .markdown .mdown .mdwn .mkdn .mkd .md
      --mason               .mas .mhtml .mpl .mtxt
      --matlab              .m
      --mathematica         .m .wl
      --mercury             .m .moo
      --nim                 .nim
      --objc                .m .h
      --objcpp              .mm .h
      --ocaml               .ml .mli .mll .mly
      --octave              .m
      --parrot              .pir .pasm .pmc .ops .pod .pg .tg
      --perl                .pl .pm .pm6 .pod .t
      --php                 .php .phpt .php3 .php4 .php5 .phtml
      --pike                .pike .pmod
      --plone               .pt .cpt .metadata .cpy .py .xml .zcml
      --proto               .proto
      --puppet              .pp
      --python              .py
      --qml                 .qml
      --racket              .rkt .ss .scm
      --rake                .Rakefile
      --restructuredtext    .rst
      --rs                  .rs
      --r                   .R .Rmd .Rnw .Rtex .Rrst
      --rdoc                .rdoc
      --ruby                .rb .rhtml .rjs .rxml .erb .rake .spec
      --rust                .rs
      --salt                .sls
      --sass                .sass .scss
      --scala               .scala
      --scheme              .scm .ss
      --shell               .sh .bash .csh .tcsh .ksh .zsh .fish
      --smalltalk           .st
      --sml                 .sml .fun .mlb .sig
      --sql                 .sql .ctl
      --stylus              .styl
      --swift               .swift
      --tcl                 .tcl .itcl .itk
      --tex                 .tex .cls .sty
      --tt                  .tt .tt2 .ttml
      --toml                .toml
      --ts                  .ts .tsx
      --vala                .vala .vapi
      --vb                  .bas .cls .frm .ctl .vb .resx
      --velocity            .vm .vtl .vsl
      --verilog             .v .vh .sv
      --vhdl                .vhd .vhdl
      --vim                 .vim
      --wix                 .wxi .wxs
      --wsdl                .wsdl
      --wadl                .wadl
      --xml                 .xml .dtd .xsl .xslt .ent .tld
      --yaml                .yaml .yml
`)
}
