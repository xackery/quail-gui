package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/xackery/quail-gui/config"
	"github.com/xackery/quail-gui/gui/handler"
	qlog "github.com/xackery/quail/log"
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
	openPath      string
	fileName      string
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

	qlog.SetLogLevel(2)

	handler.ArchiveOpenSubscribe(c.onArchiveOpen)
	handler.ArchiveSaveSubscribe(c.onArchiveSave)
	handler.ArchiveExportAllSubscribe(c.onArchiveExportAll)
	handler.ArchiveExportFileSubscribe(c.onArchiveExportFile)
	handler.ArchiveRefreshSubscribe(c.onArchiveRefresh)

	return c, nil
}

func (c *Client) Done() error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}
