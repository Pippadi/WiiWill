package wiimote

import (
	actor "gitlab.com/prithvivishak/goactor"
	"tinygo.org/x/bluetooth"
)

func SendConnect(dest actor.Inbox, btAddr bluetooth.Addresser) {
	dest <- func(a actor.Actor) error {
		return a.(*Finder).connectToWiimote(btAddr)
	}
}

type Manager interface {
	AddCandidateWiimote(btAddr bluetooth.Addresser)
	SetDevice(dev *bluetooth.Device, eventPath string)
	HandleKeyEvent(Keycode, KeyState)

	HandleConnectError(err error)
}

func sendDevice(dest actor.Inbox, dev *bluetooth.Device, eventPath string) {
	dest <- func(a actor.Actor) error {
		a.(Manager).SetDevice(dev, eventPath)
		return nil
	}
}

func sendCandidateWiimote(dest actor.Inbox, btAddr bluetooth.Addresser) {
	dest <- func(a actor.Actor) error {
		a.(Manager).AddCandidateWiimote(btAddr)
		return nil
	}
}

func sendConnectError(dest actor.Inbox, err error) {
	dest <- func(a actor.Actor) error {
		a.(Manager).HandleConnectError(err)
		return nil
	}
}

func sendKeyEvent(dest actor.Inbox, btn Keycode, state KeyState) {
	dest <- func(a actor.Actor) error {
		a.(Manager).HandleKeyEvent(btn, state)
		return nil
	}
}
