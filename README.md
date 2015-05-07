![deadman](github/logo.png)

a usb-based dead man's switch for your computer.

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
deadman is intended as an anti-forensics/compulsion tool that will prevent or limit access to your machine if you are compelled away from it, or a device is attached to it without your consent. While it is running, if any USB device is attached or removed, it will execute a forced shut down. Possible use cases could involve a usb fob attached to the user with a lanyard as a kill switch, or as a countermeasure to devices like mouse jigglers or programmable HID devices.

## Inspiration
I recently came across [heaphaest0s](https://github.com/hephaest0s)' cool project, [usbkill](https://github.com/hephaest0s/usbkill), which is written in Python. I thought I might be able to improve it somewhat by making an alternative Go version that would have no external dependencies and also would work on Windows 7/8. It currently has feature parity with USBKill, though new features are being developed all the time. 

## TODO
* Moar testing
* Moar docs
* Add support for FreeBSD
* Device whitelisting
* Better logging. You won't see much at the moment as it shuts down as soon as it can.
* Hook into system calls using shared libs as much as possible. For now, all systems parse ```lsusb``` or its equivalent every second
* Make check interval configurable
* Make commands configurable
