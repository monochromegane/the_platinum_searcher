package option

type Option struct {
	NoColor   bool     `long:"nocolor" description:"Don't print color codes in results (Disabled by default)"`
	NoGroup   bool     `long:"nogroup" description:"Don't print file name at header (Disabled by default)"`
	VcsIgnore []string `long:"vsc-ignore" description:"VCS ignore files (Default: .gitignore)"`
}

func (self *Option) VcsIgnores() []string {
	if len(self.VcsIgnore) == 0 {
		self.VcsIgnore = []string{".gitignore"}
	}
	return self.VcsIgnore
}
