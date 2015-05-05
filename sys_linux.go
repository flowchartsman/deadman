package main

import (
	"os/exec"
	"regexp"
)

var parseRE = regexp.MustCompile(`(?m)^.*ID (\S+) (.*)$`)

func enumerateDevices() ([]device, error) {
	if err := checkExe("lsusb"); err != nil {
		return nil, err
	}
	out, err := exec.Command("lsusb").Output()
	if err != nil {
		return nil, err
	}

	var deviceList []device
	devices := parseRE.FindAllStringSubmatch(string(out), -1)

	for _, d := range devices {
		deviceList = append(deviceList, device{d[2], d[1]})
	}
	return deviceList, nil
}

func shutdownNow() error {
	if err := checkExe("shutdown"); err != nil {
		return err
	}
	err := exec.Command("shutdown", "-h", "now").Run()

	//Not that this matters
	return err
}
