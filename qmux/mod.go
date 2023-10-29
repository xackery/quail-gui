package qmux

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/mesh/mod"
)

func ModDecode(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	model := common.NewModel(name)
	err := mod.Decode(model, r)
	if err != nil {
		return nil, fmt.Errorf("mod.Decode: %w", err)
	}
	treeModel := component.NewTreeModel()
	treeModel.SetRef(model)

	root := treeModel.RootAdd(ico.Grab(".mod"), "Model", model, model)

	root.ChildAdd(ico.Grab("header"), "Header", model, model.Header)

	materialRoot := root.ChildAdd(ico.Grab(".mat"), fmt.Sprintf("Materials (%d)", len(model.Materials)), model, model.Materials)
	for _, material := range model.Materials {
		matNode := materialRoot.ChildAdd(ico.Grab(".mat"), material.Name, model, material)
		for _, property := range material.Properties {
			matNode.ChildAdd(ico.Grab(".mat"), property.Name, model, property)
		}
	}

	root.ChildAdd(ico.Grab(".tri"), fmt.Sprintf("Triangles (%d)", len(model.Triangles)), model, model.Triangles)
	//for _, triangle := range model.Triangles {
	//		triangleRoot.ChildAdd(ico.Grab(".tri"), fmt.Sprintf("%d triangle", triangle.Index), model, triangle)
	//	}

	root.ChildAdd(ico.Grab(".ver"), fmt.Sprintf("Vertices (%d)", len(model.Vertices)), model, model.Vertices)
	/* for i, vert := range model.Vertices {
		vertexRoot.ChildAdd(ico.Grab(".ver"), fmt.Sprintf("%d vertex", i), model, vert)
	}*/

	boneRoot := root.ChildAdd(ico.Grab(".bon"), fmt.Sprintf("Bones (%d)", len(model.Bones)), model, model.Bones)
	for _, bone := range model.Bones {
		boneRoot.ChildAdd(ico.Grab(".bon"), bone.Name, model, bone)
	}
	return treeModel, nil
}
