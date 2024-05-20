package popup

import (
	"fmt"
	"strings"

	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func Errorf(wnd walk.Form, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	slog.Printf("%s\n", msg)
	sections := strings.Split(msg, ": ")
	if len(sections) < 2 {
		MessageBox(wnd, "Error", msg, true)
		return
	}

	for i := 1; i < len(sections); i++ {
		sections[i] = strings.ToUpper(sections[i][0:1]) + sections[i][1:]
	}

	MessageBox(wnd, "Failed to "+sections[0], strings.Join(sections[1:], "\n"), true)
}

func Open(wnd walk.Form, title string, filter string, initialDirPath string) (string, error) {
	if wnd == nil {
		return "", fmt.Errorf("gui not initialized")
	}
	dialog := walk.FileDialog{
		Title:          title,
		Filter:         filter,
		InitialDirPath: initialDirPath,
	}
	ok, err := dialog.ShowOpen(wnd)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("cancelled")
	}
	return dialog.FilePath, nil
}

func MessageBox(wnd walk.Form, title string, message string, isError bool) {

	icon := walk.MsgBoxIconInformation
	if isError {
		icon = walk.MsgBoxIconError
	}
	if wnd != nil {

		wnd.SetEnabled(false)
		defer wnd.SetEnabled(true)
	}
	walk.MsgBox(wnd, title, message, icon)
}

func MessageBoxYesNo(wnd walk.Form, title string, message string) bool {
	icon := walk.MsgBoxIconInformation
	result := walk.MsgBox(wnd, title, message, icon|walk.MsgBoxYesNo)
	return result == walk.DlgCmdYes
}

func MessageBoxf(wnd walk.Form, title string, format string, a ...interface{}) {
	icon := walk.MsgBoxIconInformation
	walk.MsgBox(wnd, title, fmt.Sprintf(format, a...), icon)
}

func InputBox(wnd walk.Form, title string, description string, label string, value string) (string, error) {
	if wnd == nil {
		return "", fmt.Errorf("gui not initialized")
	}
	var okPB, cancelPB *walk.PushButton

	var descriptionLabel *walk.Label
	var inputLabel *walk.Label
	var inputValueEdit *walk.LineEdit

	finalValue := ""
	var dlg *walk.Dialog
	dia := cpl.Dialog{
		Title:         title,
		AssignTo:      &dlg,
		DefaultButton: &okPB,
		CancelButton:  &cancelPB,
		MinSize:       cpl.Size{Width: 300, Height: 100},
		Layout:        cpl.VBox{},
		Children: []cpl.Widget{
			cpl.Composite{
				Layout: cpl.VBox{},
				Children: []cpl.Widget{
					cpl.Label{Text: description, AssignTo: &descriptionLabel},
				},
			},
			cpl.Composite{
				Layout: cpl.HBox{},
				Children: []cpl.Widget{
					cpl.Label{Text: label, AssignTo: &inputLabel},
					cpl.LineEdit{Text: value, AssignTo: &inputValueEdit},
				},
			},
			cpl.Composite{
				Layout: cpl.HBox{},
				Children: []cpl.Widget{
					cpl.HSpacer{},
					cpl.PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
					cpl.PushButton{
						AssignTo: &okPB,
						Text:     "OK",
						OnClicked: func() {
							if len(inputValueEdit.Text()) == 0 {
								descriptionLabel.SetText(fmt.Sprintf("%s\n%s", description, "Input is required"))
								return
							}
							finalValue = inputValueEdit.Text()
							dlg.Accept()
						},
					},
				},
			},
		},
	}
	result, err := dia.Run(wnd)
	if err != nil {
		return "", fmt.Errorf("run dialog: %w", err)
	}

	if result != walk.DlgCmdOK {
		return "", fmt.Errorf("cancelled")
	}
	return finalValue, nil
}
