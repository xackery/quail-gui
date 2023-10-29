package gui

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/malashin/dds"
	"github.com/sergeymakinen/go-bmp"
	"github.com/xackery/wlk/walk"
)

func imagePreview(name string, ref interface{}) error {
	data, ok := ref.([]byte)
	if !ok {
		return fmt.Errorf("ref is not []byte, instead %T", ref)
	}

	r := bytes.NewReader(data)

	var dst *walk.Bitmap
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".dds":
		img, err := dds.Decode(r)
		if err != nil {
			return fmt.Errorf("dds decode: %w", err)
		}
		dst, err = walk.NewBitmapFromImageForDPI(img, 96)
		if err != nil {
			return fmt.Errorf("walk new bitmap from dds: %w", err)
		}
	case ".png":
		img, err := png.Decode(r)
		if err != nil {
			return fmt.Errorf("png decode: %w", err)
		}
		dst, err = walk.NewBitmapFromImageForDPI(img, 96)
		if err != nil {
			return fmt.Errorf("walk new bitmap from png: %w", err)
		}
	case ".bmp":
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(r)
		if err != nil {
			return fmt.Errorf("buf read from: %w", err)
		}
		var img image.Image
		if string(buf.Bytes()[0:3]) == "DDS" {
			img, err = dds.Decode(r)
			if err != nil {
				return fmt.Errorf("dds.Decode %s: %w", name, err)
			}
		} else {
			img, err = bmp.Decode(r)
			if err != nil {
				return fmt.Errorf("bmp.Decode %s: %w", name, err)
			}
		}

		dst, err = walk.NewBitmapFromImageForDPI(img, 96)
		if err != nil {
			return fmt.Errorf("new bitmap from image for dpi: %w", err)
		}
	default:
		return fmt.Errorf("unknown extension %s", ext)
	}
	SetImage(dst)
	page := gui.pageView.Pages()
	if page == nil {
		return fmt.Errorf("page not found")
	}
	page.At(1).SetFocus()

	return nil
}
