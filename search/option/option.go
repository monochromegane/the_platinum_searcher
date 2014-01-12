package option

type Option struct {
	NoColor   bool     `long:"nocolor" description:"Don't print color codes in results (Disabled by default)"`
	NoGroup   bool     `long:"nogroup" description:"Don't print file name at header (Disabled by default)"`
	VcsIgnore []string `long:"vcs-ignore" description:"VCS ignore files (Default: .gitignore, .hgignore)"`
        Ignore    []string `long:"ignore" description:"Ignore files/directories matching pattern"`
}

func (self *Option) VcsIgnores() []string {
	if len(self.VcsIgnore) == 0 {
		self.VcsIgnore = []string{".gitignore", ".hgignore"}
	}
	return self.VcsIgnore
}
