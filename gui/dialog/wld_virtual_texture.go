package dialog

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/vwld"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualTexturePage(data *vwld.VWld, page *cpl.TabPage) error {
	bitmaps := []string{}
	for _, bitmap := range data.Bitmaps {
		bitmaps = append(bitmaps, bitmap.Tag)
	}
	onBitmapNew := func() {
		slog.Println("new bitmap")
	}
	onBitmapEdit := func() {
		slog.Println("edit bitmap")
	}
	onBitmapDelete := func() {
		slog.Println("delete bitmap")
	}

	var cmbBitmap *walk.ComboBox
	defaultBitmap := ""
	if len(bitmaps) > 0 {
		defaultBitmap = bitmaps[0]
	}

	page.Title = "Texture"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{
		cpl.GroupBox{
			Title:  "Bitmaps (BMInfo)",
			Layout: cpl.HBox{},
			Children: []cpl.Widget{
				cpl.ComboBox{
					AssignTo: &cmbBitmap,
					Editable: false,
					Model:    bitmaps,
					Value:    defaultBitmap,
				},
				cpl.PushButton{Text: "Add", OnClicked: onBitmapNew},
				cpl.PushButton{Text: "Edit", OnClicked: onBitmapEdit},
				cpl.PushButton{Text: "Delete", OnClicked: onBitmapDelete},
			},
		},
	}
	return nil
}
