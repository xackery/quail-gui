package dialog

import (
	"fmt"

	"github.com/xackery/quail-gui/config"
	"github.com/xackery/quail-gui/popup"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func ShowPreferences(mw *walk.MainWindow) error {
	var savePB, cancelPB *walk.PushButton

	formElements := cpl.Composite{
		Layout:   cpl.VBox{},
		Children: []cpl.Widget{},
	}

	cfg := config.Instance()

	var chkIsVirtualWld *walk.CheckBox
	formElements.Children = append(formElements.Children, cpl.GroupBox{
		Title:  "WLD Options",
		Layout: cpl.Grid{Columns: 2},
		Children: []cpl.Widget{
			cpl.Label{Text: "Load Virtual Wld:", ToolTipText: "If checked, the WLD will be loaded into virtual properties"},
			cpl.CheckBox{
				AssignTo:    &chkIsVirtualWld,
				Checked:     cfg.IsVirtualWld,
				ToolTipText: "If checked, the WLD will be loaded into virtual properties",
			},
		},
	})

	var dlg *walk.Dialog
	onSave := func() {
		cfg := config.Instance()
		cfg.IsVirtualWld = chkIsVirtualWld.Checked()
		err := cfg.Save()
		if err != nil {
			popup.Errorf(mw, "save config: %s", err)
			return
		}

		slog.Printf("Saved preferences\n")

		dlg.Accept()
	}

	dia := cpl.Dialog{
		AssignTo:      &dlg,
		Title:         "Preferences",
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
						OnClicked: onSave,
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
		return fmt.Errorf("cancelled")
	}
	return nil
}
