# Content expansion

This document records what the corpus is, how a new piece of content is added, and where the
corpus is thin enough that adding content would change what a learner actually meets. It is an
**analytical-tier** doc (AGENTS.md section 6): it carries measurements and the method behind
them, so `scripts/check_ste.py` does not cover it and would report violations if run against it.
That is the expected result, not a defect.

Every figure below is read from `data/`, from `internal/`, or from a script in `scripts/`. None
is hand-written. Section 7 gives the command that reproduces each one, because a number in a doc
and the artifact it describes are one fact (AGENTS.md section 7).

## 1. What the corpus contains

| Layer | Owner | Tool or origin |
|---|---|---|
| Concepts (DAG nodes) | `data/concepts/<domain>.json`, one file per domain | hand-authored |
| Enrichment overlays | `data/concepts/enrichment.json` | `encompasses`, `key_prerequisites`, `interference_group`, `variants` |
| Problems | `internal/generator/<domain>/` | generated on demand, never stored |
| Lessons | `data/lessons/{teaching,algebrica,authored}/**.md` | imports plus authored files |
| Lesson index | `data/lessons/lessons.json` | concept id to lesson source |
| KP shards | `data/lessons/kp/<concept_id>.json` | `scripts/gen_kp_shards.py` from an inline spec |
| Solution schemas | `data/lessons/solutions/<concept_id>.json` | hand-authored, learner domains only |
| Diagrams | `data/diagrams/meta.json` plus assets | Algebrica imports and generated figures |
| Courses | `data/courses.json` | named target sets |
| Destinations | `data/destinations.json` | named bundles of courses |

Two numbers are useful to hold together. There are 666 concepts, and there are 666 KP shard files
holding 1998 knowledge points, three per concept.

## 2. The pipeline

Each stage has one owner file, and one guard that fails when the stage is wrong.

| # | Stage | Owner | Guard |
|---|---|---|---|
| 1 | Concept node | `data/concepts/<domain>.json` | `scripts/validate_graph.go`: no cycles over `prerequisites` and `encompasses` (Kahn's), no orphans, no duplicate ids, valid grading enum, `streak` positive, `avg_time_seconds` in range, no singleton interference group |
| 2 | Enrichment | `data/concepts/enrichment.json` | the same loader: a `key_prerequisite` must already be a prerequisite, variant difficulty in range |
| 3 | Generator | `internal/generator/<domain>/`, registered in `internal/generator/all/all.go` | `Registry.Register` self-check (a graded generator must accept its own answer and reject an empty one), `TestFuzz`, and the pool ratchet in `internal/generator/pools_test.go` |
| 4 | Lesson | lesson markdown plus `data/lessons/lessons.json` | `scripts/audit_lessons.py`: every concept has a lesson, every id is live, every file is on disk, no orphan source |
| 5 | KP shard | `data/lessons/kp/`, generated from the `SPEC` in `scripts/gen_kp_shards.py` | `section` must name a real heading in the lesson; the three-KP rule is read from the shards by `scripts/check_docs.py` |
| 6 | Solution schema | `data/lessons/solutions/` | `TestLearnerDomainsHaveSchemas` over the six learner domains, plus `solutions_contract_test` |
| 7 | Diagram | `data/diagrams/meta.json` | `scripts/audit_lessons.py`: a mapping must point at an existing asset and an existing concept |
| 8 | Course | `data/courses.json` | `TestAllCoursesResolve`: every target resolves in the DAG |
| 9 | Documented counts | `README.md` | `scripts/counts.py --check` and `scripts/check_docs.py` |

The order is a dependency order. A concept with no generator cannot serve a problem. A generator
with no lesson has nothing to teach before practice. A KP shard whose `section` does not resolve
leaves the learner a placeholder instead of a worked example, and the audit refuses it.

## 3. What complete means today

All four content guards pass on the current corpus:

```
OK: 666 concepts across 17 domains                    scripts/validate_graph.go
OK: 666 concepts mapped, 666 entries, 565 sources,
    666 kp files, 1998 KPs, kp shards resolved,
    0 issues                                          scripts/audit_lessons.py
counts ok: 666 concepts, 666 KP files, 22 courses     scripts/counts.py --check
OK: docs agree with the code and the corpus           scripts/check_docs.py
```

`scripts/audit_baseline.json` is empty on every key it holds, including `stub_lessons`,
`placeholder_kps`, `thin_pools_single_answer` and `thin_pools_yesno_only`. That file is a
ratchet: known violations live there and shrink to empty as content lands, and a new violation
fails the build. Empty means the corpus is clean and the next wave starts from zero, but it also
means a new thin concept cannot be added quietly.

## 4. Where the corpus is thin

### 4.1 The course catalog reaches 247 of the 666 concepts

Courses carry target sets, and a target expands to its transitive prerequisite closure in
`internal/planning/planner.go:64` (`collectChain`). Across the 22 courses those closures reach
247 of the 666 concepts: 419 of them have a node, a generator, a lesson and three KP shards and
still sit outside every named course path. A learner reaches them through the diagnostic or a
free-choice goal, not through the catalog.

The count of unreachable concepts by domain:

| Domain | Concepts | Unreachable |
|---|---|---|
| calculus | 128 | 100 |
| algebra | 89 | 53 |
| topology | 30 | 27 |
| statistics | 37 | 24 |
| precalculus | 25 | 22 |
| fractions | 26 | 22 |
| machine_learning | 25 | 21 |
| complex_numbers | 26 | 20 |
| number_theory | 25 | 20 |
| differential_equations | 25 | 19 |
| prealgebra | 46 | 23 |
| trigonometry | 32 | 17 |
| arithmetic | 46 | 14 |
| linear_algebra | 26 | 12 |
| abstract_algebra | 32 | 11 |
| discrete_math | 25 | 7 |
| geometry | 23 | 7 |

Two structural facts sit behind that table. First, a course is a *sampled* path, not a domain
cover: `calc1` reaches 15 calculus concepts and `calc2` reaches 24, and their targets differ, so
the two courses are stages rather than a partition. Second, one whole domain has no course of its
own: `complex_numbers` is reached only through the `pc` course, which pulls in 6 of its 26
concepts.

### 4.2 Whole domains ask a learner to pick rather than to produce

The `grading_type` mix shows where the corpus is recognition-shaped. These domains are almost
entirely `multiple_choice`:

| Domain | Concepts | `multiple_choice` |
|---|---|---|
| abstract_algebra | 32 | 32 |
| topology | 30 | 29 |
| differential_equations | 25 | 24 |
| machine_learning | 25 | 24 |
| discrete_math | 25 | 19 |
| number_theory | 25 | 17 |
| complex_numbers | 26 | 14 |

Against that, `calculus` carries 86 `symbolic` concepts and `algebra` carries 36. The pool
ratchet in `internal/generator/pools_test.go` exists for exactly this surface: an easy pool that
yields a single distinct answer, or that is yes/no only, is a violation. It is empty today
because every registered generator clears the bar for its current grading type, not because a
multiple-choice pool produces much. Turning one of these concepts into `numeric` or `symbolic`
means a real generator and a solution schema, because the explanation stops being one fixed
string.

**The first conversion landed.** `topo.dimension` was a table of yes/no rows; it now asks for the
Lebesgue covering dimension of a Euclidean space, a sphere, a torus or the Cantor set, computes
the number from the ambient dimension, and publishes `space`, `dim` and `rule` as facts for the
authored schema in `data/lessons/solutions/topo.dimension.json`. Its `grading_type` is `numeric`,
so topology is 29 `multiple_choice` concepts rather than 30, and the worked example of the path
is the four files it touched: `data/concepts/topology.json`,
`internal/generator/topology/generators.go`, `internal/generator/topology/generators_test.go` and
the schema.

### 4.3 Lesson depth is thin at the bottom

Lessons come from three places. The `teaching/` tree holds 329 mappings, the Algebrica import
holds 328, and 9 are authored. The median lesson is 8648 bytes, and seven lessons are under 1 KB:

| Concept | Lesson | Bytes |
|---|---|---|
| `calc.series.ratio_test` | `data/lessons/authored/calc.series.ratio.md` | 671 |
| `nt.basics.crt` | `data/lessons/teaching/nt.crt.md` | 853 |
| `nt.adv.wilson` | `data/lessons/teaching/nt.wilson.md` | 865 |
| `nt.basics.bezout` | `data/lessons/teaching/nt.bezout.md` | 878 |
| `nt.basics.prime_inf` | `data/lessons/teaching/nt.prime_inf.md` | 920 |
| `abstract.group.order` | `data/lessons/teaching/abstract.group.order.md` | 960 |
| `linalg.decomp.lu` | `data/lessons/teaching/linalg.lu.md` | 1018 |

Four of the seven sit in `number_theory`, the same domain whose course reaches 5 of its 25
concepts.

### 4.4 Ingestion dropped figures and kept their captions

`docs/lesson-media.md` is the owner of this measurement, and it reports 239 paragraphs across
100 lessons that read as a caption whose image was dropped: 135 certain, 33 strong, 71 likely. A
`certain` case opens with an auto-numbered `Figure 3.2` label, which the ingestion pipeline emits
for a figure and which survives when the image does not. `scripts/audit_lesson_media.py`
regenerates the list, and the remaining work is per instance: each candidate is either matched to
an existing asset in `data/diagrams/meta.json` or removed.

### 4.5 Diagrams and solution schemas are complete against their own rules

`data/diagrams/meta.json` maps 183 assets, each to at least one concept and a license.
`data/lessons/solutions/` holds 230 schemas, and the six learner domains require one per
computed concept; the table-driven exceptions are an explicit allowlist with a reason per entry
in `internal/generator/solutions_coverage_test.go`. Growing either is additive work, not repair.

### 4.6 Application siblings: the same technique inside a story

Twenty-one concepts apply a technique their parent teaches rather than teach a new one:
`arith.{add,sub,mult,div}.word`, `frac.{add,sub,mult,div,compare,mixed}.word`,
`dec.ops.{add,sub,mult,div}.word`, `ratio.ops.{rate,proportion}.word`, `prealg.eq.word`,
`alg.eq.vars_both_sides.word`, `alg.ineq.multi_step.word`, `alg.systems.word` and
`geo.triangle.pythagorean.word`. They cover 17 subdomains.

Each one is a *dependent*: a parent technique is one of its prerequisites. Twenty of the
twenty-one are also leaves, so nothing else depends on them and no course closure needs them —
they are inside the 419 unreachable concepts of section 4.1. `prealg.eq.word` is the exception on
both counts: something depends on it, and that something is what pulls it into a course closure.

The pattern is cheap because it writes no lesson prose. A sibling shares its parent's lesson file
— `data/lessons/lessons.json` points both ids at one source — so its three KP sections come from
headings that already exist, and the new files are the shard and the authored solution schema. The
generator and its test extend package files that are already there, and the concept node joins an
existing domain file.

Where it stops, measured against the same two conditions: the parent lesson must carry an
application-shaped heading to draw three KPs from. `trig.basics.sin_cos_def` has one and it is
"Applications in integration", which is not a word problem; `geo.quad.perimeter` and
`geo.solid.volume_rect` have none; `prealg.ineq.two_step` and `alg.quad.discriminant` have one each
and carry two empty headings, so a shard there would need written prose first. Past this point the
work is lesson authoring, not generator writing.

## 5. Options, and what each one costs

### 5.1 Course wiring: sample more, or cover the domain

The gap in section 4.1 has three plausible policies.

- **Add a domain's sinks as targets.** A sink is a concept with no dependents in its domain, and
  every concept reaches some sink, so the sink set covers the domain by construction. This is
  cheap and mechanical, and it is wrong on its own: the calculus sinks are 77 concepts, so
  `calc1` and `calc2` would both become the whole domain and stop being stages.
- **Add the sinks of the depth band a course occupies.** This preserves staging and is the shape
  `cmd/masteryunits` already uses for its draft partition. It needs the band boundary to be a
  decision first.
- **Add one course per uncovered domain.** The only domain with no course at all is
  `complex_numbers`. This is a small, additive change, but it changes the course count, and the
  course count is read from `data/courses.json` by `scripts/counts.py` and quoted in the README.

None of the three is safe to apply silently. All three change the size of a course, and course
size drives the mastered/total figure and the days-remaining estimate a learner sees.

### 5.2 Producing an answer in a multiple-choice domain

This is the highest-value per-concept change and it needs no catalog decision. It is a generator
plus a solution schema for one concept: change `grading_type`, write a `GradedGenerator` whose
easy pool discriminates, add the schema so a missed answer carries an authored explanation, and
let the pool ratchet and the schema contract test say whether it holds.

### 5.3 Filling a thin lesson

A lesson under 1 KB is a lesson with one heading and little else, and its three KP shards must
name three different headings in it, which is where the placeholder risk lives. Filling one means
writing the body first and regenerating the shards second, because the shard sections resolve
against the body.

## 6. A sequenced plan

1. Decide the course-wiring policy in section 5.1, then apply it with `TestAllCoursesResolve`
   as the guard and `scripts/counts.py --write` for the README.
2. Convert one `multiple_choice` concept in an upper domain to a produced answer, as the worked
   example that proves the generator and schema path there. Done: `topo.dimension`, section 4.2.
3. Fill the seven sub-1 KB lessons, regenerating each shard after its body grows.
4. Reclaim the `certain` dropped figures from `docs/lesson-media.md` into
   `data/diagrams/meta.json`, one at a time.
5. Add a domain only after 1 to 4, using the topology addition as the copy target, because a
   domain costs a concept file, a generator package, a registration, lessons, shards, a course
   target and the domain label in `scripts/counts.py`.

## 7. Method

Reproduce each measurement. All commands run from the repository root.

```bash
go run scripts/validate_graph.go
python3 scripts/audit_lessons.py
python3 scripts/counts.py --check
python3 scripts/check_docs.py
```

Course reach and the grading-type mix need the DAG and the course targets read together; the
reach number is the union of each course's prerequisite closure, computed the way
`internal/planning/planner.go:64` computes it, and the checker in `scripts/check_docs.py`
recomputes it rather than trusting this file.

```bash
python3 - <<'PY'
import collections, glob, json
concepts, prereq = {}, {}
for path in glob.glob("data/concepts/*.json"):
    if path.endswith("enrichment.json"):
        continue
    for c in json.load(open(path)):
        concepts[c["id"]] = c
        prereq[c["id"]] = set(c.get("prerequisites") or [])

def closure(ids):
    seen, stack = set(), list(ids)
    while stack:
        cid = stack.pop()
        if cid in seen:
            continue
        seen.add(cid)
        stack += list(prereq.get(cid, ()))
    return seen

reach = set()
for course in json.load(open("data/courses.json")):
    reach |= closure(course["targets"])
print("reach", len(reach), "of", len(concepts))
by_domain = collections.Counter(concepts[c]["domain"] for c in set(concepts) - reach)
print("unreachable", dict(by_domain.most_common()))
mix = collections.Counter(concepts[c].get("grading_type") for c in concepts)
print("grading types", dict(mix))
PY
```

Lesson sizes and sources come from `data/lessons/lessons.json`; figure candidates come from
`scripts/audit_lesson_media.py`. The application siblings of section 4.6 are the concept ids ending
in `.word` in `data/concepts/*.json`; whether anything depends on one, and whether a course closure
reaches it, are read from those files and from `data/courses.json`.
