#!/usr/bin/env python3
"""
check_docs.py — find claims in the engineering docs that the code or the corpus contradicts.

Why this exists, in the shape of the bug that motivated it: README stated a 150 XP
mastery-check gate twice while `internal/xp/policy.go` said 50. Nothing was broken, no test
failed, and the landing page carried the same stale figure until it was fixed separately. A
documented number and the constant it describes are the same fact written down twice, and
nothing noticed.

The design rule that follows from that: **this script never hardcodes a fact it can read.** The
XP gate comes from the Go source, the corpus counts come from the corpus, and the KP-per-shard
rule comes from the shards. A checker that pins its own copy of the number goes stale in exactly
the way it was written to prevent.

Scope is the engineering docs only. `data/lessons/**` is learner-facing corpus content, is
covered by `audit_lessons.py`, and is deliberately not subject to American spelling or to the
link rules here.

Usage:
    python3 scripts/check_docs.py
Exit code 0 = no errors. Exit code 1 = at least one error.
"""

from __future__ import annotations

import json
import os
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

# Built once per run by file_index(); see check_symbol_citations.
_FILE_INDEX = None

# ---------------------------------------------------------------- what we check

# Engineering docs. Deliberately explicit rather than a glob, so a new doc is a deliberate act
# rather than something that silently falls under or out of the rules.
DOCS = [
    "README.md",
    "DESIGN.md",
    "CONTRIBUTING.md",
    "AGENTS.md",
    "CONTEXT.md",
    "docs/architecture.md",
    "docs/administration.md",
    "docs/code-defects.md",
    "docs/efficacy.md",
    "docs/glossary-provenance.md",
    "docs/lesson-media.md",
    "docs/mastery-assessments.md",
    "docs/mastery-semantics.md",
    "docs/math-input.md",
    "docs/mobile-device-checklist.md",
    "docs/system-design.md",
    "docs/system-design-slides.md",
    "docs/wasc-uc-ag-research.md",
]
# The ADRs are a generated, uniform set, so a glob is the honest description rather than 54
# lines that each need editing when a decision is added. They were inline in CONTEXT.md and
# were covered by the rules below until the split, so they stay covered after it.
DOCS += [str(p.relative_to(ROOT)) for p in sorted((ROOT / "docs/adr").glob("ADR-*.md"))]

# Files exempt from the spelling rule, each with its reason. An exemption without a reason is
# indistinguishable from a mistake.
SPELLING_EXEMPT = {
    "THIRD_PARTY_NOTICES.md": "names third-party license instruments (MIT, CC BY-NC 4.0); "
                              "their wording is not ours to rewrite",
}

# Paths under the Next.js app, which resolve as routes at runtime and have no file of that name.
APP_ROUTES = ("/docs", "/learn", "/study", "/review", "/onboard", "/profile", "/graph",
              "/plan", "/goals", "/history", "/leaderboard", "/settings", "/admin",
              "/transcript", "/home", "/rewards")

BRITISH = {
    "behaviour": "behavior", "behaviours": "behaviors", "behavioural": "behavioral",
    "colour": "color", "colours": "colors", "coloured": "colored", "colouring": "coloring",
    "labelled": "labeled", "labelling": "labeling",
    "licence": "license", "licences": "licenses",
    "centre": "center", "centres": "centers",
    "defence": "defense",
    "analyse": "analyze", "analysed": "analyzed", "analyses": "analyzes", "analysing": "analyzing",
    "organisation": "organization", "organisations": "organizations",
    "recognise": "recognize", "recognised": "recognized", "recognises": "recognizes",
    "normalise": "normalize", "normalised": "normalized", "normalises": "normalizes",
    "canonicalise": "canonicalize", "canonicalised": "canonicalized",
    "serialise": "serialize", "serialised": "serialized",
    "prioritise": "prioritize", "prioritised": "prioritized",
    "whilst": "while", "amongst": "among",
}


# Files allowed to cite a gate figure other than the current one, each with its reason. Both cite
# the old 150 XP value *as the example of the bug*, which is the only place that value should
# still appear. A whole-file exemption is a hole, so the list is short, explicit, and justified.
HISTORICAL_GATE_CITERS = {
    "CONTEXT.md": "ADR-020 records the 150 -> 50 rescale as history",
    "AGENTS.md": "quotes the stale 150 XP figure as the motivating example for this rule",
    "docs/adr/ADR-003-q3.md": "the superseded decision whose title names ADR-020 as its replacement",
}


class Report:
    def __init__(self) -> None:
        self.errors: list[tuple[str, str, str]] = []

    def error(self, where: str, rule: str, msg: str) -> None:
        self.errors.append((where, rule, msg))


# ---------------------------------------------------------------- helpers

def read(rel: str) -> str:
    return (ROOT / rel).read_text(encoding="utf-8")


def strip_code(text: str) -> str:
    """Blank out regions where a word is an identifier, not prose.

    Without this, a doc that mentions `behaviour` inside a fenced block or inline code trips the
    spelling rule on a token the author does not control.
    """
    out = re.sub(r"```.*?```", lambda m: "\n" * m.group(0).count("\n"), text, flags=re.S)
    out = re.sub(r"`[^`\n]*`", lambda m: " " * len(m.group(0)), out)
    return out


def load_kps(path: Path) -> list:
    """A shard is a bare JSON list in some files and an object with a "kps" key in others."""
    data = json.loads(path.read_text(encoding="utf-8"))
    if isinstance(data, dict):
        return data.get("kps") or []
    return data if isinstance(data, list) else []


def go_const(path: str, name: str) -> int | None:
    """Read an integer constant out of Go source. Returns None if absent."""
    m = re.search(rf"^\s*(?:const\s+)?{re.escape(name)}\s*=\s*(\d+)", read(path), re.M)
    return int(m.group(1)) if m else None


# ---------------------------------------------------------------- checks

def check_xp_gate(rep: Report) -> None:
    """Any mastery-check gate figure in the docs must equal the constant in the code.

    The gate word and the figure must be close together. An earlier version flagged any "N XP"
    on a line that mentioned "quiz" or "mastery", which fired on legitimate statements of other
    XP amounts -- "diagnostic 0 XP", a historical "+3 XP" bug fix -- because those words are
    everywhere in this corpus. Requiring the figure within 40 characters of the gate word keeps
    the rule on the claim it was written for.
    """
    gate = go_const("internal/xp/policy.go", "QuizGateXP")
    if gate is None:
        rep.error("internal/xp/policy.go", "quiz-gate",
                  "QuizGateXP not found; this checker cannot verify the docs")
        return
    pat = re.compile(r"(\d+)\s*XP", re.I)
    gateword = re.compile(r"\bgate\b|mastery[- ]check", re.I)
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        for i, line in enumerate(read(doc).split("\n"), 1):
            positions = [g.start() for g in gateword.finditer(line)]
            if not positions:
                continue
            for m in pat.finditer(line):
                if not any(abs(m.start() - g) <= 40 for g in positions):
                    continue
                n = int(m.group(1))
                if n == gate or doc in HISTORICAL_GATE_CITERS:
                    continue
                rep.error(f"{doc}:{i}", "xp-gate",
                          f'"{m.group(0)}" but xp.QuizGateXP is {gate}')


def sentences(text: str) -> list[tuple[int, str]]:
    """Split into (line number, sentence). A rule is stated per sentence, not per line."""
    out = []
    for i, line in enumerate(text.split("\n"), 1):
        s = line.strip()
        if not s or s.startswith(("|", "#", "-", "*", ">", "`")):
            continue
        for part in re.split(r"(?<=[.!?:])\s+", s):
            if part.strip():
                out.append((i, part.strip()))
    return out


def check_kp_rule(rep: Report) -> None:
    """The per-concept KP rule must match the corpus, and must be stated one way in the docs."""
    counts: dict[int, int] = {}
    for f in sorted((ROOT / "data/lessons/kp").glob("*.json")):
        kps = load_kps(f)
        if kps:
            counts[len(kps)] = counts.get(len(kps), 0) + 1
    if not counts:
        rep.error("data/lessons/kp", "kp-rule", "no KP shards found")
        return
    if len(counts) > 1:
        rep.error("data/lessons/kp", "kp-rule",
                  f"shards hold different KP counts: {counts}; the documented rule cannot be "
                  "one number, so state the real distribution")
        return
    kps_per = next(iter(counts))
    shards = sum(counts.values())
    total = sum(k * v for k, v in counts.items())

    # A rule claim is a number that is explicitly scoped: either "exactly N KPs", or "N KPs per
    # concept". Corpus totals such as "657 KP shard files x 3 = 1971 KPs" are not claims about
    # the per-shard rule, and reading them as such flags the arithmetic that proves the rule.
    claim = re.compile(r"([0-9]+)(?:\s*[-\u2013]\s*([0-9]+))?\s*(?:KPs?\b|knowledge points?\b)",
                       re.I)
    scoped = re.compile(r"\bper\s+(?:concept|shard|lesson|topic)\b|\bexactly\b", re.I)
    for doc in ("README.md", "CONTRIBUTING.md"):
        for line_no, sent in sentences(read(doc)):
            if not scoped.search(sent):
                continue
            for m in claim.finditer(sent):
                lo = int(m.group(1))
                hi = m.group(2)
                if lo != kps_per or (hi is not None and int(hi) != kps_per):
                    got = f"{lo}-{hi}" if hi else str(lo)
                    rep.error(f"{doc}:{line_no}", "kp-rule",
                              f"docs say {got} KPs per concept but all {shards} shards hold "
                              f"exactly {kps_per} ({total} KPs total)")


# A corpus-scale claim. Anything smaller is a per-concept parameter (a streak, a threshold, a
# percentage) that happens to sit next to the word "concept", and flagging those would bury the
# real findings in noise.
CORPUS_SCALE = 100

def check_corpus_figures(rep: Report) -> None:
    """Corpus counts quoted in docs must match the corpus. Reads the corpus, never a literal."""
    # Count concepts, not domain files: data/concepts holds one file per domain, and 17 files is
    # not the number any doc means by "concepts".
    n_concepts = 0
    for f in sorted((ROOT / "data/concepts").glob("*.json")):
        if f.name == "enrichment.json":
            continue
        data = json.loads(f.read_text(encoding="utf-8"))
        n_concepts += len(data if isinstance(data, list) else list(data.values()))
    facts = {
        "concepts": n_concepts,
        "shards": len(list((ROOT / "data/lessons/kp").glob("*.json"))),
        "kps": sum(len(load_kps(f)) for f in (ROOT / "data/lessons/kp").glob("*.json")),
        "courses": len(json.loads(read("data/courses.json"))),
    }
    nouns = {
        "concepts": r"concepts?",
        "shards": r"(?:KP\s+)?shards?",
        "kps": r"KPs|knowledge points",
        "courses": r"courses",
    }
    pat = re.compile(r"([0-9][0-9,]*)\s+(" + "|".join(f"(?:{v})" for v in nouns.values()) + r")\b",
                     re.I)
    # "128 concepts" can be a per-domain subset ("calculus is 128 concepts") rather than a corpus
    # total. Only a sentence that scopes the claim to the whole corpus is checked, because a
    # subset count is legitimate and cannot be verified without knowing which subset is meant.
    scope = re.compile(r"\b(?:corpus|graph|catalogue|catalog|total|all|across|every|whole)\b", re.I)
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, sent in sentences(read(doc)):
            if sent.startswith("|") or not scope.search(sent):
                continue
            for m in pat.finditer(sent):
                n = int(m.group(1).replace(",", ""))
                if n < CORPUS_SCALE:
                    continue
                noun = m.group(2).lower()
                key = ("shards" if "shard" in noun
                       else "kps" if noun.startswith("kp") or "knowledge" in noun
                       else "courses" if noun.startswith("course") else "concepts")
                if n != facts[key]:
                    rep.error(f"{doc}:{line_no}", "corpus-count",
                              f'"{m.group(0)}" but the corpus has {facts[key]} {key}')


def defined_adrs() -> set:
    """Every ADR, read from docs/adr/ one file each.

    The ADRs lived inline in CONTEXT.md as one bullet per decision until the file reached
    140 KB. They are now one file per decision, so the source of truth is the directory: a file
    named ADR-NNN-*.md defines ADR-NNN. Reading the directory rather than a single file is what
    lets the check survive the split.
    """
    out = set()
    for f in sorted((ROOT / "docs/adr").glob("ADR-*.md")):
        m = re.match(r"(ADR-\d{3})-", f.name)
        if m:
            out.add(m.group(1))
    return out


def check_adr_refs(rep: Report) -> None:
    """Every ADR a doc cites must be defined in docs/adr/."""
    defined = defined_adrs()
    if not defined:
        rep.error("docs/adr", "adr", "no ADR files found; the pattern changed?")
        return
    for doc in DOCS:
        if doc == "CONTEXT.md":
            continue
        for i, line in enumerate(read(doc).split("\n"), 1):
            for m in re.finditer(r"\bADR-\d{3}\b", line):
                if m.group(0) not in defined:
                    rep.error(f"{doc}:{i}", "adr-ref",
                              f"{m.group(0)} is cited but not defined in docs/adr/")


def check_links(rep: Report) -> None:
    """Relative links must resolve, as a file or as an app route."""
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        text = read(doc)
        for i, line in enumerate(text.split("\n"), 1):
            for m in re.finditer(r"\[[^\]]*\]\(([^)\s]+)\)", line):
                tgt = m.group(1)
                if tgt.startswith(("http://", "https://", "mailto:", "#")):
                    continue
                if tgt.startswith("<"):
                    # A placeholder such as <../parabola/>, used inside audit tables.
                    continue
                if tgt.startswith("/"):
                    if any(tgt == r or tgt.startswith(r + "/") for r in APP_ROUTES):
                        continue
                    rep.error(f"{doc}:{i}", "link", f"{tgt} looks like a route but matches no "
                              f"known route, and no such file exists either")
                    continue
                target = tgt.split("#")[0]
                if not target:
                    continue
                if not (p.parent / target).resolve().exists():
                    rep.error(f"{doc}:{i}", "link", f"{tgt} does not exist")


def check_spelling(rep: Report) -> None:
    """Engineering docs use American spelling. Exempt files state why."""
    targets = [d for d in DOCS if d not in SPELLING_EXEMPT]
    targets += [f for f in SPELLING_EXEMPT if (ROOT / f).is_file()]
    for doc in targets:
        p = ROOT / doc
        if not p.is_file():
            continue
        exempt = doc in SPELLING_EXEMPT
        for i, line in enumerate(read(doc).split("\n"), 1):
            if exempt:
                continue
            for w, american in BRITISH.items():
                for m in re.finditer(rf"\b{w}\b", line, re.I):
                    if m.group(0) == american:
                        continue
                    in_code = False
                    # skip occurrences inside inline code on this line
                    for cm in re.finditer(r"`[^`\n]*`", line):
                        if cm.start() <= m.start() < cm.end():
                            in_code = True
                    if in_code:
                        continue
                    rep.error(f"{doc}:{i}", "spelling",
                              f'"{m.group(0)}" is British English; use "{american}"')


# ---------------------------------------------------------------- schema and policy facts
#
# Each rule below reads its value from the code that owns it, for the reason in the module
# docstring: a checker that pins its own copy of the number goes stale in exactly the way it
# was written to prevent. Every one of these is a fact the docs already stated wrongly at some
# point, so each has already paid for itself once.

def go_roles() -> list[str]:
    """The four role values, read from the storage type that orders them."""
    m = re.search(r"const \(\s*(RoleStudent\s+Role\s*=\s*\"(\w+)\")", read("internal/storage/store.go"))
    src = read("internal/storage/store.go")
    return re.findall(r'Role\w+\s+Role\s+=\s*"(\w+)"', src)


def check_table_count(rep: Report) -> None:
    """"Five core tables" must match the CREATE TABLE statements in the migration."""
    src = read("internal/storage/migrate.go")
    tables = set(re.findall(r"CREATE TABLE IF NOT EXISTS (\w+)", src))
    if not tables:
        rep.error("internal/storage/migrate.go", "tables",
                  "no CREATE TABLE found; this checker cannot verify the docs")
        return
    pat = re.compile(r"\b(five|six|seven|eight|\d+)\s+(?:core\s+)?tables?\b", re.I)
    words = {"five": 5, "six": 6, "seven": 7, "eight": 8}
    for doc in DOCS:
        for line_no, sent in sentences(read(doc)):
            for m in pat.finditer(sent):
                got = words.get(m.group(1).lower(), None)
                if got is None:
                    got = int(m.group(1))
                if got != len(tables):
                    rep.error(f"{doc}:{line_no}", "tables",
                              f'"{m.group(0)}" but the migration creates {len(tables)}')


def check_roles(rep: Report) -> None:
    """A doc that enumerates the roles must enumerate the ones the code has."""
    roles = go_roles()
    if not roles:
        rep.error("internal/storage/store.go", "roles",
                  "no role constants found; this checker cannot verify the docs")
        return
    n = len(roles)
    # "admin|student" in a PATCH body is the two-role model the code left behind.
    pat = re.compile(r'"(admin)"\s*\|\s*"(student)"', re.I)
    for doc in DOCS:
        for line_no, line in enumerate(read(doc).split("\n"), 1):
            if pat.search(line):
                rep.error(f"{doc}:{line_no}", "roles",
                          f'a two-role body, but the code has {n}: {", ".join(roles)}')
    for doc in DOCS:
        for line_no, sent in sentences(read(doc)):
            m = re.search(r"\b(?:two|three|four|five)\s+roles?\b", sent, re.I)
            if m:
                got = {"two": 2, "three": 3, "four": 4, "five": 5}[m.group(1).lower()]
                if got != n:
                    rep.error(f"{doc}:{line_no}", "roles",
                              f'"{m.group(0)}" but the code has {n} roles')


def check_daily_goal(rep: Report) -> None:
    """A default daily XP goal in the docs must equal the constant the code owns."""
    goal = go_const("internal/xp/policy.go", "DefaultDailyGoal")
    if goal is None:
        rep.error("internal/xp/policy.go", "daily-goal",
                  "DefaultDailyGoal not found; this checker cannot verify the docs")
        return
    pat = re.compile(r"daily[_ ]xp[_ ]goal[^.\n]{0,60}?default[^.\n]{0,20}?(\d+)", re.I)
    for doc in DOCS:
        for line_no, sent in sentences(read(doc)):
            for m in pat.finditer(sent):
                if int(m.group(1)) != goal:
                    rep.error(f"{doc}:{line_no}", "daily-goal",
                              f'"{m.group(0)}" but xp.DefaultDailyGoal is {goal}')


def check_removed_xp_terms(rep: Report) -> None:
    """The time and streak multipliers were removed on purpose. No doc may state them.

    internal/xp/policy.go explains the removal and why, so a doc that reintroduces the formula
    contradicts a decision that is already recorded in the code. The XP gate check cannot see
    this: it matches "N XP" near a gate, and these terms are bare coefficients in a code fence.
    """
    banned = ("timeMultiplier", "streakMultiplier", "time multiplier", "streak multiplier")
    # A doc may name a removed term while explaining the removal. The package comment in
    # internal/xp/policy.go does exactly this, and a rule that flagged it would push authors to
    # delete the only record of why the term is gone. Require the sentence to be about removal.
    about_removal = re.compile(
        r"\b(?:remov\w+|revert\w+|no longer|used to|previously|was |were |deleted|"
        r"deliberate|does not|not rewarded)\b", re.I)
    for doc in DOCS:
        # Check the sentence, not the line: these docs wrap at ~100 columns, so the removal
        # language ("the old time multiplier paid for guessing") sits on the next line from the
        # term itself. A line-scoped check flags correct prose for a formatting reason.
        for line_no, sent in sentences(read(doc)):
            for b in banned:
                if re.search(rf"\b{re.escape(b)}\b", sent, re.I) and not about_removal.search(sent):
                    rep.error(f"{doc}:{line_no}", "xp-terms",
                              f'"{b}" was removed; see the package comment in '
                              f"internal/xp/policy.go")


def check_grading_timeout(rep: Report) -> None:
    """The SymPy subprocess timeout must match the alarm the service raises."""
    m = re.search(r"signal\.alarm\((\d+)\)", read("grading/sympy_service.py"))
    if not m:
        rep.error("grading/sympy_service.py", "grading-timeout",
                  "no signal.alarm found; this checker cannot verify the docs")
        return
    secs = int(m.group(1))
    pat = re.compile(r"(\d+)[- ]second\s+(?:SIGALRM\s+)?timeout", re.I)
    for doc in DOCS:
        for line_no, sent in sentences(read(doc)):
            for hm in pat.finditer(sent):
                if int(hm.group(1)) != secs:
                    rep.error(f"{doc}:{line_no}", "grading-timeout",
                              f'"{hm.group(0)}" but sympy_service.py raises signal.alarm({secs})')


def check_diagnostic_bounds(rep: Report) -> None:
    """The diagnostic question bounds must match the constants the engine enforces.

    The in-app docs said "binary search ... approximately 20-35" while the code bounds an
    info-gain CAT at 25-45. The wrong figure was also visible to learners, because the
    onboarding briefing quoted the real one.
    """
    src = read("internal/diagnostic/cat.go")
    lo = re.search(r"minTotalQuestions\s*=\s*(\d+)", src)
    hi = re.search(r"maxTotalQuestions\s*=\s*(\d+)", src)
    if not (lo and hi):
        rep.error("internal/diagnostic/cat.go", "diagnostic-bounds",
                  "bounds not found; this checker cannot verify the docs")
        return
    lo_n, hi_n = int(lo.group(1)), int(hi.group(1))
    # A range written as "20-35", "20–35" or "20 to 35".
    rng = re.compile(r"\b(\d+)\s*(?:-|–|to)\s*(\d+)\s*questions?\b", re.I)
    ctx = re.compile(r"diagnostic|assessment|CAT|info[- ]gain|test", re.I)
    for doc in DOCS + ["web/next-app/app/docs/architecture/page.tsx",
                       "web/next-app/app/how-it-works/sections.tsx"]:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, line in enumerate(p.read_text(encoding="utf-8").split("\n"), 1):
            if not ctx.search(line):
                continue
            for m in rng.finditer(line):
                a, b = int(m.group(1)), int(m.group(2))
                if (a, b) != (lo_n, hi_n):
                    rep.error(f"{doc}:{line_no}", "diagnostic-bounds",
                              f'"{m.group(0)}" but cat.go bounds it at {lo_n}-{hi_n}')
    # "binary search" describes an algorithm the code does not use. Diagram sources are
    # included because they are documentation that regenerates into docs/diagrams/*.svg, and
    # a wrong claim there is published to every reader of the rendered diagram.
    algo_docs = DOCS + ["web/next-app/app/docs/architecture/page.tsx",
                        "web/next-app/app/how-it-works/sections.tsx",
                        "web/next-app/app/diagnose/page.tsx"]
    algo_docs += [f"web/next-app/diagrams/{p.name}" for p in sorted((ROOT / "web/next-app/diagrams").glob("*.d2"))]
    for doc in algo_docs:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, line in enumerate(p.read_text(encoding="utf-8").split("\n"), 1):
            if re.search(r"binary search", line, re.I):
                rep.error(f"{doc}:{line_no}", "diagnostic-algo",
                          '"binary search" but cat.go runs an info-gain CAT over a covering set')


def check_adr_headers(rep: Report) -> None:
    """Every ADR a source file or a doc cites must have a file in docs/adr/.

    check_adr_refs scans markdown only. ADR-050 had no header and was cited by three Go files,
    so nothing failed: the decision existed, was load-bearing, and was unrecorded as a decision.
    This walks internal/ as well.
    """
    defined = defined_adrs()
    if not defined:
        rep.error("docs/adr", "adr", "no ADR files found; the pattern changed?")
        return
    for f in sorted((ROOT / "internal").rglob("*.go")):
        if f.name.endswith("_test.go"):
            continue
        for line_no, line in enumerate(f.read_text(encoding="utf-8").split("\n"), 1):
            for m in re.finditer(r"\bADR-\d{3}\b", line):
                if m.group(0) not in defined:
                    rel = f.relative_to(ROOT)
                    rep.error(f"{rel}:{line_no}", "adr-ref",
                              f"{m.group(0)} is cited in code but not defined in docs/adr/")


def git_ignored(paths: list) -> set:
    """Return the subset of paths that git ignores.

    A citation resolves only if the file is in the repository. A path under a gitignored tree
    exists on a developer's machine and nowhere else, so a check that reads the filesystem
    passes locally and fails on a clean clone -- which is what a citation check must not do.
    CONTEXT.md cited `skills/skills/engineering/tdd/SKILL.md`, and `skills/` is ignored, so the
    check passed here and failed in CI. Consulting git closes the gap.
    """
    if not paths:
        return set()
    try:
        import subprocess
        r = subprocess.run(["git", "check-ignore", "--stdin"], cwd=ROOT,
                           input="\n".join(str(p) for p in paths),
                           capture_output=True, text=True, timeout=30)
        return {line for line in r.stdout.splitlines() if line}
    except Exception:
        # git absent or the call failed: fall back to the filesystem alone rather than
        # reporting every citation as broken.
        return set()


def file_index() -> dict:
    """Map a basename to every file with that name, walking the tree once.

    The citation check looks up files by basename, because the docs cite "engine.go" and
    "main.go" with no directory. A rglob per citation cost about eight seconds across the
    corpus; one walk costs a fraction of that.
    """
    global _FILE_INDEX
    if _FILE_INDEX is not None:
        return _FILE_INDEX
    skip = {"node_modules", ".next", ".git", "dist", "build"}
    idx: dict = {}
    for dirpath, dirnames, filenames in os.walk(ROOT):
        dirnames[:] = [d for d in dirnames if d not in skip]
        for fn in filenames:
            idx.setdefault(fn, []).append(Path(dirpath) / fn)
    _FILE_INDEX = idx
    return idx


def check_symbol_citations(rep: Report) -> None:
    """A `file.go:LINE` citation must point at a line that still exists.

    The docs cite source locations, and the locations rot: a short file grows, the line numbers
    do not move, and the citation quietly names unrelated code. No rule read them, because
    every other check in this file looks at the claim rather than the pointer.
    """
    pat = re.compile(r"`([\w./-]+\.(?:go|py|tsx|ts|mjs|md)):(\d+)(?:-(\d+))?`")
    idx = file_index()

    # First pass: resolve every citation to the file it means, then ask git once which of
    # those files are ignored. Calling git per citation cost about ten seconds; one call over
    # the whole set costs one.
    resolved: list[tuple] = []
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, line in enumerate(p.read_text(encoding="utf-8").split("\n"), 1):
            for m in pat.finditer(line):
                rel, start = m.group(1), int(m.group(2))
                # The exact cited path, then the same basename anywhere in the tree.
                cands = [ROOT / rel] + list(idx.get(rel.split("/")[-1], []))
                # A gitignored file is not in the repository, so a citation to it does not
                # resolve for anyone who clones. Filtered below so local and CI agree.
                resolved.append((doc, line_no, rel, start, m.group(3), cands))

    present = [c for _, _, _, _, _, cs in resolved for c in cs if c.is_file()]
    ignored = git_ignored(present)

    for doc, line_no, rel, start, end_s, cands in resolved:
        cands = [c for c in cands if str(c) not in ignored]
        target = next((c for c in cands if c.is_file()), None)
        if target is None:
            rep.error(f"{doc}:{line_no}", "citation", f"`{rel}` does not exist")
            continue
        total = len(target.read_text(encoding="utf-8").split("\n"))
        end = int(end_s or start)
        if end > total:
            rep.error(f"{doc}:{line_no}", "citation",
                      f"`{rel}:{start}` is past the end of the file ({total} lines)")
            continue
        # The shape test below is for code. A citation into a markdown file is checked for
        # existence only: a blank line in prose is normal, and the rule exists to catch a code
        # anchor that slid onto whitespace.
        if rel.endswith(".md"):
            continue
        # A citation that lands on a blank line, a closing brace, or an import names nothing.
        # These are the three shapes the line numbers drift into: a file grows and the anchor
        # slides onto whitespace or into a neighbouring block. It is a mechanical test, so it
        # costs nothing and needs no knowledge of the symbol.
        content = target.read_text(encoding="utf-8").split("\n")[start - 1]
        if not content.strip() or re.match(r'^\s*[\}\)\]],?\s*$', content) \
                or re.match(r'\s*(?:import|"[\w./-]+")\s', content):
            rep.error(f"{doc}:{line_no}", "citation",
                      f"`{rel}:{start}` lands on {content.strip()[:30] or 'a blank line'!r}"
                      f", which names no symbol")


def check_no_stub_claims(rep: Report) -> None:
    """A doc may not call an implemented store a stub, or promise a fallback that is absent.

    Three claims rotted together. PostgresStore was described as "a stub with all methods
    returning 'not implemented'" in the in-app docs while postgres.go carried 2,141 lines and
    no such string. The grading diagram drew a "Symbolic Go fallback" that router.go does not
    branch to, and the endpoint count said 25+ against 88 registered routes. All three are
    reads of the code, so they are reads of the code here too.
    """
    pg = read("internal/storage/postgres.go")
    if re.search(r"not implemented", pg, re.I):
        rep.error("internal/storage/postgres.go", "postgres-stub",
                  "the doc says this store is a stub and the file now says 'not implemented'")
    routes = len(re.findall(r"mux\.HandleFunc", read("internal/server/server.go")))
    if routes == 0:
        rep.error("internal/server/server.go", "routes",
                  "no mux.HandleFunc found; this checker cannot verify the docs")
        return
    ctx = re.compile(r"REST endpoints?|routes?|API exposes", re.I)
    loose = re.compile(r"\b(\d+)\+?\s+(?:REST\s+)?(?:endpoints?|routes?)\b", re.I)
    stub_ctx = re.compile(r"stub|not implemented|awaiting implementation", re.I)
    for doc in DOCS + ["web/next-app/app/docs/architecture/page.tsx",
                       "web/next-app/components/SystemDesignFlow/ArchitectureFlow.tsx",
                       "web/next-app/components/GradingFlow.tsx"]:
        p = ROOT / doc
        if not p.is_file():
            continue
        # Paragraph-scoped: these pages wrap prose across lines, so "postgres.go" and
        # "a stub awaiting implementation" routinely land on different lines. A line-scoped
        # check found one of the two claims and let the other through.
        # The claim has two parts and both must be about PostgresStore in the same sentence:
        # "the PostgresStore is a stub" / "currently a stub awaiting implementation". Matching
        # the pair anywhere in a paragraph produced a false positive in CONTEXT.md, where the
        # whole ADR block is one paragraph and an unrelated line supplies the other half.
        text = p.read_text(encoding="utf-8")
        claim = re.compile(
            r"(?:Postgres\w*|PostgreSQL\w*)\s+(?:is|are)\s+(?:\w+\s+){0,4}?"
            r"(?:a\s+)?(?:stub|not implemented|awaiting implementation)", re.I)
        for line_no, line in enumerate(text.split("\n"), 1):
            for sent in re.split(r"(?<=[.!?:])\s+", line):
                if claim.search(sent):
                    rep.error(f"{doc}:{line_no}", "postgres-stub",
                              "postgres.go is implemented and contains no 'not implemented'")
        for line_no, line in enumerate(text.split("\n"), 1):
            if ctx.search(line):
                for m in loose.finditer(line):
                    n = int(m.group(1))
                    # A "+" floor is not automatically wrong -- "25+" is true of 88. But a
                    # floor that understates by more than half reads as a total, and that is
                    # what a reader takes from it. Require the floor to be within a factor of
                    # the real count, so the number has to be maintained.
                    if "+" in m.group(0) and n * 2 < routes:
                        rep.error(f"{doc}:{line_no}", "routes",
                                  f'"{m.group(0)}" understates the {routes} routes the router '
                                  f"registers")
                    elif "+" not in m.group(0) and n != routes:
                        rep.error(f"{doc}:{line_no}", "routes",
                                  f'"{m.group(0)}" but the router registers {routes}')
    # A grading fallback the router does not branch to.
    if "symbolicGrader" not in read("internal/grader/router.go") + read("internal/grader/grader.go"):
        for doc in ["web/next-app/components/GradingFlow.tsx",
                    "web/next-app/components/SystemDesignFlow/ArchitectureFlow.tsx"]:
            p = ROOT / doc
            if not p.is_file():
                continue
            for line_no, line in enumerate(p.read_text(encoding="utf-8").split("\n"), 1):
                if re.search(r"symbolic.{0,20}fallback|fallback.{0,20}symbolic", line, re.I):
                    rep.error(f"{doc}:{line_no}", "grading-fallback",
                              "the router has no pure-Go symbolic fallback branch")


def check_toolchain(rep: Report) -> None:
    """The README's Go badge must agree with go.mod.

    The badge said 1.21+ while the same file's requirements line said 1.26+ and go.mod
    targeted 1.26.6. A badge is the first thing a prospective contributor reads, and a wrong
    one sends them to install a toolchain the module will not build under.
    """
    m = re.search(r"^go\s+([0-9.]+)", read("go.mod"), re.M)
    if not m:
        rep.error("go.mod", "go-version", "no go directive found")
        return
    floor = ".".join(m.group(1).split(".")[:2])  # 1.26.6 -> 1.26
    for doc in ("README.md", "CONTRIBUTING.md"):
        for line_no, line in enumerate(read(doc).split("\n"), 1):
            # shields.io encodes the "+" as %2B, but the badge may carry either form.
            for hm in re.finditer(r"badge/Go-([0-9.]+)(?:\+|%2B)", line):
                if hm.group(1) != floor:
                    rep.error(f"{doc}:{line_no}", "go-version",
                              f'badge says Go {hm.group(1)}+ but go.mod targets {floor}+')
            for hm in re.finditer(r"Go\s+([0-9.]+)\+", line):
                if hm.group(1) != floor:
                    rep.error(f"{doc}:{line_no}", "go-version",
                              f'"{hm.group(0)}" but go.mod targets {floor}+')


def check_sm2_intervals(rep: Report) -> None:
    """A quoted first-or-second SM-2 interval must match the code.

    README said "1 day -> 3 days -> 1 week". sm2.go uses 1 day then 6 days, and every interval
    is additionally scaled by the learner's per-topic speed, which the same sentence did not
    mention. system-design.md had the 6 right, so the two documents disagreed.
    """
    src = read("internal/scheduler/sm2.go")
    first = re.search(r"case 0:\s*\n\s*interval = int\(math\.Round\((\d+)", src)
    second = re.search(r"case 1:\s*\n\s*interval = int\(math\.Round\((\d+)", src)
    if not (first and second):
        rep.error("internal/scheduler/sm2.go", "sm2",
                  "interval constants not found; this checker cannot verify the docs")
        return
    f_n, s_n = first.group(1), second.group(1)
    pat = re.compile(r"\b1 day\s*(?:→|->|→|to)\s*(\d+)\s*days?\b", re.I)
    for doc in DOCS:
        for line_no, line in enumerate(read(doc).split("\n"), 1):
            for m in pat.finditer(line):
                if m.group(1) != s_n:
                    rep.error(f"{doc}:{line_no}", "sm2",
                              f'"1 day -> {m.group(1)} days" but sm2.go uses {s_n} for rep 1')


def check_diagram_citations(rep: Report) -> None:
    """A cited diagram must exist, including the ones cited from a Go doc comment.

    check_links only sees markdown `[text](target)`. Three docs and internal/mastery/machine.go
    cited `docs/diagrams/student-model.svg`, which does not exist -- the files are
    student-model-dark.svg and student-model-light.svg -- and no rule noticed, because the
    citation is in backticks, not a link.
    """
    pat = re.compile(r"`((?:docs|web/next-app/public)/diagrams/[\w.-]+\.svg)`")
    for doc in DOCS + [str(f.relative_to(ROOT)) for f in
                       sorted((ROOT / "internal").rglob("*.go")) if not f.name.endswith("_test.go")]:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, line in enumerate(p.read_text(encoding="utf-8").split("\n"), 1):
            for m in pat.finditer(line):
                rel = m.group(1)
                if not (ROOT / rel).is_file():
                    # The docs are written from the repo root; a path under web/next-app/public
                    # is the same file mirrored. Accept either root.
                    alt = ROOT / "web/next-app/public" / rel.split("diagrams/", 1)[1]
                    if not alt.is_file():
                        rep.error(f"{doc}:{line_no}", "diagram",
                                  f"`{rel}` does not exist")


def check_design_tokens(rep: Report) -> None:
    """Every token DESIGN.md documents must exist in globals.css.

    DESIGN.md described a hex palette with `--accent-gold` while globals.css had moved to OKLCH
    and had no such token. The stale palette had also propagated into two other docs before the
    move. Nothing ever read the design doc, so nothing noticed. This reads both.
    """
    # Font variables are declared by next/font in app/layout.tsx, not in CSS, so both files
    # are sources of a legitimate token name.
    css = read("web/next-app/app/globals.css") + read("web/next-app/app/layout.tsx")
    defined = set(re.findall(r"(--[\w-]+)\s*:", css)) | set(re.findall(r"variable:\s*'(--[\w-]+)'", css))
    if not defined:
        rep.error("web/next-app/app/globals.css", "design-tokens",
                  "no custom properties found; this checker cannot verify DESIGN.md")
        return
    doc = ROOT / "DESIGN.md"
    if not doc.is_file():
        return
    for line_no, line in enumerate(doc.read_text(encoding="utf-8").split("\n"), 1):
        # A token named in a table cell or inline code, not a CSS snippet's own declaration.
        for m in re.finditer(r"`(--[\w-]+)`", line):
            if m.group(1) not in defined:
                rep.error(f"DESIGN.md:{line_no}", "design-tokens",
                          f"`{m.group(1)}` is documented but not defined in globals.css")


def main() -> int:
    rep = Report()
    check_xp_gate(rep)
    check_kp_rule(rep)
    check_corpus_figures(rep)
    check_adr_refs(rep)
    check_links(rep)
    check_spelling(rep)
    check_table_count(rep)
    check_roles(rep)
    check_daily_goal(rep)
    check_removed_xp_terms(rep)
    check_grading_timeout(rep)
    check_diagnostic_bounds(rep)
    check_adr_headers(rep)
    check_symbol_citations(rep)
    check_no_stub_claims(rep)
    check_toolchain(rep)
    check_sm2_intervals(rep)
    check_diagram_citations(rep)
    check_design_tokens(rep)

    for where, rule, msg in rep.errors:
        print(f"ERROR   {where} [{rule}] {msg}", file=sys.stderr)
    if rep.errors:
        print(f"\nFAIL: {len(rep.errors)} documentation inconsistency/inconsistencies.",
              file=sys.stderr)
        return 1
    print("OK: docs agree with the code and the corpus "
          "(XP gate, KP rule, corpus counts, ADR refs, links, spelling, tables, roles, "
          "daily goal, removed XP terms, grading timeout, diagnostic bounds, "
          "citations, stub claims, route counts, toolchain, SM-2 intervals, "
          "diagram citations, design tokens)")
    return 0


if __name__ == "__main__":
    sys.exit(main())