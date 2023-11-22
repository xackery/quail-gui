package gui

import (
	"bytes"

	"github.com/xackery/quail-gui/op"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/quail"
	"github.com/xackery/wlk/walk"
)

func (w *widgetBind) onFileChange() {
	slog.Println("onFileChange")

}

func (w *widgetBind) onFileActivated() {
	var err error
	slog.Println("onFileActivated")
	op.Clear()

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

	if archive == nil {
		slog.Println("archive is nil")
		return
	}

	data, err := archive.File(item.Name)
	if err != nil {
		slog.Printf("Failed to open file %s: %s\n", item.Name, err.Error())
		return
	}

	value, err := quail.Open(item.Name, bytes.NewReader(data))
	if err != nil {
		slog.Printf("Failed to open file %s: %s\n", item.Name, err.Error())
		return
	}

	node := op.NewNode(item.Name, value)
	op.SetRoot(node)
	op.SetFocus(node)

	slog.Printf("Selected file: %s\n", item.Name)
	viewSet(currentViewElement)
}

func (m *menuBind) onFileNew() {
	slog.Println("new triggered")

}

func (m *menuBind) onFileOpen() {
	path, err := ShowOpen("Open EQ Archive", "All Archives|*.pfs;*.eqg;*.s3d;*.pak|PFS Files (*.pfs)|*.pfs|EQG Files (*.eqg)|*.eqg|S3D Files (*.s3d)|*.s3d|PAK Files (*.pak)|*.pak", ".")
	if err != nil {
		slog.Printf("Failed to open: %s\n", err)
		return
	}
	slog.Printf("Menu Opening %s\n", path)
	err = Open(path, "", "")
	if err != nil {
		slog.Printf("Failed to open: %s\n", err)
		return
	}
}

func (m *menuBind) onFileOpenRecent() {
	slog.Println("open recent triggered")
}

func (m *menuBind) onFileRefresh() {
	slog.Println("file refresh triggered")
}

func (m *menuBind) onFileDelete() {
	slog.Println("delete triggered")
}

func (m *menuBind) onFileSave() {
	slog.Println("save triggered")
}

func (m *menuBind) onFileExit() {
	slog.Println("File Exit triggered")
	err := mw.Close()
	if err != nil {
		slog.Printf("Failed to close: %s\n", err.Error())
	}
	walk.App().Exit(0)
	slog.Dump()
}
