package dialog

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/wld/virtual"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualRegionPage(data *virtual.Wld, page *cpl.TabPage) error {

	regions := []string{}
	for _, region := range data.Regions {
		regions = append(regions, region.Tag)
	}
	onRegionNew := func() {
		slog.Println("new region")
	}
	onRegionEdit := func() {
		slog.Println("edit region")
	}
	onRegionDelete := func() {
		slog.Println("delete region")
	}

	var cmbRegion *walk.ComboBox
	defaultRegion := ""
	if len(regions) > 0 {
		defaultRegion = regions[0]
	}

	regionGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	regionGroup.Children = append(regionGroup.Children, cpl.GroupBox{
		Title:  "Regions (Region)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbRegion,
				Editable: false,
				Model:    regions,
				Value:    defaultRegion,
			},
			cpl.PushButton{Text: "Add", OnClicked: onRegionNew},
			cpl.PushButton{Text: "Edit", OnClicked: onRegionEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onRegionDelete},
		},
	})

	regionInstances := []string{}
	for _, regionInstance := range data.RegionInstances {
		regionInstances = append(regionInstances, regionInstance.Tag)
	}
	onRegionInstanceNew := func() {
		slog.Println("new regionInstance")
	}
	onRegionInstanceEdit := func() {
		slog.Println("edit regionInstance")
	}
	onRegionInstanceDelete := func() {
		slog.Println("delete regionInstance")
	}

	var cmbRegionInstance *walk.ComboBox
	defaultRegionInstance := ""
	if len(regionInstances) > 0 {
		defaultRegionInstance = regionInstances[0]
	}

	regionGroup.Children = append(regionGroup.Children, cpl.GroupBox{
		Title:  "RegionInstances (Zone)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbRegionInstance,
				Editable: false,
				Model:    regionInstances,
				Value:    defaultRegionInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onRegionInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onRegionInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onRegionInstanceDelete},
		},
	})

	page.Title = "Region"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{regionGroup}
	return nil
}
