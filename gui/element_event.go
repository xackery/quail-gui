package gui

import (
	"fmt"
	"path/filepath"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/op"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/common"
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

	entries := []*component.ElementViewEntry{}
	switch val := focus.(type) {
	case *common.Model:
		entries = append(entries, &component.ElementViewEntry{
			Icon: ico.Grab("header"),
			Name: fmt.Sprintf("Header version %d", val.Header.Version),
		})

		switch filepath.Ext(node.Name()) {
		case ".mod":
			entries = append(entries, &component.ElementViewEntry{
				Icon: ico.Grab("material"),
				Name: fmt.Sprintf("Material List (%d)", len(val.Materials)),
			})

			entries = append(entries, &component.ElementViewEntry{
				Icon: ico.Grab("bone"),
				Name: fmt.Sprintf("Bone List (%d)", len(val.Bones)),
			})

			entries = append(entries, &component.ElementViewEntry{
				Icon: ico.Grab("triangle"),
				Name: fmt.Sprintf("Triangle List (%d)", len(val.Triangles)),
			})

			entries = append(entries, &component.ElementViewEntry{
				Icon: ico.Grab("vertex"),
				Name: fmt.Sprintf("Vertex List (%d)", len(val.Vertices)),
			})

		}
	}
	widget.elementView.SetItems(entries)
}

func (m *menuBind) onElementDelete() {
	slog.Println("delete triggered")
}
