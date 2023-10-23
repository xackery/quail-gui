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
	src          *common.Header
	node         *component.TreeNode
	firstError   error
	name         *walk.LineEdit
	nameError    *walk.Label
	version      *walk.LineEdit
	versionError *walk.Label
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
	err = intValidate(value, 0, 3)
	if err != nil {
		return
	}
}

func HeaderEditWidgets() []cpl.Widget {
	if headerEditor == nil {
		headerEditor = &HeaderEditor{}
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
