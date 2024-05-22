package dialog

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/wld/virtual"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualBspPage(data *virtual.Wld, page *cpl.TabPage) error {

	bsps := []string{}
	for _, bsp := range data.BspTrees {
		bsps = append(bsps, bsp.Tag)
	}
	onBspNew := func() {
		slog.Println("new bsp")
	}
	onBspEdit := func() {
		slog.Println("edit bsp")
	}
	onBspDelete := func() {
		slog.Println("delete bsp")
	}

	var cmbBsp *walk.ComboBox
	defaultBsp := ""
	if len(bsps) > 0 {
		defaultBsp = bsps[0]
	}

	bspGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	bspGroup.Children = append(bspGroup.Children, cpl.GroupBox{
		Title:  "BspTrees (WorldTree)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbBsp,
				Editable: false,
				Model:    bsps,
				Value:    defaultBsp,
			},
			cpl.PushButton{Text: "Add", OnClicked: onBspNew},
			cpl.PushButton{Text: "Edit", OnClicked: onBspEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onBspDelete},
		},
	})

	page.Title = "Bsp"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{bspGroup}
	return nil
}
