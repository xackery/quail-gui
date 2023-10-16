//go:build windows
// +build windows

package gui

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/xackery/quail-gui/config"
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

type Gui struct {
	ctx                   context.Context
	cancel                context.CancelFunc
	mw                    *walk.MainWindow
	progress              *walk.ProgressBar
	log                   *walk.TextEdit
	table                 *walk.TableView
	exportSelected        *walk.Action
	newHandler            []func()
	openHandler           []func(path string, file string) error
	savePFSHandler        []func(path string) error
	saveContentHandler    []func(path string, file string) error
	saveAllContentHandler []func(path string) error
	refreshHandler        []func()
	fileView              *FileView
	fileEntries           []*FileViewEntry
	sectionLabel          *walk.Label
	sectionList           *walk.ListBox
	sections              map[string]*Section
	contentsLabel         *walk.Label
	contents              *walk.TextEdit
	image                 *walk.ImageView
	statusBar             *walk.StatusBarItem
}

var (
	gui *Gui
)

// NewMainWindow creates a new main window
func NewMainWindow(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, version string) error {
	gui = &Gui{
		ctx:    ctx,
		cancel: cancel,
	}

	var err error

	cmw := cpl.MainWindow{
		Name: "quail-gui",
		MenuItems: []cpl.MenuItem{
			cpl.Menu{
				Text: "&Archive",
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
								fn(path, "")
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
							for _, fn := range gui.savePFSHandler {
								err = fn("")
								if err != nil {
									slog.Printf("Failed to save: %s\n", err)
									return
								}
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
							for _, fn := range gui.savePFSHandler {
								err = fn(path)
								if err != nil {
									slog.Printf("Failed to save: %s\n", err)
									return
								}
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
							//lastSection := gui.sectionList.CurrentIndex()
							for _, fn := range gui.refreshHandler {
								fn()
							}
							gui.table.SetCurrentIndex(lastSel)
							//gui.sectionList.SetCurrentIndex(lastSection)
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
							for _, fn := range gui.saveContentHandler {
								err = fn(path, entry)
								if err != nil {
									slog.Printf("Failed to save %s: %s\n", entry, err)
									return
								}
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
							for _, fn := range gui.saveAllContentHandler {
								err = fn(path)
								if err != nil {
									slog.Printf("Failed to save content: %s\n", err)
									return
								}
							}
						},
					},
				},
			},
		},
		AssignTo: &gui.mw,
		Visible:  false,
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

	gui.mw.SetTitle("quail-gui v" + version)
	gui.mw.SetMinMaxSize(walk.Size{Width: 405, Height: 371}, walk.Size{Width: 0, Height: 0})
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
	gui.table.CurrentIndexChanged().Attach(onTableSelect)
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
	gui.sectionList.ItemActivated().Attach(onSectionListSelect)
	gui.sectionList.CurrentIndexChanged().Attach(onSectionListSelect)
	gui.sectionList.SetWidth(200)

	gui.contents, err = walk.NewTextEdit(gui.mw)
	if err != nil {
		return fmt.Errorf("new text edit: %w", err)
	}
	gui.contents.SetReadOnly(true)
	gui.contents.SetEnabled(false)

	comp, err := walk.NewComposite(gui.mw)
	if err != nil {
		return fmt.Errorf("new composite: %w", err)
	}
	comp.SetLayout(walk.NewHBoxLayout())

	entry, err := walk.NewComposite(gui.mw)
	if err != nil {
		return fmt.Errorf("new composite: %w", err)
	}
	entry.SetLayout(walk.NewVBoxLayout())

	entry.Children().Add(newLabel("Files"))
	entry.Children().Add(gui.table)
	comp.Children().Add(entry)

	entry, err = walk.NewComposite(gui.mw)
	if err != nil {
		return fmt.Errorf("new composite: %w", err)
	}
	entry.SetLayout(walk.NewVBoxLayout())
	gui.sectionLabel = newLabel("Category")
	entry.Children().Add(gui.sectionLabel)
	entry.Children().Add(gui.sectionList)
	comp.Children().Add(entry)

	gui.image, err = walk.NewImageView(gui.mw)
	if err != nil {
		return fmt.Errorf("new image view: %w", err)
	}
	gui.image.SetVisible(false)

	entry, err = walk.NewComposite(gui.mw)
	if err != nil {
		return fmt.Errorf("new composite: %w", err)
	}
	entry.SetLayout(walk.NewVBoxLayout())
	gui.contentsLabel = newLabel("Contents")
	entry.Children().Add(gui.contentsLabel)
	entry.Children().Add(gui.contents)
	entry.Children().Add(gui.image)
	comp.Children().Add(entry)

	gui.progress, err = walk.NewProgressBar(gui.mw)
	if err != nil {
		return fmt.Errorf("new progress bar: %w", err)
	}

	//gui.progress.SetMinMaxSize(walk.Size{Width: 400, Height: 39}, walk.Size{Width: 400, Height: 39})
	gui.progress.SetValue(0)
	//gui.progress.SetMinMaxSize(walk.Size{Width: 400, Height: 39}, walk.Size{Width: 400, Height: 39})

	gui.mw.Children().Add(gui.progress)
	gui.mw.SetSize(walk.Size{Width: 405, Height: 371})

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

	line := fmt.Sprintf(format, a...)
	if strings.Contains(line, "\n") {
		line = line[0:strings.Index(line, "\n")]
	}
	gui.statusBar.SetText(line)

	//convert \n to \r\n
	format = strings.ReplaceAll(format, "\n", "\r\n")
	gui.log.AppendText(fmt.Sprintf(format, a...))

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

func SubscribeNew(fn func()) {
	if gui == nil {
		return
	}
	gui.newHandler = append(gui.newHandler, fn)
}

func SubscribeOpen(fn func(path string, file string) error) {
	if gui == nil {
		return
	}
	gui.openHandler = append(gui.openHandler, fn)
}

func SubscribeRefresh(fn func()) {
	if gui == nil {
		return
	}
	gui.refreshHandler = append(gui.refreshHandler, fn)
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

func SubscribeSavePFS(fn func(path string) error) {
	if gui == nil {
		return
	}
	gui.savePFSHandler = append(gui.savePFSHandler, fn)
}

func SubscribeSaveContent(fn func(path string, file string) error) {
	if gui == nil {
		return
	}
	gui.saveContentHandler = append(gui.saveContentHandler, fn)
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

func SubscribeSaveAllContent(fn func(path string) error) {
	if gui == nil {
		return
	}
	gui.saveAllContentHandler = append(gui.saveAllContentHandler, fn)
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

func SetFileViewItems(items []*FileViewEntry) {
	if gui == nil {
		return
	}
	gui.fileEntries = items
	gui.fileView.SetItems(items)
}

func SetSections(sections map[string]*Section) {
	if gui == nil {
		return
	}
	gui.sections = sections
	sectionList := []string{}
	for _, v := range sections {
		sectionList = append(sectionList, fmt.Sprintf("%s (%d)", v.Name, v.Count))
	}
	sort.Strings(sectionList)
	gui.sectionList.SetModel(sectionList)
}

func onTableSelect() {
	if len(gui.fileEntries) == 0 {
		slog.Printf("No files to open")
		return
	}

	if gui.table.CurrentIndex() < 0 || gui.table.CurrentIndex() >= len(gui.fileEntries) {
		//slog.Printf("Invalid file index %d", gui.table.CurrentIndex())
		return
	}
	SetProgress(0)
	name := gui.fileEntries[gui.table.CurrentIndex()].Name
	slog.Printf("Selected %s\n", name)
	gui.exportSelected.SetText("&Export " + name)
	for _, fn := range gui.openHandler {
		err := fn("", name)
		if err != nil {
			slog.Printf("Failed to open: %s\n", err)
			gui.sectionList.SetModel([]string{})
			gui.sectionList.SetEnabled(false)
			gui.contents.SetEnabled(false)
			return
		}
	}
	gui.sectionList.SetEnabled(true)
	gui.contents.SetEnabled(true)
}

func onSectionListSelect() {
	SetProgress(0)
	fmt.Printf("Activated: %v\n", gui.sectionList.CurrentIndex())
	if gui.sectionList.CurrentIndex() < 0 || gui.sectionList.CurrentIndex() >= len(gui.sectionList.Model().([]string)) {
		slog.Printf("Invalid section index %d\n", gui.sectionList.CurrentIndex())
		return
	}
	//get current sectionlist
	name := gui.sectionList.Model().([]string)[gui.sectionList.CurrentIndex()]

	gui.sectionList.SetEnabled(true)
	gui.contents.SetEnabled(true)
	gui.contents.SetText(gui.sections[name].Content)
	slog.Printf("Selected %s\n", name)
}

func SetImage(image walk.Image) {
	if gui == nil {
		return
	}
	if image == nil {
		gui.image.SetVisible(false)
		gui.contents.SetVisible(true)
		gui.sectionList.SetVisible(true)
		gui.sectionLabel.SetVisible(true)
		gui.contentsLabel.SetVisible(true)
		return
	}
	gui.image.SetImage(image)
	gui.image.SetVisible(true)
	gui.contents.SetVisible(false)
	gui.sectionList.SetVisible(false)
	gui.sectionLabel.SetVisible(false)
	gui.contentsLabel.SetVisible(false)
}
