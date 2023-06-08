package mapeditor

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type KeyChooser struct {
	widget.BaseWidget

	parentWindow fyne.Window

	inputSelector    *widget.RadioGroup
	mouseBtnSelector *widget.Select
	keybdBtn         *widget.Button

	val fyne.KeyName

	OnChanged func(fyne.KeyName)
}

func NewKeyChooser(w fyne.Window) *KeyChooser {
	k := new(KeyChooser)
	k.ExtendBaseWidget(k)

	k.parentWindow = w
	k.inputSelector = widget.NewRadioGroup([]string{"Mouse", "Keyboard"}, k.conditionallyDisableInputs)
	k.mouseBtnSelector = widget.NewSelect([]string{
		string(MouseLeft),
		string(MouseMiddle),
		string(MouseRight),
	}, k.setSelectedMouseBtn)
	k.keybdBtn = widget.NewButton("None", k.setKeybdBtn)

	k.inputSelector.SetSelected("Keyboard")
	k.mouseBtnSelector.SetSelected(string(MouseLeft))

	k.val = desktop.KeyNone

	return k
}

func (k *KeyChooser) conditionallyDisableInputs(opt string) {
	switch opt {
	case "Mouse":
		k.keybdBtn.Disable()
		k.mouseBtnSelector.Enable()
	case "Keyboard":
		k.keybdBtn.Enable()
		k.mouseBtnSelector.Disable()
	default:
	}
}

func (k *KeyChooser) setSelectedMouseBtn(opt string) {
	k.val = fyne.KeyName(opt)
	if k.OnChanged != nil {
		k.OnChanged(k.val)
	}
}

func (k *KeyChooser) setKeybdBtn() {
	k.keybdBtn.SetText("Waiting for keypress...")
	k.parentWindow.Canvas().SetOnTypedKey(
		func(e *fyne.KeyEvent) {
			k.keybdBtn.SetText(string(e.Name))
			k.val = e.Name
			k.parentWindow.Canvas().SetOnTypedKey(nil)
			if k.OnChanged != nil {
				k.OnChanged(k.val)
			}
		},
	)
}

func (k *KeyChooser) Value() fyne.KeyName {
	return k.val
}

func (k *KeyChooser) SetValue(val fyne.KeyName) {
	k.val = val
	switch val {
	case MouseLeft, MouseRight, MouseMiddle:
		k.inputSelector.SetSelected("Mouse")
		k.mouseBtnSelector.SetSelected(string(val))
	default:
		k.inputSelector.SetSelected("Keyboard")
		if val == desktop.KeyNone {
			k.keybdBtn.SetText("None")
		} else {
			k.keybdBtn.SetText(string(val))
		}
	}

	if k.OnChanged != nil {
		k.OnChanged(k.val)
	}
}

func (k *KeyChooser) CreateRenderer() fyne.WidgetRenderer {
	return &keyChooserRenderer{selector: k.inputSelector, mouseInput: k.mouseBtnSelector, keybdInput: k.keybdBtn}
}

type keyChooserRenderer struct {
	selector   *widget.RadioGroup
	mouseInput *widget.Select
	keybdInput *widget.Button
}

func (r *keyChooserRenderer) Layout(sz fyne.Size) {
	r.selector.Move(fyne.NewPos(0, 0))
	r.selector.Resize(r.selector.MinSize())
	r.mouseInput.Move(fyne.NewPos(r.selector.MinSize().Width, 0))
	r.mouseInput.Resize(fyne.NewSize(sz.Width-r.selector.MinSize().Width, r.mouseInput.MinSize().Height))
	r.keybdInput.Move(fyne.NewPos(r.selector.MinSize().Width, r.mouseInput.MinSize().Height))
	r.keybdInput.Resize(fyne.NewSize(sz.Width-r.selector.MinSize().Width, r.keybdInput.MinSize().Height))
}

func (r *keyChooserRenderer) MinSize() fyne.Size {
	return fyne.NewSize(
		r.selector.MinSize().Width+200,
		r.selector.MinSize().Height,
	)
}

func (r *keyChooserRenderer) Refresh() {
	r.selector.Refresh()
	r.mouseInput.Refresh()
	r.keybdInput.Refresh()
}

func (r *keyChooserRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.selector, r.mouseInput, r.keybdInput}
}

func (r *keyChooserRenderer) Destroy() {}
