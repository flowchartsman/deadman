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
	devices := strings.Split(string(out), "\n\n")

	for _, d := range devices {
		if strings.HasPrefix(d, "ugen") {
			lines := strings.Split(d, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "  name:") {
					name := strings.TrimSpace(strings.TrimPrefix(line, "  name:"))
					deviceList = append(deviceList, device{name, ""})
				} else if strings.HasPrefix(line, "  devname:") {
					id := strings.TrimSpace(strings.TrimPrefix(line, "  devname:"))
					if len(deviceList) > 0 {
						deviceList[len(deviceList)-1].ID = id
					}
				}
			}
		}
	}
	if len(deviceList) == 0 {
		return nil, nil
	}
	return deviceList, nil
}
