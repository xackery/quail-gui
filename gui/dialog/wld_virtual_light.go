package dialog

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/vwld"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualLightPage(data *vwld.VWld, page *cpl.TabPage) error {

	lights := []string{}
	for _, light := range data.Lights {
		lights = append(lights, light.Tag)
	}
	onLightNew := func() {
		slog.Println("new light")
	}
	onLightEdit := func() {
		slog.Println("edit light")
	}
	onLightDelete := func() {
		slog.Println("delete light")
	}

	var cmbLight *walk.ComboBox
	defaultLight := ""
	if len(lights) > 0 {
		defaultLight = lights[0]
	}

	lightGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	lightGroup.Children = append(lightGroup.Children, cpl.GroupBox{
		Title:  "Lights (LightDef)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbLight,
				Editable: false,
				Model:    lights,
				Value:    defaultLight,
			},
			cpl.PushButton{Text: "Add", OnClicked: onLightNew},
			cpl.PushButton{Text: "Edit", OnClicked: onLightEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onLightDelete},
		},
	})

	ambientLightInstances := []string{}
	for _, ambientLightInstance := range data.AmbientLightInstances {
		ambientLightInstances = append(ambientLightInstances, ambientLightInstance.Tag)
	}
	onAmbientLightInstanceNew := func() {
		slog.Println("new ambientLightInstance")
	}
	onAmbientLightInstanceEdit := func() {
		slog.Println("edit ambientLightInstance")
	}
	onAmbientLightInstanceDelete := func() {
		slog.Println("delete ambientLightInstance")
	}

	var cmbAmbientLightInstance *walk.ComboBox
	defaultAmbientLightInstance := ""
	if len(ambientLightInstances) > 0 {
		defaultAmbientLightInstance = ambientLightInstances[0]
	}

	lightGroup.Children = append(lightGroup.Children, cpl.GroupBox{
		Title:  "AmbientLightInstances (AmbientLight)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbAmbientLightInstance,
				Editable: false,
				Model:    ambientLightInstances,
				Value:    defaultAmbientLightInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onAmbientLightInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onAmbientLightInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onAmbientLightInstanceDelete},
		},
	})

	pointLightInstances := []string{}
	for _, pointLightInstance := range data.PointLightInstances {
		pointLightInstances = append(pointLightInstances, pointLightInstance.Tag)
	}
	onPointLightInstanceNew := func() {
		slog.Println("new pointLightInstance")
	}
	onPointLightInstanceEdit := func() {
		slog.Println("edit pointLightInstance")
	}
	onPointLightInstanceDelete := func() {
		slog.Println("delete pointLightInstance")
	}

	var cmbPointLightInstance *walk.ComboBox
	defaultPointLightInstance := ""
	if len(pointLightInstances) > 0 {
		defaultPointLightInstance = pointLightInstances[0]
	}

	lightGroup.Children = append(lightGroup.Children, cpl.GroupBox{
		Title:  "PointLightInstances (PointLight)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbPointLightInstance,
				Editable: false,
				Model:    pointLightInstances,
				Value:    defaultPointLightInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onPointLightInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onPointLightInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onPointLightInstanceDelete},
		},
	})

	lightInstances := []string{}
	for _, lightInstance := range data.LightInstances {
		lightInstances = append(lightInstances, lightInstance.Tag)
	}
	onLightInstanceNew := func() {
		slog.Println("new lightInstance")
	}
	onLightInstanceEdit := func() {
		slog.Println("edit lightInstance")
	}
	onLightInstanceDelete := func() {
		slog.Println("delete lightInstance")
	}

	var cmbLightInstance *walk.ComboBox
	defaultLightInstance := ""
	if len(lightInstances) > 0 {
		defaultLightInstance = lightInstances[0]
	}

	lightGroup.Children = append(lightGroup.Children, cpl.GroupBox{
		Title:  "LightInstances (Light)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbLightInstance,
				Editable: false,
				Model:    lightInstances,
				Value:    defaultLightInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onLightInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onLightInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onLightInstanceDelete},
		},
	})

	page.Title = "Light"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{lightGroup}
	return nil
}
