package mapeditor

import (
	"encoding/json"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/WiiWill/wiimote"
	"github.com/Pippadi/loggo"
)

type MapEditor struct {
	parentApp     fyne.App
	parentWindow  fyne.Window
	mainContainer fyne.CanvasObject

	loadBtn *widget.Button
	saveBtn *widget.Button

	buttons   map[wiimote.Keycode]*widget.Button
	selectors map[wiimote.Keycode]*KeyChooser
	mapping   map[wiimote.Keycode]fyne.KeyName
}

func New(a fyne.App, w fyne.Window) *MapEditor {
	m := new(MapEditor)
	m.parentApp = a
	m.parentWindow = w
	m.buttons = make(map[wiimote.Keycode]*widget.Button)
	m.selectors = make(map[wiimote.Keycode]*KeyChooser)
	m.mapping = make(map[wiimote.Keycode]fyne.KeyName)

	m.loadBtn = widget.NewButtonWithIcon("Load", theme.UploadIcon(), m.loadMap)
	m.saveBtn = widget.NewButtonWithIcon("Save", theme.FileIcon(), m.saveMap)

	for c, _ := range wiimote.KeyMap {
		m.buttons[c] = m.remapBtnFor(c)
		m.selectors[c] = NewKeyChooser(m.parentWindow)
		m.selectors[c].OnChanged = m.chooserChangeHandler(c)
		m.mapping[c] = desktop.KeyNone
	}

	// Map stores elements in arbitrary order, so order buttons manually
	f := widget.NewForm(
		widget.NewFormItem(wiimote.KeyMap[wiimote.BtnUp].PrettyName, m.buttons[wiimote.BtnUp]),
		widget.NewFormItem(wiimote.KeyMap[wiimote.BtnDown].PrettyName, m.buttons[wiimote.BtnDown]),
		widget.NewFormItem(wiimote.KeyMap[wiimote.BtnLeft].PrettyName, m.buttons[wiimote.BtnLeft]),
		widget.NewFormItem(wiimote.KeyMap[wiimote.BtnRight].PrettyName, m.buttons[wiimote.BtnRight]),
		widget.NewFormItem(wiimote.KeyMap[wiimote.BtnA].PrettyName, m.buttons[wiimote.BtnA]),
		widget.NewFormItem(wiimote.KeyMap[wiimote.BtnB].PrettyName, m.buttons[wiimote.BtnB]),
		widget.NewFormItem(wiimote.KeyMap[wiimote.Btn1].PrettyName, m.buttons[wiimote.Btn1]),
		widget.NewFormItem(wiimote.KeyMap[wiimote.Btn2].PrettyName, m.buttons[wiimote.Btn2]),
		widget.NewFormItem(wiimote.KeyMap[wiimote.BtnPlus].PrettyName, m.buttons[wiimote.BtnPlus]),
		widget.NewFormItem(wiimote.KeyMap[wiimote.BtnMinus].PrettyName, m.buttons[wiimote.BtnMinus]),
		widget.NewFormItem(wiimote.KeyMap[wiimote.BtnHome].PrettyName, m.buttons[wiimote.BtnHome]),
	)

	nf := widget.NewForm(
		widget.NewFormItem(wiimote.KeyMap[wiimote.BtnZ].PrettyName, m.buttons[wiimote.BtnZ]),
		widget.NewFormItem(wiimote.KeyMap[wiimote.BtnC].PrettyName, m.buttons[wiimote.BtnC]),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Wiimote", f),
		container.NewTabItem("Nunchuk", nf),
	)
	m.mainContainer = container.NewVBox(
		tabs,
		container.NewHBox(
			widget.NewLabel("Directions as seen when Wiimote held vertically"),
			layout.NewSpacer(),
			m.loadBtn,
			m.saveBtn,
		),
	)

	return m
}

func (m *MapEditor) UI() fyne.CanvasObject {
	return m.mainContainer
}

func (m *MapEditor) chooserChangeHandler(c wiimote.Keycode) func(fyne.KeyName) {
	return func(n fyne.KeyName) {
		m.mapping[c] = n
		if n == desktop.KeyNone {
			m.buttons[c].SetText("None")
		} else {
			m.buttons[c].SetText(string(n))
		}
	}
}

func (m *MapEditor) loadMap() {
	dialog.ShowFileOpen(m.loadMapFromFile, m.parentWindow)
}

func (m *MapEditor) loadMapFromFile(file fyne.URIReadCloser, err error) {
	if err != nil {
		loggo.Error(err)
		dialog.ShowError(err, m.parentWindow)
		return
	}
	if file == nil {
		return
	}
	defer file.Close()

	jsonBytes := make([]byte, 2048)
	n, err := file.Read(jsonBytes)
	if err != nil {
		loggo.Error(err)
		dialog.ShowError(errors.New("Unable to read file"), m.parentWindow)
		return
	}

	jsonBytes = jsonBytes[:n]
	err = json.Unmarshal(jsonBytes, &(m.mapping))
	if err != nil {
		loggo.Error(err)
		dialog.ShowError(errors.New("Unable to parse file"), m.parentWindow)
		return
	}

	m.updateButtonsFromMap()
}

func (m *MapEditor) saveMap() {
	dialog.ShowFileSave(m.saveMapToFile, m.parentWindow)
}

func (m *MapEditor) saveMapToFile(file fyne.URIWriteCloser, err error) {
	if err != nil {
		loggo.Error(err)
		dialog.ShowError(errors.New("Unable to save file"), m.parentWindow)
		return
	}
	if file == nil {
		return
	}
	defer file.Close()

	jsonBytes, err := json.MarshalIndent(m.mapping, "", "\t")
	if err != nil {
		loggo.Error(err)
		dialog.ShowError(errors.New("Unable to assemble file contents"), m.parentWindow)
		return
	}

	_, err = file.Write(jsonBytes)
	if err != nil {
		loggo.Error(err)
		dialog.ShowError(errors.New("Unable to write to file"), m.parentWindow)
	}
}

func (m *MapEditor) remapBtnFor(c wiimote.Keycode) *widget.Button {
	return widget.NewButton("None", func() {
		dialog.ShowCustom("Remap button", "OK", m.selectors[c], m.parentWindow)
	})
}

func (m *MapEditor) updateButtonsFromMap() {
	for wiicode, keyname := range m.mapping {
		m.selectors[wiicode].SetValue(keyname)
	}
}

func (m *MapEditor) KeyFor(code wiimote.Keycode) fyne.KeyName {
	name, ok := m.mapping[code]
	if !ok {
		return desktop.KeyNone
	}
	return name
}
