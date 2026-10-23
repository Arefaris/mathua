package i18n_test

// These tests read every shipped overlay in data/i18n/<locale> and data/lessons/kp/<locale>
// and hold it to the English originals it translates. They are the check that a translation
// cannot silently drift from the corpus: an id that no longer exists, a label that was left
// blank, a step with a different number of subgoals, or a formula that did not survive.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"unicode"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/i18n"
	"github.com/chuma-beep/mathua/internal/planning"
)

const (
	dataDir = "../../data"
	kpDir   = dataDir + "/lessons/kp"
)

// overlayLocales are the supported translations that have started: a locale without a
// data/i18n/<locale> directory is not begun yet, and is skipped rather than failed, so adding a
// language to i18n.Supported and translating it can land as separate commits.
func overlayLocales(t *testing.T) []i18n.Locale {
	t.Helper()
	var out []i18n.Locale
	for _, loc := range i18n.Overlays() {
		if _, err := os.Stat(dataDir + "/i18n/" + string(loc)); err != nil {
			t.Logf("%s: no overlays yet, skipped", loc)
			continue
		}
		out = append(out, loc)
	}
	return out
}

// forbiddenLetters are letters that belong to a neighbouring language and so mark a slip, not
// a translation: Russian ы, э, ъ, ё in Ukrainian text, and Ukrainian ґ, є, і, ї in Russian.
var forbiddenLetters = map[i18n.Locale]string{
	i18n.RU: "ґєіїҐЄІЇ",
	i18n.UK: "ыэъёЫЭЪЁ",
}

// checkLetters reports text in the wrong alphabet for its locale.
func checkLetters(t *testing.T, loc i18n.Locale, where, text string) {
	t.Helper()
	if !hasCyrillic(text) {
		t.Errorf("%s: %q has no %s text", where, text, loc)
	}
	if bad := forbiddenLetters[loc]; bad != "" && strings.ContainsAny(text, bad) {
		t.Errorf("%s: %q uses letters that are not %s (%s)", where, text, loc, bad)
	}
}

// schoolDomains are the domains whose lesson steps every locale must translate. Any other domain
// may be translated too, but only whole (see checkKPShards).
var schoolDomains = []string{"arithmetic", "fractions", "prealgebra", "algebra", "geometry"}

func readJSON(t *testing.T, path string, v interface{}) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
}

// corpus returns every concept id in the English corpus, keyed to its domain, using the same
// loader the server uses, so enrichment.json and any future non-concept file are excluded
// exactly as they are in production.
func corpus(t *testing.T) map[string]string {
	t.Helper()
	dag, err := concepts.LoadDir(dataDir + "/concepts")
	if err != nil {
		t.Fatalf("load corpus: %v", err)
	}
	out := make(map[string]string, dag.Count())
	for id, c := range dag.Concepts() {
		out[id] = c.Domain
	}
	return out
}

var (
	// A formula or a placeholder the translator must carry over untouched: a $...$ span, a
	// LaTeX delimiter or command, or a token with a percent, brace, caret or underscore in
	// it. Plain words and signs such as "+" and "=" are translated around.
	protectedRe = regexp.MustCompile(`\$[^$]*\$|\\\\?[()\[\]]|\\[A-Za-z]+|\S*[%{}^_]\S*`)
	digitsRe    = regexp.MustCompile(`\d+`)
	trimRe      = regexp.MustCompile(`^[\s.,;:!?"'«»]+|[\s.,;:!?"'«»]+$`)
	// An English word of four letters or more inside a would-be formula ("base^(ratio)")
	// is prose in code clothing: it is translated, so it is not checked.
	englishWordRe = regexp.MustCompile(`[A-Za-z]{4,}`)
	// SI unit abbreviations that a Russian text writes in Russian letters (m^2 becomes м^2).
	unitRe = regexp.MustCompile(`^[cdk]?m\^[23]$`)
)

// checkCarried reports what a translation lost from its English source: a protected token
// (LaTeX, placeholder, code-like term) or a number.
func checkCarried(t *testing.T, where, en, tr string) {
	t.Helper()
	for _, tok := range protectedRe.FindAllString(en, -1) {
		tok = trimRe.ReplaceAllString(tok, "")
		if tok == "" || unitRe.MatchString(tok) {
			continue
		}
		// Math spans and LaTeX commands are always checked; a bare token is skipped when it
		// is mostly English words.
		if !strings.HasPrefix(tok, "$") && !strings.HasPrefix(tok, `\`) && englishWordRe.MatchString(tok) {
			continue
		}
		if !strings.Contains(tr, tok) {
			t.Errorf("%s: %q lost %q\n  en: %s\n  tr: %s", where, tr, tok, en, tr)
		}
	}
	for _, d := range digitsRe.FindAllString(en, -1) {
		if !strings.Contains(tr, d) {
			t.Errorf("%s: the number %s is missing\n  en: %s\n  tr: %s", where, d, en, tr)
		}
	}
}

func hasCyrillic(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Cyrillic, r) {
			return true
		}
	}
	return false
}

func TestConceptLabelOverlays(t *testing.T) {
	ids := corpus(t)

	// The English label comes from the corpus files, not from the loaded DAG, which is
	// the same data but keeps this test about the files.
	en := map[string]string{}
	files, _ := filepath.Glob(dataDir + "/concepts/*.json")
	for _, f := range files {
		if filepath.Base(f) == "enrichment.json" {
			continue
		}
		var batch []struct{ ID, Label string }
		readJSON(t, f, &batch)
		for _, c := range batch {
			en[c.ID] = c.Label
		}
	}

	for _, loc := range overlayLocales(t) {
		t.Run(string(loc), func(t *testing.T) {
			var tr map[string]string
			readJSON(t, dataDir+"/i18n/"+string(loc)+"/concepts.json", &tr)
			if len(tr) != len(ids) {
				t.Errorf("concepts.json has %d labels for %d concepts", len(tr), len(ids))
			}
			for id := range tr {
				if _, ok := ids[id]; !ok {
					t.Errorf("concepts.json names %q, which is not a concept in the corpus", id)
				}
			}
			for id := range ids {
				label, ok := tr[id]
				switch {
				case !ok:
					t.Errorf("no %s label for %s", loc, id)
				case strings.TrimSpace(label) == "":
					t.Errorf("empty %s label for %s", loc, id)
				case label != strings.TrimSpace(label):
					t.Errorf("%s: label %q has stray whitespace", id, label)
				default:
					checkLetters(t, loc, id, label)
					checkCarried(t, id, en[id], label)
				}
			}
		})
	}
}

func TestCourseAndDestinationOverlays(t *testing.T) {
	check := func(t *testing.T, loc i18n.Locale, file string, ids []string) {
		var tr map[string]planning.Text
		readJSON(t, dataDir+"/i18n/"+string(loc)+"/"+file, &tr)
		if len(tr) != len(ids) {
			t.Errorf("%s has %d entries for %d ids", file, len(tr), len(ids))
		}
		for _, id := range ids {
			e, ok := tr[id]
			if !ok {
				t.Errorf("%s: no entry for %q", file, id)
				continue
			}
			if strings.TrimSpace(e.Name) == "" || strings.TrimSpace(e.Description) == "" {
				t.Errorf("%s: %q has an empty name or description", file, id)
			}
			checkLetters(t, loc, file+" "+id+" name", e.Name)
			checkLetters(t, loc, file+" "+id+" description", e.Description)
		}
		for id := range tr {
			found := false
			for _, want := range ids {
				found = found || want == id
			}
			if !found {
				t.Errorf("%s names %q, which is not in the English file", file, id)
			}
		}
	}

	var courses []struct{ ID string }
	readJSON(t, dataDir+"/courses.json", &courses)
	var courseIDs []string
	for _, c := range courses {
		courseIDs = append(courseIDs, c.ID)
	}
	if len(courseIDs) != 28 {
		t.Fatalf("expected 28 courses, found %d", len(courseIDs))
	}

	var dests []struct{ ID string }
	readJSON(t, dataDir+"/destinations.json", &dests)
	var destIDs []string
	for _, d := range dests {
		destIDs = append(destIDs, d.ID)
	}
	for _, loc := range overlayLocales(t) {
		t.Run(string(loc), func(t *testing.T) {
			check(t, loc, "courses.json", courseIDs)
			check(t, loc, "destinations.json", destIDs)
		})
	}
}

type kpStep struct {
	Label    string   `json:"label"`
	Section  string   `json:"section"`
	Subgoals []string `json:"subgoals"`
}

func TestKPShardOverlaysMirrorTheEnglishOnes(t *testing.T) {
	ids := corpus(t)
	for _, loc := range overlayLocales(t) {
		t.Run(string(loc), func(t *testing.T) { checkKPShards(t, loc, ids) })
	}
}

func checkKPShards(t *testing.T, loc i18n.Locale, ids map[string]string) {
	t.Helper()
	dir := kpDir + "/" + string(loc)

	// Every concept in the school domains has a translated shard. A domain outside them may have
	// shards too, but only when it is complete, every one of its concepts translated, so that a
	// half-translated domain cannot ship. A shard for an id outside the corpus is always an error.
	school := map[string]bool{}
	for _, d := range schoolDomains {
		school[d] = true
	}
	total, translated := map[string]int{}, map[string]int{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s shards: %v", loc, err)
	}
	have := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			t.Errorf("unexpected entry %q in the %s shard directory", e.Name(), loc)
			continue
		}
		have[strings.TrimSuffix(e.Name(), ".json")] = true
	}
	var missing, extra, partial []string
	for id, dom := range ids {
		total[dom]++
		if have[id] {
			translated[dom]++
		} else if school[dom] {
			missing = append(missing, id)
		}
	}
	for id := range have {
		if _, ok := ids[id]; !ok {
			extra = append(extra, id)
		}
	}
	for dom, n := range translated {
		if !school[dom] && n < total[dom] {
			partial = append(partial, fmt.Sprintf("%s (%d of %d)", dom, n, total[dom]))
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)
	sort.Strings(partial)
	if len(missing) > 0 {
		t.Errorf("%d school-domain concepts have no %s shard, e.g. %v", len(missing), loc, missing[:min(5, len(missing))])
	}
	if len(extra) > 0 {
		t.Errorf("%s shards for ids that are not in the corpus: %v", loc, extra[:min(5, len(extra))])
	}
	if len(partial) > 0 {
		t.Errorf("%s shards for domains that are only partly translated, which must be finished or removed: %v", loc, partial)
	}
	if len(have) < 300 {
		t.Errorf("only %d %s shards found; the school domains have ~334", len(have), loc)
	}

	fromOpenStax := openStaxConcepts(t)
	for id := range have {
		var en, tr []kpStep
		readJSON(t, kpDir+"/"+id+".json", &en)
		readJSON(t, dir+"/"+id+".json", &tr)
		if len(tr) != len(en) {
			t.Errorf("%s: %d %s steps for %d English ones", id, len(tr), loc, len(en))
			continue
		}
		for i := range en {
			where := fmt.Sprintf("%s#%d", id, i)
			// A step's section is a heading of the lesson it points into; for an OpenStax lesson
			// that heading is OpenStax text, which is not translated (see openStaxConcepts).
			if fromOpenStax[id] && tr[i].Section != en[i].Section {
				t.Errorf("%s: section %q is OpenStax text and must stay %q", where, tr[i].Section, en[i].Section)
			}
			if strings.TrimSpace(tr[i].Label) == "" {
				t.Errorf("%s: empty label", where)
			} else {
				checkLetters(t, loc, where+" label", tr[i].Label)
			}
			if (en[i].Section == "") != (tr[i].Section == "") {
				t.Errorf("%s: section %q became %q; an empty section must stay empty and a set one must not be dropped", where, en[i].Section, tr[i].Section)
			}
			if len(tr[i].Subgoals) != len(en[i].Subgoals) {
				t.Errorf("%s: %d %s subgoals for %d English ones", where, len(tr[i].Subgoals), loc, len(en[i].Subgoals))
				continue
			}
			checkCarried(t, where+" label", en[i].Label, tr[i].Label)
			checkCarried(t, where+" section", en[i].Section, tr[i].Section)
			for j := range en[i].Subgoals {
				if strings.TrimSpace(tr[i].Subgoals[j]) == "" {
					t.Errorf("%s subgoal %d: empty", where, j)
				}
				checkCarried(t, where+" subgoal", en[i].Subgoals[j], tr[i].Subgoals[j])
			}
		}
	}
}

// openStaxRe matches the attribution line an OpenStax-derived lesson opens with.
var openStaxRe = regexp.MustCompile(`(?m)^> Content (sourced|adapted) from \[OpenStax`)

// openStaxConcepts are the concepts whose lesson text comes from an OpenStax book: those
// data/lessons/mapping.json attributes to OpenStax, and those whose teaching lesson carries the
// OpenStax attribution line.
//
// OpenStax text is not translated in this fork, by machine or AI. Every OpenStax title states
// that the book "may not be used in the training of large language models or otherwise be
// ingested into large language models or generative AI offerings without OpenStax's prior
// written permission" (recorded in THIRD_PARTY_NOTICES.md), and the translation here is made
// by AI. The lessons stay in English until OpenStax grants permission.
func openStaxConcepts(t *testing.T) map[string]bool {
	t.Helper()
	out := map[string]bool{}
	var mapping map[string]struct{ Source string }
	readJSON(t, dataDir+"/lessons/mapping.json", &mapping)
	for id, m := range mapping {
		if m.Source == "openstax" {
			out[id] = true
		}
	}
	for _, f := range openStaxLessonFiles(t) {
		out[strings.TrimSuffix(filepath.Base(f), ".md")] = true
	}
	if len(out) == 0 {
		t.Fatal("found no OpenStax lessons; the mapping or the attribution line has changed")
	}
	return out
}

// openStaxLessonFiles are the English lesson files that carry the OpenStax attribution line,
// relative to data/lessons.
func openStaxLessonFiles(t *testing.T) []string {
	t.Helper()
	var out []string
	root := dataDir + "/lessons"
	for _, corpus := range []string{"teaching", "algebrica", "authored"} {
		err := filepath.WalkDir(root+"/"+corpus, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() && i18n.Locale(d.Name()).Valid() {
				return filepath.SkipDir // a translation, not an English original
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if openStaxRe.Match(data) {
				rel, _ := filepath.Rel(root, path)
				out = append(out, filepath.ToSlash(rel))
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return out
}

// No translation of an OpenStax lesson exists in any language; see openStaxConcepts for why.
func TestOpenStaxLessonsAreNotTranslated(t *testing.T) {
	files := openStaxLessonFiles(t)
	if len(files) < 90 {
		t.Fatalf("found %d OpenStax lessons, expected about 99; the attribution line may have changed", len(files))
	}
	root := dataDir + "/lessons"
	for _, rel := range files {
		corpus, rest, _ := strings.Cut(rel, "/")
		for _, loc := range i18n.Overlays() {
			p := root + "/" + corpus + "/" + string(loc) + "/" + rest
			if _, err := os.Stat(p); err == nil {
				t.Errorf("%s translates an OpenStax lesson; OpenStax text stays in English", p)
			}
		}
	}
}

// Every overlay directory belongs to a supported locale; a directory the server does not know
// would be translated work that no learner can ever see.
func TestOverlayDirectoriesAreSupportedLocales(t *testing.T) {
	for _, root := range []string{dataDir + "/i18n", kpDir} {
		entries, err := os.ReadDir(root)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range entries {
			if e.IsDir() && !i18n.Locale(e.Name()).Valid() {
				t.Errorf("%s/%s is not a supported locale", root, e.Name())
			}
		}
	}
	if len(overlayLocales(t)) == 0 {
		t.Error("no overlay locale has started; the Russian overlays at least should ship")
	}
}
