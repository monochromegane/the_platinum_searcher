package option

type Option struct {
	NoColor          bool     `long:"nocolor" description:"Don't print color codes in results (Disabled by default)"`
	NoGroup          bool     `long:"nogroup" description:"Don't print file name at header (Disabled by default)"`
	FilesWithMatches bool     `short:"l" long:"files-with-matches" description:"Only print filenames that don't contain matches"`
	VcsIgnore        []string `long:"vcs-ignore" description:"VCS ignore files (Default: .gitignore, .hgignore)"`
	Ignore           []string `long:"ignore" description:"Ignore files/directories matching pattern"`
	Depth            int      `long:"depth" description:"Search up to NUM derectories deep (Default: 25)"`
}

func (self *Option) VcsIgnores() []string {
	if len(self.VcsIgnore) == 0 {
		self.VcsIgnore = []string{".gitignore", ".hgignore"}
	}
	return self.VcsIgnore
}

func (self *Option) MaxDepth() int {
	if self.Depth == 0 {
		return 25
	}
	return self.Depth
}
