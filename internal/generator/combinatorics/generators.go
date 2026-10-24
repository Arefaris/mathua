// Package combinatorics generates the counting problems of the `discrete_math` combinatorics
// subdomain.
//
// The six basic counts live in internal/generator/discrete/combinatorics_batch.go: counting
// principle, factorial, multinomial, circular arrangements, derangements and Catalan numbers.
// They were registered here as well until 23 October 2026, and because the `all` package
// registers `discrete` first, the second registration returned "generator already registered",
// an error every caller ignores. This package's versions were therefore unreachable while their
// tests kept passing, because those tests built the generators directly. One concept has one
// generator: the batch owns those six now.
//
// The counting techniques below stay here. Every concept is `numeric`: a counting answer is an
// exact integer that a learner produces rather than an option they pick, which is the change
// section 4.2 of docs/content-expansion.md argues for in the domains that are entirely
// `multiple_choice`. The explanation is derived from the same integers as the answer, so the two
// cannot disagree.
package combinatorics

import (
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Register adds this subdomain's generators to reg.
func Register(reg *generator.Registry) {
	registerCountingTechniques(reg)
}

// pick returns a random int in [0, n) and never panics on n < 1, so a difficulty of 0 (which the
// easy-pool ratchet samples at 0.12 but a caller may pass as 0) yields a one-element range rather
// than a panic inside rand.Intn.
func pick(n int) int {
	if n < 1 {
		n = 1
	}
	return rand.Intn(n)
}

// choose is the binomial coefficient nCr, computed multiplicatively so it stays exact.
func choose(n, r int) int {
	if r < 0 || r > n {
		return 0
	}
	if r > n-r {
		r = n - r
	}
	c := 1
	for i := 0; i < r; i++ {
		c = c * (n - i) / (i + 1)
	}
	return c
}
