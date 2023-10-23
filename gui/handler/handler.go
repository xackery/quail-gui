package handler

import "github.com/xackery/quail-gui/gui/component"

var (
	windowCloseHandler       []func(cancelled *bool, reason byte)
	archiveNewHandler        []func()
	archiveOpenHandler       []func(path string, file string, isSelect bool) error
	archiveRefreshHandler    []func()
	archiveSaveHandler       []func(path string) error
	archiveExportFileHandler []func(path string, file string) error
	archiveExportAllHandler  []func(path string) error
	fileViewRefreshHandler   []func(items []*component.FileViewEntry)
	editResetHandler         []func()
	editSaveHandler          []func()
	previewHandler           []func()
)

// WindowCloseSubscribe allows subscribing to close events
func WindowCloseSubscribe(fn func(cancelled *bool, reason byte)) {
	windowCloseHandler = append(windowCloseHandler, fn)
}

// WindowCloseInvoke invokes close events on the window
func WindowCloseInvoke(cancelled *bool, reason byte) {
	for _, fn := range windowCloseHandler {
		fn(cancelled, reason)
	}
}

// SubscribeNewArchive allows subscribing to new archve creation events
func ArchiveNewSubscribe(fn func()) {
	archiveNewHandler = append(archiveNewHandler, fn)
}

// ArchiveNewInvoke invokes new archive creation events
func ArchiveNewInvoke() {
	for _, fn := range archiveNewHandler {
		fn()
	}
}

// SubscribeOpen allows subscribing to open archve events
func ArchiveOpenSubscribe(fn func(path string, file string, isSelect bool) error) {
	archiveOpenHandler = append(archiveOpenHandler, fn)
}

// ArchiveOpenInvoke invokes open archve events
func ArchiveOpenInvoke(path string, file string, isSelect bool) error {
	for _, fn := range archiveOpenHandler {
		err := fn(path, file, isSelect)
		if err != nil {
			return err
		}
	}
	return nil
}

// ArchiveSaveSubscribe allows subscribing to save archive events
func ArchiveSaveSubscribe(fn func(path string) error) {
	archiveSaveHandler = append(archiveSaveHandler, fn)
}

// ArchiveSaveInvoke invokes save archive events
func ArchiveSaveInvoke(path string) error {
	for _, fn := range archiveSaveHandler {
		err := fn(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// ArchiveExportFileSubscribe allows subscribing to save archive file events
func ArchiveExportFileSubscribe(fn func(path string, file string) error) {
	archiveExportFileHandler = append(archiveExportFileHandler, fn)
}

// ArchiveExportFileInvoke invokes save archive file events
func ArchiveExportFileInvoke(path string, file string) error {
	for _, fn := range archiveExportFileHandler {
		err := fn(path, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// FileViewRefreshSubscribe allows subscribing to file view refresh events
func FileViewRefreshSubscribe(fn func(items []*component.FileViewEntry)) {
	fileViewRefreshHandler = append(fileViewRefreshHandler, fn)
}

// FileViewRefreshInvoke invokes file view refresh events
func FileViewRefreshInvoke(items []*component.FileViewEntry) {
	for _, fn := range fileViewRefreshHandler {
		fn(items)
	}
}

// ArchiveRefreshSubscribe allows subscribing to refresh events
func ArchiveRefreshSubscribe(fn func()) {
	archiveRefreshHandler = append(archiveRefreshHandler, fn)
}

// ArchiveRefreshInvoke invokes refresh events
func ArchiveRefreshInvoke() {
	for _, fn := range archiveRefreshHandler {
		fn()
	}
}

// SubscribeExportArchive allows subscribing to export archive events
func ArchiveExportAllSubscribe(fn func(path string) error) {
	archiveExportAllHandler = append(archiveExportAllHandler, fn)
}

// ArchiveExportAllInvoke invokes export archive events
func ArchiveExportAllInvoke(path string) error {
	for _, fn := range archiveExportAllHandler {
		err := fn(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// EditResetSubscribe allows subscribing to cancel events
func EditResetSubscribe(fn func()) {
	editResetHandler = append(editResetHandler, fn)
}

// EditResetInvoke invokes cancel events
func EditResetInvoke() {
	for _, fn := range editResetHandler {
		fn()
	}
}

// EditSaveSubscribe allows subscribing to save events
func EditSaveSubscribe(fn func()) {
	editSaveHandler = append(editSaveHandler, fn)
}

// EditSaveInvoke invokes save events
func EditSaveInvoke() {
	for _, fn := range editSaveHandler {
		fn()
	}
}

// PreviewSubscribe allows subscribing to preview events
func PreviewSubscribe(fn func()) {
	previewHandler = append(previewHandler, fn)
}

// PreviewInvoke invokes preview events
func PreviewInvoke() {
	for _, fn := range previewHandler {
		fn()
	}
}
