package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Pippadi/loggo"
	"github.com/pilebones/go-udev/netlink"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

type key byte
type state byte

const (
	btn_A     key = 0x30
	btn_B         = 0x31
	btn_1         = 0x01
	btn_2         = 0x02
	btn_UP        = 0x67
	btn_RIGHT     = 0x6a
	btn_LEFT      = 0x69
	btn_DOWN      = 0x6c
	btn_PLUS      = 0x97
	btn_MINUS     = 0x9c
	btn_HOME      = 0x3c

	pressed  state = 0x01
	released state = 0x00

	btnCodeOffset int = 18
	stateOffset       = 20
)

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
		if strings.Contains(dev.LocalName(), "RVL-CNT-01") {
			loggo.Info("Found Wiimote")
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

	eventPath, err := getDevicePathFromUdev()
	if err != nil {
		loggo.Error(err)
		return
	}

	loggo.Info("Wiimote button events at", eventPath)

	file, err := os.Open(eventPath)
	if err != nil {
		loggo.Error(err)
		return
	}
	defer file.Close()

	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	buf := make([]byte, 64)
	go func() {
		for {
			n, err := file.Read(buf)
			if err != nil {
				return
			}
			if n > 0 {
				loggo.Infof("0x%02x %d", buf[btnCodeOffset], buf[stateOffset])
			}
		}
	}()

	<-osSignal
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
