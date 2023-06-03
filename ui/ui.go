package ui

import (
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

	finderInbox   actor.Inbox
	listenerInbox actor.Inbox

	wwApp      fyne.App
	mainWindow fyne.Window

	statusLbl   *widget.Label
	activityBar *widget.ProgressBarInfinite

	mapEditor *mapeditor.MapEditor

	keyboard uinput.Keyboard
}

var _ wiimote.Manager = new(UI)

func New() *UI {
	return new(UI)
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
		u.activityBar,
		container.NewCenter(u.statusLbl),
		widget.NewSeparator(),
		u.mapEditor.UI(),
	)
	u.mainWindow.SetContent(box)
	u.mainWindow.Resize(fyne.NewSize(800, 600))
	u.mainWindow.SetMaster()

	u.CreatorInbox() <- func(a actor.Actor) error {
		u.mainWindow.ShowAndRun()
		return nil
	}

	u.startFinder()
	u.keyboard, err = uinput.CreateKeyboard("/dev/uinput", []byte("WiiWill"))

	return err
}

func (u *UI) SetEventPath(eventPath string) {
	loggo.Info("Wiimote button events at", eventPath)
	u.statusLbl.SetText("Listening for events on " + eventPath)
	u.activityBar.Stop()
	u.activityBar.Hide()
	actor.SendStopMsg(u.finderInbox)

	var err error
	u.listenerInbox, err = u.SpawnNested(wiimote.NewEventReader(eventPath), "EventReader")
	if err != nil {
		dialog.ShowError(err, u.mainWindow)
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

func (u *UI) HandleLastMsg(a actor.Actor, err error) error {
	if a.Inbox() == u.listenerInbox {
		u.startFinder()
	}
	return nil
}

func (u *UI) startFinder() {
	u.finderInbox, _ = u.SpawnNested(wiimote.NewFinder(), "Finder")
	u.statusLbl.SetText("Waiting for Wiimote connection")
	u.activityBar.Show()
	u.activityBar.Start()
}
