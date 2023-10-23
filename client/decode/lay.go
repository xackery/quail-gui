package decode

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/lay"
)

func Lay(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	model := common.NewModel(name)
	err := lay.Decode(model, r)
	if err != nil {
		return nil, fmt.Errorf("lay.Decode: %w", err)
	}
	treeModel := component.NewTreeModel()
	treeModel.SetRef(model.Layers)
	root := treeModel.RootAdd(ico.Grab(".mod"), "Model", model)

	root.ChildAdd(ico.Grab("header"), "Header", model.Header)

	child := root.ChildAdd(ico.Grab(".lay"), "Layers", model.Layers)
	for _, layer := range model.Layers {
		child.ChildAdd(ico.Grab(".lay"), layer.Material, layer)
	}

	return treeModel, nil
}
