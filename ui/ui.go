package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/WiiWill/wiimote"
	"github.com/Pippadi/loggo"
	actor "gitlab.com/prithvivishak/goactor"
	"tinygo.org/x/bluetooth"
)

type UI struct {
	actor.Base

	btAdapter   *bluetooth.Adapter
	finderInbox actor.Inbox

	wwApp      fyne.App
	mainWindow fyne.Window

	dev *bluetooth.Device
}

var _ wiimote.Manager = new(UI)

func New() *UI {
	return new(UI)
}

func (u *UI) Initialize() (err error) {
	u.wwApp = app.NewWithID("com.github.Pippadi.WiiWill")
	u.wwApp.Lifecycle().SetOnStopped(func() {
		if u.dev != nil {
			u.dev.Disconnect()
		}
		actor.SendStopMsg(u.Inbox())
	})
	u.mainWindow = u.wwApp.NewWindow("WiiWill")
	u.mainWindow.SetContent(widget.NewLabel("Wiimote"))
	u.mainWindow.Resize(fyne.NewSize(400, 400))
	u.mainWindow.SetMaster()

	u.CreatorInbox() <- func(a actor.Actor) error {
		u.mainWindow.ShowAndRun()
		return nil
	}

	u.finderInbox, err = u.SpawnNested(wiimote.NewFinder(), "Finder")

	return err
}

func (u *UI) AddCandidateWiimote(btAddr bluetooth.Addresser) {
	loggo.Info("Connecting to", btAddr.String())
	wiimote.SendConnect(u.finderInbox, btAddr)
}

func (u *UI) SetDevice(dev *bluetooth.Device, eventPath string) {
	loggo.Info("Wiimote button events at", eventPath)
	actor.SendStopMsg(u.finderInbox)
	u.SpawnNested(wiimote.NewEventReader(eventPath), "EventReader")
}

func (u *UI) HandleKeyEvent(key wiimote.Key, state wiimote.KeyState) {
	loggo.Infof("0x%02x %d", key, state)
}

func (u *UI) HandleConnectError(err error) {
	dialog.ShowError(err, u.mainWindow)
}
