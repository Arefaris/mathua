package concepts

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chuma-beep/mathua/internal/i18n"
)

// ConceptsOverlayFile is the name of the label overlay inside a locale directory:
// <i18n dir>/<locale>/concepts.json, a flat {"<concept id>": "<label>"} object. The frontend
// build script reads the same file, so its path and shape are part of a contract with it.
const ConceptsOverlayFile = "concepts.json"

// SetLabels installs the translated labels for loc. The English originals in data/concepts are
// never edited: a label the overlay does not have, or has blank, is served in English, so a
// translation can be shipped in parts.
//
// An id that is not in the DAG is dropped and returned, because it means the corpus moved
// (a concept renamed or removed) and the overlay was not updated. It is reported rather than
// fatal: an out-of-date translation must not stop the server from starting. The overlay test
// is the place that fails the build.
//
// Call it while the server is starting; the overlay is read without a lock afterwards.
func (d *DAG) SetLabels(loc i18n.Locale, labels map[string]string) (unknown []string) {
	if loc == i18n.Default {
		return nil
	}
	if d.labels == nil {
		d.labels = make(map[i18n.Locale]map[string]string)
	}
	m := make(map[string]string, len(labels))
	for id, label := range labels {
		if _, ok := d.concepts[id]; !ok {
			unknown = append(unknown, id)
			continue
		}
		if label == "" {
			continue
		}
		m[id] = label
	}
	d.labels[loc] = m
	return unknown
}

// Label is the concept's display name in loc, falling back to the English label.
func (d *DAG) Label(c *Concept, loc i18n.Locale) string {
	if c == nil {
		return ""
	}
	if l, ok := d.labels[loc][c.ID]; ok {
		return l
	}
	return c.Label
}

// LabelByID is Label for a concept id. An id the DAG does not know yields fallback, which is
// what callers that echo a bare id into a response already do.
func (d *DAG) LabelByID(id string, loc i18n.Locale, fallback string) string {
	if c := d.concepts[id]; c != nil {
		return d.Label(c, loc)
	}
	return fallback
}

// LoadLabelOverlay reads one flat {"id": "label"} overlay file.
func LoadLabelOverlay(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return m, nil
}

// LoadOverlays installs every locale's label overlay found under i18nDir
// (data/i18n/<locale>/concepts.json). A locale with no file is simply English.
func (d *DAG) LoadOverlays(i18nDir string) error {
	for _, loc := range i18n.Overlays() {
		path := filepath.Join(i18nDir, string(loc), ConceptsOverlayFile)
		labels, err := LoadLabelOverlay(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("concept labels (%s): %w", loc, err)
		}
		if unknown := d.SetLabels(loc, labels); len(unknown) > 0 {
			log.Printf("warning: %s concept labels name %d id(s) that are not in the corpus (first: %q)", loc, len(unknown), unknown[0])
		}
	}
	return nil
}
