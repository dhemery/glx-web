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
	g             *glx.GLXFile
	Assertions    map[string]*Assertion
	Citations     map[string]*Citation
	Events        map[string]*Event
	Media         map[string]*Media
	Persons       map[string]*Person
	Places        map[string]*Place
	Relationships map[string]*Relationship
	Repositories  map[string]*Repository
	Sources       map[string]*Source
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
		g:             g,
		Assertions:    make(map[string]*Assertion),
		Citations:     make(map[string]*Citation),
		Events:        make(map[string]*Event),
		Media:         make(map[string]*Media),
		Persons:       make(map[string]*Person),
		Places:        make(map[string]*Place),
		Relationships: make(map[string]*Relationship),
		Repositories:  make(map[string]*Repository),
		Sources:       make(map[string]*Source),
	}

	for id, ga := range g.Assertions {
		a.Assertions[id] = &Assertion{a: a, g: ga, ID: id}
	}

	for id, gc := range g.Citations {
		a.Citations[id] = newCitation(id, gc, a)
	}

	for id, ge := range g.Events {
		a.Events[id] = newEvent(id, ge, a)
	}

	for id, gm := range g.Media {
		a.Media[id] = &Media{a: a, g: gm, ID: id}
	}

	for id, gp := range g.Persons {
		a.Persons[id] = newPerson(id, gp, a)
	}

	for id, gp := range g.Places {
		a.Places[id] = &Place{a: a, g: gp, ID: id}
	}

	for id, gr := range g.Relationships {
		a.Relationships[id] = &Relationship{a: a, g: gr, ID: id}
	}

	for id, gr := range g.Repositories {
		a.Repositories[id] = newRepository(id, gr, a)
	}

	for id, gs := range g.Sources {
		a.Sources[id] = &Source{a: a, g: gs, ID: id}
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
