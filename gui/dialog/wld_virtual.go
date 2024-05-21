package dialog

import (
	"fmt"

	"github.com/xackery/quail-gui/popup"
	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/vwld"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func ShowWldVirtualEdit(mw *walk.MainWindow, title string, src raw.ReadWriter) error {
	var savePB, cancelPB *walk.PushButton
	formElements := cpl.Composite{
		Layout:   cpl.VBox{},
		Children: []cpl.Widget{},
	}

	rawData, ok := src.(*raw.Wld)
	if !ok {
		return fmt.Errorf("cast wld")
	}

	data := &vwld.VWld{}
	err := data.Read(rawData)
	if err != nil {
		return fmt.Errorf("read wld: %w", err)
	}

	tabWidget := cpl.TabWidget{}

	headerPage := &cpl.TabPage{}
	err = virtualHeaderPage(data, headerPage)
	if err != nil {
		return fmt.Errorf("header tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *headerPage)

	texturePage := &cpl.TabPage{}
	err = virtualTexturePage(data, texturePage)
	if err != nil {
		return fmt.Errorf("texture tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *texturePage)

	spritePage := &cpl.TabPage{}
	err = virtualSpritePage(data, spritePage)
	if err != nil {
		return fmt.Errorf("sprite tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *spritePage)

	particlePage := &cpl.TabPage{}
	err = virtualParticlePage(data, particlePage)
	if err != nil {
		return fmt.Errorf("particle tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *particlePage)

	materialPage := &cpl.TabPage{}
	err = virtualMaterialPage(data, materialPage)
	if err != nil {
		return fmt.Errorf("material tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *materialPage)

	meshPage := &cpl.TabPage{}
	err = virtualMeshPage(data, meshPage)
	if err != nil {
		return fmt.Errorf("mesh tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *meshPage)

	actorPage := &cpl.TabPage{}
	err = virtualActorPage(data, actorPage)
	if err != nil {
		return fmt.Errorf("actor tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *actorPage)

	animationPage := &cpl.TabPage{}
	err = virtualAnimationPage(data, animationPage)
	if err != nil {
		return fmt.Errorf("animation tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *animationPage)

	skeletonPage := &cpl.TabPage{}
	err = virtualSkeletonPage(data, skeletonPage)
	if err != nil {
		return fmt.Errorf("skeleton tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *skeletonPage)

	lightPage := &cpl.TabPage{}
	err = virtualLightPage(data, lightPage)
	if err != nil {
		return fmt.Errorf("light tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *lightPage)

	regionPage := &cpl.TabPage{}
	err = virtualRegionPage(data, regionPage)
	if err != nil {
		return fmt.Errorf("region tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *regionPage)

	bspPage := &cpl.TabPage{}
	err = virtualBspPage(data, bspPage)
	if err != nil {
		return fmt.Errorf("bsp tab: %w", err)
	}
	tabWidget.Pages = append(tabWidget.Pages, *bspPage)

	formElements.Children = append(formElements.Children, tabWidget)

	onSave := func() error {
		/* newVersionStr := cmbVersion.Text()
		if newVersionStr == "" {
			return fmt.Errorf("version is required")
		}
		newVersion, err := strconv.Atoi(newVersionStr)
		if err != nil {
			return fmt.Errorf("parse version: %w", err)
		}
		if data.Version != uint32(newVersion) {
			data.Version = uint32(newVersion)
		}
		slog.Println("new version:", data.Version)
		*/
		return nil
	}

	var dlg *walk.Dialog
	dia := cpl.Dialog{
		AssignTo:      &dlg,
		Title:         fmt.Sprintf("%s (Virtual)", data.FileName),
		DefaultButton: &savePB,
		CancelButton:  &cancelPB,
		MinSize:       cpl.Size{Width: 300, Height: 300},
		Layout:        cpl.VBox{},
		Children: []cpl.Widget{
			formElements,
			cpl.Composite{
				Layout: cpl.HBox{},
				Children: []cpl.Widget{
					cpl.HSpacer{},
					cpl.PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
					cpl.PushButton{
						AssignTo: &savePB,
						Text:     "Save",
						OnClicked: func() {
							err := onSave()
							if err != nil {
								popup.Errorf(dlg, "save: %s", err.Error())
								return
							}
							dlg.Accept()
						},
					},
				},
			},
		},
	}
	result, err := dia.Run(mw)
	if err != nil {
		return fmt.Errorf("run dialog: %w", err)
	}

	if result != walk.DlgCmdOK {
		return fmt.Errorf("cancelled")
	}

	return nil
}
