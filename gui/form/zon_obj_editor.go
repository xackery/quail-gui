package form

import (
	"fmt"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail/common"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func showZonObjEditor(page *walk.TabPage, node *component.TreeNode) (*walk.DataBinder, error) {
	src, ok := node.Ref().(common.Light)
	if !ok {
		return nil, fmt.Errorf("node is not a layer")
	}
	dst := common.Light{}
	dst.Color = src.Color
	dst.Name = src.Name
	dst.Position = src.Position
	dst.Radius = src.Radius
	db := new(walk.DataBinder)
	err := db.SetDataSource(dst)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ZonObjEditWidgets() []cpl.Widget {
	return []cpl.Widget{
		cpl.Composite{
			Layout: cpl.Grid{Columns: 2},
			Children: []cpl.Widget{
				cpl.Label{Text: "Name"}, cpl.LineEdit{Text: cpl.Bind("Name", validateMaterial)},
			},
		},
	}
}
