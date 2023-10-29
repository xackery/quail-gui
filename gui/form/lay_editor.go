package form

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/xackery/quail-gui/archive"
	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/lay"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
	"github.com/xackery/wlk/wcolor"
)

var (
	layEditor *LayEditor
	layReg    = regexp.MustCompile(`_S\d\d_M\d\d$`)
)

type LayEditor struct {
	node          *component.TreeNode
	firstError    error
	material      *walk.LineEdit
	materialError *walk.Label
	diffuse       *walk.LineEdit
	diffuseError  *walk.Label
	normal        *walk.LineEdit
	normalError   *walk.Label
}

func showLayEditor(page *walk.TabPage, node *component.TreeNode) (Editor, error) {
	_, ok := node.Ref().(*common.Layer)
	if !ok {
		return nil, fmt.Errorf("node is not a layer")
	}
	_, ok = node.RootRef().(*common.Model)
	if !ok {
		return nil, fmt.Errorf("root ref is not a model")
	}
	e := layEditor
	e.node = node
	e.Reset()
	return e, nil
}

func (e *LayEditor) Save() error {
	slog.Println("Saving layer")
	e.ClearError()
	e.validateMaterial()
	e.validateDiffuse()
	e.validateNormal()

	if e.firstError != nil {
		return fmt.Errorf("validation failed: %w", e.firstError)
	}

	src, ok := e.node.Ref().(*common.Layer)
	if !ok {
		return fmt.Errorf("node is not a layer")
	}
	base, ok := e.node.RootRef().(*common.Model)
	if !ok {
		return fmt.Errorf("root ref is not a model")
	}

	src.Diffuse = e.diffuse.Text()
	src.Normal = e.normal.Text()
	src.Material = e.material.Text()

	slog.Printf("Saving %+v\n", src)
	slog.Printf("model: %+v\n", base)

	buf := bytes.NewBuffer(nil)
	err := lay.Encode(base, buf)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	err = archive.SetFile("", buf.Bytes())
	if err != nil {
		return fmt.Errorf("set file: %w", err)
	}

	err = archive.Save("")
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	e.node.SetName(src.Material)

	return nil
}

func (e *LayEditor) Reset() {
	e.ClearError()
	src, ok := e.node.Ref().(*common.Layer)
	if !ok {
		return
	}

	e.diffuse.SetText(src.Diffuse)
	e.normal.SetText(src.Normal)
	e.material.SetText(src.Material)
}

func (e *LayEditor) Node() *component.TreeNode {
	return e.node
}

func (e *LayEditor) Ext() string {
	return ".lay"
}

func (e *LayEditor) ClearError() {
	e.firstError = nil
}

func (e *LayEditor) SetFirstError(err error) {
	if e.firstError != nil {
		return
	}
	e.firstError = err
}

func (e *LayEditor) validateMaterial() {
	var err error
	prop := "Material"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.materialError.SetText(err.Error())
			return
		}
		e.materialError.SetText("")
	}()
	value := e.material.Text()
	err = strValidate(value)
	if err != nil {
		return
	}
	if !layReg.MatchString(value) {
		err = fmt.Errorf("is not a valid material, needs _S##_M## pattern suffix")
		return
	}
}

func (e *LayEditor) validateDiffuse() {
	var err error
	prop := "Diffuse"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.diffuseError.SetText(err.Error())
			return
		}
		e.diffuseError.SetText("")
	}()
	value := e.diffuse.Text()
	err = strValidate(value)
	if err != nil {
		return
	}
}

func (e *LayEditor) validateNormal() {
	prop := "Normal"
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.normalError.SetText(err.Error())
			return
		}
		e.normalError.SetText("")
	}()
	value := e.normal.Text()
	err = strValidate(value)
	if err != nil {
		return
	}
}

func LayEditWidgets() []cpl.Widget {
	if layEditor == nil {
		layEditor = &LayEditor{}
	}
	return []cpl.Widget{
		cpl.Composite{
			Layout: cpl.Grid{Columns: 2},
			Children: []cpl.Widget{
				cpl.Label{Text: "Material"},
				cpl.LineEdit{
					AssignTo:      &layEditor.material,
					OnTextChanged: layEditor.validateMaterial,
				},
				cpl.Label{Text: ""},
				cpl.Label{Text: "", AssignTo: &layEditor.materialError, TextColor: wcolor.Red},
				cpl.Label{Text: "Diffuse"},
				cpl.LineEdit{
					AssignTo:      &layEditor.diffuse,
					OnTextChanged: layEditor.validateDiffuse,
				},
				cpl.Label{Text: ""},
				cpl.Label{Text: "", AssignTo: &layEditor.diffuseError, TextColor: wcolor.Red},
				cpl.Label{Text: "Normal"},
				cpl.LineEdit{
					AssignTo:      &layEditor.normal,
					OnTextChanged: layEditor.validateNormal,
				},
				cpl.Label{Text: ""},
				cpl.Label{Text: "", AssignTo: &layEditor.normalError, TextColor: wcolor.Red},
			},
		},
	}
}

func (e *LayEditor) New(src interface{}) (*component.TreeNode, error) {
	layer := &common.Layer{
		Material: "New Material",
	}
	srcLayer, ok := src.(*common.Layer)
	if ok {
		layer.Material = srcLayer.Material
		layer.Diffuse = srcLayer.Diffuse
		layer.Normal = srcLayer.Normal
	}

	base, ok := e.node.RootRef().(*common.Model)
	if !ok {
		return nil, fmt.Errorf("root ref is not a model")
	}

	base.Layers = append(base.Layers, layer)
	slog.Printf("layers: %+v\n", base.Layers)
	node := e.node.Parent().(*component.TreeNode).ChildAdd(ico.Grab(".lay"), layer.Material, base, layer)
	return node, nil
}

func (e *LayEditor) Name() string {
	return "Layer"
}
