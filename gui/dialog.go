package gui

import (
	"fmt"
	"strings"

	"github.com/xackery/wlk/walk"
)

func ShowError(err error) {
	sections := strings.Split(err.Error(), ": ")
	if len(sections) < 2 {
		ShowMessageBox("Error", err.Error(), true)
		return
	}

	for i := 1; i < len(sections); i++ {
		sections[i] = strings.ToUpper(sections[i][0:1]) + sections[i][1:]
	}

	ShowMessageBox("Failed to "+sections[0], strings.Join(sections[1:], "\n"), true)
}

func ShowOpen(title string, filter string, initialDirPath string) (string, error) {
	if mw == nil {
		return "", fmt.Errorf("gui not initialized")
	}
	dialog := walk.FileDialog{
		Title:          title,
		Filter:         filter,
		InitialDirPath: initialDirPath,
	}
	ok, err := dialog.ShowOpen(mw)
	if err != nil {
		return "", fmt.Errorf("show open: %w", err)
	}
	if !ok {
		return "", fmt.Errorf("show open: cancelled")
	}
	return dialog.FilePath, nil
}

func ShowMessageBox(title string, message string, isError bool) {
	if mw == nil {
		return
	}
	// convert style to msgboxstyle
	icon := walk.MsgBoxIconInformation
	if isError {
		icon = walk.MsgBoxIconError
	}
	mw.SetEnabled(false)
	walk.MsgBox(mw, title, message, icon)
	mw.SetEnabled(true)
}

func ShowMessageBoxYesNo(title string, message string) bool {
	if mw == nil {
		return false
	}
	// convert style to msgboxstyle
	icon := walk.MsgBoxIconInformation
	result := walk.MsgBox(mw, title, message, icon|walk.MsgBoxYesNo)
	return result == walk.DlgCmdYes
}

func ShowMessageBoxf(title string, format string, a ...interface{}) {
	if mw == nil {
		return
	}
	// convert style to msgboxstyle
	icon := walk.MsgBoxIconInformation
	walk.MsgBox(mw, title, fmt.Sprintf(format, a...), icon)
}
