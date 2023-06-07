package wiimote

import actor "gitlab.com/prithvivishak/goactor"

type Manager interface {
	AddDevice(dev Device, eventPath string)
	HandleKeyEvent(Keycode, KeyState)
}

func addDevice(dest actor.Inbox, device Device, eventPath string) {
	dest <- func(a actor.Actor) error {
		a.(Manager).AddDevice(device, eventPath)
		return nil
	}
}

func sendKeyEvent(dest actor.Inbox, btn Keycode, state KeyState) {
	dest <- func(a actor.Actor) error {
		a.(Manager).HandleKeyEvent(btn, state)
		return nil
	}
}
