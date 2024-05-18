package dialog

import (
	"fmt"

	"github.com/xackery/quail-gui/op"
	"github.com/xackery/quail/raw"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func ShowZonEdit(mw *walk.MainWindow, node *op.Node) error {
	var savePB, cancelPB *walk.PushButton
	zd, ok := node.Value().(*raw.Zon)
	if !ok {
		return fmt.Errorf("failed to cast zon")
	}

	defaultModel := ""
	if len(zd.Models) > 0 {
		defaultModel = zd.Models[0]
	}
	models := []string{}
	for _, model := range zd.Models {
		models = append(models, model)
	}

	defaultObject := ""
	if len(zd.Objects) > 0 {
		defaultObject = zd.Objects[0].ModelName
	}
	objects := []string{}
	for _, object := range zd.Objects {
		objects = append(objects, object.ModelName)
	}

	defaultRegion := ""
	if len(zd.Regions) > 0 {
		defaultRegion = zd.Regions[0].Name
	}
	regions := []string{}
	for _, region := range zd.Regions {
		regions = append(regions, region.Name)
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
			cpl.Label{Text: "Models:"},
			cpl.ListBox{
				Model: models,
				OnItemActivated: func() {
					fmt.Println("item activated: ", zd.Models)
				},
			},
			cpl.Label{Text: "Models:"},
			cpl.GenerateComposite(cpl.Grid{Columns: 3},
				cpl.ComboBox{
					Editable: false,
					Value:    defaultModel,
					Model:    models,
				},
				cpl.PushButton{
					Text: "Add",
				},
				cpl.PushButton{
					Text: "Edit",
				},
			),
			cpl.Label{Text: "Objects:"},
			cpl.GenerateComposite(cpl.Grid{Columns: 3},
				cpl.ComboBox{
					Editable: false,
					Value:    defaultObject,
					Model:    objects,
				},
				cpl.PushButton{
					Text: "Add",
				},
				cpl.PushButton{
					Text: "Edit",
				},
			),
			cpl.Label{Text: "Regions:"},
			cpl.GenerateComposite(cpl.Grid{Columns: 3},
				cpl.ComboBox{
					Editable: false,
					Value:    defaultRegion,
					Model:    regions,
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
