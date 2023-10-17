package gui

import (
	"sort"

	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/walk"
)

type FileViewEntry struct {
	Icon    *walk.Bitmap
	Name    string
	Ext     string
	Size    string
	checked bool
}

type FileView struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*FileViewEntry
}

func NewFileView() *FileView {
	m := new(FileView)
	m.ResetRows()
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *FileView) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *FileView) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case -1:
		return nil
	case 0:
		return ""
	case 1:
		return item.Name
	case 2:
		return item.Ext
	case 3:
		return item.Size
	}

	slog.Printf("invalid col: %d\n", col)
	return nil
}

// Called by the TableView to retrieve if a given row is checked.
func (m *FileView) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *FileView) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

// Called by the TableView to sort the model.
func (m *FileView) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case -1:
			return false
		case 0:
			/* if a.Ext == ".dds" {
				return true
			}
			if a.Ext == ".bmp" {
				return true
			}
			if a.Ext == ".png" {
				return true
			} */

			return c(a.Ext < b.Ext)
		case 1:
			return c(a.Name < b.Name)
		case 2:
			return c(a.Ext < b.Ext)
		case 3:
			return c(a.Size < b.Size)
		}

		slog.Printf("invalid sort col: %d", m.sortColumn)
		return false
	})

	return m.SorterBase.Sort(col, order)
}

func (m *FileView) ResetRows() {
	m.items = nil

	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}

func (m *FileView) SetItems(items []*FileViewEntry) {
	m.items = items

	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}

type fileViewStyler struct {
}

func (fv *fileViewStyler) StyleCell(style *walk.CellStyle) {
	if style.Col() != 0 {
		return
	}

	if style.Row() >= len(gui.fileView.items) {
		return
	}

	item := gui.fileView.items[style.Row()]
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
