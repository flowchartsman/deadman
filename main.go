package main

import (
	"log"
)

func main() {
	devices, err := enumerateDevices()
	if err != nil {
		log.Fatal(err)
	}
	if len(devices) > 0 {
		for _, d := range devices {
			log.Printf("Device: %s, ID: %s", d.Name, d.ID)
		}
	} else {
		log.Fatal("No devices found")
	}
}
