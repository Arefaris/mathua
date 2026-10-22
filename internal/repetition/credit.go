package repetition

import (
	"errors"
	"fmt"
)

// Epoch accounting for one (student, target concept).
//
// The units are integer micro-repetitions. Three quantities matter:
//
//	E   retained credit admitted in the epoch (sum of non-voided events)
//	A_b ordinals banked in the epoch, in every state
//	C   E - A_b*Micro, the fractional credit that has not become an ordinal
//
// Banking is what turns credit into a repetition: whenever the fractional
// balance reaches Micro, an ordinal is created and the balance drops by Micro.
// Projection later turns an ordinal into a counterfactual due date; it consumes
// no credit, and it does not restore capacity. Staling an ordinal retires it
// without projecting it.
//
// The invariant that makes reconstruction possible is 0 <= C < Micro: there is
// never a full repetition sitting in the fractional balance, because banking
// would have taken it.
var (
	ErrNegativeCredit = errors.New("repetition: negative credit")
	ErrCreditTooLarge = errors.New("repetition: credit exceeds one repetition")
	ErrBankLimit      = errors.New("repetition: epoch bank limit reached")
	ErrOrdinalTotal   = errors.New("repetition: ordinal counts disagree")
	ErrBalance        = errors.New("repetition: fractional balance out of range")
)

// Account is the reconstructable state of one target in one epoch.
type Account struct {
	RetainedUnits int64
	Ordinals      OrdinalCounts
}

// OrdinalCounts separates the lifecycle states. A_b is their sum: a stale
// ordinal remains banked for historical accounting even though it will never
// be projected, which is why it is counted rather than forgotten.
type OrdinalCounts struct {
	Banked    int
	Projected int
	Stale     int
}

// Total is A_b.
func (c OrdinalCounts) Total() int { return c.Banked + c.Projected + c.Stale }

// Pending is U, the ordinals waiting to be projected.
func (c OrdinalCounts) Pending() int { return c.Banked }

// Balance is C.
func (a Account) Balance() int64 {
	return a.RetainedUnits - int64(a.Ordinals.Total())*Micro
}

// Capacity is how many more units this target may retain in the epoch.
func (a Account) Capacity() int64 {
	remaining := int64(MaxBank)*Micro - a.RetainedUnits
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Admit applies one allocation to the account.
//
// It returns the retained units and the discarded units separately. Discarding
// is not an error: capacity exhaustion is the documented policy, and the caller
// persists both the amount and the reason so a shadow report can show what the
// cap actually cost. Banking happens here, atomically with admission, because a
// separate banking step could be lost between transactions and leave the
// fractional balance above Micro.
func (a Account) Admit(scaled int64) (next Account, retained, discarded int64, err error) {
	if scaled < 0 {
		return a, 0, 0, ErrNegativeCredit
	}
	if scaled > Micro {
		return a, 0, 0, fmt.Errorf("%w: %d", ErrCreditTooLarge, scaled)
	}
	if err := a.Validate(); err != nil {
		return a, 0, 0, err
	}

	retained = scaled
	if capacity := a.Capacity(); retained > capacity {
		retained = capacity
	}
	discarded = scaled - retained

	next = a
	next.RetainedUnits += retained
	// Banking is derived, not incremental: the number of ordinals is exactly
	// the number of whole repetitions the retained total can fund. Deriving it
	// is what makes a retried admission idempotent once the event row is
	// deduplicated by its unique key.
	wantBanked := int(next.RetainedUnits / Micro)
	if wantBanked > a.Ordinals.Total() {
		next.Ordinals.Banked += wantBanked - a.Ordinals.Total()
	}
	if err := next.Validate(); err != nil {
		return a, 0, 0, err
	}
	return next, retained, discarded, nil
}

// Project converts one pending ordinal into a projected one. It changes no
// credit: projection is a claim about a due date, not a new repetition.
func (a Account) Project() (Account, error) {
	if a.Ordinals.Banked <= 0 {
		return a, ErrOrdinalTotal
	}
	next := a
	next.Ordinals.Banked--
	next.Ordinals.Projected++
	return next, nil
}

// Stale retires one pending ordinal without projecting it. Reset uses this for
// the old epoch.
func (a Account) Stale() (Account, error) {
	if a.Ordinals.Banked <= 0 {
		return a, ErrOrdinalTotal
	}
	next := a
	next.Ordinals.Banked--
	next.Ordinals.Stale++
	return next, nil
}

// Validate enforces the ledger invariants.
func (a Account) Validate() error {
	if a.RetainedUnits < 0 {
		return ErrNegativeCredit
	}
	if a.Ordinals.Banked < 0 || a.Ordinals.Projected < 0 || a.Ordinals.Stale < 0 {
		return ErrOrdinalTotal
	}
	total := a.Ordinals.Total()
	if int64(total) > MaxBank {
		return fmt.Errorf("%w: %d ordinals", ErrBankLimit, total)
	}
	balance := a.Balance()
	if balance < 0 || balance >= Micro {
		return fmt.Errorf("%w: %d", ErrBalance, balance)
	}
	return nil
}

// NewAccount builds an account from its persisted columns and validates it.
func NewAccount(retainedUnits int64, counts OrdinalCounts) (Account, error) {
	a := Account{RetainedUnits: retainedUnits, Ordinals: counts}
	return a, a.Validate()
}
