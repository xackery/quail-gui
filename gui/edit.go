package gui

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail-gui/gui/dialog"
	"github.com/xackery/quail-gui/popup"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/raw"
	"github.com/xackery/wlk/walk"
)

var (
	menuEditPreferences *walk.Action
)

func onEditPreferences() {
	err := dialog.ShowPreferences(mw)
	if err != nil {
		if err.Error() == "cancelled" {
			return
		}
		popup.Errorf(mw, "preferences: %s", err)
	}

}

func DialogEdit(itemName string, value raw.ReadWriter) ([]byte, error) {

	valueType := value.Identity()

	slog.Println("Asserted type:", valueType)

	extFuncs := map[string]func(*walk.MainWindow, string, raw.ReadWriter) error{
		"mod":       dialog.ShowModEdit,
		"zon":       dialog.ShowZonEdit,
		"wld":       dialog.ShowWldEdit,
		"wld.ascii": dialog.ShowWldAsciiEdit,
		"mds":       dialog.ShowMdsEdit,
		"txt":       dialog.ShowTxtEdit,
	}
	for funcType, f := range extFuncs {
		if valueType != funcType {
			continue
		}
		err := f(mw, itemName, value)
		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer([]byte{})
		err = value.Write(buf)
		if err != nil {
			return nil, fmt.Errorf("write %s: %w", itemName, err)
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("unsupported file type: %s", valueType)
}
