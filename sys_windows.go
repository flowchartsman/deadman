package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"os/exec"
	"syscall"
)

//TODO: Call DLL directly
func enumerateDevices() ([]device, error) {
	if err := checkExe("powershell"); err != nil {
		return nil, err
	}
	cmd := exec.Command("powershell", "-command", `gwmi Win32_USBControllerDevice|%{[wmi]($_.Dependent)}|Select-Object Description,DeviceID|ConvertTo-CSV -notypeinformation|select -skip 1`)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(bytes.NewReader(out))
	deviceList := []device{}
	parsed, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, d := range parsed {
		if d[0] != "" && d[1] != "" {
			deviceList = append(deviceList, device{d[0], d[1]})
		} else {
			return nil, errors.New("incorrect device output")
		}
	}
	return deviceList, nil
}

func shutdownNow() error {
	if err := checkExe("shutdown"); err != nil {
		return err
	}
	err := exec.Command("shutdown", "/f", "/p").Run()

	//Not that this matters
	return err
}
