package wiimote

import (
	"os"

	"github.com/Pippadi/loggo"
	actor "gitlab.com/prithvivishak/goactor"
)

type Keycode byte
type KeyState byte
type Key struct {
	Code       Keycode
	PrettyName string
}

const (
	// Wiimote held vertically
	BtnA     Keycode = 0x30
	BtnB             = 0x31
	Btn1             = 0x01
	Btn2             = 0x02
	BtnUp            = 0x67
	BtnRight         = 0x6a
	BtnLeft          = 0x69
	BtnDown          = 0x6c
	BtnPlus          = 0x97
	BtnMinus         = 0x9c
	BtnHome          = 0x3c

	Pressed  KeyState = 0x01
	Released          = 0x00

	btnCodeOffset int = 18
	stateOffset       = 20
)

var KeyMap = map[Keycode]Key{
	BtnA:     Key{BtnA, "A"},
	BtnB:     Key{BtnB, "B"},
	Btn1:     Key{Btn1, "1"},
	Btn2:     Key{Btn2, "2"},
	BtnUp:    Key{BtnUp, "D-pad Up"},
	BtnDown:  Key{BtnDown, "D-pad Down"},
	BtnLeft:  Key{BtnLeft, "D-pad Left"},
	BtnRight: Key{BtnRight, "D-pad Right"},
	BtnPlus:  Key{BtnPlus, "+"},
	BtnMinus: Key{BtnMinus, "-"},
	BtnHome:  Key{BtnHome, "Home"},
}

type EventReader struct {
	actor.Base
	path string
}

func NewEventReader(path string) *EventReader {
	return &EventReader{path: path}
}

func (e *EventReader) Initialize() error {
	file, err := os.Open(e.path)
	if err != nil {
		return err
	}

	go func() {
		defer file.Close()
		buf := make([]byte, 64)
		for {
			n, err := file.Read(buf)
			if err != nil {
				loggo.Error(err)
				loggo.Info("Wiimote disconnected")
				actor.SendStopMsg(e.Inbox())
				return
			}
			if n > 0 {
				sendKeyEvent(
					e.CreatorInbox(),
					Keycode(buf[btnCodeOffset]),
					KeyState(buf[stateOffset]),
				)
			}
		}
	}()

	return nil
}
