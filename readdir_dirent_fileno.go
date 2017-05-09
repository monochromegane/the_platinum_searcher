// +build freebsd openbsd netbsd

package the_platinum_searcher

import "syscall"

func direntInode(dirent *syscall.Dirent) uint64 {
	return uint64(dirent.Fileno)
}
