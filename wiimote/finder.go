package wiimote

import (
	"strings"

	"github.com/Pippadi/loggo"
	"github.com/pilebones/go-udev/netlink"
	actor "gitlab.com/prithvivishak/goactor"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

const wiimoteName = "RVL-CNT-01"

type Finder struct {
	actor.Base
}

func NewFinder() *Finder {
	return new(Finder)
}

func (f *Finder) Initialize() (err error) {
	err = adapter.Enable()
	if err != nil {
		return
	}

	// Must be launched in goroutine. UI won't start otherwise. Have no clue why.
	go adapter.Scan(func(adapter *bluetooth.Adapter, dev bluetooth.ScanResult) {
		if strings.Contains(dev.LocalName(), wiimoteName) {
			sendCandidateWiimote(f.CreatorInbox(), dev.Address)
		}
	})

	return
}

func (f *Finder) Finalize() {
	adapter.StopScan()
}

func (f *Finder) connectToWiimote(btAddr bluetooth.Addresser) error {
	loggo.Infof("Connecting to %s", btAddr.String())

	dev, err := adapter.Connect(btAddr, bluetooth.ConnectionParams{})
	if err != nil {
		sendConnectError(f.CreatorInbox(), err)
	}

	path, err := getDevicePathFromUdev()
	if err != nil {
		sendConnectError(f.CreatorInbox(), err)
	}

	sendDevice(f.CreatorInbox(), dev, path)
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
