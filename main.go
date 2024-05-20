package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "embed"

	"github.com/xackery/quail-gui/config"
	"github.com/xackery/quail-gui/gui"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/popup"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/raw"
)

var (
	Version string
)

func main() {
	start := time.Now()
	if Version == "" {
		Version = "0.0.1"
	}

	_, err := config.New(context.Background(), "quail-gui")
	if err != nil {
		popup.Errorf(gui.MainWindow(), "config new: %s", err)
		os.Exit(1)
	}

	err = ico.Init()
	if err != nil {
		popup.Errorf(gui.MainWindow(), "ico init: %s", err)
		os.Exit(1)
	}

	exeName, err := os.Executable()
	if err != nil {
		popup.MessageBox(gui.MainWindow(), "Error", "Failed to get executable name", true)
		os.Exit(1)
	}
	baseName := filepath.Base(exeName)
	if strings.Contains(baseName, ".") {
		baseName = baseName[0:strings.Index(baseName, ".")]
	}

	fileToOpen := ""
	if len(os.Args) > 1 {
		fileToOpen = os.Args[1]
	}

	err = gui.New()
	if err != nil {
		slog.Printf("Failed to create main window: %s", err.Error())
		os.Exit(1)
	}

	defer slog.Dump()

	/*gui.SubscribeClose(func(canceled *bool, reason byte) {
		if ctx.Err() != nil {
			fmt.Println("Accepting exit")
			return
		}
		*canceled = true
		fmt.Println("Got close message")
		gui.SetTitle("Closing...")
		cancel()
	})

	/*
		go func() {
			<-ctx.Done()
			fmt.Println("Doing clean up process...")
			gui.Close()
			walk.App().Exit(0)
			fmt.Println("Done, exiting")
			slog.Dump(baseName + ".txt")
			os.Exit(0)
		}() */

	if len(fileToOpen) > 1 {
		go func() {
			time.Sleep(10 * time.Millisecond)
			ext := strings.ToLower(filepath.Ext(fileToOpen))
			if !isArchive(ext) {
				err = quickEditFile(fileToOpen)
				if err != nil {
					if err.Error() == "cancelled" {
						slog.Printf("Cancelled edit %s\n", baseName)
						return
					}

					popup.Errorf(gui.MainWindow(), "show edit: %s", err)
					os.Exit(1)
				}
				os.Exit(0)
			}

			err = gui.Open(fileToOpen)
			if err != nil {
				popup.Errorf(gui.MainWindow(), "gui open: %s", err)
				return
			}
		}()
	}

	slog.Printf("Started in %s\n", time.Since(start).String())
	errCode := gui.Run()
	if errCode != 0 {
		fmt.Println("Failed to run:", errCode)
		os.Exit(1)
	}

}

func isArchive(ext string) bool {
	switch ext {
	case ".pfs", ".eqg", ".s3d", ".pak":
		return true
	}
	return false
}

// open a non-archive file
func quickEditFile(path string) error {

	ext := strings.ToLower(filepath.Ext(path))

	slog.Printf("Opening path: %s\n", path)

	r, err := os.Open(path)
	if err != nil {
		popup.Errorf(gui.MainWindow(), "os open: %s", err)
		os.Exit(1)
	}
	defer r.Close()

	value, err := raw.Read(ext, r)
	if err != nil {
		popup.Errorf(gui.MainWindow(), "raw read: %s", err)
		os.Exit(1)
	}

	data, err := gui.DialogEdit(path, value)
	if err != nil {
		if err.Error() == "cancelled" {
			slog.Printf("Cancelled without saving\n")
			return nil
		}
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		popup.Errorf(gui.MainWindow(), "os write: %s", err)
		os.Exit(1)
	}

	slog.Printf("Saved %s\n", path)

	return nil

}
