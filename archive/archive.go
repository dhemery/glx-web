// Package archive represents a GLX archive in a form suitable for Go templates to consume.
package archive

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/genealogix/glx/go-glx"
)

type Archive struct {
	Events  map[string]*Event
	Persons map[string]*Person
	Places  map[string]*Place
}

func Load(path string) (*Archive, error) {
	files, err := readGLXFiles(path)
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

	return compile(g), nil
}

func compile(g *glx.GLXFile) *Archive {
	a := &Archive{
		Events:  make(map[string]*Event),
		Persons: make(map[string]*Person),
		Places:  make(map[string]*Place),
	}

	for id, ge := range g.Events {
		a.Events[id] = newEvent(id, ge, a)
	}

	for id, gp := range g.Persons {
		a.Persons[id] = newPerson(id, gp)
	}

	for id, gp := range g.Places {
		a.Places[id] = newPlace(id, gp, a)
	}

	return a
}

func readGLXFiles(rootDir string) (map[string][]byte, error) {
	files := make(map[string][]byte)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
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

		relPath, err := filepath.Rel(rootDir, cleanPath)
		if err != nil {
			return fmt.Errorf("getting relative path: %w", err)
		}

		files[relPath] = data

		return nil
	})

	return files, err
}
