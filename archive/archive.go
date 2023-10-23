package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/pfs"
)

var (
	openPath string
	fileName string
	archive  *pfs.PFS
)

func Open(path string, file string) error {
	if path == "" {
		path = openPath
	}
	openPath = path
	fileName = file

	slog.Printf("Opening %s %s\n", path, file)
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path check: %w", err)
	}
	if fi.IsDir() {
		return fmt.Errorf("inspect requires a target file, directory provided")
	}

	ext := strings.ToLower(filepath.Ext(path))
	isValidExt := false
	exts := []string{".eqg", ".s3d", ".pfs", ".pak"}
	for _, ext := range exts {
		if strings.HasSuffix(path, ext) {
			isValidExt = true
			break
		}
	}

	if !isValidExt {
		return fmt.Errorf("invalid extension %s", ext)
	}
	archive, err = pfs.NewFile(path)
	if err != nil {
		return fmt.Errorf("%s load: %w", ext, err)
	}
	return nil
}

func Save(path string) error {
	if archive == nil {
		return fmt.Errorf("no archive loaded")
	}
	if path == "" {
		path = openPath
	}
	slog.Printf("Saving %s\n", path)
	w, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer w.Close()
	err = archive.Encode(w)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	slog.Printf("Saved %s\n", path)
	return nil
}

func ExportAll(dir string) error {
	slog.Printf("Exporting %d files to %s\n", len(archive.Files()), dir)
	if archive == nil {
		return fmt.Errorf("no archive loaded")
	}
	for _, fe := range archive.Files() {
		err := writeFile(dir, fe.Name(), fe.Data())
		if err != nil {
			return fmt.Errorf("export %s: %w", fe.Name(), err)
		}
	}
	slog.Printf("Exported %d files to %s\n", len(archive.Files()), dir)
	return nil
}

func ExportFile(dir string, file string) error {
	if archive == nil {
		return fmt.Errorf("no archive loaded")
	}
	slog.Printf("Exporting %s to %s", file, dir)
	data, err := archive.File(file)
	if err != nil {
		return fmt.Errorf("file %s: %w", file, err)
	}
	err = writeFile(dir, file, data)
	if err != nil {
		return fmt.Errorf("write %s: %w", file, err)
	}
	slog.Printf("Exported %s to %s", file, dir)
	return nil
}

func writeFile(dir string, fileName string, data []byte) error {
	w, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		return fmt.Errorf("create %s: %w", fileName, err)
	}
	defer w.Close()
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("write %s: %w", fileName, err)
	}
	return nil
}

func Files() []pfs.FileEntry {
	if archive == nil {
		return nil
	}
	return archive.Files()
}

func File(name string) ([]byte, error) {
	if archive == nil {
		return nil, fmt.Errorf("no archive loaded")
	}
	return archive.File(name)
}

func SetFile(name string, data []byte) error {
	if archive == nil {
		return fmt.Errorf("no archive loaded")
	}
	if name == "" {
		name = fileName
	}
	return archive.SetFile(name, data)
}
