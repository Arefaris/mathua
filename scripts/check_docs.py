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
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

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
    "docs/code-defects.md",
    "docs/efficacy.md",
    "docs/lesson-media.md",
    "docs/mastery-assessments.md",
    "docs/mastery-semantics.md",
    "docs/math-input.md",
    "docs/mobile-device-checklist.md",
    "docs/system-design.md",
    "docs/system-design-slides.md",
    "docs/wasc-uc-ag-research.md",
]

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
    "behaviour": "behavior", "behaviours": "behaviors",
    "colour": "color", "colours": "colors", "coloured": "colored", "colouring": "coloring",
    "labelled": "labeled", "labelling": "labeling",
    "licence": "license", "licences": "licenses",
    "centre": "center", "centres": "centers",
    "defence": "defense", "analyse": "analyze", "analysed": "analyzed",
    "organisation": "organization", "organisations": "organizations",
    "recognise": "recognize", "normalise": "normalize", "canonicalise": "canonicalize",
    "whilst": "while", "amongst": "among",
}


# Files allowed to cite a gate figure other than the current one, each with its reason. Both cite
# the old 150 XP value *as the example of the bug*, which is the only place that value should
# still appear. A whole-file exemption is a hole, so the list is short, explicit, and justified.
HISTORICAL_GATE_CITERS = {
    "CONTEXT.md": "ADR-020 records the 150 -> 50 rescale as history",
    "AGENTS.md": "quotes the stale 150 XP figure as the motivating example for this rule",
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
    """Any mastery-check gate figure in the docs must equal the constant in the code."""
    gate = go_const("internal/xp/policy.go", "QuizGateXP")
    if gate is None:
        rep.error("internal/xp/policy.go", "quiz-gate",
                  "QuizGateXP not found; this checker cannot verify the docs")
        return
    # "150 XP" and "50 XP" near a gate/mastery-check/quiz mention. Historical references to the
    # old value are allowed only in CONTEXT.md, which records the rescale as history.
    pat = re.compile(r"(\d+)\s*XP", re.I)
    ctx = re.compile(r"gate|mastery[- ]check|quiz", re.I)
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        for i, line in enumerate(read(doc).split("\n"), 1):
            if not ctx.search(line):
                continue
            for m in pat.finditer(line):
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


def check_adr_refs(rep: Report) -> None:
    """Every ADR a doc cites must be defined in CONTEXT.md."""
    defined = set(re.findall(r"^- (ADR-\d{3})\b", read("CONTEXT.md"), re.M))
    if not defined:
        rep.error("CONTEXT.md", "adr", "no ADR definitions found; the pattern changed?")
        return
    for doc in DOCS:
        if doc == "CONTEXT.md":
            continue
        for i, line in enumerate(read(doc).split("\n"), 1):
            for m in re.finditer(r"\bADR-\d{3}\b", line):
                if m.group(0) not in defined:
                    rep.error(f"{doc}:{i}", "adr-ref",
                              f"{m.group(0)} is cited but not defined in CONTEXT.md")


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
    for doc in DOCS:
        for line_no, line in enumerate(read(doc).split("\n"), 1):
            for b in banned:
                if re.search(rf"\b{re.escape(b)}\b", line, re.I):
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
    # "binary search" describes an algorithm the code does not use.
    for doc in DOCS + ["web/next-app/app/docs/architecture/page.tsx",
                       "web/next-app/app/how-it-works/sections.tsx",
                       "web/next-app/app/diagnose/page.tsx"]:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, line in enumerate(p.read_text(encoding="utf-8").split("\n"), 1):
            if re.search(r"binary search", line, re.I):
                rep.error(f"{doc}:{line_no}", "diagnostic-algo",
                          '"binary search" but cat.go runs an info-gain CAT over a covering set')


def check_adr_headers(rep: Report) -> None:
    """Every ADR a source file or a doc cites must have a header in CONTEXT.md.

    check_adr_refs scans markdown only. ADR-050 had no header and was cited by three Go files,
    so nothing failed: the decision existed, was load-bearing, and was unrecorded as a
    decision. This walks internal/ as well.
    """
    defined = set(re.findall(r"^- (ADR-\d{3})\b", read("CONTEXT.md"), re.M))
    if not defined:
        rep.error("CONTEXT.md", "adr", "no ADR definitions found; the pattern changed?")
        return
    for f in sorted((ROOT / "internal").rglob("*.go")):
        if f.name.endswith("_test.go"):
            continue
        for line_no, line in enumerate(f.read_text(encoding="utf-8").split("\n"), 1):
            for m in re.finditer(r"\bADR-\d{3}\b", line):
                if m.group(0) not in defined:
                    rel = f.relative_to(ROOT)
                    rep.error(f"{rel}:{line_no}", "adr-ref",
                              f"{m.group(0)} is cited in code but not defined in CONTEXT.md")


def check_symbol_citations(rep: Report) -> None:
    """A `file.go:LINE` citation must point at a line that still exists.

    The docs cite source locations, and the locations rot: a short file grows, the line numbers
    do not move, and the citation quietly names unrelated code. No rule read them, because
    every other check in this file looks at the claim rather than the pointer.
    """
    pat = re.compile(r"`([\w./-]+\.(?:go|py|tsx|ts|mjs)):(\d+)(?:-(\d+))?`")
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, line in enumerate(p.read_text(encoding="utf-8").split("\n"), 1):
            for m in pat.finditer(line):
                rel, start = m.group(1), int(m.group(2))
                # Try the cited path, then the same basename anywhere in the tree. Docs cite
                # "engine.go", "main.go" and "audit_lessons.py" without a directory.
                roots = ["internal", "cmd", "scripts", "grading", "web/next-app"]
                cands = [ROOT / rel]
                for r in roots:
                    cands += sorted((ROOT / r).rglob(rel))
                target = next((c for c in cands if c.is_file()), None)
                if target is None:
                    rep.error(f"{doc}:{line_no}", "citation",
                              f"`{rel}` does not exist")
                    continue
                total = len(target.read_text(encoding="utf-8").split("\n"))
                end = int(m.group(3) or start)
                if end > total:
                    rep.error(f"{doc}:{line_no}", "citation",
                              f"`{rel}:{start}` is past the end of the file ({total} lines)")


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

    for where, rule, msg in rep.errors:
        print(f"ERROR   {where} [{rule}] {msg}", file=sys.stderr)
    if rep.errors:
        print(f"\nFAIL: {len(rep.errors)} documentation inconsistency/inconsistencies.",
              file=sys.stderr)
        return 1
    print("OK: docs agree with the code and the corpus "
          "(XP gate, KP rule, corpus counts, ADR refs, links, spelling, tables, roles, "
          "daily goal, removed XP terms, grading timeout, diagnostic bounds, citations)")
    return 0


if __name__ == "__main__":
    sys.exit(main())