package gui

import (
	"fmt"
	"strings"

	_ "embed"

	"github.com/xackery/quail-gui/gui/component"
	"github.com/xackery/quail-gui/ico"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

const (
	currentViewArchiveFiles = iota
	currentViewContext
)

var (
	mw          *walk.MainWindow
	statusBar   *walk.StatusBarItem
	currentView int
)

func New() error {
	if mw != nil {
		return fmt.Errorf("main window already created")
	}

	slog.AddHandler(logf)

	widget.fileView = component.NewFileView()
	fvs := component.NewFileViewStyler(widget.fileView)
	currentView = currentViewContext

	cmw := cpl.MainWindow{
		AssignTo:      &mw,
		Title:         "quail-gui",
		MinSize:       cpl.Size{Width: 300, Height: 300},
		Visible:       false,
		Name:          "quail-gui",
		OnSizeChanged: widget.onSizeChanged,
		MenuItems: []cpl.MenuItem{
			cpl.Menu{
				Text: "&File",
				Items: []cpl.MenuItem{
					cpl.Action{Text: " &New", AssignTo: &menu.fileNew, OnTriggered: menu.onFileNew},
					cpl.Separator{},
					cpl.Action{Text: "&Open", AssignTo: &menu.fileOpen, OnTriggered: menu.onFileOpen},
					cpl.Action{Text: "Open &Recent", AssignTo: &menu.fileOpenRecent, OnTriggered: menu.onFileOpenRecent},
					cpl.Separator{},
					cpl.Action{Text: "E&xit", AssignTo: &menu.fileExit, OnTriggered: menu.onFileExit},
				},
			},
			cpl.Menu{
				Text: "&Help",
				Items: []cpl.MenuItem{
					cpl.Action{Text: "&About", AssignTo: &menu.helpAbout, OnTriggered: menu.onHelpAbout},
				},
			},
		},
		ToolBar: cpl.ToolBar{
			ButtonStyle: cpl.ToolBarButtonImageBeforeText,
			Items: []cpl.MenuItem{
				cpl.Action{Image: ico.Grab("back"), AssignTo: &toolbar.back, OnTriggered: toolbar.onBack},
				cpl.Action{Text: " &New", Image: ico.Grab("open"), AssignTo: &menu.fileNew, OnTriggered: menu.onFileNew},
				cpl.Action{Image: ico.Grab("delete"), AssignTo: &menu.fileDelete, OnTriggered: menu.onFileDelete},
			},
		},
		OnDropFiles: onDrop,
		Layout:      cpl.VBox{},
		Children: []cpl.Widget{
			cpl.Label{Text: "", AssignTo: &widget.breadcrumb, Font: cpl.Font{PointSize: 10}, TextAlignment: cpl.AlignFar},
			cpl.TableView{
				AssignTo:              &widget.file,
				AlternatingRowBG:      true,
				ColumnsOrderable:      true,
				MultiSelection:        false,
				OnCurrentIndexChanged: widget.onFileChange,
				OnItemActivated:       widget.onFileActivated,
				StyleCell:             fvs.StyleCell,
				Model:                 widget.fileView,
				ContextMenuItems: []cpl.MenuItem{
					cpl.ActionRef{Action: &menu.fileNew},
					cpl.Action{Text: " Refresh", Image: ico.Grab("refresh"), AssignTo: &menu.fileRefresh, OnTriggered: menu.onFileRefresh},
					cpl.Separator{},
					cpl.Action{Text: " Delete", Image: ico.Grab("delete"), AssignTo: &menu.fileDelete, OnTriggered: menu.onFileDelete},
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
	return nil
}

func Run() int {
	if mw == nil {
		return 1
	}

	mw.SetSize(walk.Size{Width: 300, Height: 300})
	walk.CenterWindowOnScreen(mw)

	mw.SetVisible(true)
	return mw.Run()
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

func viewSetBack() {
	switch currentView {
	case currentViewArchiveFiles:
		return
	case currentViewContext:
		viewSet(currentViewArchiveFiles)
	}
}

func viewSet(view int) {
	if currentView == view {
		return
	}

	switch view {
	case currentViewArchiveFiles:
		widget.file.SetVisible(true)
	case currentViewContext:
		widget.file.SetVisible(false)
	}
	widget.breadcrumbRefresh()
	currentView = view
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
