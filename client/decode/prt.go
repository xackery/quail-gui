package decode

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/prt"
)

func Prt(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	render := common.NewParticleRender(name)
	err := prt.Decode(render, r)
	if err != nil {
		return nil, fmt.Errorf("prt.Decode: %w", err)
	}
	treeModel := component.NewTreeModel()
	treeModel.SetRef(render)

	treeModel.RootAdd(ico.Grab("header"), "Header", render.Header)
	root := treeModel.RootAdd(ico.Grab(".prt"), "Renders", render.Entries)
	for _, entry := range render.Entries {
		root.ChildAdd(ico.Grab(".prt"), entry.ParticlePoint, entry)
	}

	return treeModel, nil
}
