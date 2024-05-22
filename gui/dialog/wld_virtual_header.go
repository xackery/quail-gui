package dialog

import (
	"fmt"
	"strings"

	"github.com/xackery/quail/wld/virtual"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualHeaderPage(data *virtual.Wld, page *cpl.TabPage) error {

	headerGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	var cmbVersion *walk.ComboBox
	versions := []string{"0x00015500 (OldWorld)", "0x1000C800"}

	defaultValue := ""
	originalValue := fmt.Sprintf("0x%08X", data.Version)
	for _, version := range versions {
		if strings.HasPrefix(version, originalValue) {
			defaultValue = version
			break
		}
	}
	if defaultValue == "" {
		return fmt.Errorf("Unknown version: %s", originalValue)
	}

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
				Value:    defaultValue,
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
