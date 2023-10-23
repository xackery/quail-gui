package decode

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/mesh/mod"
)

func Mod(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	model := common.NewModel(name)
	err := mod.Decode(model, r)
	if err != nil {
		return nil, fmt.Errorf("mod.Decode: %w", err)
	}
	treeModel := component.NewTreeModel()
	treeModel.SetRef(model)

	treeModel.RootAdd(ico.Grab("preview"), "Preview", model)
	treeModel.RootAdd(ico.Grab("header"), "Header", model.Header)

	root := treeModel.RootAdd(ico.Grab(".mat"), fmt.Sprintf("Materials (%d)", len(model.Materials)), model.Materials)
	for _, material := range model.Materials {
		root.ChildAdd(ico.Grab(".mat"), material.Name, material)
	}

	root = treeModel.RootAdd(ico.Grab(".tri"), fmt.Sprintf("Triangles (%d)", len(model.Triangles)), model.Triangles)
	for _, triangle := range model.Triangles {
		root.ChildAdd(ico.Grab(".tri"), fmt.Sprintf("%d triangle", triangle.Index), triangle)
	}

	root = treeModel.RootAdd(ico.Grab(".ver"), fmt.Sprintf("Vertices (%d)", len(model.Vertices)), model.Vertices)
	for i, vert := range model.Vertices {
		root.ChildAdd(ico.Grab(".ver"), fmt.Sprintf("%d vertex", i), vert)
	}

	root = treeModel.RootAdd(ico.Grab(".bon"), fmt.Sprintf("Bones (%d)", len(model.Bones)), model.Bones)
	for _, bone := range model.Bones {
		root.ChildAdd(ico.Grab(".bon"), bone.Name, bone)
	}
	return treeModel, nil
}
