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
	BtnA     Keycode = 0x0304
	BtnB             = 0x0305
	Btn1             = 0x0257
	Btn2             = 0x0258
	BtnUp            = 0x0103
	BtnRight         = 0x0106
	BtnLeft          = 0x0105
	BtnDown          = 0x0108
	BtnPlus          = 0x0407
	BtnMinus         = 0x0412
	BtnHome          = 0x0316

	BtnZ = 0x0135
	BtnC = 0x0132

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
				loggo.Info("Wiimote disconnected")
				return
			}
			if ev.Type == Key {
				loggo.Debug("%+v", ev)
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
