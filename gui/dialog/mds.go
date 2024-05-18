package dialog

import (
	"fmt"

	"github.com/xackery/quail-gui/op"
	"github.com/xackery/quail/raw"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func ShowMdsEdit(mw *walk.MainWindow, node *op.Node) error {
	var savePB, cancelPB *walk.PushButton
	zd, ok := node.Value().(*raw.Mds)
	if !ok {
		return fmt.Errorf("failed to cast mds")
	}

	materials := []string{}
	for _, material := range zd.Materials {
		materials = append(materials, material.Name)
	}

	defaultBone := ""
	if len(zd.Bones) > 0 {
		defaultBone = zd.Bones[0].Name
	}
	bones := []string{}
	for _, bone := range zd.Bones {
		bones = append(bones, bone.Name)
	}

	defaultVertex := ""
	if len(zd.Vertices) > 0 {
		defaultVertex = "1"
	}
	vertices := []string{}
	for i := 0; i < len(zd.Vertices); i++ {
		vertices = append(vertices, fmt.Sprintf("%d", i+1))
	}

	defaultTriangle := ""
	if len(zd.Triangles) > 0 {
		defaultTriangle = "1"
	}
	triangles := []string{}
	for i := 0; i < len(zd.Triangles); i++ {
		triangles = append(triangles, fmt.Sprintf("%d", i+1))
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
			cpl.Label{Text: "Materials:"},
			cpl.ListBox{
				Model: materials,
				OnItemActivated: func() {
					fmt.Println("item activated: ", zd.Materials)
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
