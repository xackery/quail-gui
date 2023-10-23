//go:build windows
// +build windows

package gui

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/xackery/quail-gui/config"
	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/gui/form"
	"github.com/xackery/quail-gui/gui/handler"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
	"github.com/xackery/wlk/wcolor"
	"gopkg.in/yaml.v2"
)

type Gui struct {
	ctx            context.Context
	cancel         context.CancelFunc
	mw             *walk.MainWindow
	progress       *walk.ProgressBar
	log            *walk.TextEdit
	table          *walk.TableView
	exportSelected *walk.Action
	fileView       *component.FileView
	fileEntries    []*component.FileViewEntry
	contentsLabel  *walk.Label
	contents       *walk.TextEdit
	image          *walk.ImageView
	imageLabel     *walk.Label
	statusBar      *walk.StatusBarItem
	treeView       *walk.TreeView
	treeModel      *component.TreeModel
	treeCopy       *component.TreeNode
	editView       *walk.TabPage
	pageView       *walk.TabWidget
	yamlView       *walk.TabPage
	titleEdit      *walk.Label
	editSave       *walk.PushButton
}

var (
	gui       *Gui
	editViews = map[string]*walk.Composite{
		".lay":    nil,
		".pts":    nil,
		"header":  nil,
		"zon_obj": nil,
		".unk":    nil,
		".mod":    nil,
	}
	activeEditor form.Editor
)

// NewMainWindow creates a new main window
func NewMainWindow(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, version string) error {
	//walk.SetDarkModeAllowed(true)
	gui = &Gui{
		ctx:    ctx,
		cancel: cancel,
	}

	var err error
	gui.fileView = component.NewFileView()
	fvs := component.NewFileViewStyler(gui.fileView)
	gui.treeModel = component.NewTreeModel()

	layEdit := editViews[".lay"]
	ptsEdit := editViews[".pts"]
	headerEdit := editViews["header"]
	zonObjEdit := editViews["zon_obj"]
	unkEdit := editViews[".unk"]
	modEdit := editViews[".mod"]

	slog.AddHandler(Logf)
	handler.EditSaveSubscribe(onEditSave)
	handler.EditResetSubscribe(onEditReset)

	var filesSplit *walk.Splitter

	cmw := cpl.MainWindow{
		Title:   "quail-gui v" + version,
		MinSize: cpl.Size{Width: 300, Height: 300},
		Layout:  cpl.VBox{},
		Visible: false,
		Name:    "quail-gui",
		MenuItems: []cpl.MenuItem{
			cpl.Menu{
				Text: "&Archive",
				Items: []cpl.MenuItem{
					cpl.Action{
						Text: "&New",
						OnTriggered: func() {
							handler.ArchiveNewInvoke()
						},
						Shortcut: cpl.Shortcut{
							Key:       walk.KeyN,
							Modifiers: walk.ModControl,
						},
					},
					cpl.Action{
						Text: "&Open",
						OnTriggered: func() {

							path, err := ShowOpen("Open EQ Archive", "All Archives|*.pfs;*.eqg;*.s3d;*.pak|PFS Files (*.pfs)|*.pfs|EQG Files (*.eqg)|*.eqg|S3D Files (*.s3d)|*.s3d|PAK Files (*.pak)|*.pak", ".")
							if err != nil {
								slog.Printf("Failed to open: %s\n", err)
								return
							}
							slog.Printf("Menu Opening %s\n", path)
							err = handler.ArchiveOpenInvoke(path, "", true)
							if err != nil {
								slog.Printf("Failed to open: %s\n", err)
								return
							}

						},
						Shortcut: cpl.Shortcut{
							Key:       walk.KeyO,
							Modifiers: walk.ModControl,
						},
					},
					cpl.Action{
						Text: "&Save",
						OnTriggered: func() {
							err = handler.ArchiveSaveInvoke("")
							if err != nil {
								slog.Printf("Failed to save: %s\n", err)
								return
							}
						},
						Shortcut: cpl.Shortcut{
							Key:       walk.KeyS,
							Modifiers: walk.ModControl,
						},
					},
					cpl.Action{
						Text: "&Save As...",
						OnTriggered: func() {
							path, err := ShowSave("Save EQ Archive", "All Archives|*.pfs;*.eqg;*.s3d;*.pak|PFS Files (*.pfs)|*.pfs|EQG Files (*.eqg)|*.eqg|S3D Files (*.s3d)|*.s3d|PAK Files (*.pak)|*.pak", ".")
							if err != nil {
								slog.Printf("Failed to save: %s\n", err)
								return
							}
							slog.Printf("Saving %s\n", path)
							err = handler.ArchiveSaveInvoke(path)
							if err != nil {
								slog.Printf("Failed to save: %s\n", err)
								return
							}
						},
						Shortcut: cpl.Shortcut{
							Key:       walk.KeyS,
							Modifiers: walk.ModControl,
						},
					},
					cpl.Action{
						Text: "Refresh",
						OnTriggered: func() {
							lastSel := gui.table.CurrentIndex()
							handler.ArchiveRefreshInvoke()
							gui.table.SetCurrentIndex(lastSel)
						},
						Shortcut: cpl.Shortcut{
							Key: walk.KeyF5,
						},
					},
					cpl.Action{
						Text: "E&xit",
						OnTriggered: func() {
							gui.mw.Close()
						},
					},
				},
			},
			cpl.Menu{
				Text: "&File",
				Items: []cpl.MenuItem{
					cpl.Action{
						Text: "&Export Selected",
						OnTriggered: func() {
							entry := gui.fileEntries[gui.table.CurrentIndex()].Name
							slog.Printf("Exporting %s\n", entry)

							path, err := ShowSave("Export "+entry, entry, ".")
							if err != nil {
								slog.Printf("Failed to save: %s\n", err)
								return
							}
							err = handler.ArchiveExportFileInvoke(path, entry)
							if err != nil {
								slog.Printf("Failed to save: %s\n", err)
								return
							}
						},
						AssignTo: &gui.exportSelected,
					},
					cpl.Action{
						Text: "Export &All",
						OnTriggered: func() {
							path, err := ShowDirSave("Export All Contents", "All Files|*.*", ".")
							if err != nil {
								slog.Printf("Failed to save: %s\n", err)
								return
							}
							slog.Printf("Exporting all to %s\n", path)
							err = handler.ArchiveExportAllInvoke(path)
							if err != nil {
								slog.Printf("Failed to save: %s\n", err)
								return
							}

						},
					},
				},
			},
		},
		Children: []cpl.Widget{
			cpl.HSplitter{Children: []cpl.Widget{
				cpl.VSplitter{
					AssignTo:           &filesSplit,
					AlwaysConsumeSpace: false,
					Children: []cpl.Widget{
						cpl.Label{Text: "Files"},
						cpl.TableView{
							AssignTo:              &gui.table,
							Name:                  "tableView",
							AlternatingRowBG:      true,
							ColumnsOrderable:      true,
							MultiSelection:        false,
							OnCurrentIndexChanged: onFileViewSelect,
							StyleCell:             fvs.StyleCell,
							Columns: []cpl.TableViewColumn{
								{Name: "Name", Width: 160},
								{Name: "Ext", Width: 40},
								{Name: "Size", Width: 80},
							},
						},
					},
				},
				cpl.VSplitter{AlwaysConsumeSpace: true, Children: []cpl.Widget{
					cpl.Label{Text: "", AssignTo: &gui.imageLabel},
					cpl.ImageView{
						AssignTo: &gui.image,
						Visible:  false,
						Mode:     cpl.ImageViewModeZoom,
					},
				}},

				cpl.VSplitter{AlwaysConsumeSpace: true, Children: []cpl.Widget{
					cpl.Label{Text: "Tree"},
					cpl.TreeView{
						AssignTo: &gui.treeView,
						Visible:  true,
						Model:    gui.treeModel,
						OnCurrentItemChanged: func() {
							node := gui.treeView.CurrentItem().(*component.TreeNode)

							buf := bytes.NewBuffer(nil)
							enc := yaml.NewEncoder(buf)
							defer enc.Close()

							err = enc.Encode(node.Ref())
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
								unkEdit.SetVisible(true)
								gui.titleEdit.SetText(fmt.Sprintf("Select a valid node to edit on left. %T is not yet supported", node.Ref()))

								return
							}
							ext := editor.Ext()
							editView, ok := editViews[ext]
							if !ok {
								slog.Printf("Failed finding edit view %s\n", ext)
								return
							}
							if editView == nil {
								slog.Printf("Failed edit view is nil %s\n", ext)
								unkEdit.SetVisible(true)
								return
							}
							for _, view := range editViews {
								if view == nil {
									continue
								}
								view.SetVisible(false)
							}
							activeEditor = editor
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
							if key == walk.KeyInsert {
								node := gui.treeView.CurrentItem().(*component.TreeNode)
								if node == nil {
									slog.Printf("Failed to insert: node is nil\n")
									return
								}
								newNode := component.NewTreeNode(node, ico.Grab(".lay"), "New Node", nil)
								slog.Printf("Inserted %+v\n", newNode)
								gui.treeModel.PublishItemsReset(node)
								return
							}
							if key == walk.KeyV && walk.ControlDown() {
								if gui.treeCopy == nil {
									slog.Printf("Failed to paste: copy is nil\n")
									return
								}

								node := gui.treeView.CurrentItem().(*component.TreeNode)
								if node == nil {
									slog.Printf("Failed to insert: node is nil\n")
									return
								}

								newNode := component.NewTreeNode(node, gui.treeCopy.Icon(), gui.treeCopy.Text(), gui.treeCopy.DuplicateRef())

								slog.Printf("Pasted %+v\n", newNode)
								gui.treeModel.PublishItemsReset(node)
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
					},
				}},
				cpl.TabWidget{
					AssignTo: &gui.pageView,
					Pages: []cpl.TabPage{
						{
							Title:    "Edit",
							Layout:   cpl.VBox{},
							AssignTo: &gui.editView,
							Children: []cpl.Widget{
								cpl.VSplitter{
									Visible: true,
									Children: []cpl.Widget{
										cpl.Composite{
											Visible:  true,
											AssignTo: &unkEdit,
											Layout:   cpl.Grid{Columns: 2},
											Children: []cpl.Widget{
												cpl.Label{Text: "Select a valid node to edit on left", AssignTo: &gui.titleEdit},
											},
										},
										cpl.Composite{
											Visible:  false,
											AssignTo: &layEdit,
											Layout:   cpl.Grid{Columns: 2},
											Children: form.LayEditWidgets(),
										},
										cpl.Composite{
											Visible:  false,
											AssignTo: &ptsEdit,
											Layout:   cpl.Grid{Columns: 2},
											Children: form.PtsEditWidgets(),
										},
										cpl.Composite{
											Visible:  false,
											AssignTo: &headerEdit,
											Layout:   cpl.Grid{Columns: 2},
											Children: form.HeaderEditWidgets(),
										},
										cpl.Composite{
											Visible:  false,
											AssignTo: &zonObjEdit,
											Layout:   cpl.Grid{Columns: 2},
											Children: form.ZonObjEditWidgets(),
										},
										cpl.Composite{
											Visible:  false,
											AssignTo: &modEdit,
											Layout:   cpl.Grid{Columns: 2},
											Children: form.ModEditWidgets(),
										},
										cpl.Composite{
											Layout: cpl.HBox{},
											Children: []cpl.Widget{
												cpl.HSpacer{},
												/*cpl.PushButton{
													Text:      "Reset",
													OnClicked: handler.EditResetInvoke,
												},*/
												cpl.PushButton{
													Text:      "Save",
													OnClicked: handler.EditSaveInvoke,
													AssignTo:  &gui.editSave,
												},
											},
										},
									},
								},
							},
						},
						{
							Title:    "Yaml",
							AssignTo: &gui.yamlView,
							Layout:   cpl.VBox{},
							Children: []cpl.Widget{
								cpl.Label{
									Text: "A read only view of the selected node's contents. Useful for copy pasting to discord.",
								},
								cpl.TextEdit{
									AssignTo:   &gui.contents,
									ReadOnly:   true,
									Enabled:    false,
									VScroll:    true,
									Background: cpl.SolidColorBrush{Color: wcolor.RGB(255, 255, 255)},
								},
							},
						},
					},
				},
				cpl.ProgressBar{
					AssignTo: &gui.progress,
					Visible:  false,
					MaxValue: 100,
					MinValue: 0,
				},
			}},
		},
		AssignTo: &gui.mw,
		StatusBarItems: []cpl.StatusBarItem{
			{
				AssignTo: &gui.statusBar,
				Text:     "Ready",
				OnClicked: func() {
					fmt.Println("status bar clicked")
				},
			},
		},
	}
	err = cmw.Create()
	if err != nil {
		return fmt.Errorf("create main window: %w", err)
	}

	editViews[".lay"] = layEdit
	editViews[".pts"] = ptsEdit
	editViews["header"] = headerEdit
	editViews["zon_obj"] = zonObjEdit
	editViews[".mod"] = modEdit

	gui.contents.SetText("Init")

	filesSplit.SetMinMaxSizePixels(walk.Size{Width: 250, Height: 0}, walk.Size{Width: 300, Height: 0})

	gui.table.SetModel(gui.fileView)

	return nil
}

func Run() int {
	if gui == nil {
		return 1
	}

	gui.mw.SetVisible(true)
	return gui.mw.Run()
}

func SubscribeClose(fn func(cancelled *bool, reason byte)) {
	if gui == nil {
		return
	}
	gui.mw.Closing().Attach(fn)
}

// Logf logs a message to the gui
func Logf(format string, a ...interface{}) {
	if gui == nil {
		return
	}

	line := fmt.Sprintf(format, a...)
	if strings.Contains(line, "\n") {
		line = "  " + line[0:strings.Index(line, "\n")]
	}
	gui.statusBar.SetText(line)

	//convert \n to \r\n
	//format = strings.ReplaceAll(format, "\n", "\r\n")
	//gui.log.AppendText(fmt.Sprintf(format, a...))

}

func LogClear() {
	if gui == nil {
		return
	}
	gui.log.SetText("")
}

func SetMaxProgress(value int) {
	if gui == nil {
		return
	}
	gui.progress.SetRange(0, value)
}

func SetProgress(value int) {
	if gui == nil {
		return
	}
	gui.progress.SetValue(value)
	gui.progress.SetVisible(value > 0)
}

func MessageBox(title string, message string, isError bool) {
	if gui == nil {
		return
	}
	// convert style to msgboxstyle
	icon := walk.MsgBoxIconInformation
	if isError {
		icon = walk.MsgBoxIconError
	}
	walk.MsgBox(gui.mw, title, message, icon)
}

func MessageBoxYesNo(title string, message string) bool {
	if gui == nil {
		return false
	}
	// convert style to msgboxstyle
	icon := walk.MsgBoxIconInformation
	result := walk.MsgBox(gui.mw, title, message, icon|walk.MsgBoxYesNo)
	return result == walk.DlgCmdYes
}

func MessageBoxf(title string, format string, a ...interface{}) {
	if gui == nil {
		return
	}
	// convert style to msgboxstyle
	icon := walk.MsgBoxIconInformation
	walk.MsgBox(gui.mw, title, fmt.Sprintf(format, a...), icon)
}

func SetTitle(title string) {
	if gui == nil {
		return
	}
	gui.mw.SetTitle(title)
}

func Close() {
	if gui == nil {
		return
	}
	gui.mw.Close()
}

func ShowOpen(title string, filter string, initialDirPath string) (string, error) {
	if gui == nil {
		return "", fmt.Errorf("gui not initialized")
	}
	dialog := walk.FileDialog{
		Title:          title,
		Filter:         filter,
		InitialDirPath: initialDirPath,
	}
	ok, err := dialog.ShowOpen(gui.mw)
	if err != nil {
		return "", fmt.Errorf("show open: %w", err)
	}
	if !ok {
		return "", fmt.Errorf("show open: cancelled")
	}
	return dialog.FilePath, nil
}

func ShowSave(title string, fileName string, initialDirPath string) (string, error) {
	if gui == nil {
		return "", fmt.Errorf("gui not initialized")
	}
	dialog := walk.FileDialog{
		Title:          title,
		FilePath:       fileName,
		InitialDirPath: initialDirPath,
	}
	ok, err := dialog.ShowSave(gui.mw)
	if err != nil {
		return "", fmt.Errorf("show save: %w", err)
	}
	if !ok {
		return "", fmt.Errorf("show save: cancelled")
	}
	return dialog.FilePath, nil
}

func ShowDirSave(title string, filter string, initialDirPath string) (string, error) {
	if gui == nil {
		return "", fmt.Errorf("gui not initialized")
	}
	dialog := walk.FileDialog{
		Title:          title,
		Filter:         filter,
		InitialDirPath: initialDirPath,
	}
	ok, err := dialog.ShowBrowseFolder(gui.mw)
	if err != nil {
		return "", fmt.Errorf("show save: %w", err)
	}
	if !ok {
		return "", fmt.Errorf("show save: cancelled")
	}
	return dialog.FilePath, nil
}

func SetFileViewItems(items []*component.FileViewEntry) {
	if gui == nil {
		return
	}
	gui.fileEntries = items
	gui.fileView.SetItems(items)
	if len(items) > 0 {
		gui.table.SetCurrentIndex(0)
		onFileViewSelect()
	}

	handler.FileViewRefreshInvoke(items)

}

func onFileViewSelect() {
	if len(gui.fileEntries) == 0 {
		return
	}

	if gui.table.CurrentIndex() < 0 || gui.table.CurrentIndex() >= len(gui.fileEntries) {
		//slog.Printf("Invalid file index %d", gui.table.CurrentIndex())
		return
	}
	SetProgress(0)
	name := gui.fileEntries[gui.table.CurrentIndex()].Name
	slog.Printf("FileView Selected %s\n", name)
	gui.exportSelected.SetText("&Export " + name)

	err := handler.ArchiveOpenInvoke("", name, false)
	if err != nil {
		slog.Printf("Failed to open: %s\n", err)
		return
	}

	gui.contents.SetEnabled(true)
}

func SetImage(image walk.Image) {
	if gui == nil {
		return
	}
	fmt.Println("image change", image)
	if image == nil {
		gui.image.SetImage(nil)
		gui.imageLabel.SetText("")
		gui.image.SetVisible(false)
		gui.contents.SetVisible(true)
		gui.contentsLabel.SetVisible(true)
		return
	}
	gui.image.SetImage(image)
	gui.imageLabel.SetText("Image")
	gui.image.SetVisible(true)
	gui.contents.SetVisible(false)
	gui.contentsLabel.SetVisible(false)
}

func SetTreeModel(model *component.TreeModel) {
	if gui == nil {
		return
	}
	gui.treeCopy = nil
	gui.treeView.SetModel(model)
	if model.RootCount() < 3 {

		err := gui.treeView.SetExpanded(model.RootAt(0), true)
		if err != nil {
			slog.Printf("Failed to expand root 0: %s\n", err)
		}
		if model.RootCount() == 1 {
			root := model.RootAt(0)
			if root.ChildCount() < 3 {
				err = gui.treeView.SetExpanded(root.ChildAt(0), true)
				if err != nil {
					slog.Printf("Failed to child 0: %s\n", err)
				}

				if root.ChildCount() > 1 {
					err = gui.treeView.SetExpanded(root.ChildAt(1), true)
					if err != nil {
						slog.Printf("Failed to child 1: %s\n", err)
					}
				}
				if root.ChildCount() > 2 {
					err = gui.treeView.SetExpanded(root.ChildAt(2), true)
					if err != nil {
						slog.Printf("Failed to child 2: %s\n", err)
					}
				}
			}
		}

		if model.RootCount() > 1 {
			err = gui.treeView.SetExpanded(model.RootAt(1), true)
			if err != nil {
				slog.Printf("Failed to expand root 1: %s\n", err)
			}
		}

	}
}

func onEditReset() {
	if gui == nil {
		return
	}

	activeEditor.Reset()

	slog.Println("Edit restored")
}

func onEditSave() {
	if gui == nil {
		return
	}

	err := activeEditor.Save()
	if err != nil {
		slog.Printf("Failed to save: %s\n", err)
		return
	}

	//gui.treeView.UpdateItem(activeEditor.Node())
	gui.treeModel.PublishItemChanged(activeEditor.Node())
}
