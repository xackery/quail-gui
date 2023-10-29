package form

import (
	"fmt"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/gui/handler"
	"github.com/xackery/quail/common"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

var (
	modEditor *ModEditor
)

type ModEditor struct {
	node *component.TreeNode
}

func showModEditor(page *walk.TabPage, node *component.TreeNode) (Editor, error) {
	_, ok := node.Ref().(*common.Model)
	if !ok {
		return nil, fmt.Errorf("node is not a model")
	}

	_, ok = node.RootRef().(*common.Model)
	if !ok {
		return nil, fmt.Errorf("root ref is not a model")
	}

	e := modEditor

	e.node = node
	return e, nil
}

func (e *ModEditor) Save() error {
	return nil
}

func (e *ModEditor) Reset() {

}

func (e *ModEditor) Node() *component.TreeNode {
	return e.node
}

func (e *ModEditor) Ext() string {
	return ".mod"
}

func ModEditWidgets() []cpl.Widget {
	if modEditor == nil {
		modEditor = &ModEditor{}
	}
	return []cpl.Widget{
		cpl.Composite{
			Layout: cpl.Grid{Columns: 2},
			Children: []cpl.Widget{
				cpl.PushButton{
					Text:      "Preview 3D Model",
					OnClicked: handler.PreviewInvoke,
				},
			},
		},
	}
}

func (e *ModEditor) New(src interface{}) (*component.TreeNode, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e *ModEditor) Name() string {
	return "Model"
}
