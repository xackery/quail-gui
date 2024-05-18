package dialog

import (
	"fmt"

	"github.com/xackery/quail-gui/op"
	"github.com/xackery/quail/raw"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func ShowWldEdit(mw *walk.MainWindow, node *op.Node) error {
	var savePB, cancelPB *walk.PushButton
	zd, ok := node.Value().(*raw.Wld)
	if !ok {
		return fmt.Errorf("failed to cast wld")
	}

	defaultFragment := ""
	if len(zd.Fragments) > 0 {
		defaultFragment = fmt.Sprintf("%d: %s", 1, raw.FragName(zd.Fragments[1].FragCode()))
	}
	fragments := []string{}
	for i := 1; i < len(zd.Fragments); i++ {
		fragment, ok := zd.Fragments[i]
		if !ok {
			continue
		}
		fragments = append(fragments, fmt.Sprintf("%d: %s", i+1, raw.FragName(fragment.FragCode())))
	}

	formElements := cpl.Composite{
		Layout: cpl.VBox{},
		Children: []cpl.Widget{
			cpl.GenerateComposite(cpl.Grid{Columns: 2},
				cpl.Label{Text: "Version:"},
				cpl.ComboBox{
					Editable: false,
					Value:    fmt.Sprintf("%d", zd.Version),
					Model:    []string{"1", "2", "3"},
				},
			),
			cpl.Label{Text: "Fragments:"},
			cpl.GenerateComposite(cpl.Grid{Columns: 3},
				cpl.ComboBox{
					Editable: false,
					Value:    defaultFragment,
					Model:    fragments,
				},
				cpl.PushButton{
					Text: "Add",
				},
				cpl.PushButton{
					Text: "Edit",
				},
			),
		},
	}

	var dlg *walk.Dialog
	dia := cpl.Dialog{
		AssignTo:      &dlg,
		Title:         zd.FileName(),
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
						AssignTo:  &savePB,
						Text:      "Save",
						OnClicked: func() { dlg.Accept() },
					},
					cpl.PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
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
		return nil
	}

	node.SetIsEdited(true)

	return nil
}
