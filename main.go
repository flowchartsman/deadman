package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
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

	initialDevices, err := enumerateDevices()
	if err != nil {
		log.Fatal(err)
	}

	deviceMap := map[device]int64{}

	if len(initialDevices) > 0 {
		now := time.Now().UnixNano()
		for _, d := range initialDevices {
			deviceMap[d] = now
		}
	} else {
		log.Fatal("No devices found")
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	check := time.NewTicker(time.Duration(conf.PollInterval) * time.Millisecond)

	log.Println("deadman started (press ctrl-c to exit)")

	for {
		select {
		case <-check.C:
			now := time.Now().UnixNano()
			devices, err := enumerateDevices()
			if err != nil {
				log.Fatal(err)
			}
			// Look to see if there are any new devicies
			for _, d := range devices {
				if _, ok := deviceMap[d]; !ok {
					log.Printf("New device: %s [%s]\n", d.Name, d.ID)
					shutdownSequence(conf)
				}
				deviceMap[d] = now
			}

			// Look to see if any devices have been removed
			for d, t := range deviceMap {
				if t != now {
					log.Printf("Device %s [%s]  has been removed\n", d.Name, d.ID)
					shutdownSequence(conf)
				}
			}
		case <-sigint:
			log.Println("SIGINT received")
			os.Exit(0)
		}
	}
}
