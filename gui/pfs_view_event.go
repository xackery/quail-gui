package gui

import (
	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/op"
	"github.com/xackery/quail-gui/slog"
)

func (w *widgetBind) onElementChange() {
	slog.Println("onElementChange")
}

func (w *widgetBind) onElementActivated() {
	slog.Println("onElementActivated")
}

func (m *menuBind) onElementRefresh() {
	slog.Println("element refresh triggered")
	node := op.Focus()
	if node == nil {
		slog.Println("node is nil")
		return
	}
	focus := node.Value()
	if focus == nil {
		slog.Println("focus is nil")
		return
	}

	entries := []*component.PfsViewEntry{}

	widget.pfsList.SetItems(entries)
}

func (m *menuBind) onElementDelete() {
	slog.Println("delete triggered")
}
