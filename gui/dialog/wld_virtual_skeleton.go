package dialog

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/wld/virtual"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualSkeletonPage(data *virtual.Wld, page *cpl.TabPage) error {

	skeletons := []string{}
	for _, skeleton := range data.Skeletons {
		skeletons = append(skeletons, skeleton.Tag)
	}
	onSkeletonNew := func() {
		slog.Println("new skeleton")
	}
	onSkeletonEdit := func() {
		slog.Println("edit skeleton")
	}
	onSkeletonDelete := func() {
		slog.Println("delete skeleton")
	}

	var cmbSkeleton *walk.ComboBox
	defaultSkeleton := ""
	if len(skeletons) > 0 {
		defaultSkeleton = skeletons[0]
	}

	skeletonGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	skeletonGroup.Children = append(skeletonGroup.Children, cpl.GroupBox{
		Title:  "Skeletons (HierarchialSpriteDef)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbSkeleton,
				Editable: false,
				Model:    skeletons,
				Value:    defaultSkeleton,
			},
			cpl.PushButton{Text: "Add", OnClicked: onSkeletonNew},
			cpl.PushButton{Text: "Edit", OnClicked: onSkeletonEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onSkeletonDelete},
		},
	})

	skeletonInstances := []string{}
	for _, skeletonInstance := range data.SkeletonInstances {
		skeletonInstances = append(skeletonInstances, skeletonInstance.Tag)
	}
	onSkeletonInstanceNew := func() {
		slog.Println("new skeletonInstance")
	}
	onSkeletonInstanceEdit := func() {
		slog.Println("edit skeletonInstance")
	}
	onSkeletonInstanceDelete := func() {
		slog.Println("delete skeletonInstance")
	}

	var cmbSkeletonInstance *walk.ComboBox
	defaultSkeletonInstance := ""
	if len(skeletonInstances) > 0 {
		defaultSkeletonInstance = skeletonInstances[0]
	}

	skeletonGroup.Children = append(skeletonGroup.Children, cpl.GroupBox{
		Title:  "SkeletonInstances (HierarchialSprite)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbSkeletonInstance,
				Editable: false,
				Model:    skeletonInstances,
				Value:    defaultSkeletonInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onSkeletonInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onSkeletonInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onSkeletonInstanceDelete},
		},
	})

	page.Title = "Skeleton"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{skeletonGroup}
	return nil
}
