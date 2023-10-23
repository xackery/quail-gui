package main

import (
	"fmt"
	"os"

	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Failed to run: %s\n", err)
		os.Exit(1)

	}
}

func run() error {

	var canvas *walk.Label
	var mw *walk.MainWindow
	cmw := cpl.MainWindow{
		AssignTo: &mw,
		Title:    "embed",
		MinSize:  cpl.Size{Width: 300, Height: 300},
		Layout:   cpl.VBox{},
		Visible:  false,
		Name:     "quail-gui",
		Children: []cpl.Widget{
			cpl.TextEdit{
				Text: "Text Edit Box",
			},
			cpl.Label{
				AssignTo: &canvas,
				Text:     "Canvas",
				MinSize:  cpl.Size{Width: 150, Height: 150},
			},
		},
	}
	err := cmw.Create()
	if err != nil {
		return fmt.Errorf("create main window: %s", err)
	}

	mw.SetVisible(true)
	code := mw.Run()
	if code != 0 {
		return fmt.Errorf("run main window: %d", code)
	}

	return nil
}
