// Package load loads a glx.GLXFile from a GLX archive.
package load

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/genealogix/glx/go-glx"
)

func GLXFile(archiveDir string) (*glx.GLXFile, error) {
	files, err := readGLXFiles(archiveDir)
	if err != nil {
		return nil, err
	}

	s := glx.NewSerializer(nil)

	g, dups, err := s.DeserializeMultiFileFromMap(files)
	if err != nil {
		return nil, err
	}

	if len(dups) > 0 {
		return nil, fmt.Errorf("duplicates: %s", dups)
	}

	return g, nil
}

func readGLXFiles(archiveDir string) (map[string][]byte, error) {
	files := make(map[string][]byte)

	err := filepath.WalkDir(archiveDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and dirs.
		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		// Nothing to do for dirs except walk into them.
		if d.IsDir() {
			return nil
		}

		// Ignore files other than .glx files.
		if filepath.Ext(d.Name()) != glx.FileExtGLX {
			return nil
		}

		cleanPath := filepath.Clean(path)

		data, err := os.ReadFile(cleanPath)
		if err != nil {
			return fmt.Errorf("reading %s: %w", cleanPath, err)
		}

		relPath, err := filepath.Rel(archiveDir, cleanPath)
		if err != nil {
			return fmt.Errorf("getting relative path: %w", err)
		}

		files[relPath] = data

		return nil
	})

	return files, err
}
