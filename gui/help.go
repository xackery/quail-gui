package gui

import (
	"github.com/xackery/quail-gui/popup"
	"github.com/xackery/wlk/walk"
)

var (
	menuHelpAbout *walk.Action
)

func onHelpAbout() {
	popup.MessageBoxf(mw, "About", "Quail GUI\n\nVersion: %s\n\nAuthor: Xackery\n\nLicense: MIT", "0.0.1")
}
