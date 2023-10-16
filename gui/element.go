package gui

import "github.com/xackery/wlk/walk"

func newLabel(text string) *walk.Label {
	label, err := walk.NewLabel(gui.mw)
	if err != nil {
		panic(err)
	}
	label.SetText(text)
	return label
}
