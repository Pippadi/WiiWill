package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/bendahl/uinput"
)

var fyneToUinputKey = map[fyne.KeyName]int{
	fyne.Key1:               uinput.Key1,
	fyne.Key2:               uinput.Key2,
	fyne.Key3:               uinput.Key3,
	fyne.Key4:               uinput.Key4,
	fyne.Key5:               uinput.Key5,
	fyne.Key6:               uinput.Key6,
	fyne.Key7:               uinput.Key7,
	fyne.Key8:               uinput.Key8,
	fyne.Key9:               uinput.Key9,
	fyne.Key0:               uinput.Key0,
	fyne.KeyMinus:           uinput.KeyMinus,
	fyne.KeyEqual:           uinput.KeyEqual,
	fyne.KeyBackspace:       uinput.KeyBackspace,
	fyne.KeyTab:             uinput.KeyTab,
	fyne.KeyQ:               uinput.KeyQ,
	fyne.KeyW:               uinput.KeyW,
	fyne.KeyE:               uinput.KeyE,
	fyne.KeyR:               uinput.KeyR,
	fyne.KeyT:               uinput.KeyT,
	fyne.KeyY:               uinput.KeyY,
	fyne.KeyU:               uinput.KeyU,
	fyne.KeyI:               uinput.KeyI,
	fyne.KeyO:               uinput.KeyO,
	fyne.KeyP:               uinput.KeyP,
	fyne.KeyEnter:           uinput.KeyEnter,
	fyne.KeyA:               uinput.KeyA,
	fyne.KeyS:               uinput.KeyS,
	fyne.KeyD:               uinput.KeyD,
	fyne.KeyF:               uinput.KeyF,
	fyne.KeyG:               uinput.KeyG,
	fyne.KeyH:               uinput.KeyH,
	fyne.KeyJ:               uinput.KeyJ,
	fyne.KeyK:               uinput.KeyK,
	fyne.KeyL:               uinput.KeyL,
	fyne.KeySemicolon:       uinput.KeySemicolon,
	fyne.KeyApostrophe:      uinput.KeyApostrophe,
	fyne.KeyBackslash:       uinput.KeyBackslash,
	fyne.KeyZ:               uinput.KeyZ,
	fyne.KeyX:               uinput.KeyX,
	fyne.KeyC:               uinput.KeyC,
	fyne.KeyV:               uinput.KeyV,
	fyne.KeyB:               uinput.KeyB,
	fyne.KeyN:               uinput.KeyN,
	fyne.KeyM:               uinput.KeyM,
	fyne.KeyComma:           uinput.KeyComma,
	fyne.KeySlash:           uinput.KeySlash,
	fyne.KeySpace:           uinput.KeySpace,
	fyne.KeyF1:              uinput.KeyF1,
	fyne.KeyF2:              uinput.KeyF2,
	fyne.KeyF3:              uinput.KeyF3,
	fyne.KeyF4:              uinput.KeyF4,
	fyne.KeyF5:              uinput.KeyF5,
	fyne.KeyF6:              uinput.KeyF6,
	fyne.KeyF7:              uinput.KeyF7,
	fyne.KeyF8:              uinput.KeyF8,
	fyne.KeyF9:              uinput.KeyF9,
	fyne.KeyF10:             uinput.KeyF10,
	fyne.KeyF11:             uinput.KeyF11,
	fyne.KeyF12:             uinput.KeyF12,
	fyne.KeyHome:            uinput.KeyHome,
	fyne.KeyUp:              uinput.KeyUp,
	fyne.KeyLeft:            uinput.KeyLeft,
	fyne.KeyRight:           uinput.KeyRight,
	fyne.KeyEnd:             uinput.KeyEnd,
	fyne.KeyDown:            uinput.KeyDown,
	fyne.KeyInsert:          uinput.KeyInsert,
	fyne.KeyDelete:          uinput.KeyDelete,
	fyne.KeyLeftBracket:     uinput.KeyLeftbrace,
	fyne.KeyRightBracket:    uinput.KeyRightbrace,
	desktop.KeyAltLeft:      uinput.KeyLeftalt,
	desktop.KeyControlLeft:  uinput.KeyLeftctrl,
	desktop.KeySuperLeft:    uinput.KeyLeftmeta,
	desktop.KeyShiftLeft:    uinput.KeyLeftshift,
	desktop.KeyAltRight:     uinput.KeyRightalt,
	desktop.KeyControlRight: uinput.KeyRightctrl,
	desktop.KeySuperRight:   uinput.KeyRightmeta,
	desktop.KeyShiftRight:   uinput.KeyRightshift,
	desktop.KeyCapsLock:     uinput.KeyCapslock,
	desktop.KeyMenu:         uinput.KeyMenu,
}
