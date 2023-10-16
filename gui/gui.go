//go:build windows
// +build windows

package gui

import (
	"context"
	"fmt"
	"strings"

	"github.com/xackery/quail-gui/config"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

type Gui struct {
	ctx         context.Context
	cancel      context.CancelFunc
	mw          *walk.MainWindow
	progress    *walk.ProgressBar
	log         *walk.TextEdit
	table       *walk.TableView
	newHandler  []func()
	openHandler []func(path string)
	fileView    *FileView
	sectionList *walk.ListBox
	treeView    *walk.TreeView
}

var (
	gui        *Gui
	isAutoMode bool
)

// NewMainWindow creates a new main window
func NewMainWindow(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, version string) error {
	gui = &Gui{
		ctx:    ctx,
		cancel: cancel,
	}
	isAutoMode = true

	var err error

	cmw := cpl.MainWindow{
		Name: "quail-gui",
		MenuItems: []cpl.MenuItem{
			cpl.Menu{
				Text: "&File",
				Items: []cpl.MenuItem{
					cpl.Action{
						Text: "&New",
						OnTriggered: func() {
							for _, fn := range gui.newHandler {
								fn()
							}
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
							slog.Printf("Opening %s\n", path)
							for _, fn := range gui.openHandler {
								fn(path)
							}
						},
						Shortcut: cpl.Shortcut{
							Key:       walk.KeyO,
							Modifiers: walk.ModControl,
						},
					},
					cpl.Action{
						Text: "&Save As...",
						OnTriggered: func() {
							fmt.Println("Save triggered")
						},
						Shortcut: cpl.Shortcut{
							Key:       walk.KeyS,
							Modifiers: walk.ModControl,
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
		},
		AssignTo: &gui.mw,
		Visible:  false,
	}
	err = cmw.Create()
	if err != nil {
		return fmt.Errorf("create main window: %w", err)
	}

	gui.mw.SetTitle("quail-gui v" + version)
	gui.mw.SetMinMaxSize(walk.Size{Width: 305, Height: 371}, walk.Size{Width: 305, Height: 371})
	gui.mw.SetLayout(walk.NewVBoxLayout())
	gui.mw.SetVisible(false)

	gui.log, err = walk.NewTextEdit(gui.mw)
	if err != nil {
		return fmt.Errorf("new text edit: %w", err)
	}
	gui.log.SetReadOnly(true)
	gui.log.SetVisible(false)
	gui.log.SetMinMaxSize(walk.Size{Width: 400, Height: 400}, walk.Size{Width: 400, Height: 400})
	slog.AddHandler(Logf)
	gui.mw.Children().Add(gui.log)

	gui.table, err = walk.NewTableView(gui.mw)
	if err != nil {
		return fmt.Errorf("new table view: %w", err)
	}
	gui.fileView = NewFileView()
	gui.table.SetModel(gui.fileView)
	gui.table.SetName("tableView")
	gui.table.SetAlternatingRowBG(true)
	gui.table.SetColumnsOrderable(true)
	gui.table.SetMultiSelection(false)
	gui.table.ItemActivated().Attach(func() {
		fmt.Printf("Activated: %v\n", gui.table.SelectedIndexes())
	})
	gui.table.SetMinMaxSize(walk.Size{Width: 250, Height: 0}, walk.Size{Width: 0, Height: 0})

	col := walk.NewTableViewColumn()
	col.SetDataMember("Name")
	col.SetName("Name")
	col.SetWidth(160)
	err = gui.table.Columns().Add(col)
	if err != nil {
		return fmt.Errorf("add column: %w", err)
	}

	col = walk.NewTableViewColumn()
	col.SetDataMember("Ext")
	col.SetName("Ext")
	col.SetWidth(40)
	err = gui.table.Columns().Add(col)
	if err != nil {
		return fmt.Errorf("add column: %w", err)
	}

	col = walk.NewTableViewColumn()
	col.SetDataMember("Size")
	col.SetName("Size")
	col.SetWidth(80)
	err = gui.table.Columns().Add(col)
	if err != nil {
		return fmt.Errorf("add column: %w", err)
	}

	gui.sectionList, err = walk.NewListBox(gui.mw)
	if err != nil {
		return fmt.Errorf("new list box: %w", err)
	}
	gui.sectionList.SetName("Section")
	gui.sectionList.ItemActivated().Attach(func() {
		fmt.Printf("Activated: %v\n", gui.sectionList.CurrentIndex())
	})

	gui.treeView, err = walk.NewTreeView(gui.mw)
	if err != nil {
		return fmt.Errorf("new tree view: %w", err)
	}
	gui.treeView.SetName("treeView")

	comp, err := walk.NewComposite(gui.mw)
	if err != nil {
		return fmt.Errorf("new composite: %w", err)
	}
	comp.SetLayout(walk.NewHBoxLayout())
	comp.Children().Add(gui.table)
	comp.Children().Add(gui.sectionList)
	comp.Children().Add(gui.treeView)

	gui.progress, err = walk.NewProgressBar(gui.mw)
	if err != nil {
		return fmt.Errorf("new progress bar: %w", err)
	}

	gui.progress.SetMinMaxSize(walk.Size{Width: 400, Height: 39}, walk.Size{Width: 400, Height: 39})
	gui.progress.SetValue(50)
	gui.progress.SetMinMaxSize(walk.Size{Width: 400, Height: 39}, walk.Size{Width: 400, Height: 39})

	gui.mw.Children().Add(gui.progress)
	gui.mw.SetSize(walk.Size{Width: 305, Height: 371})

	return nil
}

func Run() int {
	if gui == nil {
		return 1
	}
	gui.mw.SetVisible(true)
	return gui.mw.Run()
}

func SubscribeClose(fn func(cancelled *bool, reason walk.CloseReason)) {
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

	if !isAutoMode {
		//convert \n to \r\n
		format = strings.ReplaceAll(format, "\n", "\r\n")
		gui.log.AppendText(fmt.Sprintf(format, a...))
	}
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

func SubscribeNew(fn func()) {
	if gui == nil {
		return
	}
	gui.newHandler = append(gui.newHandler, fn)
}

func SubscribeOpen(fn func(path string)) {
	if gui == nil {
		return
	}
	gui.openHandler = append(gui.openHandler, fn)
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

func SetFileViewItems(items []*FileViewEntry) {
	if gui == nil {
		return
	}
	gui.fileView.SetItems(items)
}
