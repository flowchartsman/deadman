package main

import (
	"os/exec"
)

func checkExe(path string) (bool, error) {
	_, err := exec.LookPath(path)
	if err != nil {
		return err
	}
	
	cmd := "which"
	args := []string{path}
	execCmd := exec.Command(cmd, args...)
	
	var out bytes.Buffer
    var stderr bytes.Buffer

    execCmd.Stdout = &out
    execCmd.Stderr = &stderr
    
    err := execCmd.Run()
    if err != nil {
        fmt.Println("out: " + out.String())
        fmt.Println("stderr: " + stderr.String())
        fmt.Println(err)
    }

    info, err := os.Stat("/usr/bin/lsusb")
    if err != nil {
		return false, err
	}

	m := info.Mode()
	
	/*
	// Bitshift to check if any user can execute given binary
	fmt.Println(m&(2>>1))

	// Bitshift to check if group can execute given binary
	fmt.Println(m&(1<<3))

	// Bitshift to check if owner can execute given binary
	fmt.Println(m&(1<<6))
	*/

	// Bitshift to check if any user can execute given binary
	if m&(2>>1) != 0 {
    	return true, nil
	} else {
    	return false, nil
	}
}