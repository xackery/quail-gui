package decode

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/zon"
)

func Zon(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	zone := common.NewZone(name)
	err := zon.Decode(zone, r)
	if err != nil {
		return nil, fmt.Errorf("zon.Decode: %w", err)
	}
	treeModel := component.NewTreeModel()
	treeModel.SetRef(zone)

	treeModel.RootAdd(ico.Grab("header"), "Header", zone.Header)

	root := treeModel.RootAdd(ico.Grab(".lit"), fmt.Sprintf("Lights (%d)", len(zone.Lights)), zone.Lights)
	for _, light := range zone.Lights {
		root.ChildAdd(ico.Grab(".lit"), light.Name, light)
	}

	root = treeModel.RootAdd(ico.Grab(".mod"), fmt.Sprintf("Models (%d)", len(zone.Models)), zone.Models)
	for _, model := range zone.Models {
		root.ChildAdd(ico.Grab(".mod"), model, model)
	}

	root = treeModel.RootAdd(ico.Grab(".obj"), fmt.Sprintf("Objects (%d)", len(zone.Objects)), zone.Objects)
	for _, object := range zone.Objects {
		root.ChildAdd(ico.Grab(".obj"), object.Name, object)
	}

	root = treeModel.RootAdd(ico.Grab("region"), fmt.Sprintf("Regions (%d)", len(zone.Regions)), zone.Regions)
	for _, region := range zone.Regions {
		root.ChildAdd(ico.Grab("region"), region.Name, region)
	}

	return treeModel, nil
}
