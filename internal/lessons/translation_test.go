package lessons

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/i18n"
)

// translationFixture is a lessons directory with four lessons:
//
//	a  translated, and its Russian steps point at Russian headings of the Russian body
//	b  not translated
//	c  translated, but the Russian step names a heading the Russian body does not have
//	d  translated, and shared by two concepts, one of which has no Russian steps
func translationFixture(t *testing.T) *Loader {
	t.Helper()
	dir := t.TempDir()
	w := func(rel, content string) { writeFile(t, filepath.Join(dir, filepath.FromSlash(rel)), content) }

	w("lessons.json", `[
		{"concept_id":"a","source":"teaching/a.md"},
		{"concept_id":"b","source":"teaching/b.md"},
		{"concept_id":"c","source":"teaching/c.md"},
		{"concept_id":"d1","source":"teaching/d.md"},
		{"concept_id":"d2","source":"teaching/d.md"}
	]`)

	w("teaching/a.md", "> Hand-authored lesson for the Mathua project.\n\n# Adding\n\n## First Idea\n\nEnglish idea.\n\n![A figure](svg/fig-a.svg)\n\n## Second Idea\n\nEnglish second.\n")
	w("teaching/ru/a.md", "> Hand-authored lesson for the Mathua project.\n\n# Сложение\n\n## Первая мысль\n\nРусская мысль.\n\n![Рисунок](svg/fig-a.svg)\n\n## Вторая мысль\n\nРусское второе.\n")
	w("kp/a.json", `[
		{"label":"One","section":"First Idea","subgoals":["x","y","z"]},
		{"label":"Two","section":"","subgoals":["p","q","r"],"slice":{"rule":"An English rule that is long enough to count as a rule for a slice."}},
		{"label":"Three","section":"Second Idea","subgoals":["s","t","u"]}
	]`)
	w("kp/ru/a.json", `[
		{"label":"Раз","section":"Первая мысль","subgoals":["икс","игрек","зет"]},
		{"label":"Два","section":"","subgoals":["пэ","ку","эр"]},
		{"label":"Три","section":"Вторая мысль","subgoals":["эс","тэ","у"]}
	]`)

	w("teaching/b.md", "# Bee\n\n## Only English\n\nEnglish bee.\n")
	w("kp/b.json", `[{"label":"Bee","section":"Only English","subgoals":["s","t","u"]}]`)
	w("kp/ru/b.json", `[{"label":"Пчела","section":"Только английский","subgoals":["с","т","у"]}]`)

	w("teaching/c.md", "# Cee\n\n## Real Heading\n\nEnglish cee.\n")
	w("teaching/ru/c.md", "# Цэ\n\n## Настоящий заголовок\n\nРусское цэ.\n")
	w("kp/c.json", `[{"label":"Cee","section":"Real Heading","subgoals":["s","t","u"]}]`)
	w("kp/ru/c.json", `[{"label":"Цэ","section":"Заголовок, которого нет","subgoals":["с","т","у"]}]`)

	w("teaching/d.md", "# Dee\n\n## Shared Heading\n\nEnglish dee.\n")
	w("teaching/ru/d.md", "# Дэ\n\n## Общий заголовок\n\nРусское дэ.\n")
	w("kp/d1.json", `[{"label":"D one","section":"Shared Heading","subgoals":["s","t","u"]}]`)
	w("kp/d2.json", `[{"label":"D two","section":"Shared Heading","subgoals":["s","t","u"]}]`)
	w("kp/ru/d1.json", `[{"label":"Дэ раз","section":"Общий заголовок","subgoals":["с","т","у"]}]`)

	l, err := Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	return l
}

func TestLessonIsServedInTheRequestedLanguage(t *testing.T) {
	l := translationFixture(t)

	ru := l.LessonIn("a", i18n.RU)
	if ru.Locale() != i18n.RU || ru.Title != "Сложение" || !strings.Contains(ru.Body, "Русская мысль") {
		t.Fatalf("ru lesson = %q %q (%s)", ru.Title, ru.Body, ru.Locale())
	}
	if ru.File() != "teaching/a.md" {
		t.Errorf("a translation reports the English file it translates, got %q", ru.File())
	}
	if len(ru.Assets) != 1 || ru.Assets[0].ID != "fig-a.svg" || ru.Assets[0].Section != "Первая мысль" {
		t.Errorf("the translated body has its own figure inventory: %+v", ru.Assets)
	}

	// English, an unknown language and the language of no translation are the original, and
	// the original is untouched by loading a translation next to it.
	en := l.Lesson("a")
	for _, loc := range []i18n.Locale{i18n.EN, i18n.UK, "xx"} {
		if got := l.LessonIn("a", loc); got != en || got.Locale() != i18n.EN {
			t.Errorf("%s: got %q (%s), want the English lesson", loc, got.Title, got.Locale())
		}
	}
	if en.Title != "Adding" || !strings.Contains(en.Body, "English idea") {
		t.Errorf("English lesson changed: %q", en.Title)
	}
}

func TestUntranslatedLessonStaysEnglishPerLesson(t *testing.T) {
	l := translationFixture(t)

	got := l.LessonIn("b", i18n.RU)
	if got != l.Lesson("b") || got.Locale() != i18n.EN {
		t.Errorf("b has no translation and must be English, got %s %q", got.Locale(), got.Title)
	}
	// Its steps are still shown in Russian, which is a separate overlay and was already so.
	if kps := l.KPsIn("b", i18n.RU); kps[0].Label != "Пчела" {
		t.Errorf("step label = %q", kps[0].Label)
	}
	// Lookups into its English body use the English section, exactly as before.
	kps := l.KPsForBody("b", i18n.RU)
	if kps[0].Section != "Only English" {
		t.Errorf("lookup section = %q, want the English heading", kps[0].Section)
	}
	body, _ := l.TeachingSliceIn("b", kps[0], i18n.RU)
	if !strings.Contains(body, "English bee") {
		t.Errorf("slice = %q", body)
	}
	if got := l.Localize(nil, i18n.RU); got != nil {
		t.Error("Localize(nil) must stay nil")
	}
	if got := l.LessonIn("nope", i18n.RU); got != nil {
		t.Error("an unknown concept has no lesson")
	}
}

func TestSectionsAreLookedUpInTheBodyBeingServed(t *testing.T) {
	l := translationFixture(t)

	// Translated lesson: the Russian section finds the Russian body.
	ru := l.KPsForBody("a", i18n.RU)
	if ru[0].Section != "Первая мысль" || ru[2].Section != "Вторая мысль" {
		t.Fatalf("sections = %q, %q", ru[0].Section, ru[2].Section)
	}
	body, ok := l.TeachingSliceIn("a", ru[0], i18n.RU)
	if !ok || !strings.Contains(body, "Русская мысль") || strings.Contains(body, "Русское второе") {
		t.Errorf("ru step 0 = %q", body)
	}
	body, _ = l.TeachingSliceIn("a", ru[2], i18n.RU)
	if !strings.Contains(body, "Русское второе") {
		t.Errorf("ru step 2 = %q", body)
	}

	// The English section does not exist in the Russian body, and the Russian one does not
	// exist in the English body: the two are never mixed.
	if _, ok := l.KPSectionBodyIn("a", "First Idea", i18n.RU); ok {
		t.Error("an English heading resolved in the Russian body")
	}
	if _, ok := l.KPSectionBodyIn("a", "Первая мысль", i18n.EN); ok {
		t.Error("a Russian heading resolved in the English body")
	}

	// English requests behave exactly as they did before translations existed.
	en := l.KPsForBody("a", i18n.EN)
	if en[0].Section != "First Idea" || en[0].Label != "One" {
		t.Errorf("en step = %+v", en[0])
	}
	body, _ = l.TeachingSlice("a", en[0])
	if !strings.Contains(body, "English idea") {
		t.Errorf("en step 0 = %q", body)
	}
	if body, ok := l.KPSectionBody("a", "First Idea"); !ok || !strings.Contains(body, "English idea") {
		t.Errorf("KPSectionBody = %q, %v", body, ok)
	}

	// A step with no section still gets its authored slice, in both languages: the slice text
	// is carried on the English shard.
	body, _ = l.TeachingSliceIn("a", ru[1], i18n.RU)
	if !strings.Contains(body, "An English rule") || !strings.Contains(body, "Два") {
		t.Errorf("authored slice in ru = %q", body)
	}
	body, _ = l.TeachingSlice("a", en[1])
	if !strings.Contains(body, "An English rule") || !strings.Contains(body, "Two") {
		t.Errorf("authored slice in en = %q", body)
	}
}

func TestFiguresFollowTheServedBody(t *testing.T) {
	l := translationFixture(t)

	ru := l.KPsForBody("a", i18n.RU)
	figs := l.KPAssetsIn("a", ru[0], i18n.RU)
	if len(figs) != 1 || figs[0].ID != "fig-a.svg" || figs[0].Src != "svg/fig-a.svg" {
		t.Errorf("ru figures = %+v", figs)
	}
	if got := l.KPAssetsIn("a", ru[2], i18n.RU); len(got) != 0 {
		t.Errorf("the figure is not in the second section: %+v", got)
	}
	en := l.KPsForBody("a", i18n.EN)
	figs = l.KPAssets("a", en[0])
	if len(figs) != 1 || figs[0].Alt != "A figure" {
		t.Errorf("en figures = %+v", figs)
	}
	if figs := l.KPAssetsIn("a", en[0], i18n.RU); len(figs) != 0 {
		t.Errorf("an English section must not select figures from the Russian body: %+v", figs)
	}
}

// A translation that its steps cannot find their headings in is not served at all: the learner
// reads the English lesson whole, and the steps keep their English sections to find it with.
func TestTranslationThatTheStepsCannotUseIsNotServed(t *testing.T) {
	l := translationFixture(t)

	if got := l.LessonIn("c", i18n.RU); got.Locale() != i18n.EN || got != l.Lesson("c") {
		t.Fatalf("c was served in %s although its Russian step names no heading of it", got.Locale())
	}
	kps := l.KPsForBody("c", i18n.RU)
	if kps[0].Section != "Real Heading" {
		t.Errorf("section = %q", kps[0].Section)
	}
	body, _ := l.TeachingSliceIn("c", kps[0], i18n.RU)
	if !strings.Contains(body, "English cee") {
		t.Errorf("slice = %q", body)
	}
}

// One translated file serves every concept that reads the lesson, so it is served only if it
// fits all of them. d2 has no Russian steps, so d1 gets English too.
func TestSharedLessonIsTranslatedOnlyIfEveryConceptFits(t *testing.T) {
	l := translationFixture(t)
	for _, id := range []string{"d1", "d2"} {
		if got := l.LessonIn(id, i18n.RU); got.Locale() != i18n.EN {
			t.Errorf("%s: served %s, want English because d2 has no Russian steps", id, got.Locale())
		}
	}
}

// A locale directory is not English content: it adds no lesson, no concept and no step.
func TestLocaleDirectoryIsNeverEnglishContent(t *testing.T) {
	l := translationFixture(t)
	if n := l.Count(); n != 5 {
		t.Errorf("Count = %d, want the 5 registered concepts", n)
	}
	for id, lesson := range l.All() {
		if lesson.Locale() != i18n.EN {
			t.Errorf("%s: the English catalog holds a %s lesson", id, lesson.Locale())
		}
		if strings.Contains(lesson.Body, "Русск") {
			t.Errorf("%s: Russian text in an English body", id)
		}
	}
	if got := l.LessonByTitle("Сложение"); got != nil {
		t.Error("lessons are found by their English title only")
	}
	if got := l.LessonByTitle("Adding"); got == nil || got.Locale() != i18n.EN {
		t.Error("the English title finds the English lesson")
	}
}
