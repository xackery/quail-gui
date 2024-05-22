package gui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/popup"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/walk"
)

var (
	menuFileNew        *walk.Action
	menuFileOpen       *walk.Action
	menuFileOpenRecent *walk.Action
	menuFileRefresh    *walk.Action
	menuFileDelete     *walk.Action
	menuFileSave       *walk.Action
	menuFileSaveAs     *walk.Action
	menuFileClose      *walk.Action
	menuFileExit       *walk.Action
)

func onFileNew() {
	slog.Println("new triggered")
}

func onFileOpen() {
	path, err := popup.Open(mw, "Open EQ Archive", "All Archives|*.pfs;*.eqg;*.s3d;*.pak|PFS Files (*.pfs)|*.pfs|EQG Files (*.eqg)|*.eqg|S3D Files (*.s3d)|*.s3d|PAK Files (*.pak)|*.pak|All Files (*.*)|*.*", ".")
	if err != nil {
		if err.Error() == "cancelled" {
			return
		}
		popup.Errorf(mw, "open: %s", err)
		return
	}
	slog.Printf("Menu Opening %s\n", path)
	err = Open(path)
	if err != nil {
		popup.Errorf(mw, "gui open: %s", err)
		return
	}
}

func onFileOpenRecent() {
	slog.Println("open recent triggered")
}

func onFileRefresh() {
	err := Open(archivePath)
	if err != nil {
		popup.Errorf(mw, "refresh: %s", err)
		return
	}
}

func onFileDelete() {
	slog.Println("delete triggered")
}

func onFileSave() {
	if archive == nil {
		popup.Errorf(mw, "file save: nothing to save")
		return
	}

	w, err := os.Create(archivePath)
	if err != nil {
		popup.Errorf(mw, "create: %s", err)
		return
	}
	defer w.Close()

	err = archive.Write(w)
	if err != nil {
		popup.Errorf(mw, "save: %s", err)
		return
	}

	isEdited = false
	fileName := filepath.Base(archivePath)
	mw.SetTitle(fileName)

	lastSelection := file.CurrentIndex()

	entries := []*component.FileViewEntry{}
	files := archive.Files()

	for _, fe := range files {
		ext := strings.ToLower(filepath.Ext(fe.Name()))
		img, err := ico.Generate(ext, fe.Data())
		if err != nil {
			slog.Printf("Failed to generate icon for %s: %s\n", fe.Name(), err.Error())
			img = ico.Grab("unk")
		}

		fve := &component.FileViewEntry{
			Icon:    img,
			Name:    fe.Name(),
			Ext:     ext,
			Size:    generateSize(len(fe.Data())),
			RawSize: len(fe.Data()),
		}

		entries = append(entries, fve)
	}
	fileView.SetItems(entries)
	if lastSelection >= len(entries) {
		lastSelection = len(entries) - 1
	}
	file.SetCurrentIndex(lastSelection)

	slog.Printf("Saved %s\n", fileName)
}

func onFileSaveAs() {
	slog.Println("save as triggered")
}

func onFileExit() {
	slog.Println("File Exit triggered")
	err := mw.Close()
	if err != nil {
		popup.Errorf(mw, "close: %s", err)
	}
	walk.App().Exit(0)
	slog.Dump()
}

func onFileClose() {
	if archive == nil {
		slog.Println("No archive to close")
		return
	}
	err := archive.Close()
	if err != nil {
		popup.Errorf(mw, "archive close: %s", err)
	}
	archive = nil
	archivePath = ""

	fileView.ResetRows()
	entrySetActive(false)

	slog.Printf("Closed %s\n", archivePath)
}
