package dialog

import (
	"fmt"
	"strconv"

	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/popup"
	"github.com/xackery/quail/raw"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func ShowMdsEdit(mw *walk.MainWindow, title string, src raw.ReadWriter) error {

	data, ok := src.(*raw.Mds)
	if !ok {
		return fmt.Errorf("cast mds")
	}

	var savePB, cancelPB *walk.PushButton
	formElements := cpl.Composite{
		Layout:   cpl.VBox{},
		Children: []cpl.Widget{},
	}

	var cmbVersion *walk.ComboBox
	versions := []string{"1", "2", "3"}
	formElements.Children = append(formElements.Children, cpl.GroupBox{
		Title:  "Header",
		Layout: cpl.Grid{Columns: 2},
		Children: []cpl.Widget{
			cpl.Label{Text: "Version:"},
			cpl.ComboBox{
				AssignTo: &cmbVersion,
				Editable: false,
				Value:    fmt.Sprintf("%d", data.Version),
				Model:    versions,
			},
		},
	})

	materials := []string{}
	for _, material := range data.Materials {
		materials = append(materials, material.Name)
	}
	formElements.Children = append(formElements.Children, cpl.GroupBox{
		Title:  "Materials",
		Layout: cpl.Grid{Columns: 3},
		Children: []cpl.Widget{
			cpl.ListBox{
				Model: materials,
				OnItemActivated: func() {
					fmt.Println("item activated: ", data.Materials)
				},
			},
			cpl.PushButton{
				Image: ico.Grab("new"),
			},
			cpl.PushButton{
				Image: ico.Grab("edit"),
			},
		},
	})

	defaultBone := ""
	if len(data.Bones) > 0 {
		defaultBone = data.Bones[0].Name
	}
	bones := []string{}
	for _, bone := range data.Bones {
		bones = append(bones, bone.Name)
	}
	formElements.Children = append(formElements.Children, cpl.GroupBox{
		Title:  "Bones",
		Layout: cpl.Grid{Columns: 3},
		Children: []cpl.Widget{
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
		},
	})

	defaultVertex := ""
	if len(data.Vertices) > 0 {
		defaultVertex = "1"
	}
	vertices := []string{}
	for i := 0; i < len(data.Vertices); i++ {
		vertices = append(vertices, fmt.Sprintf("%d", i+1))
	}
	formElements.Children = append(formElements.Children, cpl.GroupBox{
		Title:  "Vertices",
		Layout: cpl.Grid{Columns: 3},
		Children: []cpl.Widget{
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
		},
	})

	defaultTriangle := ""
	if len(data.Triangles) > 0 {
		defaultTriangle = "1"
	}
	triangles := []string{}
	for i := 0; i < len(data.Triangles); i++ {
		triangles = append(triangles, fmt.Sprintf("%d", i+1))
	}

	formElements.Children = append(formElements.Children, cpl.GroupBox{
		Title:  "Triangles",
		Layout: cpl.Grid{Columns: 3},
		Children: []cpl.Widget{
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
		},
	})

	onSave := func() error {
		newVersionStr := cmbVersion.Text()
		if newVersionStr == "" {
			return fmt.Errorf("version is required")
		}
		newVersion, err := strconv.Atoi(newVersionStr)
		if err != nil {
			return fmt.Errorf("parse version: %w", err)
		}
		if data.Version != uint32(newVersion) {
			data.Version = uint32(newVersion)
		}
		fmt.Println("new version:", data.Version)
		return nil
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
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
					cpl.PushButton{
						AssignTo: &savePB,
						Text:     "Save",
						OnClicked: func() {
							err := onSave()
							if err != nil {
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
