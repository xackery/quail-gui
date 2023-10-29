package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/xackery/quail-gui/archive"
	"github.com/xackery/quail-gui/gui"
	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/qmux"
	"github.com/xackery/quail-gui/slog"
)

func (c *Client) onArchiveOpen(path string, file string, isSelect bool) error {
	return c.open(path, file, isSelect)
}

func (c *Client) open(path string, file string, isSelect bool) error {
	var err error

	slog.Printf("Client Opening %s %s\n", path, file)

	err = archive.Open(path, file)
	if err != nil {
		return fmt.Errorf("archive open: %w", err)
	}

	if path == "" {
		path = c.openPath
	}

	if c.openPath != path {
		slog.Printf("Path changed from %s to %s\n", c.openPath, path)
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
		gui.SetFileViewItems(entries)
	}
	c.openPath = path
	if c.fileName == file {
		return nil
	}
	c.fileName = file

	if file == "" {
		return nil
	}

	fileExt := strings.ToLower(filepath.Ext(file))
	data, err := archive.File(file)
	if err != nil {
		return fmt.Errorf("file %s: %w", file, err)
	}

	type validExt struct {
		ext    string
		decode func(name string, r io.ReadSeeker) (*component.TreeModel, error)
	}

	validExts := []validExt{
		{".lay", qmux.LayDecode},
		{".mod", qmux.ModDecode},
		{".ter", qmux.TerDecode},
		{".pts", qmux.PtsDecode},
		{".prt", qmux.PrtDecode},
		{".ani", qmux.AniDecode},
		{".zon", qmux.ZonDecode},
		{".wld", qmux.WldDecode},
	}

	for _, ve := range validExts {
		if fileExt != ve.ext {
			continue
		}
		treeModel, err := ve.decode(file, bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("decode %s: %w", file, err)
		}

		gui.SetTreeModel(treeModel)
		return nil
	}

	return nil
}

func generateSize(in int) string {
	val := float64(in)
	if val < 1024 {
		return fmt.Sprintf("%0.0f bytes", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.0f KB", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.0f MB", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.0f GB", val)
	}
	val /= 1024
	return fmt.Sprintf("%0.0f TB", val)
}

func (c *Client) onArchiveSave(path string) error {
	files := archive.Files()
	totalSize := 0
	for _, file := range files {
		totalSize += len(file.Data())
	}
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()
	go func() {
		totalMB := float64(totalSize) / 1024 / 1024
		if totalMB < 1 {
			return
		}

		gui.SetProgress(1)
		defer gui.SetProgress(0)

		// set sleep 100ms for every mb
		sleep := time.Duration(totalMB) * 100 * time.Millisecond

		// every 100ms set progress 10
		for i := 0; i < 10; i++ {
			time.Sleep(sleep)
			select {
			case <-ctx.Done():
				return
			default:
			}
			gui.SetProgress(10 * (i + 1))
		}
	}()

	return archive.Save(path)
}

func (c *Client) onArchiveRefresh() {
	err := c.open(c.openPath, "", false)
	if err != nil {
		slog.Print("Failed to refresh: %s", err)
	}
}

func (c *Client) onArchiveExportAll(path string) error {
	return archive.ExportAll(path)
}

func (c *Client) onArchiveExportFile(path string, file string) error {
	return archive.ExportFile(path, file)
}
