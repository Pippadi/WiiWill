package mapeditor

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MapEditor struct {
	parentWindow  fyne.Window
	mainContainer *fyne.Container

	loadBtn *widget.Button
	saveBtn *widget.Button

	aBtn *widget.Button
}

func New(w fyne.Window) *MapEditor {
	m := new(MapEditor)
	m.parentWindow = w

	m.aBtn = widget.NewButton("A", func() {
		m.aBtn.SetText("Waiting for keypress...")
		m.parentWindow.Canvas().SetOnTypedKey(
			func(e *fyne.KeyEvent) {
				m.aBtn.SetText(string(e.Name))
				m.parentWindow.Canvas().SetOnTypedKey(nil)
			})
	})
	f := widget.NewForm(
		widget.NewFormItem("A", m.aBtn),
	)

	m.mainContainer = container.NewVBox(f)

	return m
}

func (m *MapEditor) UI() fyne.CanvasObject {
	return m.mainContainer
}
