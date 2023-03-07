package main

import (
	"fmt"
	"os/exec"
	"regexp"
)

var parseRE = regexp.MustCompile(`(?m)^.*ID (\S+) (.*)$`)

func enumerateDevices() ([]device, error) {

	application := "lsusb"

	// Get the path to the lsusb executable
	path, err := exec.LookPath(application)
	if err != nil {
		fmt.Printf("%s not found\n", application)
	}

	if err := checkExe(path); err != nil {
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
