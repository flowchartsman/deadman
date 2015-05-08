package main

import (
	"os/exec"
)

func shutdownNow() error {
	if err := checkExe("shutdown"); err != nil {
		return err
	}
	// kill loginwindow and finder to limit access as much as possible quickly
	// thanks, pwnsdx!
	if err := exec.Command("killall", "loginwindow", "Finder").Run(); err != nil {
		return err
	}
	if err := exec.Command("halt", "-q").Run(); err != nil {
		return err
	}
	//Not that this matters
	return nil
}
