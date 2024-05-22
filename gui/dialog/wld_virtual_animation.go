package dialog

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/wld/virtual"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualAnimationPage(data *virtual.Wld, page *cpl.TabPage) error {

	animations := []string{}
	for _, animation := range data.Animations {
		animations = append(animations, animation.Tag)
	}
	onAnimationNew := func() {
		slog.Println("new animation")
	}
	onAnimationEdit := func() {
		slog.Println("edit animation")
	}
	onAnimationDelete := func() {
		slog.Println("delete animation")
	}

	var cmbAnimation *walk.ComboBox
	defaultAnimation := ""
	if len(animations) > 0 {
		defaultAnimation = animations[0]
	}

	animationGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	animationGroup.Children = append(animationGroup.Children, cpl.GroupBox{
		Title:  "Animations (TrackDef)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbAnimation,
				Editable: false,
				Model:    animations,
				Value:    defaultAnimation,
			},
			cpl.PushButton{Text: "Add", OnClicked: onAnimationNew},
			cpl.PushButton{Text: "Edit", OnClicked: onAnimationEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onAnimationDelete},
		},
	})

	animationInstances := []string{}
	for _, animationInstance := range data.AnimationInstances {
		animationInstances = append(animationInstances, animationInstance.Tag)
	}
	onAnimationInstanceNew := func() {
		slog.Println("new animationInstance")
	}
	onAnimationInstanceEdit := func() {
		slog.Println("edit animationInstance")
	}
	onAnimationInstanceDelete := func() {
		slog.Println("delete animationInstance")
	}

	var cmbAnimationInstance *walk.ComboBox
	defaultAnimationInstance := ""
	if len(animationInstances) > 0 {
		defaultAnimationInstance = animationInstances[0]
	}

	animationGroup.Children = append(animationGroup.Children, cpl.GroupBox{
		Title:  "AnimationInstances (Track)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbAnimationInstance,
				Editable: false,
				Model:    animationInstances,
				Value:    defaultAnimationInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onAnimationInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onAnimationInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onAnimationInstanceDelete},
		},
	})

	page.Title = "Animation"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{animationGroup}
	return nil
}
