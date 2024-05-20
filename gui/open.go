package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/raw"
)

func Open(path string) error {
	if mw == nil {
		return fmt.Errorf("main window not created")
	}

	if path == "" {
		return fmt.Errorf("cancelled")
	}

	slog.Printf("Opening path: %s\n", path)

	r, err := os.Open(path)
	if err != nil {
		return err
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
		name := filepath.Base(path)
		ext := strings.ToLower(filepath.Ext(name))
		value, err := raw.Read(ext, r)
		if err != nil {
			return fmt.Errorf("quail.Open: %w", err)
		}
		value.SetFileName(name)

		return nil
	}

	archive, err = pfs.New(filepath.Base(path))
	if err != nil {
		return fmt.Errorf("pfs.New: %w", err)
	}

	err = archive.Read(r)
	if err != nil {
		archive = nil
		return fmt.Errorf("decode: %w", err)
	}

	archivePath = path

	isWorldFile := false
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
		if ext == ".wld" || ext == ".zon" {
			isWorldFile = true
		}

		entries = append(entries, fve)
	}
	fileView.SetItems(entries)
	file.SetLastColumnStretched(true)
	slog.Printf("Loaded %d files\n", len(entries))
	if len(entries) > 0 {
		file.SetCurrentIndex(0)
		entrySetActive(true)
	}
	menuEntryEditWorld.SetEnabled(isWorldFile)

	fileName := filepath.Base(path)

	mw.SetTitle(fileName)

	return nil
}
