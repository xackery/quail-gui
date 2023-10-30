package ico

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	_ "embed"

	ico "github.com/biessek/golang-ico"
	"github.com/malashin/dds"
	"github.com/sergeymakinen/go-bmp"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/walk"
	"golang.org/x/image/draw"
)

var (
	//go:embed assets
	assets embed.FS
	icos   map[string]*walk.Icon
)

func init() {
	icos = make(map[string]*walk.Icon)

	dir, err := assets.ReadDir("assets")
	if err != nil {
		slog.Printf("Failed to read assets: %s\n", err.Error())
		return
	}
	for _, fi := range dir {
		name := fi.Name()
		r, err := assets.Open(fmt.Sprintf("assets/%s", name))
		if err != nil {
			slog.Printf("Failed to open %s: %s\n", name, err.Error())
			continue
		}
		icoData, err := io.ReadAll(r)
		if err != nil {
			slog.Printf("Failed to read %s: %s\n", name, err.Error())
			continue
		}

		icon := Generate(name, icoData)
		if strings.Contains(name, ".") {
			name = name[0:strings.Index(name, ".")]
		}
		name = strings.ToLower(name)
		if len(name) == 3 {
			name = "." + name
		}

		icos[name] = icon
	}
}

// Grab returns a walk.Icon for a given icon
func Grab(name string) *walk.Icon {
	icon, ok := icos[name]
	if !ok {
		return icos[".unk"]
	}
	return icon
}

func Generate(name string, data []byte) *walk.Icon {
	var err error

	icon, ok := icos[name]
	if ok {
		return icon
	}

	ext := strings.ToLower(filepath.Ext(name))
	defer func() {
		if err != nil {
			slog.Printf("GenerateIcon: %s", err)
			return
		}
	}()

	icon = Grab(ext)
	if icon != nil && ext != ".unk" {
		return icon
	}

	unkImg := Grab(".unk")

	var img image.Image
	switch ext {
	case ".ico":
		img, err = ico.Decode(bytes.NewReader(data))
		if err != nil {
			err = fmt.Errorf("ico.Decode %s: %w", name, err)
			return unkImg
		}
		if img.Bounds().Max.X > 16 || img.Bounds().Max.Y > 16 {
			dst := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X/2, img.Bounds().Max.Y/2))
			draw.NearestNeighbor.Scale(dst, image.Rect(0, 0, 16, 16), img, img.Bounds(), draw.Over, nil)
			img = dst
		}

		icon, err = walk.NewIconFromImageForDPI(img, 96)
		if err != nil {
			err = fmt.Errorf("new icon from image for dpi: %s", err)
			return unkImg
		}
		return icon
	case ".dds":
		img, err = dds.Decode(bytes.NewReader(data))
		if err != nil {
			err = fmt.Errorf("dds.Decode %s: %w", name, err)
			return unkImg
		}
		dst := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X/2, img.Bounds().Max.Y/2))
		draw.NearestNeighbor.Scale(dst, image.Rect(0, 0, 16, 16), img, img.Bounds(), draw.Over, nil)

		icon, err = walk.NewIconFromImageForDPI(dst, 96)
		if err != nil {
			err = fmt.Errorf("new icon from image for dpi: %s", err)
			return unkImg
		}
		return icon
	case ".png":
		img, err = png.Decode(bytes.NewReader(data))
		if err != nil {
			err = fmt.Errorf("png.Decode %s: %w", name, err)
			return unkImg
		}
		if img.Bounds().Max.X > 16 || img.Bounds().Max.Y > 16 {
			dst := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X/2, img.Bounds().Max.Y/2))
			draw.NearestNeighbor.Scale(dst, image.Rect(0, 0, 16, 16), img, img.Bounds(), draw.Over, nil)
			img = dst
		}

		icon, err = walk.NewIconFromImageForDPI(img, 96)
		if err != nil {
			err = fmt.Errorf("new icon from image for dpi: %s", err)
			return unkImg
		}
		return icon
	case ".bmp":
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(bytes.NewReader(data))
		if err != nil {
			err = fmt.Errorf("buf read from: %w", err)
			return unkImg
		}
		var img image.Image
		if string(buf.Bytes()[0:3]) == "DDS" {
			img, err = dds.Decode(bytes.NewReader(data))
			if err != nil {
				err = fmt.Errorf("dds.Decode %s: %w", name, err)
				return unkImg
			}
		} else {
			img, err = bmp.Decode(bytes.NewReader(data))
			if err != nil {
				err = fmt.Errorf("bmp.Decode %s: %w", name, err)
				return unkImg
			}
		}
		dst := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X/2, img.Bounds().Max.Y/2))
		draw.NearestNeighbor.Scale(dst, image.Rect(0, 0, 16, 16), img, img.Bounds(), draw.Over, nil)

		icon, err = walk.NewIconFromImageForDPI(dst, 96)
		if err != nil {
			err = fmt.Errorf("new icon from image for dpi: %w", err)
			return unkImg
		}
		return icon
	}

	fmt.Println("unhandled extension", ext, unkImg)

	return unkImg
}

// Clear is used to flush an ico or generate cache
func Clear(name string) {
	_, err := assets.Open(fmt.Sprintf("assets/%s.ico", name))
	if err != nil {
		delete(icos, name)
		return
	}
}
