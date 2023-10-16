package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/xackery/quail-gui/client"
	"github.com/xackery/quail-gui/config"
	"github.com/xackery/quail-gui/gui"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/walk"
)

var (
	Version    string
	PatcherUrl string
)

func main() {
	if Version == "" {
		Version = "0.0.1"
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exeName, err := os.Executable()
	if err != nil {
		gui.MessageBox("Error", "Failed to get executable name", true)
		os.Exit(1)
	}
	baseName := filepath.Base(exeName)
	if strings.Contains(baseName, ".") {
		baseName = baseName[0:strings.Index(baseName, ".")]
	}
	cfg, err := config.New(context.Background(), baseName)
	if err != nil {
		slog.Printf("Failed to load config: %s", err.Error())
		os.Exit(1)
	}

	err = gui.NewMainWindow(ctx, cancel, cfg, Version)
	if err != nil {
		slog.Printf("Failed to create main window: %s", err.Error())
		os.Exit(1)
	}

	c, err := client.New(ctx, cancel, cfg, Version)
	if err != nil {
		gui.MessageBox("Error", "Failed to create client: "+err.Error(), true)
		os.Exit(1)
	}
	defer slog.Dump(baseName + ".txt")
	defer c.Done()

	gui.SubscribeClose(func(canceled *bool, reason walk.CloseReason) {
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
		c.Done() // close client
		gui.Close()
		walk.App().Exit(0)
		fmt.Println("Done, exiting")
		slog.Dump(baseName + ".txt")
		os.Exit(0)
	}()

	errCode := gui.Run()
	if errCode != 0 {
		fmt.Println("Failed to run:", errCode)
		os.Exit(1)
	}

}
