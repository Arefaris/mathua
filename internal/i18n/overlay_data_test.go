package i18n_test

// These tests read the shipped Russian overlays in data/i18n and data/lessons/kp/ru and
// hold them to the English originals they translate. They are the check that a translation
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
	ruDir   = dataDir + "/i18n/ru"
	kpDir   = dataDir + "/lessons/kp"
)

// schoolDomains are the domains whose lesson steps are translated in this round.
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
func checkCarried(t *testing.T, where, en, ru string) {
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
		if !strings.Contains(ru, tok) {
			t.Errorf("%s: %q lost %q\n  en: %s\n  ru: %s", where, ru, tok, en, ru)
		}
	}
	for _, d := range digitsRe.FindAllString(en, -1) {
		if !strings.Contains(ru, d) {
			t.Errorf("%s: the number %s is missing\n  en: %s\n  ru: %s", where, d, en, ru)
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

func TestRussianConceptLabels(t *testing.T) {
	ids := corpus(t)
	var ru map[string]string
	readJSON(t, ruDir+"/concepts.json", &ru)

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

	if len(ru) != len(ids) {
		t.Errorf("concepts.json has %d labels for %d concepts", len(ru), len(ids))
	}
	for id := range ru {
		if _, ok := ids[id]; !ok {
			t.Errorf("concepts.json names %q, which is not a concept in the corpus", id)
		}
	}
	for id := range ids {
		label, ok := ru[id]
		switch {
		case !ok:
			t.Errorf("no Russian label for %s", id)
		case strings.TrimSpace(label) == "":
			t.Errorf("empty Russian label for %s", id)
		case !hasCyrillic(label):
			t.Errorf("%s: label %q has no Russian text", id, label)
		case label != strings.TrimSpace(label):
			t.Errorf("%s: label %q has stray whitespace", id, label)
		default:
			checkCarried(t, id, en[id], label)
		}
	}
}

func TestRussianCourseAndDestinationText(t *testing.T) {
	check := func(file string, ids []string) {
		var ru map[string]planning.Text
		readJSON(t, ruDir+"/"+file, &ru)
		if len(ru) != len(ids) {
			t.Errorf("%s has %d entries for %d ids", file, len(ru), len(ids))
		}
		for _, id := range ids {
			e, ok := ru[id]
			if !ok {
				t.Errorf("%s: no entry for %q", file, id)
				continue
			}
			if strings.TrimSpace(e.Name) == "" || strings.TrimSpace(e.Description) == "" {
				t.Errorf("%s: %q has an empty name or description", file, id)
			}
			if !hasCyrillic(e.Name) || !hasCyrillic(e.Description) {
				t.Errorf("%s: %q is not Russian: %+v", file, id, e)
			}
		}
		for id := range ru {
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
	check("courses.json", courseIDs)

	var dests []struct{ ID string }
	readJSON(t, dataDir+"/destinations.json", &dests)
	var destIDs []string
	for _, d := range dests {
		destIDs = append(destIDs, d.ID)
	}
	check("destinations.json", destIDs)
}

type kpStep struct {
	Label    string   `json:"label"`
	Section  string   `json:"section"`
	Subgoals []string `json:"subgoals"`
}

func TestRussianKPShardsMirrorTheEnglishOnes(t *testing.T) {
	ids := corpus(t)

	// Every concept in the school domains has a translated shard, and nothing else does.
	want := map[string]bool{}
	for id, dom := range ids {
		for _, d := range schoolDomains {
			if dom == d {
				want[id] = true
			}
		}
	}
	entries, err := os.ReadDir(kpDir + "/ru")
	if err != nil {
		t.Fatalf("read ru shards: %v", err)
	}
	have := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			t.Errorf("unexpected entry %q in the ru shard directory", e.Name())
			continue
		}
		have[strings.TrimSuffix(e.Name(), ".json")] = true
	}
	var missing, extra []string
	for id := range want {
		if !have[id] {
			missing = append(missing, id)
		}
	}
	for id := range have {
		if !want[id] {
			extra = append(extra, id)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)
	if len(missing) > 0 {
		t.Errorf("%d school-domain concepts have no Russian shard, e.g. %v", len(missing), missing[:min(5, len(missing))])
	}
	if len(extra) > 0 {
		t.Errorf("Russian shards for ids outside the school domains or the corpus: %v", extra[:min(5, len(extra))])
	}
	if len(have) < 300 {
		t.Errorf("only %d Russian shards found; the school domains have ~334", len(have))
	}

	for id := range have {
		var en, ru []kpStep
		readJSON(t, kpDir+"/"+id+".json", &en)
		readJSON(t, kpDir+"/ru/"+id+".json", &ru)
		if len(ru) != len(en) {
			t.Errorf("%s: %d Russian steps for %d English ones", id, len(ru), len(en))
			continue
		}
		for i := range en {
			where := fmt.Sprintf("%s#%d", id, i)
			if strings.TrimSpace(ru[i].Label) == "" {
				t.Errorf("%s: empty label", where)
			}
			if (en[i].Section == "") != (ru[i].Section == "") {
				t.Errorf("%s: section %q became %q; an empty section must stay empty and a set one must not be dropped", where, en[i].Section, ru[i].Section)
			}
			if len(ru[i].Subgoals) != len(en[i].Subgoals) {
				t.Errorf("%s: %d Russian subgoals for %d English ones", where, len(ru[i].Subgoals), len(en[i].Subgoals))
				continue
			}
			checkCarried(t, where+" label", en[i].Label, ru[i].Label)
			checkCarried(t, where+" section", en[i].Section, ru[i].Section)
			for j := range en[i].Subgoals {
				if strings.TrimSpace(ru[i].Subgoals[j]) == "" {
					t.Errorf("%s subgoal %d: empty", where, j)
				}
				checkCarried(t, where+" subgoal", en[i].Subgoals[j], ru[i].Subgoals[j])
			}
		}
	}
}

func TestOverlayLocalesAreTheSupportedOnes(t *testing.T) {
	for _, loc := range i18n.Overlays() {
		if _, err := os.Stat(dataDir + "/i18n/" + string(loc)); err != nil {
			t.Errorf("supported overlay locale %q has no directory: %v", loc, err)
		}
	}
}
