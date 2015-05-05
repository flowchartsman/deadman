package main

import (
	"os/exec"
)

func checkExe(path string) error {
	_, err := exec.LookPath(path)
	return err
	//TODO: check for executability
}
