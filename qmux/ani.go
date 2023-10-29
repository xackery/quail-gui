package qmux

import (
	"fmt"
	"io"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/ani"
)

func AniDecode(name string, r io.ReadSeeker) (*component.TreeModel, error) {
	anim := common.NewAnimation(name)
	err := ani.Decode(anim, r)
	if err != nil {
		return nil, fmt.Errorf("ani.Decode: %w", err)
	}
	treeModel := component.NewTreeModel()
	treeModel.SetRef(anim)

	root := treeModel.RootAdd(ico.Grab(".ani"), "Animation", anim, anim)

	root.ChildAdd(ico.Grab("header"), "Header", anim, anim.Header)

	boneNode := root.ChildAdd(ico.Grab(".bon"), fmt.Sprintf("Bones (%d)", len(anim.Bones)), anim, anim.Bones)
	for _, bone := range anim.Bones {
		child := boneNode.ChildAdd(ico.Grab(".bon"), bone.Name, anim, bone)
		for i, frame := range bone.Frames {
			child.ChildAdd(ico.Grab(".ani"), fmt.Sprintf("Frame %d", i), anim, frame)
		}
	}

	return treeModel, nil
}
