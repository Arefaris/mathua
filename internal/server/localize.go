package server

import (
	"fmt"
	"net/http"

	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/i18n"
	"github.com/chuma-beep/mathua/internal/scheduler"
	"github.com/chuma-beep/mathua/internal/xp"
)

// The engine builds every response in English, the language of the corpus. The helpers below
// put the request's language on the human-readable names in what it returned, keyed by the ids
// that travel beside them. Doing it here, from the id, means the result never depends on which
// language a stored session or a cached value was first built in. An id the overlay does not
// know keeps the English text the engine already put there.

// localizePrereqs relabels prerequisite rows in place.
func (s *Server) localizePrereqs(r *http.Request, ps []engine.PrereqInfo) {
	dag, loc := s.eng.GetDAG(), localeOf(r)
	for i := range ps {
		ps[i].Label = dag.LabelByID(ps[i].ID, loc, ps[i].Label)
	}
}

// localizeEligibility relabels a topic's own title and its prerequisites' titles in place.
func (s *Server) localizeEligibility(r *http.Request, el *engine.ConceptEligibility) {
	dag, loc := s.eng.GetDAG(), localeOf(r)
	el.Title = dag.LabelByID(el.ConceptID, loc, el.Title)
	for i := range el.Prerequisites {
		p := &el.Prerequisites[i]
		p.Title = dag.LabelByID(p.ConceptID, loc, p.Title)
	}
}

// recommendationRU words a recommendation card in Russian, keyed by the English badge. Each
// badge belongs to exactly one wording in scheduler and engine recommend.go, so the badge
// identifies the whole card; a badge missing here keeps its English wording.
var recommendationRU = map[string]struct{ badge, detail, cta string }{
	"Recommended":   {"Рекомендуем", "Короткая адаптивная диагностика покажет, с чего начать", "Пройти диагностику →"},
	"Mastery check": {"Проверка освоения", fmt.Sprintf("С прошлой проверки вы набрали %d XP — пройдите её, чтобы закрепить изученное", xp.QuizGateXP), "Пройти проверку →"},
	"All learned":   {"Всё изучено", "Сейчас ничего не нужно повторять — на графе видно, что строится на ваших темах", "Открыть граф →"},
	"Paused":        {"На паузе", "Повторения отложены, пока вы не продолжите, — темы к повторению скрыты", "Изменить или продолжить →"},
	"Review":        {"Повторение", "Пора проверить то, что вы уже изучили", "Повторить →"},
	"Practice":      {"Практика", "Стоит попробовать ещё раз — эта тема пока не закрепилась", "Продолжить практику →"},
	"New":           {"Новая тема", "Можно начинать — всё, что для неё нужно, вы уже прошли", "Начать →"},
}

// recommendationTitlesRU names the cards that are not a concept and so carry no concept title;
// the client would otherwise show the bare id ("diagnostic").
var recommendationTitlesRU = map[string]string{
	"diagnostic":    "Диагностика",
	"mastery-check": "Проверка освоения",
	"graph":         "Граф тем",
	"paused":        "Пауза",
}

// localizeRecommendations relabels the concept titles of a recommendation in place and, for a
// Russian request, its badge, detail and call to action.
func (s *Server) localizeRecommendations(r *http.Request, resp *scheduler.RecommendationResponse) {
	dag, loc := s.eng.GetDAG(), localeOf(r)
	localize := func(rec *scheduler.Recommendation) {
		rec.ConceptTitle = dag.LabelByID(rec.ConceptID, loc, rec.ConceptTitle)
		if loc == i18n.RU {
			wordRecommendationRU(rec)
		}
	}
	if resp.Primary != nil {
		localize(resp.Primary)
	}
	for i := range resp.Alternatives {
		localize(&resp.Alternatives[i])
	}
}

// wordRecommendationRU swaps a card's English wording for the Russian one.
func wordRecommendationRU(rec *scheduler.Recommendation) {
	if w, ok := recommendationRU[rec.Badge]; ok {
		rec.Badge, rec.Detail, rec.CTA = w.badge, w.detail, w.cta
	}
	if rec.ConceptTitle == "" {
		rec.ConceptTitle = recommendationTitlesRU[rec.ID]
	}
}

// localizeCatalog swaps each course for a copy with its name and description in the request's
// language. The courses themselves are shared by every request and are never edited.
func (s *Server) localizeCatalog(r *http.Request, catalog []engine.CourseStatus) {
	p := s.eng.GetPlanner()
	if p == nil {
		return
	}
	loc := localeOf(r)
	for i := range catalog {
		catalog[i].Course = p.CourseIn(catalog[i].Course, loc)
	}
}
