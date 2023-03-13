package main

/**

This file contains a function called shutdownNow() which initiates a forced
shutdown of the Windows operating system.
It does this by:
	Loading the kernel32.dll library
	Finding the ExitWindowsEx function using MustFindProc
	Calling it with the appropriate arguments

The function takes two parameters, both of which are uintptr values:
	The first parameter specifies the shutdown operation (EWX_SHUTDOWN and EWX_FORCE)
	The second parameter specifies the shutdown reason (a combination of SHTDN_REASON_MAJOR,
SHTDN_REASON_MINOR, and SHTDN_REASON_FLAG_PLANNED).
If the function call fails, it returns an error.
*/
import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32Dll   = syscall.MustLoadDLL("kernel32.dll")
	exitWindowsEx = kernel32Dll.MustFindProc("ExitWindowsEx")
)

func shutdownNow() error {
	// Constant for the shutdown operation
	const (
		EWX_SHUTDOWN              = 0x00000001
		EWX_FORCE                 = 0x00000004
		SHTDN_REASON_MAJOR        = 0x00020000
		SHTDN_REASON_MINOR        = 0x00000000
		SHTDN_REASON_FLAG_PLANNED = 0x80000000
	)

	// Call the ExitWindowsEx function to initiate the shutdown
	ret, _, err := exitWindowsEx.Call(uintptr(EWX_SHUTDOWN|EWX_FORCE), uintptr(SHTDN_REASON_MAJOR|SHTDN_REASON_MINOR|SHTDN_REASON_FLAG_PLANNED))
	if ret == 0 {
		return fmt.Errorf("failed to initiate shutdown: %v", err)
	}
	return nil
}
