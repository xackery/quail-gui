package qmux

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/pts"
)

func PtsDecode(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	point := common.NewParticlePoint(name)
	err := pts.Decode(point, r)
	if err != nil {
		return nil, fmt.Errorf("pts.Decode: %w", err)
	}
	treeModel := component.NewTreeModel()
	treeModel.SetRef(point)

	root := treeModel.RootAdd(ico.Grab(".pts"), "Particle Point", point, point)

	root.ChildAdd(ico.Grab("header"), "Header", point, point.Header)
	pointNode := root.ChildAdd(ico.Grab(".pts"), "Point Entries", point, point.Entries)
	for _, entry := range point.Entries {
		pointNode.ChildAdd(ico.Grab(".pts"), entry.Name, point, entry)
	}

	return treeModel, nil
}
