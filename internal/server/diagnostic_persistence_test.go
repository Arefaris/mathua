package server

import (
	"fmt"
	"testing"
)

// A diagnostic used to live only in the server's map, so a deploy answered the
// next question with 404 "diagnostic session not found", the screenshot a
// learner sent in was that message. The test simulates the restart precisely:
// the map entry disappears, the row does not.
func TestDiagnosticSession_SurvivesMemoryLoss(t *testing.T) {
	s, mux := twoConceptServer(t)

	code, start := postJSON(t, mux, "/api/goal/diagnostic", `{"name":"restart","concept_ids":["a","b"]}`)
	if code != 200 {
		t.Fatalf("start: %d %v", code, start)
	}
	dsid, _ := start["session_id"].(string)
	cid, _ := start["concept_id"].(string)
	if dsid == "" || cid == "" {
		t.Fatalf("start missing ids: %v", start)
	}

	code, ans := postJSON(t, mux, "/api/goal/diagnostic/answer",
		fmt.Sprintf(`{"session_id":"%s","concept_id":"%s","answer":"4","elapsed":1.0}`, dsid, cid))
	if code != 200 {
		t.Fatalf("first answer: %d %v", code, ans)
	}
	nextCID, _ := ans["concept_id"].(string)
	if nextCID == "" {
		t.Fatalf("diagnostic ended after one answer on a two-concept path: %v", ans)
	}

	// The restart: memory forgets, storage does not.
	s.mu.Lock()
	delete(s.diagSessions, dsid)
	delete(s.diagCreated, dsid)
	s.mu.Unlock()
	if _, found, err := s.repo.LoadDiagnosticSession(dsid); err != nil || !found {
		t.Fatalf("session was not persisted before the first answer: found=%v err=%v", found, err)
	}

	code, ans = postJSON(t, mux, "/api/goal/diagnostic/answer",
		fmt.Sprintf(`{"session_id":"%s","concept_id":"%s","answer":"4","elapsed":1.0}`, dsid, nextCID))
	if code != 200 {
		t.Fatalf("answer after simulated restart: %d %v", code, ans)
	}
}

// The finished run must outlive completion by exactly one request: the plan
// endpoint reads it and then deletes it. Deleting at completion 404'd the plan
// the learner was about to ask for, which the journey tests caught.
func TestDiagnosticSession_PlanConsumesTheFinishedRun(t *testing.T) {
	s, mux := twoConceptServer(t)
	_, dsid := journeyDiagnostic(t, mux)

	if _, found, _ := s.repo.LoadDiagnosticSession(dsid); !found {
		t.Fatal("completed run was not persisted, so a restart would lose the plan")
	}
	code, plan := postJSON(t, mux, "/api/goal/plan", `{"session_id":"`+dsid+`"}`)
	if code != 200 {
		t.Fatalf("plan: %d %v", code, plan)
	}
	if _, found, _ := s.repo.LoadDiagnosticSession(dsid); found {
		t.Error("plan left the consumed run in storage")
	}
}

// A row that no longer decodes is a dead end, not a permanent 500: drop it so
// the next request 404s once instead of failing every time.
func TestDiagnosticSession_UnreadableRunIsDropped(t *testing.T) {
	s, mux := twoConceptServer(t)
	if err := s.repo.SaveDiagnosticSession("broken", "", []byte("{")); err != nil {
		t.Fatalf("seed broken row: %v", err)
	}

	code, _ := postJSON(t, mux, "/api/goal/diagnostic/answer",
		`{"session_id":"broken","concept_id":"a","answer":"4","elapsed":1.0}`)
	if code != 404 {
		t.Fatalf("unreadable session = %d, want 404", code)
	}
	if _, found, _ := s.repo.LoadDiagnosticSession("broken"); found {
		t.Error("unreadable run was left in storage")
	}
}

// A wiped learner must not be able to resume the run from the row the wipe
// forgot: ResetProgress and account deletion both go through this helper.
func TestDiagnosticSession_DropStudentDeletesTheRow(t *testing.T) {
	s, mux := twoConceptServer(t)
	code, start := postJSON(t, mux, "/api/goal/diagnostic", `{"name":"wipe","concept_ids":["a","b"]}`)
	if code != 200 {
		t.Fatalf("start: %d %v", code, start)
	}
	sid, _ := start["student_id"].(string)
	dsid, _ := start["session_id"].(string)

	s.dropStudentLiveSessions(sid)

	if _, found, _ := s.repo.LoadDiagnosticSession(dsid); found {
		t.Error("student's diagnostic row survived the wipe")
	}
}
