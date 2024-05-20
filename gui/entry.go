package gui

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/popup"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/raw"
	"github.com/xackery/wlk/walk"
)

var (
	lastEntry          int
	menuEntryNew       *walk.Action
	menuEntryEdit      *walk.Action
	menuEntryEditWorld *walk.Action // shortcut to get world edit
	menuEntryDelete    *walk.Action
	menuEntryRename    *walk.Action
)

func onMenuEntryNew() {
	var value string
	var err error
	var data []byte
	for {
		value, err = popup.InputBox(mw, "New entry", "Create a new entry in "+filepath.Base(archivePath), "Name", value)
		if err != nil {
			if err.Error() == "cancelled" {
				return
			}
			popup.Errorf(mw, "input box: %s", err)
			return
		}

		isNew := true
		for _, entry := range archive.Files() {
			if strings.EqualFold(entry.Name(), value) {
				popup.Errorf(mw, "file %s already exists", value)
				isNew = false
				break
			}
		}
		if !isNew {
			continue
		}

		ext := filepath.Ext(value)
		if ext == "" {
			popup.Errorf(mw, "entry must have an extension")
			continue
		}

		var rawWriter raw.ReadWriter

		switch ext {
		case ".mod":
			rawWriter = &raw.Mod{MetaFileName: value}
		case ".txt":
			data = []byte(value)
		case ".wld":
			menuEntryEditWorld.SetEnabled(true)
			rawWriter = &raw.Wld{MetaFileName: value}
		case ".zon":
			menuEntryEditWorld.SetEnabled(true)
			rawWriter = &raw.Zon{MetaFileName: value}
		default:
			popup.Errorf(mw, "unsupported extension %s", ext)
			continue
		}

		if rawWriter != nil {
			buf := bytes.NewBuffer([]byte{})
			err = rawWriter.Write(buf)
			if err != nil {
				popup.Errorf(mw, "write mod: %s", err)
				return
			}
			data = buf.Bytes()
		}

		break
	}

	err = archive.SetFile(value, data)
	if err != nil {
		popup.Errorf(mw, "set file %s: %s", value, err)
		return
	}

	ext := strings.ToLower(filepath.Ext(value))
	img, err := ico.Generate(ext, data)
	if err != nil {
		slog.Printf("Failed to generate icon for %s: %s\n", value, err.Error())
		img = ico.Grab("unk")
	}

	fileView.AddItem(&component.FileViewEntry{
		Icon:    img,
		Name:    value,
		Ext:     ext,
		Size:    generateSize(len(data)),
		RawSize: len(data),
	})
}

func onMenuEntryEdit() {
	EntryEdit()
}

func onMenuEntryDelete() {
	slog.Println("entry delete triggered")
}

func onMenuEntryRename() {
	slog.Println("entry rename triggered")
}

func onEntryChange() {
	if lastEntry == file.CurrentIndex() {
		return
	}
	if file.CurrentIndex() < 0 {
		entrySetActive(false)
		return
	}
	entrySetActive(true)
	fileName := fileView.Item(file.CurrentIndex()).Name

	slog.Printf("Selected %s\n", fileName)
}

func onEntryActivate() {
	EntryEdit()
}

func onMenuEntryEditWorld() {

	idx, item := fileView.ItemByExt(".wld")
	if item == nil {
		idx, item = fileView.ItemByExt(".zon")
		if item == nil {
			slog.Printf("No world file found\n")
			menuEntryEditWorld.SetEnabled(false)
			return
		}
	}

	err := file.SetCurrentIndex(idx)
	if err != nil {
		popup.Errorf(mw, "set current index: %s", err)
		return
	}

	EntryEdit()
}

func EntryEdit() {
	var err error

	if file.CurrentIndex() < 0 {
		slog.Println("current index is less than 0")
		return
	}

	if file.CurrentIndex() >= fileView.RowCount() {
		slog.Println("current index is greater than row count")
		return
	}
	item := fileView.Item(file.CurrentIndex())
	if item == nil {
		slog.Println("item is nil")
		return
	}

	if archive == nil {
		slog.Println("archive is nil")
		return
	}

	itemName := strings.ReplaceAll(item.Name, "*", "")

	data, err := archive.File(itemName)
	if err != nil {
		popup.Errorf(mw, "open file %s: %s", itemName, err)
		return
	}

	ext := filepath.Ext(strings.ToLower(itemName))
	value, err := raw.Read(ext, bytes.NewReader(data))
	if err != nil {
		popup.Errorf(mw, "read raw %s: %s", itemName, err)
		return
	}
	value.SetFileName(itemName)

	slog.Printf("Selected file: %s\n", itemName)

	data, err = DialogEdit(itemName, value)
	if err != nil {
		if err.Error() == "cancelled" {
			return
		}

		popup.Errorf(mw, "show %s edit: %s", ext, err)
		return
	}

	err = archive.SetFile(itemName, data)
	if err != nil {
		popup.Errorf(mw, "set file %s: %s", itemName, err)
		return
	}
	slog.Printf("Edited %s\n", itemName)
	item.Size = generateSize(len(data))
	item.RawSize = len(data)
	item.Name = strings.ReplaceAll(item.Name, "*", "")
	item.Name = fmt.Sprintf("%s*", itemName)

	if !isEdited {
		isEdited = true
		err = mw.SetTitle(fmt.Sprintf("%s*", mw.Title()))
		if err != nil {
			popup.Errorf(mw, "set title: %s", err)
			return
		}
	}

}

func entrySetActive(value bool) {
	menuEntryNew.SetEnabled(value)
	menuEntryEdit.SetEnabled(value)
	menuEntryDelete.SetEnabled(value)
	menuEntryRename.SetEnabled(value)
}
