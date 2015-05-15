package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

func main() {

	var configFile string
	var conf *config

	flag.StringVar(&configFile, "configfile", "", "config file to use")
	flag.Parse()

	if configFile != "" {
		var err error
		conf, err = loadConfig(configFile)
		if err != nil {
			log.Fatal("failed to parse config file: ", err)
		}
	} else {
		log.Println("defaulting to immediate shutdown and 1 second polling interval")
		conf = &config{Shutdown: true, PollInterval: 1000, ShutdownTimeout: 10000}
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	eventChan, err := getEventChannel(conf)
	if err != nil {
		log.Fatalln("error getting events", err)
	}

	log.Println("deadman started (press ctrl-c to exit)")

	for {
		select {
		case e := <-eventChan:
			switch e.EvType {
			case evRemoved:
				log.Printf("device removed:")
			case evAdded:
				log.Printf("device added:")
			case evError:
				log.Println("ERROR:", e.Error)
				continue
			}
			log.Println("Name:", e.Device.Name, "ID:", e.Device.ID)
		case <-sigint:
			log.Println("SIGINT received")
			os.Exit(0)
		}
	}
}
