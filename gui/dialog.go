package gui

import (
	"fmt"

	"github.com/xackery/wlk/walk"
)

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
	walk.MsgBox(mw, title, message, icon)
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
