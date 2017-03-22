package the_platinum_searcher

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/monochromegane/conflag"
	"github.com/monochromegane/go-home"
	"github.com/monochromegane/terminal"
)

const version = "2.1.5"

const (
	ExitCodeOK = iota
	ExitCodeError
)

var opts Option

type PlatinumSearcher struct {
	Out, Err io.Writer
}

func (p PlatinumSearcher) Run(args []string) int {

	parser := newOptionParser(&opts)

	conflag.LongHyphen = true
	conflag.BoolValue = false
	for _, c := range [...]string{
		filepath.Join(xdgConfigHomeDir(), "pt", "config.toml"),
		filepath.Join(home.Dir(), ".ptconfig.toml"),
		".ptconfig.toml",
	} {
		if args, err := conflag.ArgsFrom(c); err == nil {
			parser.ParseArgs(args)
		}
	}

	args, err := parser.ParseArgs(args)
	if err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			return ExitCodeOK
		}
		return ExitCodeError
	}

	if opts.Version {
		fmt.Printf("pt version %s\n", version)
		return ExitCodeOK
	}

	if !terminal.IsTerminal(os.Stdout) {
		if !opts.OutputOption.ForceColor {
			opts.OutputOption.EnableColor = false
		}
		if !opts.OutputOption.ForceGroup {
			opts.OutputOption.EnableGroup = false
		}
	}

	var pathPattern, contentPattern string
	pipeMode := givenStdin()

	upp := regexp.MustCompile(`[[:upper:]]`)
	switch {
	case opts.SearchOption.EnableFilesWithRegexp && opts.SearchOption.FileSearchRegexp != "":
		fmt.Fprintf(p.Err, "ERR: (-g and -G) are exclusive!\n")
		fallthrough
	case len(args) == 0 && !opts.SearchOption.EnableFilesWithRegexp:
		parser.WriteHelp(p.Err)
		return ExitCodeError
	case opts.SearchOption.EnableFilesWithRegexp && !pipeMode: // g option
		opts.OutputOption.FilesWithMatches = true
		pathPattern = opts.SearchOption.PatternFilesWithRegexp
		if opts.SearchOption.IgnoreCase ||
			(opts.SearchOption.SmartCase && !upp.MatchString(pathPattern)) {
			opts.SearchOption.IgnoreCaseFilesWithRegexp = true
		}
	case pipeMode:
		opts.SearchOption.SearchStream = true
		contentPattern = strings.Join(args, " ")
		fallthrough
	case opts.SearchOption.FileSearchRegexp != "": // G option
		pathPattern = opts.SearchOption.FileSearchRegexp
		fallthrough
	default:
		if contentPattern == "" {
			contentPattern = args[0]
			args = args[1:]
		}
		if opts.SearchOption.SmartCase && !upp.MatchString(contentPattern) {
			opts.SearchOption.IgnoreCase = true
		}
		if opts.SearchOption.WordRegexp {
			opts.SearchOption.Regexp = true
			contentPattern = "\\b" + contentPattern + "\\b"
		}
		if opts.SearchOption.IgnoreCase {
			opts.SearchOption.Regexp = true
		}
	}

	if opts.OutputOption.Count {
		opts.OutputOption.Before = 0
		opts.OutputOption.After = 0
		opts.OutputOption.Context = 0
	} else if opts.OutputOption.Context > 0 {
		opts.OutputOption.Before = opts.OutputOption.Context
		opts.OutputOption.After = opts.OutputOption.Context
	}

	search := search{
		roots: rootsFrom(args),
		out:   p.Out,
	}
	if err = search.start(pathPattern, contentPattern); err != nil {
		fmt.Fprintf(p.Err, "%s\n", err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func rootsFrom(args []string) []string {
	if len(args) == 0 {
		return []string{"."}
	}
	return args
}

func givenStdin() bool {
	fi, err := os.Stdin.Stat()
	if runtime.GOOS == "windows" {
		if err == nil {
			return true
		}
	} else {
		if err != nil {
			return false
		}

		mode := fi.Mode()
		if (mode&os.ModeNamedPipe != 0) || mode.IsRegular() {
			return true
		}
	}
	return false
}

func xdgConfigHomeDir() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(home.Dir(), ".config")
	}
	return xdgConfigHome
}
