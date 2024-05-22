package dialog

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/wld/virtual"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualActorPage(data *virtual.Wld, page *cpl.TabPage) error {

	actors := []string{}
	for _, actor := range data.Actors {
		actors = append(actors, actor.Tag)
	}
	onActorNew := func() {
		slog.Println("new actor")
	}
	onActorEdit := func() {
		slog.Println("edit actor")
	}
	onActorDelete := func() {
		slog.Println("delete actor")
	}

	var cmbActor *walk.ComboBox
	defaultActor := ""
	if len(actors) > 0 {
		defaultActor = actors[0]
	}

	actorGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	actorGroup.Children = append(actorGroup.Children, cpl.GroupBox{
		Title:  "Actors (ActorDef)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbActor,
				Editable: false,
				Model:    actors,
				Value:    defaultActor,
			},
			cpl.PushButton{Text: "Add", OnClicked: onActorNew},
			cpl.PushButton{Text: "Edit", OnClicked: onActorEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onActorDelete},
		},
	})

	actorInstances := []string{}
	for _, actorInstance := range data.ActorInstances {
		actorInstances = append(actorInstances, actorInstance.Tag)
	}
	onActorInstanceNew := func() {
		slog.Println("new actorInstance")
	}
	onActorInstanceEdit := func() {
		slog.Println("edit actorInstance")
	}
	onActorInstanceDelete := func() {
		slog.Println("delete actorInstance")
	}

	var cmbActorInstance *walk.ComboBox
	defaultActorInstance := ""
	if len(actorInstances) > 0 {
		defaultActorInstance = actorInstances[0]
	}

	actorGroup.Children = append(actorGroup.Children, cpl.GroupBox{
		Title:  "ActorInstances (Actor)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbActorInstance,
				Editable: false,
				Model:    actorInstances,
				Value:    defaultActorInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onActorInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onActorInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onActorInstanceDelete},
		},
	})

	cameras := []string{}
	for _, camera := range data.Cameras {
		cameras = append(cameras, camera.Tag)
	}
	onCameraNew := func() {
		slog.Println("new camera")
	}
	onCameraEdit := func() {
		slog.Println("edit camera")
	}
	onCameraDelete := func() {
		slog.Println("delete camera")
	}

	var cmbCamera *walk.ComboBox
	defaultCamera := ""
	if len(cameras) > 0 {
		defaultCamera = cameras[0]
	}

	cameraGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	cameraGroup.Children = append(cameraGroup.Children, cpl.GroupBox{
		Title:  "Cameras (Sprite3DDef)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbCamera,
				Editable: false,
				Model:    cameras,
				Value:    defaultCamera,
			},
			cpl.PushButton{Text: "Add", OnClicked: onCameraNew},
			cpl.PushButton{Text: "Edit", OnClicked: onCameraEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onCameraDelete},
		},
	})

	cameraInstances := []string{}
	for _, cameraInstance := range data.CameraInstances {
		cameraInstances = append(cameraInstances, cameraInstance.Tag)
	}
	onCameraInstanceNew := func() {
		slog.Println("new cameraInstance")
	}
	onCameraInstanceEdit := func() {
		slog.Println("edit cameraInstance")
	}
	onCameraInstanceDelete := func() {
		slog.Println("delete cameraInstance")
	}

	var cmbCameraInstance *walk.ComboBox
	defaultCameraInstance := ""
	if len(cameraInstances) > 0 {
		defaultCameraInstance = cameraInstances[0]
	}

	cameraGroup.Children = append(cameraGroup.Children, cpl.GroupBox{
		Title:  "CameraInstances (Sprite3D)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbCameraInstance,
				Editable: false,
				Model:    cameraInstances,
				Value:    defaultCameraInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onCameraInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onCameraInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onCameraInstanceDelete},
		},
	})

	page.Title = "Actor"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{actorGroup, cameraGroup}
	return nil
}
