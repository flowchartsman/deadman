/*
This function uses the C library from the IOKit framework on macOS to enumerate the
USB devices connected to the system. It initializes an empty device list and error,
and then opens the I/O Kit registry entry for USB devices. It creates a matching
dictionary for USB devices and an iterator to iterate over the devices. For each
device, it gets the device name and device ID by calling the IORegistryEntryGetName
and IORegistryEntryGetParentEntry functions, respectively, and converts the
CFTypeRefs to Go strings. It then appends the device to the device list. If there
are any errors in the process, it returns an error. Otherwise, it returns the
device list.
*/
package main

import (
	"C"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"unsafe"
)

// The C library and header are both part of the IOKit framework on macOS, and are
// included by specifying -framework IOKit as a linker flag.
// #cgo LDFLAGS: -framework IOKit
// #include <IOKit/IOKitLib.h>
import "C"

func enumerateDevices() ([]device, error) {
	// Initialize device list and error
	var deviceList []device
	var err error

	// Open I/O Kit registry entry for USB devices
	masterPort, result := IOMasterPort(0)
	if result != KERN_SUCCESS {
		return nil, fmt.Errorf("failed to get master port: %d", result)
	}

	// Create a matching dictionary for USB devices
	usbClass := IOServiceMatching(kIOUSBDeviceClassName)
	if usbClass == nil {
		return nil, errors.New("failed to create matching dictionary")
	}

	// Create an iterator to iterate over USB devices
	iter := IO_OBJECT_NULL
	defer func() {
		if iter != IO_OBJECT_NULL {
			IOObjectRelease(iter)
		}
	}()

	// Iterate over USB devices and add them to the device list
	for {
		service := IOIteratorNext(iter)
		if service == IO_OBJECT_NULL {
			break
		}
		defer IOObjectRelease(service)

		// Get device name
		var deviceNameCFTypeRef C.CFTypeRef
		result := IORegistryEntryGetName(service, &deviceNameCFTypeRef)
		if result != KERN_SUCCESS {
			err = fmt.Errorf("failed to get device name: %d", result)
			break
		}

		// Convert the device name CFTypeRef to a Go string
		deviceName := C.GoString((*C.char)(unsafe.Pointer(deviceNameCFTypeRef)))
		C.CFRelease(deviceNameCFTypeRef)

		// Get device ID
		var deviceIDCFTypeRef C.CFTypeRef
		result = IORegistryEntryGetParentEntry(service, kIOServicePlane, &deviceIDCFTypeRef)
		if result != KERN_SUCCESS {
			err = fmt.Errorf("failed to get device ID: %d", result)
			break
		}

		// Convert the device ID CFTypeRef to a Go string
		deviceID := C.GoString((*C.char)(unsafe.Pointer(deviceIDCFTypeRef)))
		C.CFRelease(deviceIDCFTypeRef)

		// Add device to list
		deviceList = append(deviceList, device{deviceName, deviceID})
	}

	// Check for errors and return the device list
	if err != nil {
		return nil, err
	}

	return deviceList, nil
}
