package gui

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/op"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/quail"
	"github.com/xackery/wlk/walk"
)

var (
	menu     = &menuBind{}
	lastPath string
	archive  *pfs.PFS
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
		return fmt.Errorf("open %s: %w", path, err)
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
		value, err := quail.Open(name, r)
		if err != nil {
			return fmt.Errorf("open file %s: %w", path, err)
		}

		node := op.NewNode(name, value)
		op.SetRoot(node)
		op.SetFocus(node)
		viewSet(currentViewElement)
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

	data, err := archive.File(fileName)
	if err != nil {
		return fmt.Errorf("file %s: %w", fileName, err)
	}

	value, err := quail.Open(fileName, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("open file %s: %w", fileName, err)
	}

	node := op.NewNode(fileName, value)
	op.SetRoot(node)
	op.SetFocus(node)
	//viewSet(currentViewElement)
	return nil
}

func (m *menuBind) onHelpAbout() {
	slog.Println("about triggered")
}
