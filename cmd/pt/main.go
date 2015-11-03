package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

func main() {
	grepChan := make(chan string, 32)
	done := make(chan struct{})

	go func() {
		sem := make(chan struct{}, 16)
		for path := range grepChan {
			sem <- struct{}{}
			go func(path string) {
				readFromMmap(path)
				<-sem
			}(path)

		}
		done <- struct{}{}
	}()

	Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}
		grepChan <- path
		return nil
	})
	close(grepChan)
	<-done
}

type WalkFunc func(path string, info os.FileInfo, err error) error

func Walk(root string, walkFn WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	sem := make(chan struct{}, 8)
	return walk(root, info, walkFn, sem)
}

func walk(path string, info os.FileInfo, walkFn WalkFunc, sem chan struct{}) error {
	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return walkFn(path, info, err)
	}

	wg := &sync.WaitGroup{}
	for _, file := range files {
		select {
		case sem <- struct{}{}:
			wg.Add(1)
			go func(path string, file os.FileInfo, wg *sync.WaitGroup) {
				defer wg.Done()
				defer func() { <-sem }()
				walk(path, file, walkFn, sem)
			}(filepath.Join(path, file.Name()), file, wg)
		default:
			walk(filepath.Join(path, file.Name()), file, walkFn, sem)
		}
	}
	wg.Wait()
	return nil
}

func readFromMmap(path string) {
	pattern := []byte(os.Args[2])
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}

	fi, err := f.Stat()
	if err != nil {
		log.Fatalf("stat: %s\n", err)
	}

	if int(fi.Size()) > 0 {
		mem, err := syscall.Mmap(int(f.Fd()), 0, int(fi.Size()),
			syscall.PROT_READ, syscall.MAP_SHARED)
		if err != nil {
			log.Fatalf("Mmap: %s %s\n", path, err)
		}

		if bytes.Index(mem, pattern) >= 0 {
			scanner := bufio.NewScanner(bytes.NewReader(mem))
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), os.Args[2]) {
				}
			}
		}

		err = syscall.Munmap(mem)
		if err != nil {
			log.Fatal(err)
		}
	}
	f.Close()

}
