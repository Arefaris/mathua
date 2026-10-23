package storage

import "testing"

// TestDiagnosticAttempt_ReplaysTheFirstVerdict pins the idempotency rule that
// the diagnostic depends on: the first answer for a session+concept is the one
// that counts, and a retry reads it back instead of grading again.
func TestDiagnosticAttempt_ReplaysTheFirstVerdict(t *testing.T) {
	store := newTestStore(t)
	const key = "session-1:geo.basic.angles"

	if _, found, err := store.DiagnosticAttemptResponse(key); err != nil || found {
		t.Fatalf("unknown attempt: found=%v err=%v", found, err)
	}
	if err := store.RecordDiagnosticAttempt(key, `{"done":false,"correct":true}`); err != nil {
		t.Fatalf("record: %v", err)
	}
	// A duplicate with different content must not replace the verdict.
	if err := store.RecordDiagnosticAttempt(key, `{"done":true,"correct":false}`); err != nil {
		t.Fatalf("duplicate record: %v", err)
	}
	got, found, err := store.DiagnosticAttemptResponse(key)
	if err != nil || !found {
		t.Fatalf("read back: found=%v err=%v", found, err)
	}
	if got != `{"done":false,"correct":true}` {
		t.Errorf("stored verdict = %q, want the first one", got)
	}
}

func TestDiagnosticAttempt_RefusesEmptyValues(t *testing.T) {
	store := newTestStore(t)
	if err := store.RecordDiagnosticAttempt("", "{}"); err == nil {
		t.Error("empty attempt id was accepted")
	}
	if err := store.RecordDiagnosticAttempt("k", ""); err == nil {
		t.Error("empty response was accepted")
	}
}
