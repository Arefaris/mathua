package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/chuma-beep/mathua/internal/ai"
)

// Bounds for the AI endpoints. They are small on purpose: an assistant message
// is a question, not an essay, and the goal is a sentence. A caller that
// exceeds them is rejected rather than truncated, so a visitor never gets a
// reply to something they did not ask.
const (
	maxAIBody           = 16 << 10
	maxAssistantMsgLen  = 2000
	maxGoalLen          = 500
	maxGoalDestinations = 3
)

// aiAssistantRequest is the assistant endpoint's body.
//
// It carries one field. The assistant answers product questions and helps a
// visitor choose a starting point, so it has no use for a concept id or a
// mastery line, and the request cannot carry one.
//
// There is deliberately no expected_answer field. The graded answer lives in
// the server's own study-session record and is never a request field, so a
// client cannot send one and a future refactor cannot accidentally forward one.
type aiAssistantRequest struct {
	Message string `json:"message"`
}

// aiGoalScopeRequest is the goal-scope endpoint's body.
type aiGoalScopeRequest struct {
	Goal string `json:"goal"`
}

// POST /api/ai/assistant — an optional, best-effort product assistant. It
// answers "what is Mathua" and "where do I start", and it is not a tutor: it
// holds no concept context and never guides a graded question. The deterministic
// app does not depend on it: a disabled registry, an exhausted quota or a
// provider outage all surface as an error the client shows as "unavailable".
func (s *Server) handleAIAssistant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.aiChain == nil || !s.aiChain.Enabled() {
		writeError(w, "ai unavailable", http.StatusServiceUnavailable)
		return
	}

	var req aiAssistantRequest
	if err := decodeAIJSON(w, r, &req); err != nil {
		return
	}
	msg := strings.TrimSpace(req.Message)
	if msg == "" || len([]rune(msg)) > maxAssistantMsgLen {
		writeError(w, "message must be 1 to 2000 characters", http.StatusBadRequest)
		return
	}

	if !s.reserveAIUsage(w, r) {
		return
	}

	// The catalogue is the server's own destination list, so the assistant can
	// name a real starting point and cannot invent one. With no catalogue it
	// still answers, and asks where the visitor wants to go.
	actx := ai.AssistantContext{
		Message:      msg,
		Destinations: assistantDestinations(s.goalCandidates()),
	}
	res, err := s.aiChain.Complete(r.Context(), ai.DataPseudonymous, ai.Request{
		Messages:  actx.Messages(),
		MaxTokens: s.aiAssistantMaxTokens,
	})
	if err != nil {
		s.writeAIError(w, "assistant", err)
		return
	}
	writeJSON(w, map[string]string{"reply": res.Response.Text})
}

// Goal-scope reasons. They exist so a failure is not silently indistinguishable
// from a model that simply chose nothing: the client shows different copy, and
// the value is what a log or counter can be built on.
const (
	goalReasonDisabled     = "disabled"      // AI is off or has no eligible provider
	goalReasonNoCandidates = "no_candidates" // the server has no destinations to choose from
	goalReasonFailed       = "failed"        // every provider in the chain failed
	goalReasonNoMatch      = "no_match"      // the model answered, but chose no valid destination
)

// aiGoalScopeResponse is the goal-scope reply.
type aiGoalScopeResponse struct {
	Destinations []string `json:"destinations"`
	UsedAI       bool     `json:"used_ai"`
	// Reason is set only when UsedAI is false.
	Reason string `json:"reason,omitempty"`
}

// POST /api/ai/goal-scope — maps a learner's free-text goal to destination ids
// from the repository's own list. It is an enhancement, not the mechanism: when
// AI is unavailable it returns used_ai=false with a reason, and the client keeps
// its deterministic fuzzy match.
func (s *Server) handleAIGoalScope(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.aiChain == nil || !s.aiChain.Enabled() {
		writeJSON(w, aiGoalScopeResponse{Destinations: []string{}, Reason: goalReasonDisabled})
		return
	}

	var req aiGoalScopeRequest
	if err := decodeAIJSON(w, r, &req); err != nil {
		return
	}
	goal := strings.TrimSpace(req.Goal)
	if goal == "" || len([]rune(goal)) > maxGoalLen {
		writeError(w, "goal must be 1 to 500 characters", http.StatusBadRequest)
		return
	}

	// The candidate list comes from the server's own data, never from the
	// request, so a client cannot make the model pick an id that does not exist.
	candidates := s.goalCandidates()
	if len(candidates) == 0 {
		writeJSON(w, aiGoalScopeResponse{Destinations: []string{}, Reason: goalReasonNoCandidates})
		return
	}

	if !s.reserveAIUsage(w, r) {
		return
	}

	res, err := s.aiChain.Complete(r.Context(), ai.DataPseudonymous, ai.Request{
		Messages:  goalScopeMessages(goal, candidates),
		MaxTokens: s.aiGoalMaxTokens,
		JSON:      true,
	})
	if err != nil {
		// A provider failure is not worth an error page here: the client's
		// deterministic match is the fallback and is already on screen. It is
		// still reported, because "the model failed" and "the model chose
		// nothing" must not look the same in a response or in a log.
		log.Printf("ai goal-scope: failed: %v", err)
		writeJSON(w, aiGoalScopeResponse{Destinations: []string{}, Reason: goalReasonFailed})
		return
	}

	picked := validateDestinationIDs(parseDestinationIDs(res.Response.Text), candidates)
	if len(picked) == 0 {
		writeJSON(w, aiGoalScopeResponse{Destinations: []string{}, Reason: goalReasonNoMatch})
		return
	}
	writeJSON(w, aiGoalScopeResponse{Destinations: picked, UsedAI: true})
}

// goalCandidate is one destination offered to the model.
type goalCandidate struct {
	ID          string
	Name        string
	Description string
}

// goalCandidates reads the destination list from the engine, which loaded it
// from data/destinations.json at boot.
func (s *Server) goalCandidates() []goalCandidate {
	if s.eng == nil {
		return nil
	}
	dests, err := s.eng.Destinations("")
	if err != nil {
		log.Printf("ai goal-scope: load destinations: %v", err)
		return nil
	}
	out := make([]goalCandidate, 0, len(dests))
	for _, d := range dests {
		if d.Destination == nil {
			continue
		}
		out = append(out, goalCandidate{ID: d.ID, Name: d.Name, Description: d.Description})
	}
	return out
}

// goalScopeMessages builds the goal-scope prompt. The candidate list is authored
// course data, not learner data, so this request carries only the learner's own
// words plus public text.
func goalScopeMessages(goal string, candidates []goalCandidate) []ai.Message {
	var b strings.Builder
	b.WriteString("Destinations:\n")
	for _, c := range candidates {
		b.WriteString("- ")
		b.WriteString(c.ID)
		b.WriteString(": ")
		b.WriteString(c.Name)
		if c.Description != "" {
			b.WriteString(" — ")
			b.WriteString(c.Description)
		}
		b.WriteByte('\n')
	}
	return []ai.Message{
		{Role: ai.RoleSystem, Content: "You map a learner's goal to the destinations in the list that meet it.\n" +
			"Reply with JSON only, in the form {\"destinations\":[\"id\"]}.\n" +
			"Use only ids from the list. Choose at most 3. If none fit, reply with an empty list."},
		{Role: ai.RoleUser, Content: "Goal: " + goal + "\n\n" + b.String()},
	}
}

// parseDestinationIDs reads the model's JSON reply. A reply that is not JSON, or
// is missing the field, yields no ids rather than an error: the caller falls
// back to the deterministic path.
func parseDestinationIDs(text string) []string {
	trimmed := strings.TrimSpace(text)
	// Tolerate a fenced block, which some models add even under JSON mode.
	if i := strings.Index(trimmed, "{"); i > 0 {
		trimmed = trimmed[i:]
	}
	if j := strings.LastIndex(trimmed, "}"); j >= 0 && j < len(trimmed)-1 {
		trimmed = trimmed[:j+1]
	}
	var parsed struct {
		Destinations []string `json:"destinations"`
	}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}
	return parsed.Destinations
}

// validateDestinationIDs keeps only ids that are in the candidate set, in the
// model's order, capped so a chatty reply cannot produce a huge scope.
func validateDestinationIDs(ids []string, candidates []goalCandidate) []string {
	valid := make(map[string]bool, len(candidates))
	for _, c := range candidates {
		valid[c.ID] = true
	}
	out := make([]string, 0, maxGoalDestinations)
	seen := make(map[string]bool, maxGoalDestinations)
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if !valid[id] || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
		if len(out) == maxGoalDestinations {
			break
		}
	}
	return out
}

// assistantDestinations converts the server's destination list into the public
// catalogue the assistant may name. It is a projection, not a filter: the
// assistant needs ids, names and descriptions, and nothing about a learner's
// progress, which is why this type cannot carry any.
func assistantDestinations(candidates []goalCandidate) []ai.Destination {
	if len(candidates) == 0 {
		return nil
	}
	out := make([]ai.Destination, 0, len(candidates))
	for _, c := range candidates {
		out = append(out, ai.Destination{ID: c.ID, Name: c.Name, Description: c.Description})
	}
	return out
}

// reserveAIUsage consumes one request from the caller's daily budget and writes
// the error response itself when the budget is gone or the store failed. It
// returns false when the request must stop.
func (s *Server) reserveAIUsage(w http.ResponseWriter, r *http.Request) bool {
	if s.aiQuota == nil {
		return true
	}
	ok, err := s.aiQuota.Reserve(r.Context(), s.aiQuotaKey(r))
	if err != nil {
		log.Printf("ai quota: %v", err)
		writeError(w, "ai unavailable", http.StatusServiceUnavailable)
		return false
	}
	if !ok {
		// A refusal is the one AI outcome that leaves no trace in a provider
		// log, because no provider was called.
		if s.aiChain != nil {
			s.aiChain.Stats().Denied()
		}
		writeError(w, "you have reached today's AI limit", http.StatusTooManyRequests)
		return false
	}
	return true
}

// aiQuotaKey is the identity a daily AI budget is spent against.
//
// An authenticated caller is keyed by student id. Everyone else is keyed by IP,
// not by the guest id they sent: guest ids are generated in the browser and are
// therefore freely mintable, so an id-based guest cap is bypassed by clearing
// localStorage.
func (s *Server) aiQuotaKey(r *http.Request) string {
	if sid, _ := r.Context().Value(authStudentKey{}).(string); sid != "" {
		return "u:" + sid
	}
	return "ip:" + clientIP(r, s.ipHeader)
}

// writeAIError maps a chain failure to a response. ErrNoProvider is not the
// client's fault and is reported as unavailable, never as a data problem.
func (s *Server) writeAIError(w http.ResponseWriter, label string, err error) {
	if errors.Is(err, ai.ErrNoProvider) {
		writeError(w, "ai unavailable", http.StatusServiceUnavailable)
		return
	}
	log.Printf("ai %s: %v", label, err)
	writeError(w, "ai request failed", http.StatusBadGateway)
}

// decodeAIJSON decodes a bounded JSON body, writing a 400 on failure.
func decodeAIJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxAIBody)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return err
	}
	return nil
}
