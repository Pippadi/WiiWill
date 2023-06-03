package wiimote

import (
	"strings"

	"github.com/pilebones/go-udev/netlink"
	actor "gitlab.com/prithvivishak/goactor"
)

const wiimoteName = "RVL-CNT-01"

type Finder struct {
	actor.Base
}

func NewFinder() *Finder {
	return new(Finder)
}

func (f *Finder) Initialize() error {
	go func() {
		for {
			devname, err := getDevicePathFromUdev()
			if err == nil {
				sendEventPath(f.CreatorInbox(), devname)
				return
			}
		}
	}()
	return nil
}

func getDevicePathFromUdev() (string, error) {
	conn := new(netlink.UEventConn)
	err := conn.Connect(netlink.UdevEvent)
	if err != nil {
		return "", err
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
			if event.Action == "add" &&
				strings.Contains(event.KObj, "bluetooth") {
				maj, majOk := event.Env["MAJOR"]
				min, minOk := event.Env["MINOR"]
				keyIn, keyInOk := event.Env["ID_INPUT_KEY"]
				dev, devOk := event.Env["DEVNAME"]
				// min="0" gives us /dev/input/jsX, which doesn't register D-pad events
				if majOk && minOk && keyInOk && devOk &&
					maj == "13" && min != "0" && keyIn == "1" {
					return dev, nil
				}
			}
		}
	}
	return "", nil
}
