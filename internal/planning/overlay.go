package planning

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chuma-beep/mathua/internal/i18n"
)

// Overlay file names inside a locale directory (<i18n dir>/<locale>/...). Each is an object
// keyed by the English file's id: {"<id>": {"name": ..., "description": ...}}.
const (
	CoursesOverlayFile      = "courses.json"
	DestinationsOverlayFile = "destinations.json"
)

// Text is the translated, human-readable part of a course or destination. The ids, grades and
// targets stay in the English file; only these two fields are ever replaced.
type Text struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type overlays struct {
	courses      map[i18n.Locale]map[string]Text
	destinations map[i18n.Locale]map[string]Text
}

// LoadOverlays reads the translated course and destination text for every overlay locale
// found under i18nDir. A missing file means that locale is English for that catalog; an
// unreadable one is an error. Call it at startup, after SetDestinations.
func (p *Planner) LoadOverlays(i18nDir string) error {
	p.translated = overlays{
		courses:      map[i18n.Locale]map[string]Text{},
		destinations: map[i18n.Locale]map[string]Text{},
	}
	for _, loc := range i18n.Overlays() {
		for _, f := range []struct {
			file string
			dst  map[i18n.Locale]map[string]Text
			has  func(id string) bool
		}{
			{CoursesOverlayFile, p.translated.courses, func(id string) bool { return p.byID[id] != nil }},
			{DestinationsOverlayFile, p.translated.destinations, func(id string) bool { return p.destByID[id] != nil }},
		} {
			data, err := os.ReadFile(filepath.Join(i18nDir, string(loc), f.file))
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return fmt.Errorf("read %s overlay (%s): %w", f.file, loc, err)
			}
			var m map[string]Text
			if err := json.Unmarshal(data, &m); err != nil {
				return fmt.Errorf("parse %s overlay (%s): %w", f.file, loc, err)
			}
			for id := range m {
				if !f.has(id) {
					log.Printf("warning: %s %s overlay names %q, which is not in the catalog", loc, f.file, id)
					delete(m, id)
				}
			}
			f.dst[loc] = m
		}
	}
	return nil
}

// CourseIn returns c with its name and description in loc. The result is a copy: courses are
// shared by every request, so the one request that wants Russian must not rewrite them for the
// next. A field the overlay does not translate keeps its English text.
func (p *Planner) CourseIn(c *Course, loc i18n.Locale) *Course {
	if c == nil {
		return nil
	}
	t, ok := p.translated.courses[loc][c.ID]
	if !ok {
		return c
	}
	out := *c
	if t.Name != "" {
		out.Name = t.Name
	}
	if t.Description != "" {
		out.Description = t.Description
	}
	return &out
}

// DestinationIn is CourseIn for a destination bundle.
func (p *Planner) DestinationIn(d *Destination, loc i18n.Locale) *Destination {
	if d == nil {
		return nil
	}
	t, ok := p.translated.destinations[loc][d.ID]
	if !ok {
		return d
	}
	out := *d
	if t.Name != "" {
		out.Name = t.Name
	}
	if t.Description != "" {
		out.Description = t.Description
	}
	return &out
}
