package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	configFile := ""
	flag.StringVar(&configFile, "configfile", "", "The config file to use.")
	flag.Parse()

	conf := getConfig(configFile)

	deviceMap, err := getInitialDeviceMap(conf)
	if err != nil {
		log.Fatal(err)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	check := time.NewTicker(time.Duration(conf.PollInterval) * time.Millisecond)
	log.Println("Deadman running (Press Ctrl-C to exit.)  ...")

	for {
		select {
		case <-check.C:
			now := time.Now().UnixNano()
			checkForNewDevices(deviceMap, now, conf)
			checkForRemovedDevices(deviceMap, now)
		case <-sigint:
			log.Println("SIGINT was received, exiting.")
			os.Exit(0)
		}
	}
}

func getConfig(configFile string) *config {
	if configFile != "" {
		conf, err := loadConfig(configFile)
		if err != nil {
			log.Fatal("Failed to parse config file: ", err)
		}
		return conf
	} else {
		log.Println("Defaulting to immediate shutdown and 1 second polling interval.")
		return &config{Shutdown: true, PollInterval: 1000, ShutdownTimeout: 10000}
	}
}

func getInitialDeviceMap(conf *config) (map[device]int64, error) {
	initialDevices, err := enumerateDevices()
	if err != nil {
		return nil, err
	}

	if len(initialDevices) == 0 {
		log.Fatal("No devices were found. Is the USB subsystem running?")
	}

	deviceMap := make(map[device]int64)
	now := time.Now().UnixNano()
	for _, d := range initialDevices {
		if isExemptDevice(d, conf) {
			deviceMap[d] = now
		}
	}
	return deviceMap, nil
}

func isExemptDevice(d device, conf *config) bool {
	for _, wd := range conf.ExemptDevices {
		if wd.Name == d.Name && wd.ID == d.ID {
			return true
		}
	}
	return false
}

func checkForNewDevices(deviceMap map[device]int64, now int64, conf *config) {
	devices, err := enumerateDevices()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range devices {
		if _, ok := deviceMap[d]; !ok && isExemptDevice(d, conf) {
			log.Printf("New device: %s [%s]\n", d.Name, d.ID)
			shutdownNow()
			deviceMap[d] = now
		} else {
			deviceMap[d] = now
		}
	}
}

func checkForRemovedDevices(deviceMap map[device]int64, now int64) {
	for d, t := range deviceMap {
		if t != now {
			log.Printf("Device %s [%s] has been removed\n", d.Name, d.ID)
			shutdownNow()
		}
	}
}
