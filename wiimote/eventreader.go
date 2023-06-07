package wiimote

import (
	"encoding/binary"
	"os"
	"syscall"

	"github.com/Pippadi/loggo"
	actor "gitlab.com/prithvivishak/goactor"
)

type Keycode uint16
type KeyState uint32
type EventType uint16

type KeyInfo struct {
	Code       Keycode
	PrettyName string
}

const (
	// Wiimote held vertically
	BtnA     Keycode = 0x130
	BtnB             = 0x131
	Btn1             = 0x101
	Btn2             = 0x102
	BtnUp            = 0x67
	BtnRight         = 0x6A
	BtnLeft          = 0x69
	BtnDown          = 0x6C
	BtnPlus          = 0x197
	BtnMinus         = 0x19C
	BtnHome          = 0x13C

	BtnZ = 0x135
	BtnC = 0x132

	Pressed  KeyState = 0x01
	Released          = 0x00

	Sync     EventType = 0x00
	Key      EventType = 0x01
	Relative EventType = 0x02
	Absolute EventType = 0x03
)

var KeyMap = map[Keycode]KeyInfo{
	BtnA:     KeyInfo{BtnA, "A"},
	BtnB:     KeyInfo{BtnB, "B"},
	Btn1:     KeyInfo{Btn1, "1"},
	Btn2:     KeyInfo{Btn2, "2"},
	BtnUp:    KeyInfo{BtnUp, "D-pad Up"},
	BtnDown:  KeyInfo{BtnDown, "D-pad Down"},
	BtnLeft:  KeyInfo{BtnLeft, "D-pad Left"},
	BtnRight: KeyInfo{BtnRight, "D-pad Right"},
	BtnPlus:  KeyInfo{BtnPlus, "+"},
	BtnMinus: KeyInfo{BtnMinus, "-"},
	BtnHome:  KeyInfo{BtnHome, "Home"},

	BtnZ: KeyInfo{BtnZ, "Z"},
	BtnC: KeyInfo{BtnC, "C"},
}

// See https://www.kernel.org/doc/Documentation/input/input.txt
type InputEvent struct {
	Timestamp syscall.Timeval
	Type      EventType
	Code      uint16
	Value     uint32
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
		for {
			ev := new(InputEvent)
			err = binary.Read(file, binary.LittleEndian, ev)
			if err != nil {
				loggo.Error(err)
				actor.SendStopMsg(e.Inbox())
				return
			}
			if ev.Type == Key {
				loggo.Debugf("%x", ev.Code)
				sendKeyEvent(
					e.CreatorInbox(),
					Keycode(ev.Code),
					KeyState(ev.Value),
				)
			}
		}
	}()

	return nil
}
