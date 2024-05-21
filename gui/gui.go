package gui

import (
	"fmt"
	"strings"

	_ "embed"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/pfs"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

var (
	mw          *walk.MainWindow
	statusBar   *walk.StatusBarItem
	isEdited    bool
	archive     *pfs.Pfs // currently loaded pfs archive
	archivePath string   // used for writing archive back to file
	file        *walk.TableView
	fileView    *component.FileView
)

func New() error {
	if mw != nil {
		return fmt.Errorf("main window already created")
	}

	slog.AddHandler(logf)

	fileView = component.NewFileView()
	fvs := component.NewFileViewStyler(fileView)

	cmw := cpl.MainWindow{
		AssignTo:      &mw,
		MinSize:       cpl.Size{Width: 400, Height: 300},
		Visible:       false,
		Name:          "quail-gui",
		OnSizeChanged: onSizeChanged,
		MenuItems: []cpl.MenuItem{
			cpl.Menu{
				Text: "&File",
				Items: []cpl.MenuItem{
					cpl.Action{Text: " &New", AssignTo: &menuFileNew, OnTriggered: onFileNew},
					cpl.Separator{},
					cpl.Action{Text: "&Open", Shortcut: cpl.Shortcut{Modifiers: walk.ModControl, Key: walk.KeyO}, AssignTo: &menuFileOpen, OnTriggered: onFileOpen},
					cpl.Action{Text: "Open &Recent", AssignTo: &menuFileOpenRecent, OnTriggered: onFileOpenRecent},
					cpl.Separator{},
					cpl.Action{Text: "&Save", Shortcut: cpl.Shortcut{Modifiers: walk.ModControl, Key: walk.KeyS}, AssignTo: &menuFileSave, OnTriggered: onFileSave},
					cpl.Action{Text: "Save As...", AssignTo: &menuFileSaveAs, OnTriggered: onFileSaveAs},
					cpl.Separator{},
					cpl.Action{Text: "&Refresh", Shortcut: cpl.Shortcut{Key: walk.KeyF5}, AssignTo: &menuFileRefresh, OnTriggered: onFileRefresh},

					cpl.Action{Text: "&Close", Shortcut: cpl.Shortcut{Modifiers: walk.ModControl, Key: walk.KeyX}, AssignTo: &menuFileClose, OnTriggered: onFileClose},
					cpl.Action{Text: "E&xit", Shortcut: cpl.Shortcut{Modifiers: walk.ModControl, Key: walk.KeyQ}, AssignTo: &menuFileExit, OnTriggered: onFileExit},
				},
			},
			cpl.Menu{
				Text: "&Edit",
				Items: []cpl.MenuItem{
					cpl.Action{Text: "&Preferences", AssignTo: &menuEditPreferences, OnTriggered: onEditPreferences},
				},
			},
			cpl.Menu{
				Text: "&Entry",
				Items: []cpl.MenuItem{
					cpl.Action{Text: " &New", Shortcut: cpl.Shortcut{Modifiers: walk.ModControl, Key: walk.KeyN}, AssignTo: &menuEntryNew, OnTriggered: onMenuEntryNew},
					cpl.Separator{},
					cpl.Action{Text: " &Edit", AssignTo: &menuEntryEdit, OnTriggered: onMenuEntryEdit},
					cpl.Action{Text: " &Delete", Shortcut: cpl.Shortcut{Key: walk.KeyDelete}, AssignTo: &menuEntryDelete, OnTriggered: onMenuEntryDelete},
					cpl.Separator{},
					cpl.Action{Text: " &Rename", AssignTo: &menuEntryRename, OnTriggered: onMenuEntryRename},
				},
			},
			cpl.Menu{
				Text: "&Help",
				Items: []cpl.MenuItem{
					cpl.Action{Text: "&About", AssignTo: &menuHelpAbout, OnTriggered: onHelpAbout},
				},
			},
		},
		ToolBar: cpl.ToolBar{
			ButtonStyle: cpl.ToolBarButtonImageBeforeText,
			Items: []cpl.MenuItem{
				cpl.Action{Image: ico.Grab("open"), AssignTo: &menuFileOpen, OnTriggered: onFileOpen},
				cpl.Action{Image: ico.Grab("save"), AssignTo: &menuFileSave, OnTriggered: onFileSave},
				cpl.Separator{},
				cpl.Action{Image: ico.Grab("new"), AssignTo: &menuEntryNew, OnTriggered: onMenuEntryNew},
				cpl.Action{Image: ico.Grab("edit"), AssignTo: &menuEntryEdit, OnTriggered: onMenuEntryEdit},
				cpl.Action{Image: ico.Grab("delete"), AssignTo: &menuEntryDelete, OnTriggered: onMenuEntryDelete},
				cpl.Separator{},
				cpl.Action{Image: ico.Grab("wld"), Shortcut: cpl.Shortcut{Modifiers: walk.ModControl, Key: walk.KeyW}, AssignTo: &menuEntryEditWorld, OnTriggered: onMenuEntryEditWorld},
			},
		},
		OnDropFiles: onDrop,
		Layout:      cpl.VBox{},
		Children: []cpl.Widget{
			cpl.TableView{
				AssignTo:         &file,
				AlternatingRowBG: true,
				ColumnsOrderable: true,
				MultiSelection:   false,
				OnKeyDown: func(key walk.Key) {
					if key == walk.KeyUp || key == walk.KeyDown {
						onEntryChange()
					}
				},
				OnCurrentIndexChanged: onEntryChange,
				OnItemActivated:       onEntryActivate,
				StyleCell:             fvs.StyleCell,
				Model:                 fileView,
				ContextMenuItems: []cpl.MenuItem{
					cpl.Action{Text: "Refresh", Image: ico.Grab("refresh"), AssignTo: &menuFileRefresh, OnTriggered: onFileRefresh},
					cpl.Separator{},
					cpl.Action{Text: "Delete", Image: ico.Grab("delete"), AssignTo: &menuFileDelete, OnTriggered: onFileDelete},
				},
				//MaxSize:               cpl.Size{Width: 300, Height: 0},
				Columns: []cpl.TableViewColumn{
					{Name: "Name", Width: 160},
					{Name: "Ext", Width: 40},
					{Name: "Size", Width: 80},
				},
			},
		},
		StatusBarItems: []cpl.StatusBarItem{
			{
				AssignTo: &statusBar,
				Text:     "Ready",
				OnClicked: func() {
					fmt.Println("status bar clicked")
				},
			},
		},
	}
	err := cmw.Create()
	if err != nil {
		return fmt.Errorf("create main window: %w", err)
	}

	entrySetActive(false)

	return nil
}

func Run() int {
	if mw == nil {
		return 1
	}

	mw.SetSize(walk.Size{Width: 300, Height: 300})

	//	walk.CenterWindowOnScreen(mw)
	mw.SetBounds(walk.Rectangle{X: 400, Y: 300, Width: mw.Width(), Height: mw.Height()})

	mw.SetVisible(true)
	return mw.Run()
}

func MainWindow() *walk.MainWindow {
	return mw
}

// logf logs a message to the gui
func logf(format string, a ...interface{}) {
	if mw == nil {
		return
	}

	line := fmt.Sprintf(format, a...)
	if strings.Contains(line, "\n") {
		line = "  " + line[0:strings.Index(line, "\n")]
	}
	statusBar.SetText(line)
}

func generateSize(in int) string {
	val := float64(in)
	if val < 1024 {
		return fmt.Sprintf("%0.0f bytes", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.0f KB", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.0f MB", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.0f GB", val)
	}
	val /= 1024
	return fmt.Sprintf("%0.0f TB", val)
}

func onSizeChanged() {
	slog.Printf("Size changed: %d x %d\n", mw.Width(), mw.Height())
}
