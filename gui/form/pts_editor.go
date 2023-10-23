package form

import (
	"fmt"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail/common"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func showPtsEditor(page *walk.TabPage, node *component.TreeNode) (*walk.DataBinder, error) {
	src, ok := node.Ref().(common.ParticlePointEntry)
	if !ok {
		return nil, fmt.Errorf("node is not a particle point entry")
	}
	dst := &common.ParticlePointEntry{}
	dst.Bone = src.Bone
	dst.BoneSuffix = src.BoneSuffix
	dst.Name = src.Name
	dst.NameSuffix = src.NameSuffix
	dst.Rotation = src.Rotation
	dst.Scale = src.Scale
	dst.Translation = src.Translation
	db := new(walk.DataBinder)
	err := db.SetDataSource(dst)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func PtsEditWidgets() []cpl.Widget {
	return []cpl.Widget{
		cpl.Composite{
			Layout: cpl.Grid{Columns: 2},
			Children: []cpl.Widget{
				cpl.Label{Text: "Bone"}, cpl.LineEdit{Text: cpl.Bind("Bone")},
				cpl.Label{Text: "Name"}, cpl.LineEdit{Text: cpl.Bind("Name")},
				cpl.Label{Text: "Rotation"}, cpl.Composite{
					Layout: cpl.HBox{},
					Children: []cpl.Widget{
						cpl.LineEdit{Text: cpl.Bind("Rotation.X")},
						cpl.LineEdit{Text: cpl.Bind("Rotation.Y")},
						cpl.LineEdit{Text: cpl.Bind("Rotation.Z")},
					},
				},
				cpl.Label{Text: "Scale"}, cpl.Composite{
					Layout: cpl.HBox{},
					Children: []cpl.Widget{
						cpl.LineEdit{Text: cpl.Bind("Scale.X")},
						cpl.LineEdit{Text: cpl.Bind("Scale.Y")},
						cpl.LineEdit{Text: cpl.Bind("Scale.Z")},
					},
				},
				cpl.Label{Text: "Translation"}, cpl.Composite{
					Layout: cpl.HBox{},
					Children: []cpl.Widget{
						cpl.LineEdit{Text: cpl.Bind("Translation.X")},
						cpl.LineEdit{Text: cpl.Bind("Translation.Y")},
						cpl.LineEdit{Text: cpl.Bind("Translation.Z")},
					},
				},
			},
		},
	}
}
