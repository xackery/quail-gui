package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/op"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
	"github.com/xackery/wlk/walk"
)

var (
	menu     = &menuBind{}
	lastPath string
	archive  *pfs.Pfs
)

type menuBind struct {
	fileNew        *walk.Action
	fileOpen       *walk.Action
	fileOpenRecent *walk.Action
	fileRefresh    *walk.Action
	fileDelete     *walk.Action
	fileSave       *walk.Action
	fileExit       *walk.Action
	helpAbout      *walk.Action
	elementRefresh *walk.Action
	elementDelete  *walk.Action
}

func Open(path string, fileName string, element string) error {
	if mw == nil {
		return fmt.Errorf("main window not created")
	}

	op.Clear()

	slog.Printf("Opening path: %s, file: %s, section: %s\n", path, fileName, element)

	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	lastPath = path

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
		name := filepath.Base(path)
		ext := strings.ToLower(filepath.Ext(name))
		value, err := raw.Read(ext, r)
		if err != nil {
			return fmt.Errorf("quail.Open: %w", err)
		}
		value.SetFileName(name)

		node := op.NewNode(name, value)
		op.SetRoot(node)
		op.SetFocus(node)
		viewSet(pfsList)
		return nil
	}
	archive, err = pfs.New(filepath.Base(path))
	if err != nil {
		return fmt.Errorf("pfs.New: %w", err)
	}
	err = archive.Read(r)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
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
	/*
		data, err := archive.File(fileName)
		if err != nil {
			return fmt.Errorf("file %s: %w", fileName, err)
		}

		value, err := raw.Read(ext, bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("raw read %s: %w", fileName, err)
		}

		node := op.NewNode(fileName, value)
		op.SetRoot(node)
		op.SetFocus(node) */
	//viewSet(currentViewElement)
	return nil
}

func (m *menuBind) onHelpAbout() {
	slog.Println("about triggered")
}
