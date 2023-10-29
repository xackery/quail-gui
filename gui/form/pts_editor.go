package form

import (
	"bytes"
	"fmt"

	"github.com/xackery/quail-gui/archive"
	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/model/metadata/pts"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
	"github.com/xackery/wlk/wcolor"
)

var (
	ptsEditor *PtsEditor
)

type PtsEditor struct {
	node             *component.TreeNode
	firstError       error
	name             *walk.LineEdit
	nameError        *walk.Label
	boneName         *walk.LineEdit
	boneNameError    *walk.Label
	translationX     *walk.LineEdit
	translationY     *walk.LineEdit
	translationZ     *walk.LineEdit
	translationError *walk.Label
	rotationX        *walk.LineEdit
	rotationY        *walk.LineEdit
	rotationZ        *walk.LineEdit
	rotationError    *walk.Label
	scaleX           *walk.LineEdit
	scaleY           *walk.LineEdit
	scaleZ           *walk.LineEdit
	scaleError       *walk.Label
}

func showPtsEditor(page *walk.TabPage, node *component.TreeNode) (Editor, error) {
	_, ok := node.Ref().(*common.ParticlePointEntry)
	if !ok {
		return nil, fmt.Errorf("node is not a particle point entry")
	}

	_, ok = node.RootRef().(*common.ParticlePoint)
	if !ok {
		return nil, fmt.Errorf("root ref is not a particle point")
	}

	e := ptsEditor
	e.node = node
	e.Reset()

	return e, nil
}

func (e *PtsEditor) Save() error {
	slog.Println("Saving particle point entry")
	e.ClearError()
	e.validateName()
	e.validateBoneName()
	e.validateTranslationX()
	e.validateTranslationY()
	e.validateTranslationZ()
	e.validateRotationX()
	e.validateRotationY()
	e.validateRotationZ()
	e.validateScaleX()
	e.validateScaleY()
	e.validateScaleZ()

	if e.firstError != nil {
		return fmt.Errorf("validation failed: %w", e.firstError)
	}

	src, ok := e.node.Ref().(*common.ParticlePointEntry)
	if !ok {
		return fmt.Errorf("node is not a particle point entry")
	}

	base, ok := e.node.RootRef().(*common.ParticlePoint)
	if !ok {
		return fmt.Errorf("root ref is not a particle point")
	}

	src.Name = e.name.Text()
	src.BoneName = e.boneName.Text()
	src.Rotation.X = fastFloat32(e.rotationX.Text())
	src.Rotation.Y = fastFloat32(e.rotationY.Text())
	src.Rotation.Z = fastFloat32(e.rotationZ.Text())
	src.Scale.X = fastFloat32(e.scaleX.Text())
	src.Scale.Y = fastFloat32(e.scaleY.Text())
	src.Scale.Z = fastFloat32(e.scaleZ.Text())

	slog.Printf("Saving %+v\n", src)

	buf := bytes.NewBuffer(nil)
	err := pts.Encode(base, uint32(base.Header.Version), buf)
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

	e.node.SetName(src.Name)
	return nil
}

func (e *PtsEditor) Reset() {
	e.ClearError()
	src, ok := e.node.Ref().(*common.ParticlePointEntry)
	if !ok {
		return
	}
	e.name.SetText(src.Name)
	e.boneName.SetText(src.BoneName)
	e.translationX.SetText(fmt.Sprintf("%f", src.Translation.X))
	e.translationY.SetText(fmt.Sprintf("%f", src.Translation.Y))
	e.translationZ.SetText(fmt.Sprintf("%f", src.Translation.Z))
	e.rotationX.SetText(fmt.Sprintf("%f", src.Rotation.X))
	e.rotationY.SetText(fmt.Sprintf("%f", src.Rotation.Y))
	e.rotationZ.SetText(fmt.Sprintf("%f", src.Rotation.Z))
	e.scaleX.SetText(fmt.Sprintf("%f", src.Scale.X))
	e.scaleY.SetText(fmt.Sprintf("%f", src.Scale.Y))
	e.scaleZ.SetText(fmt.Sprintf("%f", src.Scale.Z))
}

func (e *PtsEditor) Node() *component.TreeNode {
	return e.node
}

func (e *PtsEditor) Ext() string {
	return ".pts"
}

func (e *PtsEditor) ClearError() {
	e.firstError = nil
}

func (e *PtsEditor) SetFirstError(err error) {
	if e.firstError != nil {
		return
	}
	e.firstError = err
}

func (e *PtsEditor) validateName() {
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

func (e *PtsEditor) validateBoneName() {
	var err error
	prop := "BoneName"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.boneNameError.SetText(err.Error())
			return
		}
		e.boneNameError.SetText("")
	}()
	value := e.boneName.Text()
	err = strValidate(value)
	if err != nil {
		return
	}
}

func (e *PtsEditor) validateTranslationX() {
	var err error
	prop := "Translation.X"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.translationError.SetText(err.Error())
			return
		}
		e.translationError.SetText("")
	}()
	value := e.translationX.Text()
	err = floatValidate(value)
	if err != nil {
		return
	}
}

func (e *PtsEditor) validateTranslationY() {
	var err error
	prop := "Translation.Y"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.translationError.SetText(err.Error())
			return
		}
		e.translationError.SetText("")
	}()
	value := e.translationY.Text()
	err = floatValidate(value)
	if err != nil {
		return
	}
}

func (e *PtsEditor) validateTranslationZ() {
	var err error
	prop := "Translation.Z"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.translationError.SetText(err.Error())
			return
		}
		e.translationError.SetText("")
	}()
	value := e.translationZ.Text()
	err = floatValidate(value)
	if err != nil {
		return
	}
}

func (e *PtsEditor) validateRotationX() {
	var err error
	prop := "Rotation.X"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.rotationError.SetText(err.Error())
			return
		}
		e.rotationError.SetText("")
	}()
	value := e.rotationX.Text()
	err = floatValidate(value)
	if err != nil {
		return
	}
}

func (e *PtsEditor) validateRotationY() {
	var err error
	prop := "Rotation.Y"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.rotationError.SetText(err.Error())
			return
		}
		e.rotationError.SetText("")
	}()
	value := e.rotationY.Text()
	err = floatValidate(value)
	if err != nil {
		return
	}
}

func (e *PtsEditor) validateRotationZ() {
	var err error
	prop := "Rotation.Z"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.rotationError.SetText(err.Error())
			return
		}
		e.rotationError.SetText("")
	}()
	value := e.rotationZ.Text()
	err = floatValidate(value)
	if err != nil {
		return
	}
}

func (e *PtsEditor) validateScaleX() {
	var err error
	prop := "Scale.X"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.scaleError.SetText(err.Error())
			return
		}
		e.scaleError.SetText("")
	}()
	value := e.scaleX.Text()
	err = floatValidate(value)
	if err != nil {
		return
	}
}

func (e *PtsEditor) validateScaleY() {
	var err error
	prop := "Scale.Y"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.scaleError.SetText(err.Error())
			return
		}
		e.scaleError.SetText("")
	}()
	value := e.scaleY.Text()
	err = floatValidate(value)
	if err != nil {
		return
	}
}

func (e *PtsEditor) validateScaleZ() {
	var err error
	prop := "Scale.Z"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", prop, err)
			e.SetFirstError(err)
			e.scaleError.SetText(err.Error())
			return
		}
		e.scaleError.SetText("")
	}()
	value := e.scaleZ.Text()
	err = floatValidate(value)
	if err != nil {
		return
	}
}

func PtsEditWidgets() []cpl.Widget {
	if ptsEditor == nil {
		ptsEditor = &PtsEditor{}
	}

	return []cpl.Widget{
		cpl.Composite{
			Layout: cpl.Grid{Columns: 2},
			Children: []cpl.Widget{
				cpl.Label{Text: "Name"},
				cpl.LineEdit{
					AssignTo:      &ptsEditor.name,
					OnTextChanged: ptsEditor.validateName,
				},
				cpl.Label{Text: ""},
				cpl.Label{Text: "", AssignTo: &ptsEditor.nameError, TextColor: wcolor.Red},
				cpl.Label{Text: "BoneName"},
				cpl.LineEdit{
					AssignTo:      &ptsEditor.boneName,
					OnTextChanged: ptsEditor.validateBoneName,
				},
				cpl.Label{Text: ""},
				cpl.Label{Text: "", AssignTo: &ptsEditor.boneNameError, TextColor: wcolor.Red},
				cpl.Composite{
					Layout: cpl.Grid{Columns: 4},
					Children: []cpl.Widget{
						cpl.Label{Text: "Translation"},
						cpl.LineEdit{AssignTo: &ptsEditor.translationX, OnTextChanged: ptsEditor.validateTranslationX},
						cpl.LineEdit{AssignTo: &ptsEditor.translationY, OnTextChanged: ptsEditor.validateTranslationY},
						cpl.LineEdit{AssignTo: &ptsEditor.translationZ, OnTextChanged: ptsEditor.validateTranslationZ},
						cpl.Label{Text: "", AssignTo: &ptsEditor.translationError, TextColor: wcolor.Red},
					},
				},
				cpl.Composite{
					Layout: cpl.Grid{Columns: 4},
					Children: []cpl.Widget{
						cpl.Label{Text: "Rotation"},
						cpl.LineEdit{AssignTo: &ptsEditor.rotationX, OnTextChanged: ptsEditor.validateRotationX},
						cpl.LineEdit{AssignTo: &ptsEditor.rotationY, OnTextChanged: ptsEditor.validateRotationY},
						cpl.LineEdit{AssignTo: &ptsEditor.rotationZ, OnTextChanged: ptsEditor.validateRotationZ},
						cpl.Label{Text: "", AssignTo: &ptsEditor.rotationError, TextColor: wcolor.Red},
					},
				},
				cpl.Composite{
					Layout: cpl.Grid{Columns: 4},
					Children: []cpl.Widget{
						cpl.Label{Text: "Scale"},
						cpl.LineEdit{AssignTo: &ptsEditor.scaleX, OnTextChanged: ptsEditor.validateScaleX},
						cpl.LineEdit{AssignTo: &ptsEditor.scaleY, OnTextChanged: ptsEditor.validateScaleY},
						cpl.LineEdit{AssignTo: &ptsEditor.scaleZ, OnTextChanged: ptsEditor.validateScaleZ},
						cpl.Label{Text: "", AssignTo: &ptsEditor.scaleError, TextColor: wcolor.Red},
					},
				},
			},
		},
	}
}

func (e *PtsEditor) New(src interface{}) (*component.TreeNode, error) {
	entry := &common.ParticlePointEntry{
		Name: "New Particle Point Entry",
	}
	srcEntry, ok := src.(*common.ParticlePointEntry)
	if ok {
		entry.Name = srcEntry.Name
		entry.BoneName = srcEntry.BoneName
		entry.Rotation = srcEntry.Rotation
		entry.Scale = srcEntry.Scale
	}

	base, ok := e.node.RootRef().(*common.ParticlePoint)
	if !ok {
		return nil, fmt.Errorf("root ref is not a particle point")
	}

	base.Entries = append(base.Entries, entry)
	slog.Printf("entries: %+v\n", base.Entries)
	node := e.node.Parent().(*component.TreeNode).ChildAdd(ico.Grab(".pts"), entry.Name, base, entry)
	return node, nil
}

func (e *PtsEditor) Name() string {
	return "Particle Point Entry"
}
