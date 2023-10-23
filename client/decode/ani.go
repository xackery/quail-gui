package decode

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/ani"
)

func Ani(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	anim := common.NewAnimation(name)
	err := ani.Decode(anim, r)
	if err != nil {
		return nil, fmt.Errorf("ani.Decode: %w", err)
	}
	treeModel := component.NewTreeModel()
	treeModel.SetRef(anim)

	treeModel.RootAdd(ico.Grab("header"), "Header", anim.Header)

	root := treeModel.RootAdd(ico.Grab(".bon"), fmt.Sprintf("Bones (%d)", len(anim.Bones)), anim.Bones)
	for _, bone := range anim.Bones {
		child := root.ChildAdd(ico.Grab(".bon"), bone.Name, bone)
		for i, frame := range bone.Frames {
			child.ChildAdd(ico.Grab(".ani"), fmt.Sprintf("Frame %d", i), frame)
		}
	}

	return treeModel, nil
}
