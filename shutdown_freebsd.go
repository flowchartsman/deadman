/*
This file contains a function called shutdownNow() that attempts to initiate
a forced shutdown of the FreeBSD OS using the syscall package.

The function first tries to call the Reboot system call with the
RB_POWERCYCLE command, which should initiate an immediate shutdown of the
system. If that fails, it tries to call the Reboot system call with the
RB_HALT command, which should initiate a system halt. If both of
these system calls fail, the function returns an error with the message
"Failed to initiate system shutdown."
*/
package main

import (
	"syscall"
)

func shutdownNow() error {
	// First, try to call the reboot system call with the appropriate arguments
	if err := syscall.Reboot(syscall.RB_POWERCYCLE); err == nil {
		return nil
	}

	// If that fails, try to call the halt system call with the appropriate arguments
	if err := syscall.Reboot(syscall.RB_HALT); err == nil {
		return nil
	}

	// If both system calls fail, return an error
	return fmt.Errorf("Failed to initiate system shutdown.")
}
