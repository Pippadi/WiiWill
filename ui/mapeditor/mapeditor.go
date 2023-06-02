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
	parentWindow  fyne.Window
	mainContainer fyne.CanvasObject

	loadBtn *widget.Button
	saveBtn *widget.Button

	buttons map[wiimote.Keycode]*widget.Button
	mapping map[wiimote.Keycode]fyne.KeyName
}

func New(w fyne.Window) *MapEditor {
	m := new(MapEditor)
	m.parentWindow = w
	m.buttons = make(map[wiimote.Keycode]*widget.Button)
	m.mapping = make(map[wiimote.Keycode]fyne.KeyName)

	m.loadBtn = widget.NewButtonWithIcon("Load", theme.UploadIcon(), m.loadMap)
	m.saveBtn = widget.NewButtonWithIcon("Save", theme.FileIcon(), m.saveMap)

	f := widget.NewForm()
	for c, _ := range wiimote.KeyMap {
		m.buttons[c] = widget.NewButton("None", nil)
		m.buttons[c].OnTapped = m.remapButtonHandler(m.buttons[c], c)
		m.mapping[c] = desktop.KeyNone
	}
	// Map stores elements in arbitrary order, so order buttons manually
	f.AppendItem(widget.NewFormItem(wiimote.KeyMap[wiimote.BtnUp].PrettyName, m.buttons[wiimote.BtnUp]))
	f.AppendItem(widget.NewFormItem(wiimote.KeyMap[wiimote.BtnDown].PrettyName, m.buttons[wiimote.BtnDown]))
	f.AppendItem(widget.NewFormItem(wiimote.KeyMap[wiimote.BtnLeft].PrettyName, m.buttons[wiimote.BtnLeft]))
	f.AppendItem(widget.NewFormItem(wiimote.KeyMap[wiimote.BtnRight].PrettyName, m.buttons[wiimote.BtnRight]))
	f.AppendItem(widget.NewFormItem(wiimote.KeyMap[wiimote.BtnA].PrettyName, m.buttons[wiimote.BtnA]))
	f.AppendItem(widget.NewFormItem(wiimote.KeyMap[wiimote.BtnB].PrettyName, m.buttons[wiimote.BtnB]))
	f.AppendItem(widget.NewFormItem(wiimote.KeyMap[wiimote.Btn1].PrettyName, m.buttons[wiimote.Btn1]))
	f.AppendItem(widget.NewFormItem(wiimote.KeyMap[wiimote.Btn2].PrettyName, m.buttons[wiimote.Btn2]))
	f.AppendItem(widget.NewFormItem(wiimote.KeyMap[wiimote.BtnPlus].PrettyName, m.buttons[wiimote.BtnPlus]))
	f.AppendItem(widget.NewFormItem(wiimote.KeyMap[wiimote.BtnMinus].PrettyName, m.buttons[wiimote.BtnMinus]))
	f.AppendItem(widget.NewFormItem(wiimote.KeyMap[wiimote.BtnHome].PrettyName, m.buttons[wiimote.BtnHome]))

	m.mainContainer = container.NewVBox(
		widget.NewLabel("Wiimote held vertically"),
		f,
		container.NewHBox(m.loadBtn, layout.NewSpacer(), m.saveBtn),
	)

	return m
}

func (m *MapEditor) UI() fyne.CanvasObject {
	return m.mainContainer
}

func (m *MapEditor) remapButtonHandler(b *widget.Button, c wiimote.Keycode) func() {
	return func() {
		b.SetText("Waiting for keypress...")
		m.parentWindow.Canvas().SetOnTypedKey(
			func(e *fyne.KeyEvent) {
				b.SetText(string(e.Name))
				m.mapping[c] = e.Name
				m.parentWindow.Canvas().SetOnTypedKey(nil)
			})
	}
}

func (m *MapEditor) loadMap() {
	dialog.ShowFileOpen(m.loadMapFromFile, m.parentWindow)
}

func (m *MapEditor) loadMapFromFile(file fyne.URIReadCloser, err error) {
	defer file.Close()
	if err != nil {
		loggo.Error(err)
		dialog.ShowError(err, m.parentWindow)
		return
	}

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
	defer file.Close()
	if err != nil {
		loggo.Error(err)
		dialog.ShowError(errors.New("Unable to save file"), m.parentWindow)
		return
	}

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

func (m *MapEditor) updateButtonsFromMap() {
	for wiicode, keyname := range m.mapping {
		if keyname == desktop.KeyNone {
			m.buttons[wiicode].SetText("None")
		} else {
			m.buttons[wiicode].SetText(string(keyname))
		}
	}
}

func (m *MapEditor) KeyFor(code wiimote.Keycode) fyne.KeyName {
	name, ok := m.mapping[code]
	if !ok {
		return desktop.KeyNone
	}
	return name
}
