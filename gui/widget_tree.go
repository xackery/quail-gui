package gui

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/gui/form"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
	"gopkg.in/yaml.v2"
)

func treeWidget() cpl.Widget {
	return cpl.TreeView{
		AssignTo: &gui.treeView,
		Visible:  true,
		Model:    gui.treeModel,
		ContextMenuItems: []cpl.MenuItem{
			cpl.Action{
				AssignTo:    &gui.newContext,
				Text:        "New Layer",
				OnTriggered: onNew,
			},
			cpl.Action{
				Text: "Refresh",
				OnTriggered: func() {
					gui.treeModel.ItemsReset()
				},
			},
		},
		OnCurrentItemChanged: func() {
			node := gui.treeView.CurrentItem().(*component.TreeNode)

			buf := bytes.NewBuffer(nil)
			enc := yaml.NewEncoder(buf)
			defer enc.Close()

			err := enc.Encode(node.Ref())
			if err != nil {
				slog.Printf("yaml encode: %s\n", err)
				return
			}

			gui.editSave.SetEnabled(false)
			gui.contents.SetText(strings.ReplaceAll(buf.String(), "\n", "\r\n"))
			editor, err := form.ShowEditor(gui.editView, node)
			if err != nil {
				slog.Printf("Failed editing: %s\n", err)
				for _, view := range editViews {
					if view == nil {
						continue
					}
					view.SetVisible(false)
				}
				editViews[".unk"].SetVisible(true)
				gui.titleEdit.SetText(fmt.Sprintf("Select a valid node to edit on left. %T is not yet supported", node.Ref()))
				SetPageVisible(false, false, false)
				return
			}
			ext := editor.Ext()
			editView, ok := editViews[ext]
			if !ok {
				slog.Printf("Failed finding edit view %s\n", ext)
				gui.titleEdit.SetText(fmt.Sprintf("Select a valid node to edit on left. %T is not yet supported", node.Ref()))
				editViews[".unk"].SetVisible(true)
				SetPageVisible(false, false, false)
				return
			}
			if editView == nil {
				slog.Printf("Failed edit view is nil %s\n", ext)
				gui.titleEdit.SetText(fmt.Sprintf("Select a valid node to edit on left. %T is not yet supported", node.Ref()))
				editViews[".unk"].SetVisible(true)
				SetPageVisible(false, false, false)
				return
			}
			for _, view := range editViews {
				if view == nil {
					continue
				}
				view.SetVisible(false)
			}
			activeEditor = editor
			gui.newContext.SetText(fmt.Sprintf("New %s", activeEditor.Name()))
			SetPageVisible(activeEditor.IsEdit(), activeEditor.IsYaml(), activeEditor.IsPreview())
			gui.titleEdit.SetText(fmt.Sprintf("Editing node %T. Press Save to apply changes.", node.Ref()))
			editView.SetVisible(true)
			gui.editSave.SetEnabled(true)
		},
		OnKeyDown: func(key walk.Key) {
			if key == walk.KeyDelete {
				node := gui.treeView.CurrentItem().(*component.TreeNode)
				if node == nil {
					slog.Printf("Failed to delete: node is nil\n")
					return
				}
				if node.Parent() == nil {
					slog.Printf("Failed to delete: node parent is nil\n")
					return
				}
				node.Parent().RemoveChild(node)
				//gui.treeModel.PublishItemsReset(node.Parent())
				slog.Printf("Deleted %+v\n", node)
				return
			}
			if key == walk.KeyF5 {
				node := gui.treeView.CurrentItem().(*component.TreeNode)
				if node == nil {
					slog.Printf("Failed to refresh: node is nil\n")
					return
				}
				gui.treeModel.PublishItemsReset(node)
				slog.Printf("Refreshed %+v\n", node)
				return
			}
			if key == walk.KeyInsert {
				onNew()
				return
			}
			if key == walk.KeyV && walk.ControlDown() {
				if activeEditor == nil {
					slog.Printf("New failed: editor is nil\n")
					return
				}
				node, err := activeEditor.New(gui.treeCopy.Ref())
				if err != nil {
					slog.Printf("New failed: %s\n", err)
					return
				}

				gui.treeView.SetCurrentItem(node)
				slog.Println("Added new node")

				if gui.treeCopy == nil {
					slog.Printf("Failed to paste: copy is nil\n")
					return
				}

				slog.Printf("Pasted %+v\n", node)
				return
			}
			if key == walk.KeyC && walk.ControlDown() {
				node := gui.treeView.CurrentItem().(*component.TreeNode)
				if node == nil {
					slog.Printf("Failed to copy: node is nil\n")
					return
				}
				gui.treeCopy = node
				slog.Printf("Copied %+v\n", node)
				return
			}

		},
		MinSize: cpl.Size{Width: 250, Height: 0},
	}
}
