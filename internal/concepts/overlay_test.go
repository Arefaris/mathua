package concepts

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chuma-beep/mathua/internal/i18n"
)

func overlayDAG(t *testing.T) *DAG {
	t.Helper()
	mt := MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}
	d, err := Build([]Concept{
		{ID: "a", Label: "Alpha", Domain: "d", MasteryThreshold: mt},
		{ID: "b", Label: "Beta", Domain: "d", Prerequisites: []string{"a"}, MasteryThreshold: mt},
		{ID: "c", Label: "Gamma", Domain: "d", Prerequisites: []string{"a"}, MasteryThreshold: mt},
	})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	return d
}

func TestLabelFallsBackPerItem(t *testing.T) {
	d := overlayDAG(t)
	unknown := d.SetLabels(i18n.RU, map[string]string{
		"a":    "Альфа",
		"c":    "", // blank is "not translated", never a blank label
		"gone": "Исчезнувшая",
	})
	if len(unknown) != 1 || unknown[0] != "gone" {
		t.Fatalf("unknown = %v, want [gone]", unknown)
	}
	cases := []struct {
		id   string
		loc  i18n.Locale
		want string
	}{
		{"a", i18n.RU, "Альфа"},
		{"b", i18n.RU, "Beta"}, // missing from the overlay
		{"c", i18n.RU, "Gamma"},
		{"a", i18n.EN, "Alpha"}, // English is never touched
		{"a", "", "Alpha"},      // a zero locale is English too
	}
	for _, c := range cases {
		if got := d.Label(d.Concept(c.id), c.loc); got != c.want {
			t.Errorf("Label(%s, %q) = %q, want %q", c.id, c.loc, got, c.want)
		}
	}
	if got := d.Concept("a").Label; got != "Alpha" {
		t.Errorf("the English concept was edited: %q", got)
	}
}

func TestLabelByID(t *testing.T) {
	d := overlayDAG(t)
	d.SetLabels(i18n.RU, map[string]string{"a": "Альфа"})
	if got := d.LabelByID("a", i18n.RU, "x"); got != "Альфа" {
		t.Errorf("known id = %q", got)
	}
	if got := d.LabelByID("nope", i18n.RU, "fallback"); got != "fallback" {
		t.Errorf("unknown id = %q, want the fallback", got)
	}
	if got := d.Label(nil, i18n.RU); got != "" {
		t.Errorf("nil concept = %q", got)
	}
}

func TestSetLabelsIgnoresDefaultLocale(t *testing.T) {
	d := overlayDAG(t)
	d.SetLabels(i18n.EN, map[string]string{"a": "Changed"})
	if got := d.Label(d.Concept("a"), i18n.EN); got != "Alpha" {
		t.Errorf("an English overlay replaced the original: %q", got)
	}
}

func TestLoadOverlays(t *testing.T) {
	d := overlayDAG(t)
	dir := t.TempDir()
	// No file for ru is not an error: that locale is simply English.
	if err := d.LoadOverlays(dir); err != nil {
		t.Fatalf("missing overlay: %v", err)
	}
	if got := d.Label(d.Concept("a"), i18n.RU); got != "Alpha" {
		t.Fatalf("no overlay, got %q", got)
	}

	ruDir := filepath.Join(dir, "ru")
	if err := os.MkdirAll(ruDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ruDir, ConceptsOverlayFile), []byte(`{"a":"Альфа","stale":"Старая"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := d.LoadOverlays(dir); err != nil {
		t.Fatalf("load: %v", err)
	}
	if got := d.Label(d.Concept("a"), i18n.RU); got != "Альфа" {
		t.Errorf("a = %q", got)
	}

	// A corrupt file is reported, so the caller can log it and carry on in English.
	if err := os.WriteFile(filepath.Join(ruDir, ConceptsOverlayFile), []byte(`{not json`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := overlayDAG(t).LoadOverlays(dir); err == nil {
		t.Error("a corrupt overlay was accepted silently")
	}
}
