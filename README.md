# deadman
A usb-based dead man's switch for your computer.

## Usage
```
sudo ./deadman
```
-or-
```
deadman.exe
```

Can be killed with SIGINT (ctrl-c)

## Purpose
Inspired by [usbkill](https://github.com/hephaest0s/usbkill), I thought I'd make a version in Go that has no dependencies and runs on Windows as well. It currently has feature parity with USBKill, meaning that the moment a new USB device is added or removed, your system will shut down.

## TODO
* Add support for FreeBSD
* Moar testing
* Moar docs
* Device whitelisting
* Better logging. You won't see much at the moment as it shuts down as soon as it can.
* Hook into system calls using shared libs as much as possible. For now, all systems parse ```lsusb``` or its equivalent every second
* Make check interval configurable
* Make commands configurable
