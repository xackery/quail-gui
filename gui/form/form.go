package form

import (
	"fmt"
	"strconv"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail/common"
	"github.com/xackery/wlk/walk"
)

type Editor interface {
	Name() string
	New(src interface{}) (*component.TreeNode, error)
	Save() error
	Ext() string
	Reset()
	Node() *component.TreeNode
}

// ShowEditor opens a form that let's you edit various components from the tree view
func ShowEditor(page *walk.TabPage, node *component.TreeNode) (Editor, error) {
	ref := node.Ref()
	switch ref.(type) {
	//case *common.Material:
	//	return ShowMaterialEditor(node.name, node.ref.(*common.Material))
	case *common.Layer:
		return showLayEditor(page, node)
	case *common.Header:
		return showHeaderEditor(page, node)
	case *common.Model:
		return showModEditor(page, node)
	case *common.ParticlePointEntry:
		return showPtsEditor(page, node)
	}

	return nil, fmt.Errorf("unknown type %T", ref)
}

func fastFloat32(in string) float32 {
	f, err := strconv.ParseFloat(in, 32)
	if err != nil {
		return 0
	}
	return float32(f)
}
