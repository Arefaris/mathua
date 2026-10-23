package server

import (
	"os"
	"regexp"
	"testing"

	"github.com/chuma-beep/mathua/internal/i18n"
	"github.com/chuma-beep/mathua/internal/scheduler"
)

// Every badge the scheduler and engine can put on a card must have a wording in every
// translated language, or a learner sees one English card among translated ones. The badges
// are read from source so a new card added upstream fails here instead of shipping
// half-translated.
func TestEveryRecommendationBadgeHasEveryWording(t *testing.T) {
	badge := regexp.MustCompile(`Badge\s*(?::|=)\s*"([^"]+)"`)
	seen := 0
	for _, f := range []string{"../scheduler/recommend.go", "../engine/recommend.go"} {
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range badge.FindAllStringSubmatch(string(src), -1) {
			seen++
			for _, loc := range i18n.Overlays() {
				if _, ok := recommendationWording[loc][m[1]]; !ok {
					t.Errorf("%s: badge %q has no %s wording in recommendationWording", f, m[1], loc)
				}
			}
		}
	}
	if seen == 0 {
		t.Fatal("found no badges; the pattern no longer matches the source")
	}
	for _, loc := range i18n.Overlays() {
		for _, id := range []string{"diagnostic", "mastery-check", "graph", "paused"} {
			if recommendationTitles[loc][id] == "" {
				t.Errorf("card %q has no %s title in recommendationTitles", id, loc)
			}
		}
	}
}

func TestWordRecommendation(t *testing.T) {
	rec := scheduler.Recommendation{ID: "diagnostic", Badge: "Recommended", Detail: "x", CTA: "y"}
	wordRecommendation(&rec, i18n.RU)
	if rec.Badge != "Рекомендуем" || rec.CTA != "Пройти диагностику →" || rec.ConceptTitle != "Диагностика" {
		t.Fatalf("diagnostic card not worded in Russian: %+v", rec)
	}

	uk := scheduler.Recommendation{ID: "diagnostic", Badge: "Recommended"}
	wordRecommendation(&uk, i18n.UK)
	if uk.Badge != "Рекомендуємо" || uk.ConceptTitle != "Діагностика" {
		t.Fatalf("diagnostic card not worded in Ukrainian: %+v", uk)
	}

	concept := scheduler.Recommendation{ID: "frac.add", ConceptTitle: "Сложение дробей", Badge: "New"}
	wordRecommendation(&concept, i18n.RU)
	if concept.ConceptTitle != "Сложение дробей" || concept.Badge != "Новая тема" {
		t.Fatalf("concept card: %+v", concept)
	}

	unknown := scheduler.Recommendation{ID: "x", Badge: "Something new", Detail: "d"}
	wordRecommendation(&unknown, i18n.RU)
	if unknown.Badge != "Something new" || unknown.Detail != "d" {
		t.Fatalf("unknown badge must keep its English wording: %+v", unknown)
	}

	en := scheduler.Recommendation{ID: "diagnostic", Badge: "Recommended"}
	wordRecommendation(&en, i18n.EN)
	if en.Badge != "Recommended" || en.ConceptTitle != "" {
		t.Fatalf("an English request must keep the engine's wording: %+v", en)
	}
}
