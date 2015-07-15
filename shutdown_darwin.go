package main

import (
	"os/exec"
)

func shutdownNow() error {
	if err := checkExe("shutdown"); err != nil {
		return err
	}
	// kill loginwindow and finder to limit access as much as possible quickly
	// thanks, pwnsdx! If there's an error with this, we're still going to try and
	//shut down, so  we really don't care
	_ = exec.Command("killall", "loginwindow", "Finder").Run()

	if err := exec.Command("halt", "-q").Run(); err != nil {
		return err
	}
	//Not that this matters
	return nil
}
