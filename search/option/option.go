package option

type Option struct {
	NoColor bool `long:"nocolor" description:"Don't print color codes in results (Disabled by default)"`
	NoGroup bool `long:"nogroup" description:"Don't print file name at header (Disabled by default)"`
}
