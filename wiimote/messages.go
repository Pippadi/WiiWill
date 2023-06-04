package wiimote

import actor "gitlab.com/prithvivishak/goactor"

type Manager interface {
	SetEventPath(eventPath string)
	HandleKeyEvent(Keycode, KeyState)
}

func sendEventPath(dest actor.Inbox, eventPath string) {
	dest <- func(a actor.Actor) error {
		a.(Manager).SetEventPath(eventPath)
		return nil
	}
}

func sendKeyEvent(dest actor.Inbox, btn Keycode, state KeyState) {
	dest <- func(a actor.Actor) error {
		a.(Manager).HandleKeyEvent(btn, state)
		return nil
	}
}
