package decode

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/pts"
)

func Pts(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	point := common.NewParticlePoint(name)
	err := pts.Decode(point, r)
	if err != nil {
		return nil, fmt.Errorf("pts.Decode: %w", err)
	}
	treeModel := component.NewTreeModel()
	treeModel.SetRef(point)

	root := treeModel.RootAdd(ico.Grab(".pts"), "Particle Point", point)

	root.ChildAdd(ico.Grab("header"), "Header", point.Header)
	pointNode := root.ChildAdd(ico.Grab(".pts"), "Point Entries", point.Entries)
	for _, entry := range point.Entries {
		pointNode.ChildAdd(ico.Grab(".pts"), entry.Name, entry)
	}

	return treeModel, nil
}
