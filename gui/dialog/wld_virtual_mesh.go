package dialog

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/wld/virtual"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualMeshPage(data *virtual.Wld, page *cpl.TabPage) error {

	meshes := []string{}
	for _, mesh := range data.Meshes {
		meshes = append(meshes, mesh.Tag)
	}
	onMeshNew := func() {
		slog.Println("new mesh")
	}
	onMeshEdit := func() {
		slog.Println("edit mesh")
	}
	onMeshDelete := func() {
		slog.Println("delete mesh")
	}

	var cmbMesh *walk.ComboBox
	defaultMesh := ""
	if len(meshes) > 0 {
		defaultMesh = meshes[0]
	}

	meshGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	meshGroup.Children = append(meshGroup.Children, cpl.GroupBox{
		Title:  "Meshes (SimpleMeshDef)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbMesh,
				Editable: false,
				Model:    meshes,
				Value:    defaultMesh,
			},
			cpl.PushButton{Text: "Add", OnClicked: onMeshNew},
			cpl.PushButton{Text: "Edit", OnClicked: onMeshEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onMeshDelete},
		},
	})

	altMeshes := []string{}
	for _, altMesh := range data.AlternateMeshes {
		altMeshes = append(altMeshes, altMesh.Tag)
	}
	onaltMeshNew := func() {
		slog.Println("new altMesh")
	}
	onaltMeshEdit := func() {
		slog.Println("edit altMesh")
	}
	onaltMeshDelete := func() {
		slog.Println("delete altMesh")
	}

	var cmbaltMesh *walk.ComboBox
	defaultaltMesh := ""
	if len(altMeshes) > 0 {
		defaultaltMesh = altMeshes[0]
	}

	meshGroup.Children = append(meshGroup.Children, cpl.GroupBox{
		Title:  "AltMeshes (DMSpriteDef)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbaltMesh,
				Editable: false,
				Model:    altMeshes,
				Value:    defaultaltMesh,
			},
			cpl.PushButton{Text: "Add", OnClicked: onaltMeshNew},
			cpl.PushButton{Text: "Edit", OnClicked: onaltMeshEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onaltMeshDelete},
		},
	})

	meshInstances := []string{}
	for _, meshInstance := range data.MeshInstances {
		meshInstances = append(meshInstances, meshInstance.Tag)
	}
	onMeshInstanceNew := func() {
		slog.Println("new meshInstance")
	}
	onMeshInstanceEdit := func() {
		slog.Println("edit meshInstance")
	}
	onMeshInstanceDelete := func() {
		slog.Println("delete meshInstance")
	}

	var cmbMeshInstance *walk.ComboBox
	defaultMeshInstance := ""
	if len(meshInstances) > 0 {
		defaultMeshInstance = meshInstances[0]
	}

	meshGroup.Children = append(meshGroup.Children, cpl.GroupBox{
		Title:  "MeshInstances (SimpleMesh)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbMeshInstance,
				Editable: false,
				Model:    meshInstances,
				Value:    defaultMeshInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onMeshInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onMeshInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onMeshInstanceDelete},
		},
	})

	page.Title = "Mesh"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{meshGroup}
	return nil
}
