package client

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xackery/encdec"
	"github.com/xackery/quail-gui/config"
	"github.com/xackery/quail-gui/gui"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/common"
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
	gui.SubscribeSavePFS(c.onSavePFS)
	gui.SubscribeSaveAllContent(c.onSaveAllContent)
	gui.SubscribeSaveContent(c.onSaveContent)
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

	c.sections = make(map[string]*gui.Section)
	if path != c.currentPath && file == "" {
		ext := strings.ToLower(filepath.Ext(path))
		isValidExt := false
		exts := []string{".eqg", ".s3d", ".pfs", ".pak"}
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
	fileExt := strings.ToLower(filepath.Ext(file))
	if fileExt == ".bat" {
		data, err := c.pfs.File(file)
		if err != nil {
			return fmt.Errorf("file %s: %w", file, err)
		}
		c.sections[".Info"] = &gui.Section{
			Name:    ".Info",
			Content: string(data),
			Icon:    generateIcon(file, data),
		}

		gui.SetSections(c.sections)
		return nil
	}
	if fileExt == ".lit" {
		c.sections[".Info"] = &gui.Section{
			Name:    ".Info",
			Content: "Lit files contain baked light data in a binary format.\r\nThere isn't much to show for contents",
			Icon:    generateIcon(file, nil),
		}
		gui.SetSections(c.sections)
		return nil
	}
	if c.wldInspect(file) {
		return nil
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

	gui.SetTitle(fmt.Sprintf("Archive: %s   quail-gui v%s", filepath.Base(path), c.version))
	c.sections[".Info"] = &gui.Section{
		Name: ".Info",
	}
	c.reflectTraversal(inspected, ".Info", 0, -1)

	gui.SetSections(c.sections)

	return nil
}

func (c *Client) onSavePFS(path string) error {
	return c.SavePFS(path)
}

func (c *Client) SavePFS(path string) error {
	if path == "" {
		path = c.currentPath
	}
	slog.Printf("Client saving %s\n", path)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	totalSize := 0
	for _, file := range c.pfs.Files() {
		totalSize += len(file.Data())
	}
	go func() {
		totalMB := float64(totalSize) / 1024 / 1024
		if totalMB < 1 {
			return
		}

		gui.SetProgress(1)
		defer gui.SetProgress(0)

		// set sleep 100ms for every mb
		sleep := time.Duration(totalMB) * 100 * time.Millisecond

		// every 100ms set progress 10
		for i := 0; i < 10; i++ {
			time.Sleep(sleep)
			select {
			case <-ctx.Done():
				return
			default:
			}
			gui.SetProgress(10 * (i + 1))
		}
	}()
	w, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer w.Close()
	err = c.pfs.Encode(w)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	slog.Printf("Saved %s\n", path)
	return nil
}

func (c *Client) onRefresh() {
	err := c.Open(c.openPath, "")
	if err != nil {
		slog.Print("Failed to refresh: %s", err)
	}
}

func (c *Client) onSaveAllContent(path string) error {
	return c.SaveAllContent(path)
}

func (c *Client) SaveAllContent(path string) error {
	slog.Printf("client saving all %s\n", path)
	for _, file := range c.pfs.Files() {
		err := c.saveContent(file)
		if err != nil {
			return fmt.Errorf("save %s: %w", file.Name(), err)
		}
	}
	slog.Printf("Saved %d files\n", len(c.pfs.Files()))
	return nil

}

func (c *Client) saveContent(entry pfs.FileEntry) error {
	w, err := os.Create(entry.Name())
	if err != nil {
		return fmt.Errorf("create %s: %w", entry.Name(), err)
	}
	defer w.Close()
	_, err = w.Write(entry.Data())
	if err != nil {
		return fmt.Errorf("write %s: %w", entry.Name(), err)
	}
	return nil
}

func (c *Client) onSaveContent(path string, file string) error {
	return c.SaveContent(path, file)
}

func (c *Client) SaveContent(path string, file string) error {
	slog.Printf("client saving %s %s\n", path, file)
	for _, entry := range c.pfs.Files() {
		slog.Printf("%s vs %s", entry.Name(), file)
		if !strings.EqualFold(entry.Name(), file) {
			continue
		}
		err := c.saveContent(entry)
		if err != nil {
			return fmt.Errorf("save %s: %w", entry.Name(), err)
		}
		slog.Printf("Saved %s\n", entry.Name())
		return nil
	}
	return nil

}

func (c *Client) wldInspect(file string) (isInspected bool) {

	fileExt := strings.ToLower(filepath.Ext(file))
	if fileExt != ".wld" {
		return
	}
	isInspected = true

	data, err := c.pfs.File(file)
	if err != nil {
		slog.Printf("Failed to open file %s: %s", file, err)
		return
	}
	c.sections[".Info"] = &gui.Section{
		Name:    ".Info",
		Content: "Wld Zone Data",
		Icon:    generateIcon(file, nil),
	}
	wld, err := common.WldOpen(bytes.NewReader(data))
	if err != nil {
		slog.Printf("Failed to open wld %s: %s", file, err)
		return
	}

	for i := uint32(0); i < wld.FragmentCount; i++ {
		data, err := wld.Fragment(int(i))
		if err != nil {
			slog.Printf("Failed to open fragment %d: %s", i, err)
			return
		}

		r := bytes.NewReader(data)
		dec := encdec.NewDecoder(r, binary.LittleEndian)

		fragCode := dec.Int32()
		fragName := common.FragName(int(fragCode))
		_, ok := c.sections[fragName]
		if !ok {
			c.sections[fragName] = &gui.Section{
				Name:    fragName,
				Count:   0,
				Content: "Todo: Actual content",
			}
		}
		c.sections[fragName].Count++
	}

	gui.SetSections(c.sections)

	return true
}
