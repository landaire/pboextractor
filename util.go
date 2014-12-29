package main

import (
	"os"
	"path"
	"runtime"
	"strings"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func mkdirIntermediate(dir string) error {
	// god damn Windows users
	if runtime.GOOS == "windows" {
		dir = strings.Replace(dir, "\\", "/", -1)
	}

	dir = path.Clean(dir)
	parts := strings.Split(dir, "/")

	curdir := "/"
	for _, part := range parts {
		curdir = path.Join(curdir, part)

		exist, err := exists(curdir)
		if err != nil {
			return err
		}

		if !exist {
			err := os.Mkdir(curdir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
