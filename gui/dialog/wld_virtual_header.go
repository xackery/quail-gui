package dialog

import (
	"fmt"

	"github.com/xackery/quail/vwld"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualHeaderPage(data *vwld.VWld, page *cpl.TabPage) error {

	headerGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	var cmbVersion *walk.ComboBox
	versions := []string{"1", "2", "3"}

	var cmbGlobalAmbientLight *walk.ComboBox
	globalAmbientLights := []string{"TODO"}

	headerGroup.Children = append(headerGroup.Children, cpl.GroupBox{
		Title:  "Header",
		Layout: cpl.Grid{Columns: 2},
		Children: []cpl.Widget{
			cpl.Label{Text: "Version"},
			cpl.ComboBox{
				AssignTo: &cmbVersion,
				Editable: false,
				Value:    fmt.Sprintf("%d", data.Version),
				Model:    versions,
			},
			cpl.Label{Text: "Global Ambient Light"},
			cpl.ComboBox{
				AssignTo: &cmbGlobalAmbientLight,
				Editable: false,
				Value:    "TODO", //data.GlobalAmbientLight,
				Model:    globalAmbientLights,
			},
		},
	})

	page.Title = "Header"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{headerGroup}
	return nil
}
