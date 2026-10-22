package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Durable diagnostic sessions (SQLite). The Postgres twin lives in
// diagnostic_postgres.go and mirrors this statement for statement.
//
// updated_at is Unix nanoseconds written by the server rather than by a trigger,
// so both stores age a row on the same clock and the sweeper compares integers
// instead of parsing timestamps in two dialects.

// SaveDiagnosticSession writes the encoded run, replacing any earlier copy.
// Upsert rather than insert-then-update because the first save of a session and
// its twenty-fifth answer are the same event from the caller's side.
func (s *SQLiteStore) SaveDiagnosticSession(id, studentID string, data []byte) error {
	if id == "" {
		return fmt.Errorf("save diagnostic session: empty id")
	}
	_, err := s.db.Exec(
		`INSERT INTO diagnostic_sessions (id, student_id, data, updated_at)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   student_id = excluded.student_id,
		   data       = excluded.data,
		   updated_at = excluded.updated_at`,
		id, studentID, data, time.Now().UTC().UnixNano())
	if err != nil {
		return fmt.Errorf("save diagnostic session: %w", err)
	}
	return nil
}

// LoadDiagnosticSession returns the encoded run. found is false when the row is
// gone, which is a 404, not an error: the caller distinguishes "the learner's
// run ended" from "the database is unreachable".
func (s *SQLiteStore) LoadDiagnosticSession(id string) ([]byte, bool, error) {
	var data []byte
	err := s.db.QueryRow(
		`SELECT data FROM diagnostic_sessions WHERE id = ?`, id).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("load diagnostic session: %w", err)
	}
	return data, true, nil
}

func (s *SQLiteStore) DeleteDiagnosticSession(id string) error {
	if _, err := s.db.Exec(`DELETE FROM diagnostic_sessions WHERE id = ?`, id); err != nil {
		return fmt.Errorf("delete diagnostic session: %w", err)
	}
	return nil
}

// DeleteDiagnosticSessionsForStudent backs ResetProgress and account deletion:
// a wipe that leaves the run on disk would come back the next time the tab
// answered a question.
func (s *SQLiteStore) DeleteDiagnosticSessionsForStudent(studentID string) error {
	if studentID == "" {
		return nil
	}
	if _, err := s.db.Exec(`DELETE FROM diagnostic_sessions WHERE student_id = ?`, studentID); err != nil {
		return fmt.Errorf("delete student diagnostic sessions: %w", err)
	}
	return nil
}

func (s *SQLiteStore) SweepDiagnosticSessions(before time.Time) (int, error) {
	res, err := s.db.Exec(
		`DELETE FROM diagnostic_sessions WHERE updated_at < ?`,
		before.UTC().UnixNano())
	if err != nil {
		return 0, fmt.Errorf("sweep diagnostic sessions: %w", err)
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}
