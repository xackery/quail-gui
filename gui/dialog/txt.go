package dialog

import (
	"fmt"

	"github.com/xackery/quail-gui/popup"
	"github.com/xackery/quail/raw"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func ShowTxtEdit(mw *walk.MainWindow, title string, src raw.ReadWriter) error {
	data, ok := src.(*raw.Txt)
	if !ok {
		return fmt.Errorf("cast WldAscii")
	}

	var savePB, cancelPB *walk.PushButton
	formElements := cpl.Composite{
		Layout:   cpl.VBox{},
		Children: []cpl.Widget{},
	}

	var textEdit *walk.TextEdit

	var dlg *walk.Dialog
	onSave := func() error {
		isEdited := false

		newData := textEdit.Text()
		if newData != data.Data {
			data.Data = newData
			isEdited = true
		}

		if !isEdited {
			return fmt.Errorf("no changes")
		}
		return nil
	}
	formElements.Children = append(formElements.Children, cpl.GroupBox{
		Title:  "Text",
		Layout: cpl.VBox{},
		Children: []cpl.Widget{
			cpl.TextEdit{
				AssignTo: &textEdit,
				Text:     data.Data,
				OnKeyPress: func(key walk.Key) {
					isSave := false
					if key == walk.KeyS && walk.ModifiersDown() == walk.ModControl {
						isSave = true
					}
					if key == walk.KeyReturn && walk.ModifiersDown() == walk.ModControl {
						isSave = true
					}
					if !isSave {
						return
					}

					err := onSave()
					if err != nil {
						if err.Error() == "no changes" {
							dlg.Cancel()
							return
						}
						popup.Errorf(dlg, "save: %s", err.Error())
						return

					}
					dlg.Accept()
				},
			},
		},
	})

	dia := cpl.Dialog{
		AssignTo:      &dlg,
		Title:         fmt.Sprintf("%s (ASCII)", title),
		DefaultButton: &savePB,
		CancelButton:  &cancelPB,
		MinSize:       cpl.Size{Width: 300, Height: 300},
		Layout:        cpl.VBox{},
		Children: []cpl.Widget{
			formElements,
			cpl.Composite{
				Layout: cpl.HBox{},
				Children: []cpl.Widget{
					cpl.HSpacer{},
					cpl.PushButton{
						AssignTo:  &cancelPB,
						Text:      "&Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
					cpl.PushButton{
						AssignTo: &savePB,
						Text:     "&Save",
						OnClicked: func() {
							err := onSave()
							if err != nil {
								if err.Error() == "no changes" {
									dlg.Cancel()
									return
								}
								popup.Errorf(dlg, "save: %s", err.Error())
								return

							}
							dlg.Accept()
						},
					},
				},
			},
		},
	}
	result, err := dia.Run(mw)
	if err != nil {
		return fmt.Errorf("run dialog: %w", err)
	}

	if result != walk.DlgCmdOK {
		return fmt.Errorf("cancelled")
	}

	return nil
}
