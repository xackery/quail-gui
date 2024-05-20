package dialog

import (
	"fmt"

	"github.com/xackery/quail/raw"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func ShowModEdit(mw *walk.MainWindow, title string, src raw.ReadWriter) error {
	var savePB, cancelPB *walk.PushButton
	data, ok := src.(*raw.Mod)
	if !ok {
		return fmt.Errorf("cast mod")
	}

	materials := []string{}
	for _, material := range data.Materials {
		materials = append(materials, material.Name)
	}

	defaultBone := ""
	if len(data.Bones) > 0 {
		defaultBone = data.Bones[0].Name
	}
	bones := []string{}
	for _, bone := range data.Bones {
		bones = append(bones, bone.Name)
	}

	defaultVertex := ""
	if len(data.Vertices) > 0 {
		defaultVertex = "1"
	}
	vertices := []string{}
	for i := 0; i < len(data.Vertices); i++ {
		vertices = append(vertices, fmt.Sprintf("%d", i+1))
	}

	defaultTriangle := ""
	if len(data.Triangles) > 0 {
		defaultTriangle = "1"
	}
	triangles := []string{}
	for i := 0; i < len(data.Triangles); i++ {
		triangles = append(triangles, fmt.Sprintf("%d", i+1))
	}

	formElements := cpl.Composite{
		Layout: cpl.VBox{},
		Children: []cpl.Widget{
			cpl.GenerateComposite(cpl.Grid{Columns: 2},
				cpl.Label{Text: "Version:"},
				cpl.ComboBox{
					Editable: false,
					Value:    fmt.Sprintf("%d", data.Version),
					Model:    []string{"1", "2", "3"},
				},
			),
			cpl.Label{Text: "Materials:"},
			cpl.ListBox{
				Model: materials,
				OnItemActivated: func() {
					fmt.Println("item activated: ", data.Materials)
				},
			},
			cpl.Label{Text: "Bones:"},
			cpl.GenerateComposite(cpl.Grid{Columns: 3},
				cpl.ComboBox{
					Editable: false,
					Value:    defaultBone,
					Model:    bones,
				},
				cpl.PushButton{
					Text: "Add",
				},
				cpl.PushButton{
					Text: "Edit",
				},
			),
			cpl.Label{Text: "Vertices:"},
			cpl.GenerateComposite(cpl.Grid{Columns: 3},
				cpl.ComboBox{
					Editable: false,
					Value:    defaultVertex,
					Model:    vertices,
				},
				cpl.PushButton{
					Text: "Add",
				},
				cpl.PushButton{
					Text: "Edit",
				},
			),
			cpl.Label{Text: "Triangles:"},
			cpl.GenerateComposite(cpl.Grid{Columns: 3},
				cpl.ComboBox{
					Editable: false,
					Value:    defaultTriangle,
					Model:    triangles,
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
		Title:         data.FileName(),
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
		return fmt.Errorf("cancelled")
	}

	return nil
}
