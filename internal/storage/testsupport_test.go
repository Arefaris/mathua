package storage

import (
	"database/sql"
	"testing"
)

// rawDB reaches the *sql.DB behind either store.
//
// A few tests seed rows no exported method can create: a session two hours old,
// XP backdated into last week, an attempt from yesterday. Those rows are part of
// the state the method under test reads, not part of what it does, so writing
// them directly is the honest seam. Everything else goes through Repository.
func rawDB(t *testing.T, repo Repository) *sql.DB {
	t.Helper()
	switch s := repo.(type) {
	case *SQLiteStore:
		return s.db
	case *PostgresStore:
		return s.db
	default:
		t.Fatalf("unsupported store %T", repo)
		return nil
	}
}

// execSeed runs a seeding statement written with `?` placeholders, rewriting
// them for Postgres. Keeping one dialect in the tests means a seed cannot drift
// from the statement beside it, and the rewrite is the whole difference: both
// schemas store the XP timestamp as BIGINT nanoseconds.
func execSeed(t *testing.T, repo Repository, query string, args ...any) {
	t.Helper()
	db := rawDB(t, repo)
	if _, ok := repo.(*PostgresStore); ok {
		query = questionMarksToDollars(query)
	}
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("seed %q: %v", query, err)
	}
}

// questionMarksToDollars is the production rewriter, so a seed and a shared
// statement cannot disagree about how a placeholder is spelled.
func questionMarksToDollars(query string) string { return rebindDollars(query) }

// hasIndex reports whether a named index exists on the running engine. The two
// catalogues differ, and a test that asserted on sqlite_master passed only
// because it never ran anywhere else.
func hasIndex(t *testing.T, repo Repository, name string) bool {
	t.Helper()
	var (
		query string
		found string
	)
	if _, ok := repo.(*PostgresStore); ok {
		query = "SELECT indexname FROM pg_indexes WHERE schemaname = current_schema() AND indexname = $1"
	} else {
		query = "SELECT name FROM sqlite_master WHERE type='index' AND name = ?"
	}
	err := rawDB(t, repo).QueryRow(query, name).Scan(&found)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		t.Fatalf("look up index %q: %v", name, err)
	}
	return found == name
}

// querySeed is execSeed's read half, for the counts a test asserts on rather
// than for state it drives. The caller closes the rows.
func querySeed(t *testing.T, repo Repository, query string, args ...any) *sql.Rows {
	t.Helper()
	if _, ok := repo.(*PostgresStore); ok {
		query = questionMarksToDollars(query)
	}
	rows, err := rawDB(t, repo).Query(query, args...)
	if err != nil {
		t.Fatalf("query %q: %v", query, err)
	}
	return rows
}

// seedInt runs a query that returns one number.
func seedInt(t *testing.T, repo Repository, query string, args ...any) int {
	t.Helper()
	rows := querySeed(t, repo, query, args...)
	defer rows.Close()
	var n int
	if !rows.Next() {
		t.Fatalf("query %q returned no row", query)
	}
	if err := rows.Scan(&n); err != nil {
		t.Fatalf("scan %q: %v", query, err)
	}
	return n
}
