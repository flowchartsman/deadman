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
	return fmt.Errorf("failed to initiate system shutdown")
}
