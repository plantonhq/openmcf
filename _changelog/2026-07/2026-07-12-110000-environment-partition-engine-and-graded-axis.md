# Environment partitioning: deterministic engine, persistable rule, graded axis

**Date:** 2026-07-12
**Scope:** `apis/dev/planton/iac/environmentpartitionrule/v1` (new kind), `apis/dev/planton/qa/mappingevalsuite/v1` (`grade_environment_partition`), `pkg/iac/envpartition` (new engine + `awsscan` adapter), `pkg/iac/mappingeval` (partition scoring axis, baseline integration, updated pins), `e2e/aws` (messy live lane pin)

## Summary

Assigning discovered account resources to **environments** is now a
first-class, deterministic capability — and the mapping exams grade it.
Environment is part of a resource's identity and `value_from` references
are env-scoped, so a wrong split poisons every downstream reference; this
change makes the split rule-driven, replayable, attributable, and
machine-scored, never guesswork.

## The rule and the engine

- **`EnvironmentPartitionRule`** (`iac.planton.dev/v1`, new kind): the
  persistable vocabulary — authoritative tag keys, environment definitions
  (canonical name + tokens), a containment-inheritance switch, an optional
  fallback environment, and per-resource overrides (what a review gate's
  corrections become). The **untaught default rule ships as an embedded
  document of the same kind** (`pkg/iac/envpartition/default-rule.yaml`),
  so untaught and taught states differ only in data.
- **`pkg/iac/envpartition`**: applies the rule with a strict precedence
  ladder — deterministic user intent always outranks inference:

  ```text
  taught override > authoritative tag > name token > containment > unassigned
  ```

  An authoritative tag's **value IS the environment** (normalized through
  declared tokens, accepted literally otherwise), and the default
  vocabulary includes `planton.ai/environment` — the tag the modules
  themselves write — so a previously-deployed resource self-identifies on
  re-scan and a user can force any resource's environment by stamping that
  one tag. Name tokens match **whole tokens of human-authored name
  surfaces only** (opaque cloud identifiers are never tokenized — a random
  id can embed a misleading substring). Containment inheritance resolves
  multi-hop chains (NAT → subnet → VPC) to a fixpoint. Contradictory
  signals are **flagged, never resolved silently**, and a resource no tier
  reaches stays **honestly unassigned** — the review gate's business.
  Every assignment carries signal attribution (which tier, which signal).
- **`envpartition/awsscan`**: the adapter owning all AWS property-model
  knowledge (name surfaces, container properties, IGW attachments) so
  every consumer reads the same signals the same way. The engine itself is
  provider-agnostic product code; eval machinery depends on it, never the
  reverse.

## The partition scoring axis

- The mapping eval `Report` gains a **partition** axis: each proposed
  manifest's `metadata.env` against its ground-truth counterpart's, over
  the ground-truth instances that set one — unproposed instances owe
  theirs (the refs completeness rule extended). Wrong-env is the worst
  class; honest unassignment reads as missing; proposed-env-where-none is
  informational. `Perfect()` includes the axis.
- **Grading is suite-declared** (`MappingEvalSuite.grade_environment_partition`,
  default off) — an exam-fairness call: all three suites' members carry
  `metadata.env`, but network-staples' `e2e` and identity-and-egress's
  `ops` are seeding bookkeeping with zero scan-visible signal after
  fingerprint redaction — debt no proposer could ever pay. Only
  **messy-account** grades the axis: its member names deliberately carry
  recoverable env tokens (`orders-prod-*`, `orders-stg-*`).
- The **baseline** now runs the engine (untaught default rule) over every
  scan and stamps `metadata.env` only where the engine assigned — it never
  guesses an environment, exactly as it never guesses a grouping. On
  seeded exams the authoritative-tag tier is exactly what redaction
  removes, so the graded floors exercise the inference tiers; the
  tag/override tiers are pinned by direct engine unit tests.

## The floor (pinned offline AND live, both green)

Messy-account gains the partition line: **6/11 environments assigned, 0
wrong, 5 missing, 0 extra** — prod/staging recovered from name tokens
(`stg` normalizes to staging); the token-less archive bucket plus the four
unproposed uncovered-tier instances stay honestly owed (the AI mapper's
headroom on this axis). New discrimination tests prove wrong-env and
missing-env each produce their specific finding. The network-staples
PERFECT pin and the identity-and-egress floor re-ran **byte-identical**
(neither suite grades the axis; no member name carries a vocabulary
token), so only the messy live lane re-ran:
`TestMappingEval_MessyAccount` PASS (534s, create-and-destroy, teardown
leak-checked by read-only probes).

## Verification

- `go test ./pkg/iac/envpartition/...` — engine (precedence ladder,
  determinism, conflicts, fallback, transitive containment) + adapter.
- `go test ./pkg/iac/mappingeval/...` — all suites' pins incl. the new
  partition floor and discrimination tests; `./pkg/iac/importmap/...`
  conformance guard.
- `make protos` (new kind + suite flag; generated-Java compile gate),
  `make gazelle` (BUILD enrollment incl. the embedded rule),
  `go build ./pkg/...`, `make e2e-vet`, `make e2e-build`.
- LIVE `PLANTON_E2E_MAPPING_EVAL=1` messy-account lane — PASS, report ==
  pinned floor exactly.
