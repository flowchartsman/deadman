package main

import (
	"os/exec"
)

func shutdownNow() error {
	if err := checkExe("shutdown"); err != nil {
		return err
	}
	err := exec.Command("shutdown", "/f", "/p").Run()

	//Not that this matters
	return err
}
