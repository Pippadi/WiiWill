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
	SetEventPath(path string)
	HandleKeyEvent(Key, KeyState)

	HandleConnectError(err error)
}

func sendEventPath(dest actor.Inbox, path string) {
	dest <- func(a actor.Actor) error {
		a.(Manager).SetEventPath(path)
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

func sendKeyEvent(dest actor.Inbox, btn Key, state KeyState) {
	dest <- func(a actor.Actor) error {
		a.(Manager).HandleKeyEvent(btn, state)
		return nil
	}
}
