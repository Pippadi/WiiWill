package mapeditor

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type StickConfig struct {
	AsMouse bool `json:"AsMouse"`

	// Multiplier
	Speed float64 `json:"SpeedMultiplier"`
}

type StickConfigurator struct {
	widget.BaseWidget

	parentWindow fyne.Window

	asMouseChk  *widget.Check
	speedSlider *widget.Slider

	OnChanged func(StickConfig)
}

func NewStickConfigurator(w fyne.Window) *StickConfigurator {
	k := new(StickConfigurator)
	k.ExtendBaseWidget(k)

	k.parentWindow = w

	k.asMouseChk = widget.NewCheck("Control mouse", k.conditionallyDisableInputs)
	k.speedSlider = widget.NewSlider(0.05, 1.0)
	k.speedSlider.Step = 0.05
	k.speedSlider.OnChanged = k.setSpeed

	k.speedSlider.SetValue(1.0)
	k.asMouseChk.SetChecked(true)

	return k
}

func (k *StickConfigurator) conditionallyDisableInputs(asMouse bool) {
	if asMouse {
		k.speedSlider.Show()
	} else {
		k.speedSlider.Hide()
	}

	if k.OnChanged != nil {
		k.OnChanged(k.Value())
	}
}

func (k *StickConfigurator) setSpeed(float64) {
	if k.OnChanged != nil {
		k.OnChanged(k.Value())
	}
}

func (k *StickConfigurator) Value() StickConfig {
	return StickConfig{AsMouse: k.asMouseChk.Checked, Speed: k.speedSlider.Value}
}

func (k *StickConfigurator) SetValue(val StickConfig) {
	k.asMouseChk.SetChecked(val.AsMouse)
	k.speedSlider.SetValue(val.Speed)
}

func (k *StickConfigurator) CreateRenderer() fyne.WidgetRenderer {
	return &stickConfigRenderer{check: k.asMouseChk, slider: k.speedSlider}
}

type stickConfigRenderer struct {
	check  *widget.Check
	slider *widget.Slider
}

func (r *stickConfigRenderer) Layout(sz fyne.Size) {
	r.check.Move(fyne.NewPos(sz.Width/2-r.check.MinSize().Width/2, 0))
	r.check.Resize(r.check.MinSize())
	r.slider.Move(fyne.NewPos(0, r.check.MinSize().Height))
	r.slider.Resize(fyne.NewSize(sz.Width, r.slider.MinSize().Height))
}

func (r *stickConfigRenderer) MinSize() fyne.Size {
	return fyne.NewSize(
		r.slider.MinSize().Width,
		r.check.MinSize().Height+r.slider.MinSize().Height,
	)
}

func (r *stickConfigRenderer) Refresh() {
	r.check.Refresh()
	r.slider.Refresh()
}

func (r *stickConfigRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.check, r.slider}
}

func (r *stickConfigRenderer) Destroy() {}
