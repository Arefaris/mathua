package lessons

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chuma-beep/mathua/internal/i18n"
	"github.com/chuma-beep/mathua/internal/latex"
)

// Translated lesson bodies.
//
// A translation of data/lessons/<corpus>/<path>.md lives at data/lessons/<corpus>/<locale>/<path>.md,
// next to the original and nowhere else: the progress counter (web/next-app/scripts/i18n-progress.mjs)
// and the validator (internal/i18n/lesson_body_test.go) look for it at exactly that path.
//
// Which body a learner gets is decided per lesson, not per paragraph: either the whole translated
// lesson or the whole English one. A lesson without a translation is English, so a language can
// ship one lesson at a time.
//
// The English heading in a step's `section` is the key into the English body. A translated body
// has translated headings, so for a translated lesson the key is the translated section from
// data/lessons/kp/<locale>/<id>.json instead. The two are never mixed: KPs and KPSectionBody work
// on the English lesson exactly as before, and the *In variants below look in whichever body they
// are told the learner is reading, with the sections that belong to it.

// loadTranslations reads every translated body that sits beside a loaded English lesson.
//
// Only the locale directories directly under a corpus directory are read, and only at the path
// a registered English lesson implies. Nothing walks data/lessons, so a locale directory can
// never be mistaken for English content, and a stray file in one is simply never served.
func (l *Loader) loadTranslations(lessonsDir string, bySource map[string]*Lesson) map[i18n.Locale]map[string]*Lesson {
	out := make(map[i18n.Locale]map[string]*Lesson)
	for source, en := range bySource {
		corpus, rest, ok := strings.Cut(filepath.ToSlash(source), "/")
		if !ok {
			continue
		}
		for _, loc := range i18n.Overlays() {
			path := filepath.Join(lessonsDir, corpus, string(loc), filepath.FromSlash(rest))
			raw, err := os.ReadFile(path)
			if err != nil {
				if !os.IsNotExist(err) {
					log.Printf("warning: read %s translation %q: %v", loc, path, err)
				}
				continue
			}
			// The English lesson's adapter, not one sniffed from the translation: both bodies are
			// the same document and must be canonicalised the same way.
			content := latex.Canonicalize(string(raw), latex.ForSource(source, en.Body))
			tr := &Lesson{
				Title:    extractTitle(content),
				Body:     content,
				Concepts: en.Concepts,
				Source:   en.Source,
				Assets:   ExtractAssets(content, en.Concepts),
				source:   source,
				loc:      loc,
			}
			if tr.Title == "" {
				tr.Title = en.Title
			}
			if why, ok := l.translationFits(loc, en, tr); !ok {
				log.Printf("warning: the %s translation of %q is not served, English is used instead: %s", loc, source, why)
				continue
			}
			if out[loc] == nil {
				out[loc] = make(map[string]*Lesson)
			}
			out[loc][source] = tr
		}
	}
	return out
}

// translationFits reports whether tr can stand in for en for every concept that reads en.
//
// A concept's steps point into the lesson body by heading. If a step's heading is found in the
// English body, its translated heading must be found in the translated one; otherwise the step
// would be shown with no worked example, or with the English one under a Russian page. A heading
// that does not resolve in the English body to begin with (an authored slice stands in for it)
// has nothing to preserve and is not checked.
func (l *Loader) translationFits(loc i18n.Locale, en, tr *Lesson) (string, bool) {
	for _, id := range en.Concepts {
		if l.concepts[id] != en {
			continue // another lesson wins this concept
		}
		base := l.kps[id]
		over := l.kpText[loc][id]
		for i, kp := range base {
			if kp.Section == "" {
				continue
			}
			if _, found := extractSection(en.Body, kp.Section); !found {
				continue
			}
			switch {
			case len(over) != len(base):
				return fmt.Sprintf("concept %s has %d %s steps for %d English ones", id, len(over), loc, len(base)), false
			case over[i].Section == "":
				return fmt.Sprintf("concept %s step %d has no %s section", id, i, loc), false
			}
			if _, found := extractSection(tr.Body, over[i].Section); !found {
				return fmt.Sprintf("concept %s step %d: section %q is not a heading of the translation", id, i, over[i].Section), false
			}
		}
	}
	return "", true
}

// Localize returns the translation of lesson into loc when there is a usable one, and lesson
// itself otherwise. English is never translated, and a nil lesson stays nil.
func (l *Loader) Localize(lesson *Lesson, loc i18n.Locale) *Lesson {
	if lesson == nil || loc == i18n.Default {
		return lesson
	}
	if tr, ok := l.translated[loc][lesson.source]; ok {
		return tr
	}
	return lesson
}

// LessonIn is Lesson for a learner who asked for loc: the translated lesson if there is one,
// the English lesson if not. Check Locale on the result to see which one came back.
func (l *Loader) LessonIn(conceptID string, loc i18n.Locale) *Lesson {
	return l.Localize(l.concepts[conceptID], loc)
}

// KPsForBody returns the steps to use for lookups into the lesson body LessonIn serves in loc.
//
// For an English body that is KPs. For a translated body it is KPsIn, because the headings the
// steps point at are the translated ones; the sections there are the keys into that body as
// well as display text. Never use this to look into the English body, and never use KPs to
// look into a translated one.
func (l *Loader) KPsForBody(conceptID string, loc i18n.Locale) []KP {
	if l.LessonIn(conceptID, loc).Locale() != i18n.Default {
		return l.KPsIn(conceptID, loc)
	}
	return l.KPs(conceptID)
}

// KPSectionBodyIn is KPSectionBody in the body LessonIn serves for loc. The section is a heading
// of that body, as KPsForBody returns it.
func (l *Loader) KPSectionBodyIn(conceptID, section string, loc i18n.Locale) (string, bool) {
	lesson := l.LessonIn(conceptID, loc)
	if lesson == nil || section == "" {
		return "", false
	}
	return extractSection(lesson.Body, section)
}

// TeachingSliceIn is TeachingSlice for the body LessonIn serves for loc. kp comes from
// KPsForBody.
func (l *Loader) TeachingSliceIn(conceptID string, kp KP, loc i18n.Locale) (string, bool) {
	if body, ok := l.KPSectionBodyIn(conceptID, kp.Section, loc); ok {
		return body, true
	}
	if s := kp.Slice; s != nil && strings.TrimSpace(s.Rule) != "" {
		return renderSlice(kp.Label, *s), true
	}
	return renderMinimalSlice(kp.Label, l.LessonIn(conceptID, loc)), true
}

// SectionResolves reports whether heading names a heading of body, by the same rule the loader
// uses to find a step's worked example. The translation validator uses it so that it checks what
// the server will do.
func SectionResolves(body, heading string) bool {
	if heading == "" {
		return false
	}
	_, ok := extractSection(body, heading)
	return ok
}
