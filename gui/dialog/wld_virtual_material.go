package dialog

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/vwld"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualMaterialPage(data *vwld.VWld, page *cpl.TabPage) error {

	materials := []string{}
	for _, material := range data.Materials {
		materials = append(materials, material.Tag)
	}
	onMaterialNew := func() {
		slog.Println("new material")
	}
	onMaterialEdit := func() {
		slog.Println("edit material")
	}
	onMaterialDelete := func() {
		slog.Println("delete material")
	}

	var cmbMaterial *walk.ComboBox
	defaultMaterial := ""
	if len(materials) > 0 {
		defaultMaterial = materials[0]
	}

	materialGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	materialGroup.Children = append(materialGroup.Children, cpl.GroupBox{
		Title:  "Materials (MaterialDef)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbMaterial,
				Editable: false,
				Model:    materials,
				Value:    defaultMaterial,
			},
			cpl.PushButton{Text: "Add", OnClicked: onMaterialNew},
			cpl.PushButton{Text: "Edit", OnClicked: onMaterialEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onMaterialDelete},
		},
	})

	materialInstances := []string{}
	for _, materialInstance := range data.MaterialInstances {
		materialInstances = append(materialInstances, materialInstance.Tag)
	}
	onMaterialInstanceNew := func() {
		slog.Println("new materialInstance")
	}
	onMaterialInstanceEdit := func() {
		slog.Println("edit materialInstance")
	}
	onMaterialInstanceDelete := func() {
		slog.Println("delete materialInstance")
	}

	var cmbMaterialInstance *walk.ComboBox
	defaultMaterialInstance := ""
	if len(materialInstances) > 0 {
		defaultMaterialInstance = materialInstances[0]
	}

	materialGroup.Children = append(materialGroup.Children, cpl.GroupBox{
		Title:  "MaterialInstances (MaterialPalette)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbMaterialInstance,
				Editable: false,
				Model:    materialInstances,
				Value:    defaultMaterialInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onMaterialInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onMaterialInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onMaterialInstanceDelete},
		},
	})

	page.Title = "Material"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{materialGroup}
	return nil
}
