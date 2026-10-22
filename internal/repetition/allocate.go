// Package repetition implements the accounting half of the FIRe-inspired
// implicit credit ledger.
//
// It is deliberately free of storage, time and scheduler dependencies: the
// allocator and the ledger arithmetic are pure functions so that the same
// inputs always produce the same allocation and the same balance, on every
// machine and in every transaction retry. The storage layer owns ordering,
// locking and persistence; this package owns the numbers.
//
// Nothing here touches mastery, streaks, XP, diagnostics or the production
// review schedule. Shadow accounting is the whole scope.
package repetition

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

// Micro is the number of accounting units in one implicit repetition.
//
// Credit is integer micro-units, never a float. Repeated float addition would
// make "0.2 five times" a matter of platform rounding, and the whole point of
// the ledger is that a balance can be reconstructed exactly.
const Micro int64 = 1_000_000

// MaxBank is the lifetime number of implicit repetitions one target concept
// may bank per student per epoch.
//
// It is a cap on admission, not on outstanding ordinals: projecting or staling
// an ordinal does not restore capacity. That is a product decision, recorded
// here because a reader will otherwise assume the friendlier semantics.
const MaxBank = 3

// MaxCandidates bounds one source attempt's fan-out. A concept with more than
// a thousand encompassed targets is a graph defect, not a legitimate case, and
// the bound keeps every intermediate sum inside int64 without heroics.
const MaxCandidates = 1024

var (
	ErrTooManyCandidates = errors.New("repetition: too many candidates")
	ErrNegativeRaw       = errors.New("repetition: negative raw units")
	ErrRawTooLarge       = errors.New("repetition: raw units exceed one repetition")
	ErrNegativeCap       = errors.New("repetition: negative cap")
	ErrCapTooLarge       = errors.New("repetition: cap exceeds one repetition")
	ErrDuplicateTarget   = errors.New("repetition: duplicate target concept")
	ErrOverflow          = errors.New("repetition: integer overflow")
)

// Candidate is one target concept and the raw credit it would receive from a
// single source attempt, before the aggregate per-attempt cap is applied.
type Candidate struct {
	ConceptID string
	RawUnits  int64
}

// Allocation is one emitted target allocation, after the cap and after
// zero-unit candidates have been dropped.
type Allocation struct {
	ConceptID string
	Units     int64
}

// Allocate applies the per-source-attempt cap to raw target credits.
//
// The source attempt may grant at most capUnits in total, whatever the graph
// suggests, because one exercise is one exercise however many concepts it
// touches. When the raw total fits, credit is passed through unchanged; when it
// does not, every target is scaled proportionally and the leftover units are
// given to the largest remainders, ties broken by ConceptID so the result does
// not depend on the order candidates arrived in.
func Allocate(candidates []Candidate, capUnits int64) ([]Allocation, error) {
	if len(candidates) > MaxCandidates {
		return nil, fmt.Errorf("%w: %d", ErrTooManyCandidates, len(candidates))
	}
	if capUnits < 0 {
		return nil, ErrNegativeCap
	}
	if capUnits > Micro {
		return nil, ErrCapTooLarge
	}

	seen := make(map[string]struct{}, len(candidates))
	var total int64
	for _, c := range candidates {
		if c.RawUnits < 0 {
			return nil, fmt.Errorf("%w: %s", ErrNegativeRaw, c.ConceptID)
		}
		if c.RawUnits > Micro {
			return nil, fmt.Errorf("%w: %s", ErrRawTooLarge, c.ConceptID)
		}
		if _, dup := seen[c.ConceptID]; dup {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateTarget, c.ConceptID)
		}
		seen[c.ConceptID] = struct{}{}
		var err error
		if total, err = addChecked(total, c.RawUnits); err != nil {
			return nil, err
		}
	}

	// A zero cap is valid and means "emit nothing". It is not an error: an
	// attempt can legitimately have no budget, and the caller persists the
	// empty allocation rather than treating it as a failure.
	if capUnits == 0 || total == 0 {
		return nil, nil
	}

	sorted := make([]Candidate, len(candidates))
	copy(sorted, candidates)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ConceptID < sorted[j].ConceptID })

	out := make([]Allocation, 0, len(sorted))
	if total <= capUnits {
		for _, c := range sorted {
			if c.RawUnits > 0 {
				out = append(out, Allocation{ConceptID: c.ConceptID, Units: c.RawUnits})
			}
		}
		return out, nil
	}

	type share struct {
		id      string
		units   int64
		remaind int64
	}
	shares := make([]share, 0, len(sorted))
	var assigned int64
	for _, c := range sorted {
		product, err := mulChecked(c.RawUnits, capUnits)
		if err != nil {
			return nil, err
		}
		base := product / total
		if assigned, err = addChecked(assigned, base); err != nil {
			return nil, err
		}
		shares = append(shares, share{id: c.ConceptID, units: base, remaind: product % total})
	}

	leftover := capUnits - assigned
	// Largest remainder first, then the stable concept id. Sorting by id and
	// then performing a stable sort on the remainder keeps the tie-break
	// explicit rather than incidental.
	sort.SliceStable(shares, func(i, j int) bool { return shares[i].remaind > shares[j].remaind })
	for i := 0; i < len(shares) && leftover > 0; i++ {
		shares[i].units++
		leftover--
	}

	out = out[:0]
	for _, s := range shares {
		if s.units > 0 {
			out = append(out, Allocation{ConceptID: s.id, Units: s.units})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ConceptID < out[j].ConceptID })

	var emitted int64
	for _, a := range out {
		var err error
		if emitted, err = addChecked(emitted, a.Units); err != nil {
			return nil, err
		}
	}
	if emitted > capUnits {
		return nil, fmt.Errorf("%w: emitted %d over cap %d", ErrOverflow, emitted, capUnits)
	}
	return out, nil
}

func addChecked(a, b int64) (int64, error) {
	if b > 0 && a > math.MaxInt64-b {
		return 0, ErrOverflow
	}
	if b < 0 && a < math.MinInt64-b {
		return 0, ErrOverflow
	}
	return a + b, nil
}

func mulChecked(a, b int64) (int64, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}
	if a > 0 && b > 0 && a > math.MaxInt64/b {
		return 0, ErrOverflow
	}
	return a * b, nil
}
