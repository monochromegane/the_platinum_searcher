// +build appengine !linux,!darwin,!freebsd,!openbsd,!netbsd

package the_platinum_searcher

import (
	"io/ioutil"
	"os"
)

func readDir(dirName string) ([]*fileInfo, error) {
	fis := []*fileInfo{}
	fs, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	for _, fi := range fs {
		fis = append(fis, newFileInfo(dirName, fi.Name(), fi.Mode()&os.ModeType))
	}
	return fis, nil
}
