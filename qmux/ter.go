package qmux

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/mesh/ter"
)

func TerDecode(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	model := common.NewModel(name)
	err := ter.Decode(model, r)
	if err != nil {
		return nil, fmt.Errorf("ter.Decode: %w", err)
	}
	treeModel := component.NewTreeModel()
	treeModel.SetRef(model)

	treeModel.RootAdd(ico.Grab("header"), "Header", model, model.Header)

	root := treeModel.RootAdd(ico.Grab(".mat"), fmt.Sprintf("Materials (%d)", len(model.Materials)), model, model.Materials)
	for _, material := range model.Materials {
		child := root.ChildAdd(ico.Grab(".mat"), material.Name, model, material)
		for _, property := range material.Properties {
			child2 := child.ChildAdd(ico.Grab(".mat"), fmt.Sprintf("Property %s", property.Name), model, &property)
			child2.ChildAdd(ico.Grab(".mat"), fmt.Sprintf("Name: %s", property.Name), model, &property.Name)
			child2.ChildAdd(ico.Grab(".mat"), fmt.Sprintf("Category: %d", property.Category), model, &property.Category)
			child2.ChildAdd(ico.Grab(".mat"), fmt.Sprintf("Value: %s", property.Value), model, &property.Value)
		}
	}
	treeModel.RootAdd(ico.Grab(".tri"), fmt.Sprintf("Triangles (%d)", len(model.Triangles)), model, model.Triangles)
	//for _, triangle := range model.Triangles {
	//	root.ChildAdd(ico.Grab(".tri"), fmt.Sprintf("%d", triangle.Index), model, triangle)
	//}

	treeModel.RootAdd(ico.Grab(".ver"), fmt.Sprintf("Vertices (%d)", len(model.Vertices)), model, model.Vertices)
	//for i, vert := range model.Vertices {
	//	root.ChildAdd(ico.Grab(".ver"), fmt.Sprintf("%d vertex", i), model, vert)
	//}

	return treeModel, nil
}
