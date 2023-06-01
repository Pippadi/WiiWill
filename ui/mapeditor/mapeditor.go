package mapeditor

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Pippadi/WiiWill/wiimote"
)

type MapEditor struct {
	parentWindow  fyne.Window
	mainContainer *fyne.Container

	loadBtn *widget.Button
	saveBtn *widget.Button

	buttons map[wiimote.Keycode]*widget.Button
	mapping map[wiimote.Keycode]int
}

func New(w fyne.Window) *MapEditor {
	m := new(MapEditor)
	m.parentWindow = w
	m.buttons = make(map[wiimote.Keycode]*widget.Button)
	m.mapping = make(map[wiimote.Keycode]int)

	f := widget.NewForm()
	for c, k := range wiimote.KeyMap {
		m.buttons[c] = widget.NewButton("None", nil)
		m.buttons[c].OnTapped = m.RemapButtonHandler(m.buttons[c], c)
		m.mapping[c] = -1
		f.AppendItem(widget.NewFormItem(k.PrettyName, m.buttons[c]))
	}

	m.mainContainer = container.NewVBox(f)

	return m
}

func (m *MapEditor) UI() fyne.CanvasObject {
	return m.mainContainer
}

func (m *MapEditor) RemapButtonHandler(b *widget.Button, c wiimote.Keycode) func() {
	return func() {
		b.SetText("Waiting for keypress...")
		m.parentWindow.Canvas().SetOnTypedKey(
			func(e *fyne.KeyEvent) {
				b.SetText(string(e.Name))
				m.mapping[c] = e.Physical.ScanCode
				m.parentWindow.Canvas().SetOnTypedKey(nil)
			})
	}
}
