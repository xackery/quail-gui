package dialog

import (
	"fmt"
	"strconv"

	"github.com/xackery/quail-gui/popup"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/vwld"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func ShowWldVirtualEdit(mw *walk.MainWindow, title string, src raw.ReadWriter) error {
	var savePB, cancelPB *walk.PushButton
	formElements := cpl.Composite{
		Layout:   cpl.VBox{},
		Children: []cpl.Widget{},
	}

	rawData, ok := src.(*raw.Wld)
	if !ok {
		return fmt.Errorf("cast wld")
	}

	data := &vwld.VWld{}
	err := data.Read(rawData)
	if err != nil {
		return fmt.Errorf("read wld: %w", err)
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

	bitmaps := []string{}
	for _, bitmap := range data.Bitmaps {
		bitmaps = append(bitmaps, bitmap.Name)
	}
	onBitmapNew := func() {
		slog.Println("new bitmap")
	}
	onBitmapEdit := func() {
		slog.Println("edit bitmap")
	}
	onBitmapDelete := func() {
		slog.Println("delete bitmap")
	}

	var cmbBitmap *walk.ComboBox
	defaultBitmap := ""
	if len(bitmaps) > 0 {
		defaultBitmap = bitmaps[0]
	}

	formElements.Children = append(formElements.Children, cpl.GroupBox{
		Title:  "Bitmaps",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbBitmap,
				Editable: false,
				Model:    bitmaps,
				Value:    defaultBitmap,
			},
			cpl.PushButton{Text: "Add", OnClicked: onBitmapNew},
			cpl.PushButton{Text: "Edit", OnClicked: onBitmapEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onBitmapDelete},
		},
	})

	sprites := []string{}
	for _, sprite := range data.Sprites {
		sprites = append(sprites, sprite.Name)
	}
	onSpriteNew := func() {
		slog.Println("new sprite")
	}
	onSpriteEdit := func() {
		slog.Println("edit sprite")
	}
	onSpriteDelete := func() {
		slog.Println("delete sprite")
	}

	var cmbSprite *walk.ComboBox
	defaultSprite := ""
	if len(sprites) > 0 {
		defaultSprite = sprites[0]
	}

	spriteGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	spriteGroup.Children = append(spriteGroup.Children, cpl.GroupBox{
		Title:  "Sprites",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbSprite,
				Editable: false,
				Model:    sprites,
				Value:    defaultSprite,
			},
			cpl.PushButton{Text: "Add", OnClicked: onSpriteNew},
			cpl.PushButton{Text: "Edit", OnClicked: onSpriteEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onSpriteDelete},
		},
	})

	spriteInstances := []string{}
	for _, spriteInstance := range data.SpriteInstances {
		spriteInstances = append(spriteInstances, spriteInstance.Name)
	}
	onSpriteInstanceNew := func() {
		slog.Println("new spriteInstance")
	}
	onSpriteInstanceEdit := func() {
		slog.Println("edit spriteInstance")
	}
	onSpriteInstanceDelete := func() {
		slog.Println("delete spriteInstance")
	}

	var cmbSpriteInstance *walk.ComboBox
	defaultSpriteInstance := ""
	if len(spriteInstances) > 0 {
		defaultSpriteInstance = spriteInstances[0]
	}

	spriteGroup.Children = append(spriteGroup.Children, cpl.GroupBox{
		Title:  "SpriteInstances",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbSpriteInstance,
				Editable: false,
				Model:    spriteInstances,
				Value:    defaultSpriteInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onSpriteInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onSpriteInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onSpriteInstanceDelete},
		},
	})
	formElements.Children = append(formElements.Children, spriteGroup)

	particles := []string{}
	for _, particle := range data.Particles {
		particles = append(particles, particle.Name)
	}
	onParticleNew := func() {
		slog.Println("new particle")
	}
	onParticleEdit := func() {
		slog.Println("edit particle")
	}
	onParticleDelete := func() {
		slog.Println("delete particle")
	}

	var cmbParticle *walk.ComboBox
	defaultParticle := ""
	if len(particles) > 0 {
		defaultParticle = particles[0]
	}

	particleGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	particleGroup.Children = append(particleGroup.Children, cpl.GroupBox{
		Title:  "Particles",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbParticle,
				Editable: false,
				Model:    particles,
				Value:    defaultParticle,
			},
			cpl.PushButton{Text: "Add", OnClicked: onParticleNew},
			cpl.PushButton{Text: "Edit", OnClicked: onParticleEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onParticleDelete},
		},
	})

	particleInstances := []string{}
	for _, particleInstance := range data.ParticleInstances {
		particleInstances = append(particleInstances, particleInstance.Name)
	}
	onParticleInstanceNew := func() {
		slog.Println("new particleInstance")
	}
	onParticleInstanceEdit := func() {
		slog.Println("edit particleInstance")
	}
	onParticleInstanceDelete := func() {
		slog.Println("delete particleInstance")
	}

	var cmbParticleInstance *walk.ComboBox
	defaultParticleInstance := ""
	if len(particleInstances) > 0 {
		defaultParticleInstance = particleInstances[0]
	}

	particleGroup.Children = append(particleGroup.Children, cpl.GroupBox{
		Title:  "ParticleInstances",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbParticleInstance,
				Editable: false,
				Model:    particleInstances,
				Value:    defaultParticleInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onParticleInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onParticleInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onParticleInstanceDelete},
		},
	})
	formElements.Children = append(formElements.Children, particleGroup)

	materials := []string{}
	for _, material := range data.Materials {
		materials = append(materials, material.Name)
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
		Title:  "Materials",
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
		materialInstances = append(materialInstances, materialInstance.Name)
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
		Title:  "MaterialInstances",
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
	formElements.Children = append(formElements.Children, materialGroup)

	meshes := []string{}
	for _, mesh := range data.Meshes {
		meshes = append(meshes, mesh.Name)
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
		Title:  "Meshes",
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

	meshInstances := []string{}
	for _, meshInstance := range data.MeshInstances {
		meshInstances = append(meshInstances, meshInstance.Name)
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
		Title:  "MeshInstances",
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
	formElements.Children = append(formElements.Children, meshGroup)

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
		slog.Println("new version:", data.Version)
		return nil
	}

	var dlg *walk.Dialog
	dia := cpl.Dialog{
		AssignTo:      &dlg,
		Title:         fmt.Sprintf("%s (Virtual)", data.FileName),
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
