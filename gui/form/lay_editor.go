package form

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/xackery/quail-gui/archive"
	"github.com/xackery/quail-gui/gui/component"
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
	src           *common.Layer
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
	src, ok := node.Ref().(*common.Layer)
	if !ok {
		return nil, fmt.Errorf("node is not a layer")
	}

	e := layEditor

	e.src = src
	e.node = node
	e.diffuse.SetText(src.Diffuse)
	e.normal.SetText(src.Normal)
	e.material.SetText(src.Material)
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

	parent := e.node.Parent()
	if parent == nil {
		return fmt.Errorf("parent is nil")
	}
	parentTree, ok := parent.(*component.TreeNode)
	if !ok {
		return fmt.Errorf("parent is not *component.TreeNode, instead %T", parent)
	}
	_, ok = parentTree.Ref().([]*common.Layer)
	if !ok {
		return fmt.Errorf("parent is not []*common.Layer, instead %T", parentTree.Ref())
	}
	parent = parentTree.Parent()
	parentTree, ok = parent.(*component.TreeNode)
	if !ok {
		return fmt.Errorf("parent is not *component.TreeNode, instead %T", parent)
	}
	model, ok := parentTree.Ref().(*common.Model)
	if !ok {
		return fmt.Errorf("parent is not *common.Model, instead %T", parentTree.Ref())
	}

	e.src.Diffuse = e.diffuse.Text()
	e.src.Normal = e.normal.Text()
	e.src.Material = e.material.Text()

	slog.Printf("Saving %+v\n", e.src)
	slog.Printf("model: %+v\n", model)

	buf := bytes.NewBuffer(nil)
	err := lay.Encode(model, buf)
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

	e.node.SetName(e.src.Material)

	return nil
}

func (e *LayEditor) Reset() {
	e.ClearError()
	e.diffuse.SetText(e.src.Diffuse)
	e.normal.SetText(e.src.Normal)
	e.material.SetText(e.src.Material)
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
