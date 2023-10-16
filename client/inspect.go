package client

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/xackery/quail-gui/gui"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/mesh/mds"
	"github.com/xackery/quail/model/mesh/mod"
	"github.com/xackery/quail/model/metadata/lay"
	"github.com/xackery/quail/model/metadata/prt"
	"github.com/xackery/quail/model/metadata/pts"
	"github.com/xackery/quail/model/metadata/zon"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/quail/quail"
)

func (c *Client) inspect(file string) (interface{}, error) {
	if len(file) < 2 {
		slog.Printf("Inspecting %s\n", filepath.Base(c.currentPath))
	} else {
		slog.Printf("Inspecting %s %s\n", filepath.Base(c.currentPath), filepath.Base(file))
	}

	if len(file) < 2 {
		entries := []*gui.FileViewEntry{}
		for _, fe := range c.pfs.Files() {
			entries = append(entries, &gui.FileViewEntry{
				Name: fe.Name(),
				Ext:  strings.ToLower(filepath.Ext(fe.Name())),
				Size: generateSize(len(fe.Data())),
			})
		}
		gui.SetFileViewItems(entries)
		return c.pfs, nil
	}
	return c.inspectFile(c.pfs, c.currentPath, file)
}

func (c *Client) inspectFile(pfs *pfs.PFS, path string, file string) (interface{}, error) {
	if pfs == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return c.inspectContent(filepath.Base(file), bytes.NewReader(data))
	}

	for _, fe := range pfs.Files() {
		if !strings.EqualFold(fe.Name(), file) {
			continue
		}
		return c.inspectContent(file, bytes.NewReader(fe.Data()))
	}
	return nil, fmt.Errorf("%s not found in %s", file, filepath.Base(path))
}

func (c *Client) inspectContent(file string, data *bytes.Reader) (interface{}, error) {
	var err error
	ext := strings.ToLower(filepath.Ext(file))
	switch ext {
	case ".mds":
		model := &common.Model{
			Name: strings.TrimSuffix(strings.ToUpper(file), ".MDS"),
		}
		err = mds.Decode(model, data)
		if err != nil {
			return nil, fmt.Errorf("mds.Decode %s: %w", file, err)
		}
		return model, nil
	case ".mod":
		model := &common.Model{
			Name: strings.TrimSuffix(strings.ToUpper(file), ".MOD"),
		}
		err = mod.Decode(model, data)
		if err != nil {
			return nil, fmt.Errorf("mod.Decode %s: %w", file, err)
		}
		return model, nil
	case ".pts":
		point := &common.ParticlePoint{
			Name: strings.TrimSuffix(strings.ToUpper(file), ".MDS"),
		}
		err = pts.Decode(point, data)
		if err != nil {
			return nil, fmt.Errorf("pts.Decode %s: %w", file, err)
		}
		return point, nil
	case ".prt":
		render := &common.ParticleRender{
			Name: strings.TrimSuffix(strings.ToUpper(file), ".MDS"),
		}
		err = prt.Decode(render, data)
		if err != nil {
			return nil, fmt.Errorf("prt.Decode %s: %w", file, err)
		}
		return render, nil
	case ".zon":
		zone := &common.Zone{
			Name: strings.TrimSuffix(strings.ToUpper(file), ".ZON"),
		}
		err = zon.Decode(zone, data)
		if err != nil {
			return nil, fmt.Errorf("zon.Decode %s: %w", file, err)
		}
		return zone, nil
	case ".wld":
		models, err := quail.WLDDecode(data, nil)
		if err != nil {
			return nil, fmt.Errorf("wld.Decode %s: %w", file, err)
		}
		return models, nil
	case ".lay":
		model := &common.Model{}
		err := lay.Decode(model, data)
		if err != nil {
			return nil, fmt.Errorf("lay.Decode %s: %w", file, err)
		}
		return model.Layers, nil
	default:
		return nil, fmt.Errorf("unknown file type %s", ext)
	}
}

func (c *Client) reflectTraversal(inspected interface{}, section string, nest int, index int) {
	v := reflect.ValueOf(inspected)
	tv := v.Type()

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		tv = v.Type()
	}

	if v.Kind() == reflect.Slice {

		if v.Len() == 0 {
			c.sections[section].Content += fmt.Sprintf("%s%s: (Empty)\n", strings.Repeat("  ", nest), tv.Name())
			//slog.Printf("%s%s (Empty)\n", strings.Repeat("  ", nest), tv.Name())
			return
		}
		c.sections[section].Content += fmt.Sprintf("%s%s:\n", strings.Repeat("  ", nest), tv.Name())
		//slog.Printf("%s%s:", strings.Repeat("  ", nest), tv.Name())
		for i := 0; i < v.Len(); i++ {
			c.reflectTraversal(v.Index(i).Interface(), section, nest+1, i)
		}
		return
	}

	if v.Kind() != reflect.Struct {
		c.sections[section].Content += fmt.Sprintf("%s%v\n", strings.Repeat("  ", nest), v.Interface())
		//slog.Printf("%s%v\n", strings.Repeat("  ", nest), v.Interface())
		return
	}

	for i := 0; i < v.NumField(); i++ {
		// check if field is exported
		if tv.Field(i).PkgPath != "" {
			continue
		}

		indexStr := ""
		if index >= 0 {
			indexStr = fmt.Sprintf("[%d]", index)
		}

		// is it a slice?
		if v.Field(i).Kind() == reflect.Slice {
			if nest == 0 {
				slog.Printf("Changing section to %s\n", tv.Field(i).Name)
				section = tv.Field(i).Name
				_, ok := c.sections[section]
				if !ok {
					c.sections[section] = &gui.Section{
						Name:  section,
						Count: 0,
					}
				}
			}
			s := v.Field(i)
			if s.Len() == 0 {
				c.sections[section].Content += fmt.Sprintf("%s%s %s: (Empty)\n", strings.Repeat("  ", nest), indexStr, tv.Field(i).Name)
				//slog.Printf("%s%s %s: (Empty)\n", strings.Repeat("  ", nest), indexStr, tv.Field(i).Name)
				continue
			}
			c.sections[section].Count = s.Len()
			c.sections[section].Content += fmt.Sprintf("%s%s %s:\n", strings.Repeat("  ", nest), indexStr, tv.Field(i).Name)
			//slog.Printf("%s%s %s:", strings.Repeat("  ", nest), indexStr, tv.Field(i).Name)

			for j := 0; j < s.Len(); j++ {
				if tv.Field(i).PkgPath != "" {
					continue
				}

				if s.Index(j).Kind() == reflect.Uint8 {
					if j == 0 {
						c.sections[section].Content += fmt.Sprintf("%s", strings.Repeat("  ", nest+1))
						//slog.Printf("%s", strings.Repeat("  ", nest+1))
					}
					//fmt.Printf("0x%02x ", s.Index(j).Interface())
					if j == s.Len()-1 {
						//slog.Println()
					}
					continue
				}
				c.reflectTraversal(s.Index(j).Interface(), section, nest+1, j)
				//slog.Printf("  %d %s\t %+v", j, tv.Field(i).Name, s.Index(j).Interface())
			}
			continue
		}

		if tv.Field(i).Name == "MaterialName" {
			continue
		}
		c.sections[section].Content += fmt.Sprintf("%s%s %s: %v\n", strings.Repeat("  ", nest), indexStr, tv.Field(i).Name, v.Field(i).Interface())
		//slog.Printf("%s%s %s: %v\n", strings.Repeat("  ", nest), indexStr, tv.Field(i).Name, v.Field(i).Interface())
	}
}

func generateSize(in int) string {
	val := float64(in)
	if val < 1024 {
		return fmt.Sprintf("%0.0f bytes", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.0f KB", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.0f MB", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.0f GB", val)
	}
	val /= 1024
	return fmt.Sprintf("%0.0f TB", val)
}
