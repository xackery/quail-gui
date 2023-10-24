package decode

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/wld"
)

func Wld(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	world := common.NewWld(name)
	err := wld.Decode(world, r)
	if err != nil {
		return nil, fmt.Errorf("wld.Decode: %w", err)
	}
	err = wld.Convert(world)
	if err != nil {
		return nil, fmt.Errorf("wld.Convert: %w", err)
	}

	treeModel := component.NewTreeModel()
	treeModel.SetRef(world)

	root := treeModel.RootAdd(ico.Grab(".wld"), "World", world)

	root.ChildAdd(ico.Grab("header"), "Header", world.Header)

	modelRoot := root.ChildAdd(ico.Grab(".mod"), fmt.Sprintf("Models (%d)", len(world.Models)), world.Models)
	for _, model := range world.Models {
		modelChild := modelRoot.ChildAdd(ico.Grab(".mod"), model.Header.Name, model)
		materialRoot := modelChild.ChildAdd(ico.Grab(".mat"), fmt.Sprintf("Materials (%d)", len(model.Materials)), model.Materials)
		for _, material := range model.Materials {
			matNode := materialRoot.ChildAdd(ico.Grab(".mat"), material.Name, material)
			for _, property := range material.Properties {
				matNode.ChildAdd(ico.Grab(".mat"), property.Name, property)
			}
		}

		triangleRoot := modelChild.ChildAdd(ico.Grab(".tri"), fmt.Sprintf("Triangles (%d)", len(model.Triangles)), model.Triangles)
		for _, triangle := range model.Triangles {
			triangleRoot.ChildAdd(ico.Grab(".tri"), fmt.Sprintf("%d triangle", triangle.Index), triangle)
		}

		vertexRoot := modelChild.ChildAdd(ico.Grab(".ver"), fmt.Sprintf("Vertices (%d)", len(model.Vertices)), model.Vertices)
		for i, vert := range model.Vertices {
			vertexRoot.ChildAdd(ico.Grab(".ver"), fmt.Sprintf("%d vertex", i), vert)
		}

		boneRoot := modelChild.ChildAdd(ico.Grab(".bon"), fmt.Sprintf("Bones (%d)", len(model.Bones)), model.Bones)
		for _, bone := range model.Bones {
			boneRoot.ChildAdd(ico.Grab(".bon"), bone.Name, bone)
		}
	}

	return treeModel, nil
}
