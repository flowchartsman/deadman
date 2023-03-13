package main

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func checkExe(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return errors.New("Path is a directory")
	}

	if runtime.GOOS == "windows" {
		if filepath.Ext(path) != ".exe" && filepath.Ext(path) != ".com" {
			return errors.New("Path is not an executable")
		}
	} else {
		if info.Mode()&0111 == 0 {

			return errors.New("Path is not an executable")
		}
	}

	return nil
}
