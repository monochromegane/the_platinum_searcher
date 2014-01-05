package find

import (
	"github.com/monochromegane/the_platinum_searcher/search/file"
	"github.com/monochromegane/the_platinum_searcher/search/grep"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"os"
	"path/filepath"
	"strings"
)

type Finder struct {
	Out    chan *grep.Params
	Option *option.Option
}

func (self *Finder) Find(root, pattern string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if len(info.Name()) > 1 && strings.Index(info.Name(), ".") == 0 {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
		fileType := file.IdentifyType(path)
		if fileType == file.BINARY {
			return nil
		}
		self.Out <- &grep.Params{path, pattern, fileType}
		return nil
	})
	close(self.Out)
}
