package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xackery/quail-gui/config"
	"github.com/xackery/quail-gui/gui"
	"github.com/xackery/quail-gui/slog"
	qlog "github.com/xackery/quail/log"
	"github.com/xackery/quail/pfs"
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
	sections      map[string]*gui.Section
	pfs           *pfs.PFS
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

	gui.SubscribeOpen(c.onOpen)
	gui.SubscribeSave(c.onSave)
	gui.SubscribeRefresh(c.onRefresh)

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

func (c *Client) onOpen(path string, file string) error {
	return c.Open(path, file)
}

func (c *Client) Open(path string, file string) error {
	c.openPath = path
	c.fileName = file

	if path == "" {
		path = c.currentPath
	}
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path check: %w", err)
	}
	if fi.IsDir() {
		return fmt.Errorf("inspect requires a target file, directory provided")
	}

	if path != c.currentPath && file == "" {
		isValidExt := false
		exts := []string{".eqg", ".s3d", ".pfs", ".pak"}
		ext := strings.ToLower(filepath.Ext(path))
		for _, ext := range exts {
			if strings.HasSuffix(path, ext) {
				isValidExt = true
				break
			}
		}

		if isValidExt {
			c.pfs, err = pfs.NewFile(path)
			if err != nil {
				return fmt.Errorf("%s load: %w", ext, err)
			}
		}
	}
	c.currentPath = path
	c.fileName = file

	inspected, err := c.inspect(file)
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}
	if inspected == nil {
		return nil
	}

	gui.SetTitle(filepath.Base(path))

	c.sections = make(map[string]*gui.Section)
	c.sections[".Info"] = &gui.Section{
		Name: ".Info",
	}
	c.reflectTraversal(inspected, ".Info", 0, -1)

	gui.SetSections(c.sections)

	return nil
}

func (c *Client) onSave(path string) {
	err := c.Save(path)
	if err != nil {
		slog.Print("Failed to save: %s", err)
	}
}

func (c *Client) Save(path string) error {
	slog.Printf("client saving %s\n", path)
	w, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer w.Close()
	err = c.pfs.Encode(w)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	return nil

}

func (c *Client) onRefresh() {
	err := c.Open(c.openPath, "")
	if err != nil {
		slog.Print("Failed to refresh: %s", err)
	}
}
