package main

import (
	"os/exec"
	"strings"
)

func enumerateDevices() ([]device, error) {
	if err := checkExe("usbconfig"); err != nil {
		return nil, err
	}

	out, err := exec.Command("usbconfig", "-d", "ugen*").Output()
	if err != nil {
		return nil, err
	}

	var deviceList []device
	var currentDevice device

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			if currentDevice.Name != "" {
				deviceList = append(deviceList, currentDevice)
			}
			currentDevice = device{}
			continue
		}
		if strings.HasPrefix(line, "name:") {
			currentDevice.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		} else if strings.HasPrefix(line, "devname:") {
			currentDevice.ID = strings.TrimSpace(strings.TrimPrefix(line, "devname:"))
		}
	}

	if len(deviceList) == 0 {
		return nil, nil
	}
	return deviceList, nil
}
