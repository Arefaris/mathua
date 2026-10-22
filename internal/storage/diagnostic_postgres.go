package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Durable diagnostic sessions (Postgres). Mirrors diagnostic_sqlite.go statement
// for statement so a rule one store enforces is not a rule the other lacks.

func (s *PostgresStore) SaveDiagnosticSession(id, studentID string, data []byte) error {
	if id == "" {
		return fmt.Errorf("save diagnostic session: empty id")
	}
	_, err := s.db.Exec(
		`INSERT INTO diagnostic_sessions (id, student_id, data, updated_at)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (id) DO UPDATE SET
		   student_id = EXCLUDED.student_id,
		   data       = EXCLUDED.data,
		   updated_at = EXCLUDED.updated_at`,
		id, studentID, data, time.Now().UTC().UnixNano())
	if err != nil {
		return fmt.Errorf("save diagnostic session: %w", err)
	}
	return nil
}

func (s *PostgresStore) LoadDiagnosticSession(id string) ([]byte, bool, error) {
	var data []byte
	err := s.db.QueryRow(
		`SELECT data FROM diagnostic_sessions WHERE id = $1`, id).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("load diagnostic session: %w", err)
	}
	return data, true, nil
}

func (s *PostgresStore) DeleteDiagnosticSession(id string) error {
	if _, err := s.db.Exec(`DELETE FROM diagnostic_sessions WHERE id = $1`, id); err != nil {
		return fmt.Errorf("delete diagnostic session: %w", err)
	}
	return nil
}

func (s *PostgresStore) DeleteDiagnosticSessionsForStudent(studentID string) error {
	if studentID == "" {
		return nil
	}
	if _, err := s.db.Exec(`DELETE FROM diagnostic_sessions WHERE student_id = $1`, studentID); err != nil {
		return fmt.Errorf("delete student diagnostic sessions: %w", err)
	}
	return nil
}

func (s *PostgresStore) SweepDiagnosticSessions(before time.Time) (int, error) {
	res, err := s.db.Exec(
		`DELETE FROM diagnostic_sessions WHERE updated_at < $1`,
		before.UTC().UnixNano())
	if err != nil {
		return 0, fmt.Errorf("sweep diagnostic sessions: %w", err)
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}
