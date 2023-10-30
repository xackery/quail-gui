package gui

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/quail"
	"github.com/xackery/wlk/walk"
)

var (
	menu         = &menuBind{}
	lastPath     string
	archive      *pfs.PFS
	selectedFile string
	base         interface{}
	section      interface{}
)

type menuBind struct {
	fileNew        *walk.Action
	fileOpen       *walk.Action
	fileOpenRecent *walk.Action
	fileRefresh    *walk.Action
	fileDelete     *walk.Action
	fileExit       *walk.Action
	helpAbout      *walk.Action
}

func (m *menuBind) onFileNew() {
	slog.Println("new triggered")

}

func Open(path string, fileName string, section string) error {
	if mw == nil {
		return fmt.Errorf("main window not created")
	}

	slog.Printf("Opening path: %s, file: %s, section: %s\n", path, fileName, section)

	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer r.Close()

	isPFS := false
	ext := filepath.Ext(strings.ToLower(path))
	switch ext {
	case ".pfs":
		isPFS = true
	case ".eqg":
		isPFS = true
	case ".s3d":
		isPFS = true
	case ".pak":
		isPFS = true
	}

	if !isPFS {
		base, err = quail.Open(filepath.Base(path), r)
		if err != nil {
			return fmt.Errorf("open file %s: %w", path, err)
		}
		viewSet(currentViewContext)
		return nil
	}
	archive, err = pfs.New(filepath.Base(path))
	if err != nil {
		return fmt.Errorf("open archive %s: %w", path, err)
	}
	err = archive.Decode(r)
	if err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}

	entries := []*component.FileViewEntry{}
	files := archive.Files()
	for _, fe := range files {
		ext := strings.ToLower(filepath.Ext(fe.Name()))
		fve := &component.FileViewEntry{
			Icon:    ico.Generate(strings.ToLower(filepath.Ext(fe.Name())), fe.Data()),
			Name:    fe.Name(),
			Ext:     ext,
			Size:    generateSize(len(fe.Data())),
			RawSize: len(fe.Data()),
		}

		entries = append(entries, fve)
	}
	widget.fileView.SetItems(entries)
	slog.Printf("Loaded %d files\n", len(entries))

	if len(fileName) == 0 {
		viewSet(currentViewArchiveFiles)
		return nil
	}

	selectedFile = fileName

	data, err := archive.File(fileName)
	if err != nil {
		return fmt.Errorf("file %s: %w", fileName, err)
	}

	base, err = quail.Open(fileName, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("open file %s: %w", fileName, err)
	}

	viewSet(currentViewContext)
	return nil
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
	slog.Println("refresh triggered")
}

func (m *menuBind) onFileDelete() {
	slog.Println("delete triggered")
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

func (m *menuBind) onHelpAbout() {
	slog.Println("about triggered")
}
