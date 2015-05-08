package main

import (
	"errors"
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

	err := toml.UnmarshalFile(filename, &conf)
	if err != nil {
		return nil, err
	}

	if !conf.Shutdown && conf.Commands == nil {
		return nil, errors.New(`"Shutdown" is false or unconfigured and no commands have been specified. Add "Shutdown = <shutdownmode>" to the config file`)
	}
	return &conf, nil
}
