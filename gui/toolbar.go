package gui

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/walk"
)

var (
	toolbar = &toolbarBind{}
)

type toolbarBind struct {
	back *walk.Action
}

func (t *toolbarBind) onBack() {
	slog.Println("onBack")
	viewSetBack()
}
