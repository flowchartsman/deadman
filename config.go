package main

import (
	"errors"
	"os"

	"github.com/naoina/toml"
)

type config struct {
	Shutdown        bool
	PollInterval    int
	ShutdownTimeout int
	Commands        []string
	//EmergencyFallback bool
	//EmergencyInterval int
}

type command struct {
	cmd  string
	args []string
}

func loadConfig(filename string) (*config, error) {
	conf := config{}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := toml.NewDecoder(f)
	err = dec.Decode(&conf)
	if err != nil {
		return nil, err
	}

	if !conf.Shutdown && conf.Commands == nil {
		return nil, errors.New(`"Shutdown" is false or unconfigured and no commands have been specified. Add "Shutdown = <shutdownmode>" to the config file`)
	}
	return &conf, nil
}
