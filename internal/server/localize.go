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

type cardWording struct{ badge, detail, cta string }

// recommendationWording words a recommendation card in each translated language, keyed by the
// English badge. Each badge belongs to exactly one wording in scheduler and engine
// recommend.go, so the badge identifies the whole card; a badge missing here keeps its English
// wording.
var recommendationWording = map[i18n.Locale]map[string]cardWording{
	i18n.RU: {
		"Recommended":   {"Рекомендуем", "Короткая адаптивная диагностика покажет, с чего начать", "Пройти диагностику →"},
		"Mastery check": {"Проверка освоения", fmt.Sprintf("С прошлой проверки вы набрали %d XP — пройдите её, чтобы закрепить изученное", xp.QuizGateXP), "Пройти проверку →"},
		"All learned":   {"Всё изучено", "Сейчас повторять нечего — граф показывает, что опирается на уже изученное", "Открыть граф →"},
		"Paused":        {"На паузе", "Повторения отложены, пока вы не возобновите обучение, — темы к повторению скрыты", "Изменить или возобновить →"},
		"Review":        {"Повторение", "Пора проверить то, что вы уже изучили", "Повторить →"},
		"Remediation":   {"Восполнить пробел", "Закрепите эту необходимую тему, чтобы продвинуться там, где вы застряли", "Повторить необходимую тему →"},
		"Practice":      {"Практика", "Стоит попробовать ещё раз — эта тема пока не закрепилась", "Продолжить практику →"},
		"New":           {"Новая тема", "Можно начинать — всё, что для неё нужно, вы уже прошли", "Начать →"},
	},
	i18n.UK: {
		"Recommended":   {"Рекомендуємо", "Коротка адаптивна діагностика покаже, з чого почати", "Пройти діагностику →"},
		"Mastery check": {"Перевірка засвоєння", fmt.Sprintf("Від останньої перевірки ви набрали %d XP — пройдіть її, щоб закріпити вивчене", xp.QuizGateXP), "Пройти перевірку →"},
		"All learned":   {"Усе вивчено", "Зараз нічого повторювати — граф показує, що спирається на вже вивчене", "Відкрити граф →"},
		"Paused":        {"На паузі", "Повторення відкладено, доки ви не відновите навчання, — теми для повторення приховано", "Змінити або відновити →"},
		"Review":        {"Повторення", "Час перевірити те, що ви вже вивчили", "Повторити →"},
		"Remediation":   {"Заповнити прогалину", "Закріпіть цю необхідну тему, щоб просунутися там, де ви застрягли", "Повторити необхідну тему →"},
		"Practice":      {"Практика", "Варто спробувати ще раз — ця тема поки не закріпилася", "Продовжити практику →"},
		"New":           {"Нова тема", "Можна починати — усе, що для неї потрібно, ви вже пройшли", "Почати →"},
	},
}

// recommendationTitles names the cards that are not a concept and so carry no concept title;
// the client would otherwise show the bare id ("diagnostic").
var recommendationTitles = map[i18n.Locale]map[string]string{
	i18n.RU: {
		"diagnostic":    "Диагностика",
		"mastery-check": "Проверка освоения",
		"graph":         "Граф тем",
		"paused":        "Пауза",
	},
	i18n.UK: {
		"diagnostic":    "Діагностика",
		"mastery-check": "Перевірка засвоєння",
		"graph":         "Граф тем",
		"paused":        "Пауза",
	},
}

// localizeRecommendations relabels the concept titles of a recommendation in place and, for a
// translated language, its badge, detail and call to action.
func (s *Server) localizeRecommendations(r *http.Request, resp *scheduler.RecommendationResponse) {
	dag, loc := s.eng.GetDAG(), localeOf(r)
	localize := func(rec *scheduler.Recommendation) {
		rec.ConceptTitle = dag.LabelByID(rec.ConceptID, loc, rec.ConceptTitle)
		wordRecommendation(rec, loc)
	}
	if resp.Primary != nil {
		localize(resp.Primary)
	}
	for i := range resp.Alternatives {
		localize(&resp.Alternatives[i])
	}
}

// wordRecommendation swaps a card's English wording for the one in loc, if there is one.
func wordRecommendation(rec *scheduler.Recommendation, loc i18n.Locale) {
	if w, ok := recommendationWording[loc][rec.Badge]; ok {
		rec.Badge, rec.Detail, rec.CTA = w.badge, w.detail, w.cta
	}
	if rec.ConceptTitle == "" {
		rec.ConceptTitle = recommendationTitles[loc][rec.ID]
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
