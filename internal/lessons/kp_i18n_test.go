package lessons

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/chuma-beep/mathua/internal/i18n"
)

// kpFixture writes a lessons dir with English steps for "a" and "b" and Russian
// overlays whose shapes differ in the ways KPsIn has to handle.
func kpFixture(t *testing.T) *Loader {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lessons.json"), `[{"concept_id":"a","source":"a.md"},{"concept_id":"b","source":"b.md"}]`)
	writeFile(t, filepath.Join(dir, "a.md"), "# A\n\n## Section One\n\nbody")
	writeFile(t, filepath.Join(dir, "b.md"), "# B\n\nbody")
	os.MkdirAll(filepath.Join(dir, "kp", "ru"), 0o755)
	writeFile(t, filepath.Join(dir, "kp", "a.json"), `[
		{"label":"One","section":"Section One","subgoals":["x","y","z"],"asset_ids":["fig-1"]},
		{"label":"Two","section":"","subgoals":["p","q","r"]}
	]`)
	writeFile(t, filepath.Join(dir, "kp", "b.json"), `[{"label":"Bee","section":"Bee","subgoals":["s","t","u"]}]`)
	writeFile(t, filepath.Join(dir, "kp", "ru", "a.json"), `[
		{"label":"Раз","section":"Раздел один","subgoals":["икс","игрек","зет"]},
		{"label":"Два","section":"Раздел, которого нет","subgoals":["пэ","ку"]}
	]`)
	// b has the wrong number of steps, so it must be ignored as a whole.
	writeFile(t, filepath.Join(dir, "kp", "ru", "b.json"), `[
		{"label":"Пчела","section":"Пчела","subgoals":["с","т","у"]},
		{"label":"Лишний","section":"Лишний","subgoals":["а","б","в"]}
	]`)
	l, err := Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	return l
}

func TestKPOverlaySubdirectoryIsNotEnglishContent(t *testing.T) {
	l := kpFixture(t)
	if _, ok := l.kps["ru"]; ok {
		t.Fatal("the ru directory was read as a concept's shard")
	}
	if len(l.kps) != 2 {
		t.Fatalf("English shards = %d, want 2", len(l.kps))
	}
	if got := l.KPs("a")[0].Label; got != "One" {
		t.Errorf("English label = %q", got)
	}
}

func TestKPsInTranslatesAndFallsBack(t *testing.T) {
	l := kpFixture(t)

	ru := l.KPsIn("a", i18n.RU)
	if len(ru) != 2 {
		t.Fatalf("steps = %d, want 2", len(ru))
	}
	if ru[0].Label != "Раз" || ru[0].Section != "Раздел один" || !reflect.DeepEqual(ru[0].Subgoals, []string{"икс", "игрек", "зет"}) {
		t.Errorf("step 0 = %+v", ru[0])
	}
	if !reflect.DeepEqual(ru[0].AssetIDs, []string{"fig-1"}) {
		t.Errorf("non-text fields must be copied from the English step: %+v", ru[0].AssetIDs)
	}
	// Step 1 has no English section, so it must not gain one, and its
	// two-item subgoal list does not match the English three.
	if ru[1].Label != "Два" {
		t.Errorf("label = %q", ru[1].Label)
	}
	if ru[1].Section != "" {
		t.Errorf("section = %q, want it to stay empty", ru[1].Section)
	}
	if !reflect.DeepEqual(ru[1].Subgoals, []string{"p", "q", "r"}) {
		t.Errorf("a short translated subgoal list must fall back to English, got %v", ru[1].Subgoals)
	}

	// The English steps, used for lookups, are untouched.
	en := l.KPs("a")
	if en[0].Section != "Section One" || en[0].Label != "One" {
		t.Errorf("English steps were edited: %+v", en[0])
	}
	if _, ok := l.KPSectionBody("a", en[0].Section); !ok {
		t.Error("the English section no longer resolves in the lesson body")
	}

	// A step-count mismatch ignores the overlay for the concept.
	if got := l.KPsIn("b", i18n.RU); got[0].Label != "Bee" || len(got) != 1 {
		t.Errorf("mismatched overlay = %+v", got)
	}
	// English and an unknown locale are the originals.
	if got := l.KPsIn("a", i18n.EN); got[0].Label != "One" {
		t.Errorf("en = %q", got[0].Label)
	}
	// No shard at all: the lesson-title fallback step comes back unchanged.
	if got := l.KPsIn("nope", i18n.RU); got != nil {
		t.Errorf("unknown concept = %+v", got)
	}
}
