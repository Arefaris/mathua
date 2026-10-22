package server

import (
	"net/http"

	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/scheduler"
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

// localizeRecommendations relabels the concept titles of a recommendation in place. The
// badge, detail and call-to-action strings are the scheduler's own wording and are not
// translated here; kind and reason identify them.
func (s *Server) localizeRecommendations(r *http.Request, resp *scheduler.RecommendationResponse) {
	dag, loc := s.eng.GetDAG(), localeOf(r)
	if resp.Primary != nil {
		resp.Primary.ConceptTitle = dag.LabelByID(resp.Primary.ConceptID, loc, resp.Primary.ConceptTitle)
	}
	for i := range resp.Alternatives {
		a := &resp.Alternatives[i]
		a.ConceptTitle = dag.LabelByID(a.ConceptID, loc, a.ConceptTitle)
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
