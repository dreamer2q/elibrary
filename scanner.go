package main

import (
	"io/ioutil"
	"path/filepath"
	"sync"
)

func walkDir(dir string, n *sync.WaitGroup, bookPath chan<- string, errCh chan<- error, filter map[string]bool) {
	defer n.Done()

	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		errCh <- err
		return
	}
	for _, d := range dirs {
		subDir := filepath.Join(dir, d.Name())
		if d.IsDir() {
			n.Add(1)
			walkDir(subDir, n, bookPath, errCh, filter)
		} else {
			if filter[filepath.Ext(d.Name())] {
				bookPath <- subDir
			}
		}
	}
}
