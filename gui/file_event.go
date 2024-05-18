package gui

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail-gui/gui/dialog"
	"github.com/xackery/quail-gui/op"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/raw"
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

	ext := filepath.Ext(strings.ToLower(item.Name))
	value, err := raw.Read(ext, bytes.NewReader(data))
	if err != nil {
		slog.Printf("Failed to open file %s: %s\n", item.Name, err.Error())
		return
	}
	value.SetFileName(item.Name)

	node := op.NewNode(item.Name, value)
	op.SetRoot(node)
	op.SetFocus(node)

	slog.Printf("Selected file: %s\n", item.Name)

	extFuncs := map[string]func(*walk.MainWindow, *op.Node) error{
		".mod": dialog.ShowModEdit,
		".zon": dialog.ShowZonEdit,
		".wld": dialog.ShowWldEdit,
		".mds": dialog.ShowMdsEdit,
	}

	nodeExt := strings.ToLower(filepath.Ext(item.Name))
	for ext, f := range extFuncs {
		if nodeExt != ext {
			continue
		}
		err = f(mw, node)
		if err != nil {
			slog.Printf("Failed to show %s edit: %s\n", ext, err.Error())
			ShowError(fmt.Errorf("edit %s: %w", item.Name, err))
			return
		}
	}

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
		walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxIconError)
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
