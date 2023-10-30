package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "embed"

	"github.com/xackery/quail-gui/gui"
	"github.com/xackery/quail-gui/slog"
)

var (
	Version string
)

func main() {
	start := time.Now()
	if Version == "" {
		Version = "0.0.1"
	}
	exeName, err := os.Executable()
	if err != nil {
		gui.ShowMessageBox("Error", "Failed to get executable name", true)
		os.Exit(1)
	}
	baseName := filepath.Base(exeName)
	if strings.Contains(baseName, ".") {
		baseName = baseName[0:strings.Index(baseName, ".")]
	}

	err = gui.New()
	if err != nil {
		slog.Printf("Failed to create main window: %s", err.Error())
		os.Exit(1)
	}

	defer slog.Dump()

	/* gui.SubscribeClose(func(canceled *bool, reason byte) {
		if ctx.Err() != nil {
			fmt.Println("Accepting exit")
			return
		}
		*canceled = true
		fmt.Println("Got close message")
		gui.SetTitle("Closing...")
		cancel()
	})

	go func() {
		<-ctx.Done()
		fmt.Println("Doing clean up process...")
		gui.Close()
		walk.App().Exit(0)
		fmt.Println("Done, exiting")
		slog.Dump(baseName + ".txt")
		os.Exit(0)
	}() */

	if len(os.Args) > 1 {
		path := os.Args[1]
		fileName := ""
		if len(os.Args) > 2 {
			fileName = os.Args[2]
		}
		section := ""
		if len(os.Args) > 3 {
			section = os.Args[3]
		}

		err = gui.Open(path, fileName, section)
		if err != nil {
			slog.Printf("Failed to open %s: %s", path, err.Error())
		}
	}

	slog.Printf("Started in %s\n", time.Since(start).String())
	errCode := gui.Run()
	if errCode != 0 {
		fmt.Println("Failed to run:", errCode)
		os.Exit(1)
	}

}
