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
	"github.com/xackery/wlk/wcolor"
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
	//walk.SetDarkModeAllowed(true)
	gui = &Gui{
		ctx:    ctx,
		cancel: cancel,
	}

	var err error
	fvs := &fileViewStyler{}
	gui.fileView = NewFileView()
	cmw := cpl.MainWindow{
		Title:   "quail-gui v" + version,
		MinSize: cpl.Size{Width: 405, Height: 371},
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
		Children: []cpl.Widget{
			cpl.HSplitter{Children: []cpl.Widget{
				cpl.VSplitter{Children: []cpl.Widget{
					cpl.Label{Text: "Files"},
					cpl.TableView{
						AssignTo:              &gui.table,
						Name:                  "tableView",
						AlternatingRowBG:      true,
						ColumnsOrderable:      true,
						MultiSelection:        false,
						OnCurrentIndexChanged: onFileViewSelect,
						StyleCell:             fvs.StyleCell,
						MinSize:               cpl.Size{Width: 250, Height: 0},
						Columns: []cpl.TableViewColumn{
							{Name: "Name", Width: 160},
							{Name: "Ext", Width: 40},
							{Name: "Size", Width: 80},
						},
					},
				}},
				cpl.VSplitter{Children: []cpl.Widget{
					//cpl.Label{Text: "Image"},
					cpl.ImageView{
						AssignTo: &gui.image,
						Visible:  false,
						Mode:     cpl.ImageViewModeZoom,
					},
				}},
				cpl.VSplitter{Children: []cpl.Widget{
					cpl.Label{Text: "", AssignTo: &gui.sectionLabel},
					cpl.ListBox{
						AssignTo:              &gui.sectionList,
						Name:                  "Section",
						OnCurrentIndexChanged: onSectionListSelect,
						MinSize:               cpl.Size{Width: 200, Height: 0},
					},
				}},
				cpl.VSplitter{Children: []cpl.Widget{
					cpl.Label{Text: "", AssignTo: &gui.contentsLabel},
					cpl.TextEdit{
						AssignTo:   &gui.contents,
						ReadOnly:   true,
						Enabled:    false,
						VScroll:    true,
						Background: cpl.SolidColorBrush{Color: wcolor.RGB(255, 255, 255)},
					},
				}},
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
	if len(items) > 0 {
		gui.table.SetCurrentIndex(0)
		onFileViewSelect()
	}
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
	if len(sectionList) > 0 {
		gui.sectionList.SetCurrentIndex(0)
		onSectionListSelect()
	}
}

func onFileViewSelect() {
	if len(gui.fileEntries) == 0 {
		gui.sectionLabel.SetText("")
		return
	}

	if gui.table.CurrentIndex() < 0 || gui.table.CurrentIndex() >= len(gui.fileEntries) {
		//slog.Printf("Invalid file index %d", gui.table.CurrentIndex())
		return
	}
	SetProgress(0)
	name := gui.fileEntries[gui.table.CurrentIndex()].Name
	slog.Printf("FileView Selected %s\n", name)
	gui.sectionLabel.SetText(name)
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

	if strings.Contains(name, "(") {
		name = name[0 : strings.Index(name, "(")-1]
	}

	fmt.Println("name", name)
	gui.contentsLabel.SetText(name)
	gui.sectionList.SetEnabled(true)
	gui.contents.SetEnabled(true)
	gui.contents.SetText(gui.sections[name].Content)
	gui.image.SetVisible(false)
	gui.image.SetImage(nil)
	gui.contents.SetVisible(true)
	gui.sectionList.SetVisible(true)
	gui.sectionLabel.SetVisible(true)
	gui.contentsLabel.SetVisible(true)
	slog.Printf("Section Selected %s\n", name)
}

func SetImage(image walk.Image) {
	if gui == nil {
		return
	}
	fmt.Println("image change", image)
	if image == nil {
		gui.image.SetImage(nil)
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
