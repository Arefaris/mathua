package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/storage"
)

func seedReport(t *testing.T, s *Server, conceptID, status string) int64 {
	t.Helper()
	id, err := s.repo.CreateReport(storage.QuestionReport{
		ConceptID: conceptID,
		Kind:      "question",
		Question:  "2+2=?",
		Expected:  "4",
		Reason:    "wrong_answer",
		Status:    status,
	})
	if err != nil {
		t.Fatalf("seed report: %v", err)
	}
	return id
}

func bulkPost(t *testing.T, s *Server, body map[string]interface{}) (int, map[string]interface{}) {
	t.Helper()
	raw, _ := json.Marshal(body)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/admin/reports/bulk", bytes.NewReader(raw))
	// The route sits behind requirePermission in production, which is what puts
	// the actor in the context. Calling the handler directly means supplying it
	// here, so the stamp is asserted against a real identity.
	req = req.WithContext(context.WithValue(req.Context(), authStudentKey{}, "admin-test"))
	s.handleAdminReportsBulk(rec, req)
	var res map[string]interface{}
	if rec.Body.Len() > 0 {
		_ = json.Unmarshal(rec.Body.Bytes(), &res)
	}
	return rec.Code, res
}

// Bulk close is about a bug, not a list: a fix ships, and every still-unsettled
// report about the concept it touched is stale. The test pins the two halves of
// that promise, the whole set closes together, and each report keeps its own
// decision record rather than inheriting one.
func TestAdminReportsBulk_ClosesAConceptsUnsettledReports(t *testing.T) {
	s := testServer(t)
	a1 := seedReport(t, s, "concept.x", "open")
	a2 := seedReport(t, s, "concept.x", "reviewing")
	settled := seedReport(t, s, "concept.x", "resolved")
	other := seedReport(t, s, "concept.y", "open")

	code, res := bulkPost(t, s, map[string]interface{}{
		"concept_id": "concept.x",
		"status":     "resolved",
		"resolution": "Fixed: the answer format now declares its own input.",
	})
	if code != 200 {
		t.Fatalf("bulk close: %d %v", code, res)
	}
	if res["updated"] != float64(2) {
		t.Errorf("updated = %v, want 2", res["updated"])
	}
	if failed, _ := res["failed"].([]interface{}); len(failed) != 0 {
		t.Errorf("failed = %v, want none", failed)
	}

	for _, id := range []int64{a1, a2} {
		r, err := s.repo.GetReport(id)
		if err != nil || r == nil {
			t.Fatalf("read report %d: %v", id, err)
		}
		if r.Status != "resolved" {
			t.Errorf("report %d status = %q, want resolved", id, r.Status)
		}
		if r.Resolution == "" || r.ResolvedBy == "" || r.ResolvedAt == "" {
			t.Errorf("report %d has no decision record: %+v", id, r)
		}
		// Each report owns its audit row; a single row naming the set would not
		// let a later reader see when this particular complaint was closed.
		history, err := s.repo.ListAdminAuditFor("report", strconv.FormatInt(id, 10), 10)
		if err != nil {
			t.Fatalf("audit for %d: %v", id, err)
		}
		if len(history) != 1 {
			t.Errorf("report %d audit rows = %d, want 1", id, len(history))
		}
	}

	// An already-settled report is not revisited: reopening and re-stamping it
	// would erase who actually decided it.
	unchanged, _ := s.repo.GetReport(settled)
	if unchanged == nil || unchanged.Status != "resolved" {
		t.Errorf("settled report changed: %+v", unchanged)
	}
	// A different concept is untouched.
	otherReport, _ := s.repo.GetReport(other)
	if otherReport == nil || otherReport.Status != "open" {
		t.Errorf("other concept's report changed: %+v", otherReport)
	}
}

func TestAdminReportsBulk_DismissesByIds(t *testing.T) {
	s := testServer(t)
	id := seedReport(t, s, "concept.x", "open")

	code, res := bulkPost(t, s, map[string]interface{}{
		"ids":        []int64{id},
		"status":     "dismissed",
		"resolution": "Not a bug: the decimal answer was accepted by the grader.",
	})
	if code != 200 || res["updated"] != float64(1) {
		t.Fatalf("dismiss by ids: %d %v", code, res)
	}
	r, _ := s.repo.GetReport(id)
	if r == nil || r.Status != "dismissed" {
		t.Errorf("report status = %+v, want dismissed", r)
	}
}

func TestAdminReportsBulk_RequiresAReasonToSettle(t *testing.T) {
	s := testServer(t)
	seedReport(t, s, "concept.x", "open")

	code, _ := bulkPost(t, s, map[string]interface{}{
		"concept_id": "concept.x",
		"status":     "resolved",
	})
	if code != 400 {
		t.Errorf("resolve without a resolution = %d, want 400", code)
	}
}

func TestAdminReportsBulk_RefusesAnEmptySelection(t *testing.T) {
	s := testServer(t)
	code, _ := bulkPost(t, s, map[string]interface{}{
		"status":     "reviewing",
		"resolution": "",
	})
	if code != 400 {
		t.Errorf("empty selection = %d, want 400", code)
	}
}

func TestAdminReportsBulk_CountsAnUnchangedReportAsSkipped(t *testing.T) {
	s := testServer(t)
	id := seedReport(t, s, "concept.x", "reviewing")

	code, res := bulkPost(t, s, map[string]interface{}{
		"ids":        []int64{id},
		"status":     "reviewing",
		"resolution": "",
	})
	if code != 200 {
		t.Fatalf("no-op bulk: %d %v", code, res)
	}
	if res["skipped"] != float64(1) || res["updated"] != float64(0) {
		t.Errorf("updated=%v skipped=%v, want 0/1", res["updated"], res["skipped"])
	}
}
