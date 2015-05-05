package main

import (
	"log"
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
		log.Fatal(err)
	}

	var deviceList []device
	devices := parseRE.FindAllStringSubmatch(string(out), -1)

	for _, d := range devices {
		deviceList = append(deviceList, device{d[2], d[1]})
	}
	return deviceList, nil
}
