package planning

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chuma-beep/mathua/internal/i18n"
)

func TestOverlayFallsBackAndNeverMutates(t *testing.T) {
	p := testPlanner(t)
	p.SetDestinations([]*Destination{{ID: "dest", Name: "Destination", Description: "Where to", Courses: []string{"test"}}})

	dir := t.TempDir()
	ruDir := filepath.Join(dir, "ru")
	if err := os.MkdirAll(ruDir, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(ruDir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// "Description" is left out on purpose, and "ghost" is not a course.
	write(CoursesOverlayFile, `{"test":{"name":"Тест"},"ghost":{"name":"Призрак","description":"Нет"}}`)
	write(DestinationsOverlayFile, `{"dest":{"name":"Цель","description":"Куда идти"}}`)
	if err := p.LoadOverlays(dir); err != nil {
		t.Fatalf("load: %v", err)
	}

	en := p.Course("test")
	ru := p.CourseIn(en, i18n.RU)
	if ru.Name != "Тест" {
		t.Errorf("name = %q", ru.Name)
	}
	if ru.Description != "d" {
		t.Errorf("an untranslated description must stay English, got %q", ru.Description)
	}
	if en.Name != "Test" {
		t.Errorf("the shared English course was edited: %q", en.Name)
	}
	if ru == en {
		t.Error("CourseIn returned the shared pointer for a translated course")
	}
	if got := p.CourseIn(en, i18n.EN); got != en {
		t.Error("English must return the original")
	}
	if p.CourseIn(nil, i18n.RU) != nil {
		t.Error("nil course")
	}

	d := p.DestinationIn(p.Destination("dest"), i18n.RU)
	if d.Name != "Цель" || d.Description != "Куда идти" {
		t.Errorf("destination = %+v", d)
	}
	if p.Destination("dest").Name != "Destination" {
		t.Error("the shared English destination was edited")
	}
}

func TestOverlayMissingFilesAreEnglish(t *testing.T) {
	p := testPlanner(t)
	if err := p.LoadOverlays(t.TempDir()); err != nil {
		t.Fatalf("missing files: %v", err)
	}
	c := p.Course("test")
	if got := p.CourseIn(c, i18n.RU); got != c {
		t.Error("no overlay should return the original course")
	}
}
