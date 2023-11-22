package component

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/walk"
)

type ElementViewStyler struct {
	elementView *ElementView
}

func NewElementViewStyler(elementView *ElementView) *ElementViewStyler {
	fvs := new(ElementViewStyler)
	fvs.elementView = elementView
	return fvs
}

func (fv *ElementViewStyler) StyleCell(style *walk.CellStyle) {
	if style.Col() != 0 {
		style.Image = ""
		return
	}

	if style.Row() >= len(fv.elementView.items) {
		return
	}

	item := fv.elementView.items[style.Row()]
	if item == nil {
		slog.Printf("item %d is nil\n", style.Row())
		return
	}

	if item.Icon == nil {
		slog.Printf("item %d icon is nil\n", style.Row())
		return
	}

	style.Image = item.Icon

	/* canvas := style.Canvas()
	if canvas == nil {
		return
	}
	bounds := style.Bounds()
	bounds.X += 2
	bounds.Y += 2
	bounds.Width = 16
	bounds.Height = 16
	err := canvas.DrawBitmapPartWithOpacityPixels(item.Icon, bounds, walk.Rectangle{X: 0, Y: 0, Width: 16, Height: 16}, 127)
	if err != nil {
		slog.Printf("failed to draw bitmap: %s\n", err.Error())
	} */

	/*

		switch style.Col() {
		case 1:
			if canvas := style.Canvas(); canvas != nil {
				bounds := style.Bounds()
				bounds.X += 2
				bounds.Y += 2
				bounds.Width = int((float64(bounds.Width) - 4) / 5 * float64(len(item.Bar)))
				bounds.Height -= 4
				canvas.DrawBitmapPartWithOpacity(barBitmap, bounds, walk.Rectangle{0, 0, 100 / 5 * len(item.Bar), 1}, 127)

				bounds.X += 4
				bounds.Y += 2
				canvas.DrawText(item.Bar, tv.Font(), 0, bounds, walk.TextLeft)
			}

		case 2:
			if item.Baz >= 900.0 {
				style.TextColor = walk.RGB(0, 191, 0)
				style.Image = goodIcon
			} else if item.Baz < 100.0 {
				style.TextColor = walk.RGB(255, 0, 0)
				style.Image = badIcon
			}

		case 3:
			if item.Quux.After(time.Now().Add(-365 * 24 * time.Hour)) {
				style.Font = boldFont
			}
		}
	*/
}
