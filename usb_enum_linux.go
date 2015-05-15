package main

import (
	"os"
	"strings"
	"syscall"
)

const NETLINK_KOBJECT_UEVENT = 15
const UEVENT_BUFFER_SIZE = 2048

func getEventChannel(*config) (chan *devEvent, error) {
	eventChan := make(chan *devEvent, 1)

	fd, err := syscall.Socket(
		syscall.AF_NETLINK, syscall.SOCK_RAW,
		NETLINK_KOBJECT_UEVENT,
	)
	if err != nil {
		return nil, err
	}

	nl := syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Pid:    uint32(os.Getpid()),
		Groups: 1,
	}
	err = syscall.Bind(fd, &nl)
	if err != nil {
		return nil, err
	}

	b := make([]byte, UEVENT_BUFFER_SIZE*2)
	go func() {
		for {
			_, err := syscall.Read(fd, b)
			if err != nil {
				eventChan <- ErrorEvent(err)
			}
			event := extractUSBEvent(b)
			if event != nil {
				eventChan <- event
			}
		}
	}()
	return eventChan, nil
}

func extractUSBEvent(msg []byte) *devEvent {
	j := 0
	var dev, act, subsys, product, busnum string
	for i := 0; i < len(msg)+1; i++ {
		if i == len(msg) || msg[i] == 0 {
			str := string(msg[j:i])
			a := strings.Split(str, "=")
			if len(a) == 2 {
				switch a[0] {
				case "DEVNAME":
					dev = a[1]
				case "ACTION":
					act = a[1]
				case "SUBSYSTEM":
					//TODO: return nil here if not USB?
					subsys = a[1]
				case "PRODUCT":
					product = a[1]
				case "BUSNUM":
					busnum = a[1]
				}
			}
			j = i + 1
		}
	}
	if subsys == "usb" && act != "" {
		var evType devEventType
		switch act {
		case "add":
			evType = evAdded
		case "remove":
			evType = evRemoved
		default:
			return nil
		}
		if dev == "" || product == "" || busnum == "" {
			return nil
			//return ErrorEvent(fmt.Errorf("Error parsing udev message: [%s]", string(msg)))
		}
		return &devEvent{
			EvType: evType,
			Error:  nil,
			Device: &device{
				Name: dev,
				ID:   product,
			},
		}
	}
	return nil
}
