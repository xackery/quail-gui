package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/xackery/quail-gui/config"
	"github.com/xackery/quail-gui/gui"
	"github.com/xackery/quail-gui/slog"
)

// Client wraps the entire UI
type Client struct {
	ctx           context.Context
	cancel        context.CancelFunc
	currentPath   string
	clientVersion string
	cfg           *config.Config
	version       string
	httpClient    *http.Client
}

// New creates a new client
func New(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, version string) (*Client, error) {
	var err error
	c := &Client{
		ctx:           ctx,
		cancel:        cancel,
		cfg:           cfg,
		clientVersion: "rof",
		version:       version,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}

	fmt.Printf("Starting quail-gui %s\n", c.version)
	c.currentPath, err = os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("wd invalid: %w", err)
	}

	gui.SubscribeOpen(c.onOpen)

	/*gui.SubscribePatchButton(func() {
		err := c.Patch()
		if err != nil {
			slog.Print("Failed to patch: %s", err)
		}
	})*/

	return c, nil
}

func (c *Client) Done() error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

func (c *Client) onOpen(path string) {

	err := open(path)
	if err != nil {
		slog.Print("Failed to open: %s", err)
	}
}

func open(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path check: %w", err)
	}
	if fi.IsDir() {
		return fmt.Errorf("inspect requires a target file, directory provided")
	}

	inspected, err := inspect(path, "")
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}
	if inspected == nil {
		return nil
	}

	gui.SetTitle(filepath.Base(path))

	reflectTraversal(inspected, 0, -1)
	return nil
}
