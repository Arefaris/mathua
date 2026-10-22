package storage

import (
	"bytes"
	"testing"
	"time"
)

func TestDiagnosticSessions_RoundTripAndSweep(t *testing.T) {
	s := newTestStore(t)

	data := []byte(`{"id":"run-1","state":"probing"}`)
	if err := s.SaveDiagnosticSession("run-1", "student-1", data); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, found, err := s.LoadDiagnosticSession("run-1")
	if err != nil || !found {
		t.Fatalf("load: found=%v err=%v", found, err)
	}
	if !bytes.Equal(got, data) {
		t.Errorf("load = %s, want %s", got, data)
	}

	// Saving again replaces the snapshot rather than appending a second row:
	// one run has one current state.
	updated := []byte(`{"id":"run-1","state":"done"}`)
	if err := s.SaveDiagnosticSession("run-1", "student-1", updated); err != nil {
		t.Fatalf("resave: %v", err)
	}
	got, _, _ = s.LoadDiagnosticSession("run-1")
	if !bytes.Equal(got, updated) {
		t.Errorf("after resave = %s, want %s", got, updated)
	}

	// A sweep with a cutoff in the future collects it; one in the past does not.
	if n, err := s.SweepDiagnosticSessions(time.Now().Add(-time.Hour)); err != nil || n != 0 {
		t.Fatalf("sweep with past cutoff: n=%d err=%v", n, err)
	}
	if n, err := s.SweepDiagnosticSessions(time.Now().Add(time.Minute)); err != nil || n != 1 {
		t.Fatalf("sweep with future cutoff: n=%d err=%v", n, err)
	}
	if _, found, _ := s.LoadDiagnosticSession("run-1"); found {
		t.Error("swept session is still readable")
	}
}

func TestDiagnosticSessions_StudentDeleteLeavesOthers(t *testing.T) {
	s := newTestStore(t)
	mustSave := func(id, student string) {
		t.Helper()
		if err := s.SaveDiagnosticSession(id, student, []byte(`{"id":"`+id+`"}`)); err != nil {
			t.Fatalf("save %s: %v", id, err)
		}
	}
	mustSave("a", "student-1")
	mustSave("b", "student-2")

	if err := s.DeleteDiagnosticSessionsForStudent("student-1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, found, _ := s.LoadDiagnosticSession("a"); found {
		t.Error("deleted student's run remains")
	}
	if _, found, _ := s.LoadDiagnosticSession("b"); !found {
		t.Error("another student's run was deleted")
	}
}

func TestDiagnosticSessions_EmptyIDIsRefused(t *testing.T) {
	s := newTestStore(t)
	if err := s.SaveDiagnosticSession("", "student-1", []byte(`{}`)); err == nil {
		t.Error("saving a session without an id was accepted")
	}
}
