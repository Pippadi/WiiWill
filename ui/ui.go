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
	wiimoteConnected bool
}

var _ wiimote.Manager = new(UI)

func New() *UI {
	return &UI{wiimoteConnected: false}
}

func (u *UI) Initialize() (err error) {
	u.wwApp = app.NewWithID("com.github.Pippadi.WiiWill")
	u.wwApp.Lifecycle().SetOnStopped(func() {
		actor.SendStopMsg(u.Inbox())
	})
	u.mainWindow = u.wwApp.NewWindow("WiiWill")

	u.activityBar = widget.NewProgressBarInfinite()
	u.statusLbl = widget.NewLabel("")
	u.mapEditor = mapeditor.New(u.mainWindow)

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

	u.setWiimoteDisconnected()
	u.finderIbx, _ = u.SpawnNested(wiimote.NewFinder(), "Finder")
	u.keyboard, err = uinput.CreateKeyboard("/dev/uinput", []byte("WiiWill"))

	return err
}

func (u *UI) AddDevice(dev wiimote.Device, eventPath string) {
	loggo.Infof("%s events at %s", dev, eventPath)

	// Wait for permissions to be applied on /dev/eventX
	time.Sleep(250 * time.Millisecond)
	var err error
	if dev == wiimote.Wiimote && !u.wiimoteConnected {
		u.setWiimoteConnected()

		u.moteEventIbx, err = u.SpawnNested(wiimote.NewEventReader(eventPath), "WiimoteEventReader")
		if err != nil {
			dialog.ShowError(err, u.mainWindow)
		}
	} else if dev == wiimote.Nunchuk {
		u.extEventIbx, err = u.SpawnNested(wiimote.NewEventReader(eventPath), "NunchukEventReader")
		if err != nil {
			dialog.ShowError(err, u.mainWindow)
		}
	}
}

func (u *UI) RemoveDevice(dev wiimote.Device) {
	switch dev {
	case wiimote.Wiimote:
		actor.SendStopMsg(u.moteEventIbx)
		u.setWiimoteDisconnected()
		fallthrough
	case wiimote.Nunchuk:
		actor.SendStopMsg(u.extEventIbx)
	default:
	}
}

func (u *UI) HandleKeyEvent(key wiimote.Keycode, state wiimote.KeyState) {
	loggo.Infof("0x%02x %d", key, state)
	name := u.mapEditor.KeyFor(key)
	uiKey, ok := fyneToUinputKey[name]

	if name == desktop.KeyNone || !ok {
		return
	}

	if state == wiimote.Pressed {
		u.keyboard.KeyDown(uiKey)
	} else {
		u.keyboard.KeyUp(uiKey)
	}
}

func (u *UI) setWiimoteConnected() {
	u.wiimoteConnected = true
	u.statusLbl.SetText("Wiimote connected")
	u.activityBar.Hide()
	u.activityBar.Stop()
}

func (u *UI) setWiimoteDisconnected() {
	u.wiimoteConnected = false
	u.statusLbl.SetText("Waiting for Wiimote connection")
	u.activityBar.Show()
	u.activityBar.Start()
}
