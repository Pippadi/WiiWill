package wiimote

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/pilebones/go-udev/netlink"
	actor "gitlab.com/prithvivishak/goactor"
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
			eventPath, err := getEventPathFromUdev()
			if err == nil {
				sendEventPath(f.CreatorInbox(), eventPath)
				actor.SendStopMsg(f.Inbox())
				return
			}
		}
	}()
	return nil
}

func getSysfsPathFromUdev() (string, error) {
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
				name, nameOk := event.Env["NAME"]
				devpath, devpathOk := event.Env["DEVPATH"]
				if nameOk && devpathOk && name == `"Nintendo Wii Remote"` {
					return path.Join("/sys", devpath), nil
				}
			}
		}
	}
	return "", nil
}

func getEventPathFromUdev() (string, error) {
	sysfsPath, err := getSysfsPathFromUdev()
	if err != nil {
		return "", err
	}

	files, err := ioutil.ReadDir(sysfsPath)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "event") {
			return path.Join("/dev/input", f.Name()), nil
		}
	}

	return "", errors.New("Event file not found")
}
