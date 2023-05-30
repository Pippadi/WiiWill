package wiimote

import (
	"os"

	actor "gitlab.com/prithvivishak/goactor"
)

type Key byte
type KeyState byte

const (
	BtnA     Key = 0x30
	BtnB         = 0x31
	Btn1         = 0x01
	Btn2         = 0x02
	BtnUp        = 0x67
	BtnRight     = 0x6a
	BtnLeft      = 0x69
	BtnDown      = 0x6c
	BtnPlus      = 0x97
	BtnMinus     = 0x9c
	BtnHome      = 0x3c

	Pressed  KeyState = 0x01
	Released          = 0x00

	btnCodeOffset int = 18
	stateOffset       = 20
)

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
				return
			}
			if n > 0 {
				sendKeyEvent(
					e.CreatorInbox(),
					Key(buf[btnCodeOffset]),
					KeyState(buf[stateOffset]),
				)
			}
		}
	}()

	return nil
}
