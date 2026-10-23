# Mastery semantics

Three mechanisms decide what a learner is shown and when a concept is certified.
They are easy to confuse because each of them produces evidence, and only one of
them is allowed to certify anything.

## The ladder certifies

`internal/mastery/evidence.go` defines the rungs. `RungDemonstrated` is rung 3,
and it requires variety and difficulty, not repetition. `Machine.Next` moves a
concept up a rung; `mastery.EnoughForMastery` is the only place that answers yes.

Anything else that looks like progress is an input to the ladder. None of the
three mechanisms below may return `MASTERED` on its own.

## Knowledge points pace, they do not certify

A concept may declare `knowledge_points`. They split a topic into steps so a
learner is asked one thing at a time instead of the hardest variant first.

Completing every point used to return `MASTERED` from `Machine.Next`. That
bypassed the ladder entirely: four correct answers against one instance at one
difficulty certified a concept that `RungDemonstrated` refuses. Completion now
falls through to the ladder, and a concept with unfinished points is held at its
current rung rather than being promoted.

A miss no longer leaves a finished scaffold finished. `Transition` clears
`Completed` and reopens the last point. The ordinal is not decremented, because
the streak reset already records that the point being worked on is unearned, and
stepping back on every slip pushes a learner away from material they just
reached.

Scaffold state is reconstructed on restart from the ordinal persisted with the
session (`NewKPStateMachineAt`), not from the progress row. Seeding completion
from a `MASTERED` row would restore the bypass in a different disguise, so a
mastered concept with knowledge points may still be asked scaffold questions
until the ladder itself certifies. That is the conservative reading of the split:
the scaffold decides which question, the ladder decides whether it counts.

## Placement grants reachability, not mastery

A diagnostic may conditionally complete a concept the learner barely placed out
of. `ApplyDiagnosticReport` writes the row as `LEARNING` with `PlacementSeeded`
and `IsProvisional` set. Reachability comes from the `Placement` flag in the
scheduler, which is what that flag is for.

One correct answer used to clear `IsProvisional`, because hysteresis keeps the
seeded row at `MASTERED` and the clearing condition only looked at the status.
It now requires `mastery.EnoughForMastery`. The same transition clears
`PlacementSeeded`, which is what lets a later diagnostic tell real work from a
seed instead of re-seeding demonstrated mastery.

A placement can still be demoted: a wrong answer on a dependent concept pushes a
provisional prerequisite back to `LEARNING` with a high weakness score.

## Implicit credit is review evidence, one advance per cycle

A correct answer on a downstream concept earns fractional credit toward the
prerequisites it encompasses. That credit banks an ordinal, and an ordinal can
advance a review without asking the prerequisite directly.

It used to advance once per ordinal, so ten answers in one sitting applied the
SM-2 step ten times: an interval of 3 days walked to 30,595 days, due in 2166,
and the concept was never reviewed again. Credit is now banked at most once per
`(student, concept, cycle)`. Consuming it closes the cycle by advancing
`credit_epoch`, and the next eligible review can earn another. Refused units are
recorded with `discard_reason = "cycle_already_advanced"` rather than dropped.

Credit earned while a concept has no review date stays banked. It is not applied
later as several spaced reviews, because it was never several reviews.

## What to check when changing any of this

- A concept reaches `MASTERED` only through `EnoughForMastery`.
- A rewording of progress that returns `StatusMastered` directly is a regression.
- Credit arithmetic must multiply an interval once per review, never once per
  answer.
