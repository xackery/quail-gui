package ico

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"io"
	"strings"

	ico "github.com/biessek/golang-ico"
	"github.com/malashin/dds"
	"github.com/sergeymakinen/go-bmp"
	"github.com/xackery/wlk/walk"
	"golang.org/x/image/draw"
)

var (
	//go:embed assets
	assets embed.FS
	icos   map[string]*walk.Icon
)

func Init() error {
	icos = make(map[string]*walk.Icon)

	// first load unk icon
	r, err := assets.Open("assets/unk.ico")
	if err != nil {
		return fmt.Errorf("open unk.ico: %w", err)
	}
	icoData, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read unk.ico: %w", err)
	}
	icon, err := Generate("unk", icoData)
	if err != nil {
		return fmt.Errorf("generate unk: %w", err)
	}
	icos["unk"] = icon

	dir, err := assets.ReadDir("assets")
	if err != nil {
		return fmt.Errorf("read assets: %w", err)
	}

	for _, fi := range dir {
		name := fi.Name()
		r, err := assets.Open(fmt.Sprintf("assets/%s", name))
		if err != nil {
			return fmt.Errorf("open %s: %w", name, err)
		}
		icoData, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		if strings.Contains(name, ".") {
			name = name[0:strings.Index(name, ".")]
		}
		icon, err := Generate(name, icoData)
		if err != nil {
			return fmt.Errorf("generate %s: %w", name, err)
		}

		name = strings.ToLower(name)

		icos[name] = icon
	}
	return nil
}

// Grab returns a walk.Icon for a given icon
func Grab(name string) *walk.Icon {
	icon, ok := icos[name]
	if !ok {
		return icos["unk"]
	}
	return icon
}

func Generate(name string, data []byte) (*walk.Icon, error) {
	var err error

	if len(name) == 0 {
		return nil, fmt.Errorf("name is empty")
	}
	if name[0] == '.' {
		name = name[1:]
	}

	icon, ok := icos[name]
	if ok && !isImageExt(name) && name != "dds" {
		return icon, nil
	}

	unkImg := Grab("unk")

	var img image.Image

	generators := map[string]func([]byte) (image.Image, error){
		"ico": icoGen,
		"dds": ddsGen,
		"png": pngGen,
		"bmp": bmpGen,
	}

	for _, gen := range generators {
		img, err = gen(data)
		if err == nil {
			break
		}
	}
	if err != nil {
		return unkImg, fmt.Errorf("generate %s: %w", name, err)
	}
	icon, err = walk.NewIconFromImageForDPI(img, 96)
	if err != nil {
		return nil, fmt.Errorf("new icon from image for dpi: %w", err)
	}

	return icon, nil
}

// Clear is used to flush an ico or generate cache
func Clear(name string) {
	_, err := assets.Open(fmt.Sprintf("assets/%s.ico", name))
	if err != nil {
		delete(icos, name)
		return
	}
}

func isImageExt(ext string) bool {
	switch ext {
	case "ico", "dds", "png", "bmp":
		return true
	}
	return false
}

func icoGen(data []byte) (image.Image, error) {
	icoReader := bytes.NewReader(data)
	img, err := ico.Decode(icoReader)
	if err != nil {
		return nil, fmt.Errorf("ico.Decode: %w", err)
	}
	if img.Bounds().Max.X > 16 || img.Bounds().Max.Y > 16 {
		dst := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X/2, img.Bounds().Max.Y/2))
		draw.NearestNeighbor.Scale(dst, image.Rect(0, 0, 16, 16), img, img.Bounds(), draw.Over, nil)
		img = dst
	}
	return img, nil
}

func ddsGen(data []byte) (image.Image, error) {
	img, err := dds.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("dds.Decode %w", err)
	}
	dst := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X/2, img.Bounds().Max.Y/2))
	draw.NearestNeighbor.Scale(dst, image.Rect(0, 0, 16, 16), img, img.Bounds(), draw.Over, nil)

	return dst, nil
}

func pngGen(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("image.Decode %w", err)
	}
	if img.Bounds().Max.X > 16 || img.Bounds().Max.Y > 16 {
		dst := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X/2, img.Bounds().Max.Y/2))
		draw.NearestNeighbor.Scale(dst, image.Rect(0, 0, 16, 16), img, img.Bounds(), draw.Over, nil)
		img = dst
	}

	return img, nil
}

func bmpGen(data []byte) (image.Image, error) {
	var err error
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("buf read from: %w", err)
	}
	var img image.Image
	if string(buf.Bytes()[0:3]) == "DDS" {
		img, err = dds.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("dds.Decode %w", err)
		}
	} else {
		img, err = bmp.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("bmp.Decode  %w", err)
		}
	}
	dst := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X/2, img.Bounds().Max.Y/2))
	draw.NearestNeighbor.Scale(dst, image.Rect(0, 0, 16, 16), img, img.Bounds(), draw.Over, nil)

	return dst, nil
}
