package input

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/Pippadi/WiiWill/ui/mapeditor"
	"github.com/bendahl/uinput"
	actor "gitlab.com/prithvivishak/goactor"
)

const mouseMoveInterval = 10 * time.Millisecond

type Inputter struct {
	actor.Base

	keyboard uinput.Keyboard
	mouse    uinput.Mouse

	// In pixels per mouseMoveInterval
	mouseXSpeed int32
	mouseYSpeed int32
}

func New() *Inputter {
	return &Inputter{mouseXSpeed: 0, mouseYSpeed: 0}
}

func (i *Inputter) Initialize() error {
	var err error

	i.keyboard, err = uinput.CreateKeyboard("/dev/uinput", []byte("WiiWill"))
	if err != nil {
		return err
	}

	i.mouse, err = uinput.CreateMouse("/dev/uinput", []byte("WiiWill"))
	if err != nil {
		return err
	}

	go i.moveMouse()

	return nil
}

func (i *Inputter) keyEvent(key fyne.KeyName, pressed bool) {
	switch key {
	case desktop.KeyNone:
		return
	case mapeditor.MouseLeft:
		if pressed {
			i.mouse.LeftPress()
		} else {
			i.mouse.LeftRelease()
		}
	case mapeditor.MouseMiddle:
		if pressed {
			i.mouse.MiddlePress()
		} else {
			i.mouse.MiddleRelease()
		}
	case mapeditor.MouseRight:
		if pressed {
			i.mouse.RightPress()
		} else {
			i.mouse.RightRelease()
		}
	default:
		uiKey, ok := mapeditor.FyneToUinputKey[key]
		if !ok {
			return
		}

		if pressed {
			i.keyboard.KeyDown(uiKey)
		} else {
			i.keyboard.KeyUp(uiKey)
		}
	}
}

func (i *Inputter) moveMouse() {
	for {
		i.mouse.Move(i.mouseXSpeed, i.mouseYSpeed)
		time.Sleep(mouseMoveInterval)
	}
}

func (i *Inputter) mouseSpeedX(speed int32) {
	i.mouseXSpeed = speed
}
func (i *Inputter) mouseSpeedY(speed int32) {
	i.mouseYSpeed = speed
}
