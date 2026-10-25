// Package load loads the files used by glx.web
package load

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/genealogix/glx/go-glx"
)

func GLX(fsys fs.FS) (*glx.GLXFile, error) {
	files, err := readGLXFiles(fsys)
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

func readGLXFiles(fsys fs.FS) (map[string][]byte, error) {
	files := make(map[string][]byte)

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and dirs.
		if strings.HasPrefix(path, ".") && path != "." {
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

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		files[path] = data

		return nil
	})

	return files, err
}
