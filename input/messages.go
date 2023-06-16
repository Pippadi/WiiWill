package input

import (
	"fyne.io/fyne/v2"
	actor "gitlab.com/prithvivishak/goactor"
)

func SendKeyEvent(dest actor.Inbox, key fyne.KeyName, pressed bool) {
	dest <- func(a actor.Actor) error {
		a.(*Inputter).keyEvent(key, pressed)
		return nil
	}
}

func SendMouseXSpeed(dest actor.Inbox, speed int32) {
	dest <- func(a actor.Actor) error {
		a.(*Inputter).mouseSpeedX(speed)
		return nil
	}
}

func SendMouseYSpeed(dest actor.Inbox, speed int32) {
	dest <- func(a actor.Actor) error {
		a.(*Inputter).mouseSpeedY(speed)
		return nil
	}
}
