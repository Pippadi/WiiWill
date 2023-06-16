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

	buttons   map[wiimote.Keycode]*widget.Button
	selectors map[wiimote.Keycode]*KeyChooser
	Mapping   map[wiimote.Keycode]fyne.KeyName `json:"KeyMap"`

	StickConfigs       map[wiimote.Stick]StickConfig `json:"StickConfig"`
	stickConfigurators map[wiimote.Stick]*StickConfigurator

	mapFile string
}

func New(w fyne.Window) *MapEditor {
	m := new(MapEditor)
	m.parentWindow = w
	m.buttons = make(map[wiimote.Keycode]*widget.Button)
	m.selectors = make(map[wiimote.Keycode]*KeyChooser)
	m.Mapping = make(map[wiimote.Keycode]fyne.KeyName)
	m.StickConfigs = make(map[wiimote.Stick]StickConfig)
	m.stickConfigurators = make(map[wiimote.Stick]*StickConfigurator)

	m.loadBtn = widget.NewButtonWithIcon("Load", theme.UploadIcon(), m.loadMap)
	m.saveBtn = widget.NewButtonWithIcon("Save", theme.FileIcon(), m.saveMap)

	for c, _ := range wiimote.KeyMap {
		m.buttons[c] = m.remapBtnFor(c)
		m.selectors[c] = NewKeyChooser(m.parentWindow)
		m.selectors[c].OnChanged = m.chooserChangeHandler(c)
		m.Mapping[c] = desktop.KeyNone
	}
	for s, _ := range wiimote.StickMap {
		m.StickConfigs[s] = StickConfig{}
		m.stickConfigurators[s] = NewStickConfigurator(m.parentWindow)
		m.stickConfigurators[s].OnChanged = m.updateStickCfgFor(s)
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
		widget.NewFormItem(wiimote.StickMap[wiimote.NunchukStick], m.configBtnFor(wiimote.NunchukStick)),
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
		m.Mapping[c] = n
		if n == desktop.KeyNone {
			m.buttons[c].SetText("None")
		} else {
			m.buttons[c].SetText(string(n))
		}
	}
}

func (m *MapEditor) loadMap() {
	dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
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
		jsonBytes = jsonBytes[:n]
		if err != nil {
			loggo.Error(err)
			dialog.ShowError(errors.New("Unable to read file"), m.parentWindow)
		}

		err = m.LoadMapFromBytes(jsonBytes)
		if err != nil {
			loggo.Error(err)
			dialog.ShowError(errors.New("Unable to parse file"), m.parentWindow)
		}
		m.mapFile = file.URI().Path()
	}, m.parentWindow,
	)
}

func (m *MapEditor) LoadMapFromBytes(jsonBytes []byte) error {
	err := json.Unmarshal(jsonBytes, m)
	if err != nil {
		return err
	}

	m.updateButtonsFromMap()
	return nil
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

	jsonBytes, err := json.MarshalIndent(m, "", "\t")
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

	m.mapFile = file.URI().Path()
}

func (m *MapEditor) remapBtnFor(c wiimote.Keycode) *widget.Button {
	return widget.NewButton("None", func() {
		dialog.ShowCustom("Remap button", "OK", m.selectors[c], m.parentWindow)
	})
}

func (m *MapEditor) configBtnFor(c wiimote.Stick) *widget.Button {
	return widget.NewButton("Configure", func() {
		dialog.ShowCustom("Configure stick", "OK", m.stickConfigurators[c], m.parentWindow)
	})
}

func (m *MapEditor) updateStickCfgFor(c wiimote.Stick) func(StickConfig) {
	return func(cfg StickConfig) {
		m.StickConfigs[c] = cfg
	}
}

func (m *MapEditor) updateButtonsFromMap() {
	for wiicode, keyname := range m.Mapping {
		m.selectors[wiicode].SetValue(keyname)
	}
	for wiicode, cfg := range m.StickConfigs {
		m.stickConfigurators[wiicode].SetValue(cfg)
	}
}

func (m *MapEditor) KeyFor(code wiimote.Keycode) fyne.KeyName {
	name, ok := m.Mapping[code]
	if !ok {
		return desktop.KeyNone
	}
	return name
}

func (m *MapEditor) StickConfigFor(stick wiimote.Stick) StickConfig {
	cfg, ok := m.StickConfigs[stick]
	if !ok {
		return StickConfig{false, 0}
	}
	return cfg
}

func (m *MapEditor) MapFile() string { return m.mapFile }
