package form

import (
	"fmt"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail/common"
	"github.com/xackery/wlk/walk"
)

type Editor interface {
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
		/* case common.ParticlePointEntry:
			ext := "pts"
			err := showPtsEditor(page, node)
			return ext, err
		case *common.Model:
			ext := "mod"
			err := showModEditor(page, node)
			return ext, err
		*/
	}

	return nil, fmt.Errorf("unknown type %T", ref)
}
