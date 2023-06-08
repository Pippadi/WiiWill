package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/WiiWill/ui/mapeditor"
	"github.com/Pippadi/WiiWill/wiimote"
	"github.com/Pippadi/loggo"
	"github.com/bendahl/uinput"
	actor "gitlab.com/prithvivishak/goactor"
)

type UI struct {
	actor.Base

	finderIbx    actor.Inbox
	moteEventIbx actor.Inbox
	extEventIbx  actor.Inbox

	wwApp      fyne.App
	mainWindow fyne.Window

	statusLbl   *widget.Label
	activityBar *widget.ProgressBarInfinite

	mapEditor *mapeditor.MapEditor

	keyboard         uinput.Keyboard
	mouse            uinput.Mouse
	wiimoteConnected bool
	extension        wiimote.Device
}

var _ wiimote.Manager = new(UI)

func New() *UI {
	return &UI{wiimoteConnected: false, extension: wiimote.NoDevice}
}

func (u *UI) Initialize() (err error) {
	u.wwApp = app.NewWithID("com.github.Pippadi.WiiWill")
	u.wwApp.Lifecycle().SetOnStopped(func() {
		actor.SendStopMsg(u.Inbox())
	})
	u.mainWindow = u.wwApp.NewWindow("WiiWill")

	u.activityBar = widget.NewProgressBarInfinite()
	u.statusLbl = widget.NewLabel("")
	u.mapEditor = mapeditor.New(u.wwApp, u.mainWindow)

	box := container.NewVBox(
		u.mapEditor.UI(),
		widget.NewSeparator(),
		container.NewCenter(u.statusLbl),
		u.activityBar,
	)
	u.mainWindow.SetContent(box)
	u.mainWindow.Resize(fyne.NewSize(550, 500))
	u.mainWindow.SetMaster()

	u.CreatorInbox() <- func(a actor.Actor) error {
		u.mainWindow.ShowAndRun()
		return nil
	}

	u.setStatusLbl()
	u.finderIbx, _ = u.SpawnNested(wiimote.NewFinder(), "Finder")
	u.keyboard, err = uinput.CreateKeyboard("/dev/uinput", []byte("WiiWill"))
	if err != nil {
		return err
	}
	u.mouse, err = uinput.CreateMouse("/dev/uinput", []byte("WiiWill"))

	return err
}

func (u *UI) AddDevice(dev wiimote.Device, eventPath string) {
	loggo.Infof("%s events at %s", dev, eventPath)
	defer u.setStatusLbl()

	// Wait for permissions to be applied on /dev/eventX
	time.Sleep(250 * time.Millisecond)
	var err error
	if dev == wiimote.Wiimote && !u.wiimoteConnected {
		u.moteEventIbx, err = u.SpawnNested(wiimote.NewEventReader(eventPath), "WiimoteEventReader")
		if err != nil {
			dialog.ShowError(err, u.mainWindow)
			return
		}
		u.setWiimoteConnected()
	} else if dev == wiimote.Nunchuk && u.extension == wiimote.NoDevice {
		u.extEventIbx, err = u.SpawnNested(wiimote.NewEventReader(eventPath), "NunchukEventReader")
		if err != nil {
			dialog.ShowError(err, u.mainWindow)
			return
		}
		u.setExtConnected(dev)
	}
}

func (u *UI) HandleKeyEvent(key wiimote.Keycode, state wiimote.KeyState) {
	loggo.Infof("0x%03x %d", key, state)
	name := u.mapEditor.KeyFor(key)

	switch name {
	case desktop.KeyNone:
		return
	case mapeditor.MouseLeft:
		if state == wiimote.Pressed {
			u.mouse.LeftPress()
		} else {
			u.mouse.LeftRelease()
		}
	case mapeditor.MouseMiddle:
		if state == wiimote.Pressed {
			u.mouse.MiddlePress()
		} else {
			u.mouse.MiddleRelease()
		}
	case mapeditor.MouseRight:
		if state == wiimote.Pressed {
			u.mouse.RightPress()
		} else {
			u.mouse.RightRelease()
		}
	default:
		uiKey, ok := mapeditor.FyneToUinputKey[name]
		if !ok {
			return
		}
		loggo.Debug(name, uiKey)

		if state == wiimote.Pressed {
			u.keyboard.KeyDown(uiKey)
		} else {
			u.keyboard.KeyUp(uiKey)
		}
	}
}

func (u *UI) HandleLastMsg(a actor.Actor, reason error) error {
	if !u.IsStopping() {
		defer u.setStatusLbl()
		if a == nil {
			return nil
		}

		loggo.Debug("%+v", a)
		switch a.ID() {
		case "WiimoteEventReader":
			u.setWiimoteDisconnected()
		case "NunchukEventReader":
			u.setExtDisconnected()
		default:
		}
	}
	return nil
}

func (u *UI) setWiimoteConnected() {
	loggo.Info("Wiimote connected")
	u.wiimoteConnected = true
	u.activityBar.Hide()
	u.activityBar.Stop()
}

func (u *UI) setWiimoteDisconnected() {
	loggo.Info("Wiimote disconnected")
	u.wiimoteConnected = false
	u.activityBar.Show()
	u.activityBar.Start()
}

func (u *UI) setExtConnected(ext wiimote.Device) {
	loggo.Info(ext, "connected")
	u.extension = ext
}

func (u *UI) setExtDisconnected() {
	loggo.Info(u.extension, "disconnected")
	u.extension = wiimote.NoDevice
}

func (u *UI) setStatusLbl() {
	if !u.wiimoteConnected {
		u.statusLbl.SetText("Waiting for Wiimote connection")
	} else {
		if u.extension == wiimote.NoDevice {
			u.statusLbl.SetText("Wiimote connected")
		} else {
			u.statusLbl.SetText("Wiimote, " + string(u.extension) + " connected")
		}
	}
}
