package dialog

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/wld/virtual"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualSpritePage(data *virtual.Wld, page *cpl.TabPage) error {

	sprites := []string{}
	for _, sprite := range data.Sprites {
		sprites = append(sprites, sprite.Tag)
	}
	onSpriteNew := func() {
		slog.Println("new sprite")
	}
	onSpriteEdit := func() {
		slog.Println("edit sprite")
	}
	onSpriteDelete := func() {
		slog.Println("delete sprite")
	}

	var cmbSprite *walk.ComboBox
	defaultSprite := ""
	if len(sprites) > 0 {
		defaultSprite = sprites[0]
	}

	spriteGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	spriteGroup.Children = append(spriteGroup.Children, cpl.GroupBox{
		Title:  "Sprites (SimpleSpriteDef)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbSprite,
				Editable: false,
				Model:    sprites,
				Value:    defaultSprite,
			},
			cpl.PushButton{Text: "Add", OnClicked: onSpriteNew},
			cpl.PushButton{Text: "Edit", OnClicked: onSpriteEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onSpriteDelete},
		},
	})

	spriteInstances := []string{}
	for _, spriteInstance := range data.SpriteInstances {
		spriteInstances = append(spriteInstances, spriteInstance.Tag)
	}
	onSpriteInstanceNew := func() {
		slog.Println("new spriteInstance")
	}
	onSpriteInstanceEdit := func() {
		slog.Println("edit spriteInstance")
	}
	onSpriteInstanceDelete := func() {
		slog.Println("delete spriteInstance")
	}

	var cmbSpriteInstance *walk.ComboBox
	defaultSpriteInstance := ""
	if len(spriteInstances) > 0 {
		defaultSpriteInstance = spriteInstances[0]
	}

	spriteGroup.Children = append(spriteGroup.Children, cpl.GroupBox{
		Title:  "SpriteInstances (SimpleSprite)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbSpriteInstance,
				Editable: false,
				Model:    spriteInstances,
				Value:    defaultSpriteInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onSpriteInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onSpriteInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onSpriteInstanceDelete},
		},
	})

	page.Title = "Sprite"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{spriteGroup}
	return nil
}
