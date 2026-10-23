package server

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestGoalDiagnosticAnswer_DuplicatePostReplaysTheVerdict is the duplicate
// submission case: the second POST of the same concept must return exactly the
// verdict the first one produced, not a second grading and not an error.
func TestGoalDiagnosticAnswer_DuplicatePostReplaysTheVerdict(t *testing.T) {
	_, mux := twoConceptServer(t)

	code, start := postJSON(t, mux, "/api/goal/diagnostic", `{"name":"tester","concept_ids":["a","b"]}`)
	if code != 200 {
		t.Fatalf("start: %d %v", code, start)
	}
	sid, _ := start["session_id"].(string)
	cid, _ := start["concept_id"].(string)
	body := fmt.Sprintf(`{"session_id":%q,"concept_id":%q,"answer":"4","elapsed":5.0}`, sid, cid)

	firstCode, first := postJSON(t, mux, "/api/goal/diagnostic/answer", body)
	if firstCode != 200 {
		t.Fatalf("first answer: %d %v", firstCode, first)
	}
	secondCode, second := postJSON(t, mux, "/api/goal/diagnostic/answer", body)
	if secondCode != 200 {
		t.Fatalf("replayed answer: %d %v", secondCode, second)
	}

	firstJSON, _ := json.Marshal(first)
	secondJSON, _ := json.Marshal(second)
	if string(firstJSON) != string(secondJSON) {
		t.Errorf("replay differs from the first verdict:\nfirst  %s\nsecond %s", firstJSON, secondJSON)
	}
}
