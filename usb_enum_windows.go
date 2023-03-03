// Note that this may require additional privileges to execute, depending on the system configuration.

package main

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	GUID_DEVINTERFACE_USB_DEVICE = "{A5DCBF10-6530-11D2-901F-00C04FB951ED}"
	MAX_DEVICE_ID_LEN            = 200
	MAX_CLASS_NAME_LEN           = 32
	MAX_REG_KEY_LEN              = 256
)

var (
	cfgMgr32           = syscall.MustLoadDLL("CfgMgr32.dll")
	cmEnumerateDevices = cfgMgr32.MustFindProc("CM_Enumerate_Devices")
)

type device struct {
	Name string
	ID   string
}

func enumerateDevices() ([]device, error) {
	buf := make([]byte, MAX_DEVICE_ID_LEN)
	var deviceList []device

	ret, _, _ := cmEnumerateDevices.Call(
		uintptr(unsafe.Pointer(&GUID_DEVINTERFACE_USB_DEVICE[0])),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(MAX_DEVICE_ID_LEN),
	)
	if ret != 0 {
		return nil, fmt.Errorf("failed to enumerate devices: %d", ret)
	}

	for ret == 0 {
		classBuf := make([]uint16, MAX_CLASS_NAME_LEN)
		regKeyBuf := make([]uint16, MAX_REG_KEY_LEN)

		ret, _, _ = cmEnumerateDevices.Call(
			uintptr(unsafe.Pointer(&GUID_DEVINTERFACE_USB_DEVICE[0])),
			ret,
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(MAX_DEVICE_ID_LEN),
		)

		if ret == 0 {
			cmGetClassName := cfgMgr32.MustFindProc("CM_Get_Class_NameW")
			_, _, _ = cmGetClassName.Call(
				uintptr(unsafe.Pointer(&classBuf[0])),
				uintptr(MAX_CLASS_NAME_LEN),
				uintptr(0),
			)

			cmOpenClassRegKeyEx := cfgMgr32.MustFindProc("CM_Open_Class_Key_ExW")
			ret, _, _ = cmOpenClassRegKeyEx.Call(
				uintptr(unsafe.Pointer(&GUID_DEVINTERFACE_USB_DEVICE[0])),
				uintptr(unsafe.Pointer(nil)),
				uintptr(0),
				uintptr(syscall.KEY_QUERY_VALUE|syscall.KEY_ENUMERATE_SUB_KEYS|syscall.KEY_READ),
				uintptr(unsafe.Pointer(&regKeyBuf[0])),
				uintptr(0),
			)
			if ret == 0 {
				deviceList = append(deviceList, device{syscall.UTF16ToString(classBuf), syscall.UTF16ToString(regKeyBuf)})
			}
		}
	}

	if len(deviceList) == 0 {
		return nil, nil
	}

	return deviceList, nil
}
