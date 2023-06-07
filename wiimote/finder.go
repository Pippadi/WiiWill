package wiimote

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/Pippadi/loggo"
	"github.com/pilebones/go-udev/netlink"
	actor "gitlab.com/prithvivishak/goactor"
)

type Device string
type Action string

const (
	Wiimote  Device = `"Nintendo Wii Remote"`
	Nunchuk         = `"Nintendo Wii Remote Nunchuk"`
	NoDevice        = ""

	Add    Action = "add"
	Change        = "change"
	Remove        = "remove"
)

type Finder struct {
	actor.Base
	conn     *netlink.UEventConn
	eventQ   chan netlink.UEvent
	connQuit chan struct{}
}

func NewFinder() *Finder {
	return new(Finder)
}

func (f *Finder) Initialize() error {
	f.conn = new(netlink.UEventConn)
	if err := f.conn.Connect(netlink.UdevEvent); err != nil {
		return err
	}

	f.eventQ = make(chan netlink.UEvent)
	f.connQuit = f.conn.Monitor(f.eventQ, nil, nil)

	go func() {
		for {
			event := <-f.eventQ

			name, nameOk := event.Env["NAME"]
			if !nameOk {
				continue
			}

			var dev Device
			switch Device(name) {
			case Wiimote:
				dev = Wiimote
			case Nunchuk:
				dev = Nunchuk
			default:
				continue
			}

			switch Action(event.Action) {
			case Add:
				eventPath, err := eventPathFromSysfs(path.Join("/sys", event.KObj))
				if err != nil {
					loggo.Error(err)
					continue
				}
				addDevice(f.CreatorInbox(), dev, eventPath)
			default:
			}
		}
	}()
	return nil
}

func (f *Finder) Finalize() {
	close(f.connQuit)
	f.conn.Close()
}

func eventPathFromSysfs(sysfsPath string) (string, error) {
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
