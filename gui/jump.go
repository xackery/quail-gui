package gui

import (
	"path/filepath"
	"strings"

	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/walk"
)

var (
	menuJumpToWorld     *walk.Action
	menuJumpToLight     *walk.Action
	menuJumpToObject    *walk.Action
	toolbarJumpToWorld  *walk.Action
	toolbarJumpToLight  *walk.Action
	toolbarJumpToObject *walk.Action
)

func setJumpWorldEnabled(enabled bool) {
	menuJumpToWorld.SetEnabled(enabled)
	toolbarJumpToWorld.SetEnabled(enabled)
}

func setJumpLightEnabled(enabled bool) {
	menuJumpToLight.SetEnabled(enabled)
	toolbarJumpToLight.SetEnabled(enabled)
}

func setJumpObjectEnabled(enabled bool) {
	menuJumpToObject.SetEnabled(enabled)
	toolbarJumpToObject.SetEnabled(enabled)
}

func onMenuJumpToWorld() {
	wldName := filepath.Base(archivePath)
	// replace ext with .wld
	wldName = strings.ReplaceAll(wldName, filepath.Ext(wldName), ".wld")
	idx, item := fileView.ItemByName(wldName)
	if item != nil {
		file.SetCurrentIndex(idx)
		EntryEdit()
		return
	}

	wldName = strings.ReplaceAll(wldName, ".wld", ".zon")
	idx, item = fileView.ItemByName(wldName)
	if item != nil {
		file.SetCurrentIndex(idx)
		EntryEdit()
		return
	}

	idx, item = fileView.ItemByExt(".wld")
	if item != nil {
		file.SetCurrentIndex(idx)
		EntryEdit()
		return
	}

	idx, item = fileView.ItemByExt(".zon")
	if item != nil {
		file.SetCurrentIndex(idx)
		EntryEdit()
		return
	}

	slog.Printf("No world file found\n")
	setJumpWorldEnabled(false)
}

func onMenuJumpToLight() {
	idx, item := fileView.ItemByName("lights.wld")
	if item != nil {
		file.SetCurrentIndex(idx)
		EntryEdit()
		return
	}

	slog.Printf("No light file found\n")
	setJumpLightEnabled(false)
}

func onMenuJumpToObject() {
	idx, item := fileView.ItemByName("objects.wld")
	if item != nil {
		file.SetCurrentIndex(idx)
		EntryEdit()
		return
	}

	slog.Printf("No object file found\n")
	setJumpObjectEnabled(false)
}
