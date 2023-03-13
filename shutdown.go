package main

import (
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type commandList struct {
	sync.Mutex
	commands []*exec.Cmd
}

func (cl *commandList) Add(cmd *exec.Cmd) {
	cl.Lock()
	cl.commands = append(cl.commands, cmd)
	cl.Unlock()
}

func (cl *commandList) KillAll() {
	for _, c := range cl.commands {
		if c.Process != nil {
			_ = c.Process.Kill()
		}
	}
	cl.commands = nil
}

func parallelExecute(cmd *exec.Cmd, wg *sync.WaitGroup) {
	err := cmd.Start()
	if err != nil {
		log.Printf("Error starting parallel command %s: %s", cmd.Path, err)
	}
	wg.Add(1)
	err = cmd.Wait()
	wg.Done()
	if err != nil {
		log.Printf("Error while running parallel command %s: %s", cmd.Path, err)
	}
}

func shutdownSequence(conf *config) {
	if conf.Commands == nil {
		shutdownNow()
		return
	}

	done := make(chan struct{})
	cl := new(commandList)
	wg := new(sync.WaitGroup)

	for _, command := range conf.Commands {
		cmdParts := strings.Split(command, " ")
		if string(cmdParts[0][0]) == "!" {
			osCommand := exec.Command(cmdParts[0][1:], cmdParts[1:]...)
			cl.Add(osCommand)
			go parallelExecute(osCommand, wg)
		} else {
			osCommand := exec.Command(cmdParts[0], cmdParts[1:]...)
			cl.Add(osCommand)
			err := osCommand.Run()
			if err != nil {
				log.Printf("Error while running command %s: %s", osCommand.Path, err)
			}
		}
	}

	go func() {
		wg.Wait()
		if conf.Shutdown {
			if err := shutdownNow(); err != nil {
				log.Fatal("error shutting down:", err)
			}
		}
		done <- struct{}{}
	}()

	select {
	case <-done:
		log.Println("commands have finished")
	case <-time.After(time.Duration(conf.ShutdownTimeout) * time.Millisecond):
		log.Println("timed out waiting for commands to finish")
		cl.KillAll()
	}
}
