package form

import (
	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/gui/handler"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func showModEditor(page *walk.TabPage, node *component.TreeNode) error {
	/* _, ok := node.Ref().(*common.Model)
	if !ok {
		return nil, fmt.Errorf("node is not a model")
	}
	dst := &common.Model{}
	db := new(walk.DataBinder)

	present, err := walk.NewToolTipErrorPresenter()
	if err != nil {
		return nil, err
	}
	db.SetErrorPresenter(present)

	err = db.SetDataSource(dst)
	if err != nil {
		return nil, err
	} */

	return nil
}

func ModEditWidgets() []cpl.Widget {
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
