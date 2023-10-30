package gui

import "github.com/xackery/quail-gui/slog"

func onDrop(files []string) {
	if len(files) == 0 {
		slog.Println("Ignoring file drop, no files")
		return
	}
	if len(files) > 1 {
		slog.Println("Ignoring file drop, too many files")
		return
	}
	slog.Printf("File dropped: %s\n", files[0])
}
