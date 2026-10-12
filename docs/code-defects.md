# Code defects found by a documentation audit

These are code defects, not documentation defects. A documentation sweep found both because
the prose and the code disagreed, and in each case the prose was closer to correct than a
reader would expect.

They are recorded separately from the documentation work on purpose. A fix here changes
product behaviour, so it needs its own review and its own test.

## 1. The daily XP goal fallback is a literal, and the package above it exists to prevent that

**Location:** `internal/scoring/updater.go:73`

```go
dailyGoal := 10
if st, err := u.repo.GetStudent(studentID); err == nil && st != nil && st.DailyXPGoal > 0 {
    dailyGoal = st.DailyXPGoal
}
```

`xp.DefaultDailyGoal` is 30 (`internal/xp/policy.go:61`). This fallback applies to any learner
whose stored `DailyXPGoal` is absent, zero, or on a row `GetStudent` does not return. Such a
learner is then measured against 10 XP/day.

The package comment at `internal/xp/policy.go:53-56` records this exact failure:

> Daily XP goal policy. Owned here because the server, the estimator and the client all need the
> same numbers, and they had drifted: the shipped default was 10 while the planner's fallbacks
> said 30, so a learner with no stored goal planned against 3x the XP/day they were actually
> banking.

The same field on the struct just above (`updater.go:31-34`) carries the matching history for
`QuizGateXP`. The invariant was enforced everywhere except the one line that needed it most.

**Suggested fix.** Replace the literal with `xp.DefaultDailyGoal`. Add a test that calls the
updater for a student with no stored goal and asserts the fallback equals `xp.DefaultDailyGoal`,
so the two cannot drift apart again.

**Impact.** Wrong progress bars and wrong daily-goal completion percentages for learners with
no stored goal. No data loss, and no effect on XP actually awarded.

## 2. `-repair-grading` reconstructs progress with a fallback threshold

**Location:** `internal/repair/repair.go:144-155`

```go
reqStreak := 3
thresh := 10.0
if c := dag.Concept(conceptID); c != nil {
    reqStreak = c.MasteryThreshold.Streak
    thresh = c.MasteryThreshold.AvgTimeSeconds
}
if reqStreak <= 0 {
    reqStreak = 3
}
if thresh <= 0 {
    thresh = 10.0
}
```

**This is narrower than it first appears, and the correction matters.** The DAG is consulted
first, so a concept in the graph uses its own real `required_streak` and `time_threshold`. The
literals apply in two cases only: the concept is absent from the DAG, or its stored threshold
is zero or negative.

The residual concern is therefore real but small. For those two cases the replay uses
`required_streak = 3` where the live engine would have used a different value, so the rebuilt
`concept_progress` row can differ from the row the live path would have produced. `replay` is
already documented as best-effort at `repair.go:140-142`, and it states that it does not
reproduce halt/rush XP or per-topic learning-speed evolution exactly. The threshold fallback is
a further deviation that the comment does not name.

**Suggested fix.** Either log when the fallback fires, or extend the `replay` doc comment to
name this deviation alongside the two it already names. A test covering a concept with a zero
`AvgTimeSeconds` would pin the behaviour either way.

**Impact.** Low. Affects only repair runs over concepts with a missing or zero threshold.

## What this document is not

Neither defect is a documentation error, and neither is fixed by the documentation sweep. They
are listed here so that the sweep's findings are not read as a claim that the prose is the only
thing that can be wrong.
