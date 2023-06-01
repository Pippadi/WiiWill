package ui

import (
	"errors"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/WiiWill/ui/mapeditor"
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

	candidates        map[string]bluetooth.Addresser
	candidateSelector *widget.SelectEntry
	connectBtn        *widget.Button

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

	u.candidates = make(map[string]bluetooth.Addresser)
	u.candidateSelector = widget.NewSelectEntry([]string{})
	u.candidateSelector.PlaceHolder = "Wiimote Bluetooth Address"
	u.candidateSelector.Validator = fyne.StringValidator(validBtAddr)
	u.connectBtn = widget.NewButtonWithIcon(
		"Connect",
		theme.ContentAddIcon(),
		u.connectToSelected,
	)
	connectAcc := widget.NewAccordionItem("Connect", container.NewVBox(
		u.candidateSelector,
		container.NewHBox(layout.NewSpacer(), u.connectBtn, layout.NewSpacer()),
	))

	mapAcc := widget.NewAccordionItem("Map", mapeditor.New(u.mainWindow).UI())

	mainAcc := widget.NewAccordion(connectAcc, mapAcc)
	mainAcc.OpenAll()

	u.mainWindow.SetContent(mainAcc)
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
	u.candidates[btAddr.String()] = btAddr
	addrs := make([]string, 0)
	for a, _ := range u.candidates {
		addrs = append(addrs, a)
	}
	u.candidateSelector.SetOptions(addrs)
	u.candidateSelector.SetText(btAddr.String())
}

func (u *UI) connectToSelected() {
	if err := u.candidateSelector.Validate(); err != nil {
		dialog.ShowError(err, u.mainWindow)
		return
	}
	loggo.Info("Connecting to")
	wiimote.SendConnect(u.finderInbox, u.candidates[u.candidateSelector.Text])
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

func validBtAddr(addr string) (err error) {
	err = errors.New("Invalid bluetooth address")

	chunks := strings.Split(addr, ":")
	if len(chunks) != 6 {
		return
	}
	for _, chunk := range chunks {
		_, e := strconv.ParseUint(chunk, 16, 8)
		if e != nil {
			return
		}
	}

	err = nil
	return
}
