package gui

import (
	"path/filepath"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/op"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/walk"
)

var (
	widget = &widgetBind{}
)

type widgetBind struct {
	file        *walk.TableView
	fileView    *component.FileView
	element     *walk.TableView
	elementView *component.ElementView
}

func (w *widgetBind) onSizeChanged() {
	slog.Println("onSizeChanged")
}

func (w *widgetBind) breadcrumbRefresh() {
	text := ""
	if archive != nil {
		text += filepath.Base(lastPath)
	}
	focus := op.Breadcrumb()
	for focus != "" {
		if text != "" {
			text += " > "
		}

		text += focus
	}
	slog.Println("setting title to", text)
	mw.SetTitle(text)
}
