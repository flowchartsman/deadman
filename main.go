package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	initialDevices, err := enumerateDevices()
	if err != nil {
		log.Fatal(err)
	}

	deviceMap := map[string]int64{}

	if len(initialDevices) > 0 {
		now := time.Now().UnixNano()
		for _, d := range initialDevices {
			deviceMap[d.ID] = now
		}
	} else {
		log.Fatal("No devices found")
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	check := time.NewTicker(1 * time.Second)

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
				if _, ok := deviceMap[d.ID]; !ok {
					log.Printf("New device: %s\n", d.Name)
					shutdownNow()
				}
				deviceMap[d.ID] = now
			}

			// Look to see if any devices have been removed
			for id, t := range deviceMap {
				if t != now {
					log.Printf("Device with id %s has been removed\n", id)
					shutdownNow()
				}
			}
		case <-sigint:
			log.Println("SIGINT received")
			os.Exit(0)
		}
	}
}
