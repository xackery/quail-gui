package gui

import (
	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/walk"
)

var (
	widget = &widgetBind{}
)

type widgetBind struct {
	file       *walk.TableView
	fileView   *component.FileView
	breadcrumb *walk.Label
}

func (w *widgetBind) onFileChange() {
	slog.Println("onFileChange")
}

func (w *widgetBind) onFileActivated() {
	slog.Println("onFileActivated")

	if w.file.CurrentIndex() < 0 {
		slog.Println("current index is less than 0")
		return
	}

	if w.file.CurrentIndex() >= w.fileView.RowCount() {
		slog.Println("current index is greater than row count")
		return
	}
	item := w.fileView.Item(w.file.CurrentIndex())
	if item == nil {
		slog.Println("item is nil")
		return
	}

	selectedFile = item.Name
	slog.Printf("selected file: %s\n", selectedFile)
	viewSet(currentViewContext)
}

func (w *widgetBind) onSizeChanged() {
	slog.Println("onSizeChanged")
}

func (w *widgetBind) breadcrumbRefresh() {
	text := ""
	if archive != nil {
		text += lastPath
	}
	if currentView != currentViewArchiveFiles {
		if text != "" {
			text += " > "
		}
		text += selectedFile
	}
	w.breadcrumb.SetText(text)
}
