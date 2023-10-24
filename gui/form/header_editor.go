package form

import (
	"fmt"
	"strconv"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail/common"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
	"github.com/xackery/wlk/wcolor"
)

var (
	headerEditor *HeaderEditor
)

type HeaderEditor struct {
	base         interface{}
	src          *common.Header
	node         *component.TreeNode
	firstError   error
	name         *walk.LineEdit
	nameError    *walk.Label
	version      *walk.LineEdit
	versionError *walk.Label
	versionMax   int
}

func showHeaderEditor(page *walk.TabPage, node *component.TreeNode) (Editor, error) {
	src, ok := node.Ref().(*common.Header)
	if !ok {
		return nil, fmt.Errorf("node is not a header")
	}

	e := headerEditor

	e.src = src
	e.node = node
	e.name.SetText(src.Name)
	e.version.SetText(fmt.Sprintf("%d", src.Version))

	parent := e.node.Parent()
	if parent == nil {
		return nil, fmt.Errorf("parent is nil")
	}
	parentTree, ok := parent.(*component.TreeNode)
	if !ok {
		return nil, fmt.Errorf("parent is not *component.TreeNode, instead %T", parent)
	}
	e.base = parentTree.Ref()
	switch parentTree.Ref().(type) {
	case *common.Wld:
		e.versionMax = 0x1000C800
	default:
		e.versionMax = 3
	}
	return e, nil
}

func (e *HeaderEditor) Save() error {
	e.ClearError()
	e.validateName()
	e.validateVersion()

	if e.firstError != nil {
		return fmt.Errorf("validation failed: %w", e.firstError)
	}

	e.src.Name = e.name.Text()
	val, err := strconv.Atoi(e.version.Text())
	if err != nil {
		return fmt.Errorf("version is not a number")
	}
	e.src.Version = val
	return nil
}

func (e *HeaderEditor) Reset() {
	e.ClearError()
	e.name.SetText(e.src.Name)
	e.version.SetText(fmt.Sprintf("%d", e.src.Version))
}

func (e *HeaderEditor) Node() *component.TreeNode {
	return e.node
}

func (e *HeaderEditor) Ext() string {
	return "header"
}

func (e *HeaderEditor) ClearError() {
	e.firstError = nil
}

func (e *HeaderEditor) SetFirstError(err error) {
	if e.firstError != nil {
		return
	}
	e.firstError = err
}

func (e *HeaderEditor) validateName() {
	var err error
	prop := "Name"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.nameError.SetText(err.Error())
			return
		}
		e.nameError.SetText("")
	}()
	value := e.name.Text()
	err = strValidate(value)
	if err != nil {
		return
	}
}

func (e *HeaderEditor) validateVersion() {
	var err error
	prop := "Version"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.versionError.SetText(err.Error())
			return
		}
		e.versionError.SetText("")
	}()
	value := e.version.Text()
	err = intValidate(value, 0, e.versionMax)
	if err != nil {
		return
	}
}

func HeaderEditWidgets() []cpl.Widget {
	if headerEditor == nil {
		headerEditor = &HeaderEditor{
			versionMax: 3,
		}
	}
	return []cpl.Widget{
		cpl.Composite{
			Layout: cpl.Grid{Columns: 2},
			Children: []cpl.Widget{
				cpl.Label{Text: "Name"},
				cpl.LineEdit{
					AssignTo:      &headerEditor.name,
					OnTextChanged: headerEditor.validateName,
				},
				cpl.Label{Text: ""},
				cpl.Label{Text: "", AssignTo: &headerEditor.nameError, TextColor: wcolor.Red},
				cpl.Label{Text: "Version"},
				cpl.LineEdit{
					AssignTo:      &headerEditor.version,
					OnTextChanged: headerEditor.validateVersion,
				},
				cpl.Label{Text: ""},
				cpl.Label{Text: "", AssignTo: &headerEditor.versionError, TextColor: wcolor.Red},
			},
		},
	}
}

func (e *HeaderEditor) IsPreview() bool {
	return false
}

func (e *HeaderEditor) IsYaml() bool {
	return true
}

func (e *HeaderEditor) IsEdit() bool {
	return true
}

func (e *HeaderEditor) New(src interface{}) (*component.TreeNode, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e *HeaderEditor) Name() string {
	return "Header"
}
