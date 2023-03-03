package main

import (
	"os/exec"
)

func shutdownNow() error {
	if err := checkExe("shutdown"); err != nil {
		return err
	}
	err := exec.Command("shutdown", "-p", "now").Run()
	return err
}
