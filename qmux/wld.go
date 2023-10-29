package qmux

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/wld"
	"github.com/xackery/quail/quail"
)

func WldDecode(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	world := common.NewWld(name)
	err := wld.Decode(world, r)
	if err != nil {
		return nil, fmt.Errorf("wld.Decode: %w", err)
	}
	q := quail.New()
	err = q.WldUnmarshal(world)
	if err != nil {
		return nil, fmt.Errorf("WldUnmarshal: %w", err)
	}

	treeModel := component.NewTreeModel()
	treeModel.SetRef(world)

	root := treeModel.RootAdd(ico.Grab(".wld"), "World", world, q)

	root.ChildAdd(ico.Grab("header"), "Header", world, q.Header)

	modelRoot := root.ChildAdd(ico.Grab(".mod"), fmt.Sprintf("Models (%d)", len(world.Models)), world, q.Models)
	for _, model := range world.Models {
		modelChild := modelRoot.ChildAdd(ico.Grab(".mod"), model.Header.Name, world, model)
		materialRoot := modelChild.ChildAdd(ico.Grab(".mat"), fmt.Sprintf("Materials (%d)", len(model.Materials)), world, model.Materials)
		for _, material := range model.Materials {
			matNode := materialRoot.ChildAdd(ico.Grab(".mat"), material.Name, world, material)
			for _, property := range material.Properties {
				matNode.ChildAdd(ico.Grab(".mat"), property.Name, world, property)
			}
		}

		modelChild.ChildAdd(ico.Grab(".tri"), fmt.Sprintf("Triangles (%d)", len(model.Triangles)), world, model.Triangles)
		//for _, triangle := range model.Triangles {
		//	triangleRoot.ChildAdd(ico.Grab(".tri"), fmt.Sprintf("%d triangle", triangle.Index), world, triangle)
		//}

		modelChild.ChildAdd(ico.Grab(".ver"), fmt.Sprintf("Vertices (%d)", len(model.Vertices)), world, model.Vertices)
		//for i, vert := range model.Vertices {
		//	vertexRoot.ChildAdd(ico.Grab(".ver"), fmt.Sprintf("%d vertex", i), world, vert)
		//}

		boneRoot := modelChild.ChildAdd(ico.Grab(".bon"), fmt.Sprintf("Bones (%d)", len(model.Bones)), world, model.Bones)
		for _, bone := range model.Bones {
			boneRoot.ChildAdd(ico.Grab(".bon"), bone.Name, world, bone)
		}
	}

	return treeModel, nil
}
