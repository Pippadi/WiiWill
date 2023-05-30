package main

import (
	"strings"

	"github.com/Pippadi/loggo"
	"github.com/pilebones/go-udev/netlink"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func main() {
	loggo.SetLevel(loggo.DebugLevel)
	err := adapter.Enable()
	if err != nil {
		loggo.Error(err)
		return
	}

	var addr bluetooth.Addresser
	foundChan := make(chan struct{}, 1)

	err = adapter.Scan(func(adapter *bluetooth.Adapter, dev bluetooth.ScanResult) {
		if dev.LocalName() != "" {
			loggo.Debug(dev.LocalName())
		}
		if strings.Contains(dev.LocalName(), "RVL-CNT-01") {
			loggo.Debug("Found Wiimote")
			addr = dev.Address
			foundChan <- struct{}{}
			adapter.StopScan()
		}
	})
	if err != nil {
		loggo.Error(err)
		return
	}

	<-foundChan
	close(foundChan)

	loggo.Infof("Connecting to %s", addr.String())
	_, err = adapter.Connect(addr, bluetooth.ConnectionParams{})
	if err != nil {
		loggo.Error(err)
		return
	}

	loggo.Info(getDevicePathFromUdev())
}

func getDevicePathFromUdev() (string, error) {
	conn := new(netlink.UEventConn)
	err := conn.Connect(netlink.UdevEvent)
	if err != nil {
		return "", err
	}

	eventQ := make(chan netlink.UEvent)
	quit := conn.Monitor(eventQ, nil, nil)
	defer func() { close(quit) }()

	for {
		select {
		case event := <-eventQ:
			if event.Action == "add" &&
				strings.Contains(event.KObj, "bluetooth") {
				maj, majOk := event.Env["MAJOR"]
				min, minOk := event.Env["MINOR"]
				dev, devOk := event.Env["DEVNAME"]
				if majOk && minOk && devOk && maj == "13" && min == "79" {
					return dev, nil
				}
			}
		}
	}
	return "", nil
}
