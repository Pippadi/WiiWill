package wiimote

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/pilebones/go-udev/netlink"
	actor "gitlab.com/prithvivishak/goactor"
)

type Device string

const (
	Wiimote Device = `"Nintendo Wii Remote"`
	Nunchuk        = `"Nintendo Wii Remote Nunchuk"`
)

type Finder struct {
	actor.Base
}

func NewFinder() *Finder {
	return new(Finder)
}

func (f *Finder) Initialize() error {
	go func() {
		for {
			eventPath, dev, err := getEventPathFromUdev()
			if err == nil {
				sendDevice(f.CreatorInbox(), dev, eventPath)
				actor.SendStopMsg(f.Inbox())
				return
			}
		}
	}()
	return nil
}

func getSysfsPathFromUdev() (string, Device, error) {
	conn := new(netlink.UEventConn)
	err := conn.Connect(netlink.UdevEvent)
	if err != nil {
		return "", "", err
	}

	eventQ := make(chan netlink.UEvent)
	quit := conn.Monitor(eventQ, nil, nil)
	defer func() {
		close(quit)
		conn.Close()
	}()

	for {
		select {
		case event := <-eventQ:
			name, nameOk := event.Env["NAME"]
			if event.Action == "add" && nameOk {
				switch Device(name) {
				case Wiimote:
					return path.Join("/sys", event.KObj), Wiimote, nil
				case Nunchuk:
					return path.Join("/sys", event.KObj), Nunchuk, nil
				default:
				}
			}
		}
	}
	return "", "", nil
}

func getEventPathFromUdev() (string, Device, error) {
	sysfsPath, dev, err := getSysfsPathFromUdev()
	if err != nil {
		return "", "", err
	}

	files, err := ioutil.ReadDir(sysfsPath)
	if err != nil {
		return "", "", err
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "event") {
			return path.Join("/dev/input", f.Name()), dev, nil
		}
	}

	return "", "", errors.New("Event file not found")
}
