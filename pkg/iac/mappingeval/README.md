# mappingeval — the import-mapping examination system

Import mapping is the judgment half of bringing existing cloud
infrastructure under management: deciding which discovered cloud resources
group into which component instances, what each instance's spec says, and
where `value_from` references run between them. Whatever performs that
judgment — a deterministic mapper, an AI mapping agent — its quality must be
**measured, never hoped for**. This package is the measuring instrument.

## The model

```
seed (answer key)              blind side                    grade
──────────────────             ─────────────────────────     ─────────────────
MappingEvalSuite    deploy →   read-only scan (inventory)    Score(gt, proposal)
(known manifests)              → proposer                    grouping / spec /
= GroundTruth                  → ImportMappingProposal       refs / coverage
```

- **`MappingEvalSuite`** (`qa.planton.dev/v1`,
  `{provider}/aa_eval/suites/*.yaml`): an ordered list of fixture manifests
  plus the scan scope. A member is a manifest for a **live-proven
  component** — either an existing E2E scenario or a suite-owned fixture
  (under `suites/<name>/members/`); either way the deployment code is the
  same module the component's own E2E lane proves, and the suite's live
  lane proves the composition. Suites never invent parallel deployment
  paths. Deploying one yields the **ground truth**: per instance, the
  manifest as authored (references unresolved) and the scan-visible cloud
  resources it owns (read from its IaC state, translated to scan
  coordinates through the provider import catalog's
  `cloud_control_type_name`).
- **`inventory.Scanner`**: the blind side's only window — Cloud Control
  list/get keyed by CloudFormation type names, mirroring the platform's
  inventory capability shape, with **declared per-type enrichments**
  (S3 bucket regions, route-table routes/associations, internet-gateway
  attachments) closing Cloud Control's model gaps with typed read-only SDK
  calls. Structurally read-only: no mutating code path exists.
- **`ImportMappingProposal`** (`iac.planton.dev/v1`): the proposer's answer
  — proposed **manifests** (not a parallel encoding; validation parses each
  into its typed kind with unknown fields rejected) plus per-instance
  **claims** of the discovered resources it accounts for, and the honest
  `unmapped` remainder. This contract is the seam a future AI mapper must
  emit through; the review surfaces a human approves consume the same
  shape. Nothing is ever created from a proposal without an explicit human
  gate.
- **`baseline.Propose`**: the deterministic reference mapper. Pinned to a
  **perfect score** on network-staples and to a **specific imperfect
  score** on messy-account (see "The two exams") — any drift in either pin
  is a harness, recipe, or scanner regression, never model variance. Its
  mappers are deliberately bounded to the staples' kinds; generalizing
  hand-written mapping code is exactly the infeasible work the AI proposer
  exists to replace.
- **`Score`**: the grader. Entirely structural — driven by the kinds' own
  proto schemas and the shared `StringValueOrRef` encoding — so it works
  for any component on any provider with zero per-kind grading code.

## The two exams

| Suite | Purpose | Baseline pin |
|-------|---------|--------------|
| `network-staples` | proves the **instrument**: on a clean, well-signaled account, the whole chain (seed, scan, propose, score) works end to end | PERFECT |
| `messy-account` | measures **headroom**: an account shaped like a real company's — look-alike prod/staging networks (identical CIDRs), an uncovered tier (security group, KMS key + alias, DynamoDB, ECR), cross-service `value_from` edges into the KMS key | a specific IMPERFECT report — **the floor** |

The floor's authoritative record is the pinned offline test
(`TestBaselineFloorOnMessyAccount`) and each live run's report artifact.
**Beating the floor** means strictly higher grouping/spec/refs recall
without introducing what the baseline never produces: duplicate claims,
misassignments, wrong-target edges, unaccounted resources, or
name-derivability breaks. The baseline is honest even where it is blind —
everything it cannot map is *declared* unmapped, never silently dropped —
and that honesty is part of the bar. (There is deliberately no blended
numeric score or comparator machinery yet; it becomes worth building when a
second proposer exists.)

### Exam-fairness rule for suite members

Every spec field a member sets must be reconstructable from the
scan-visible + enriched surface. A field driving scan-invisible state (the
staples lane's SNS-policy lesson) would make the exam count a leaf NO
scan-driven proposer could ever match — permanently deflating the ceiling
instead of testing judgment. When a new member's field turns out invisible
on the live lane, either drop the field from the member or add a declared
enrichment.

### Seeding-fingerprint redaction

Planton's own modules tag everything they deploy
(`planton.ai/resource-kind`, `planton.ai/environment`, ...) — on a seeded
exam account those tags ARE the answer key. Both eval lanes (offline and
live) apply `RedactSeedFingerprints` between scan and proposer, so the
account presents exactly as a stranger's: `Name` tags and realistic user
tags stay; only the deploy machinery's identity tags leave. The redaction
is a property of the PIPELINE, not of any recorded fixture, and it is never
a general tag scrubber.

## The axes

| Axis | Question | Unit |
|------|----------|------|
| grouping | did each discovered resource land in the right instance? | claims vs ground-truth ownership (precision + recall) |
| spec | do the manifests reconstruct the declared settings? | recall over the leaves the ground-truth spec sets |
| refs | are dependencies wired as references, not frozen literals? | edges at the same spec location, targets translated through instance matching; an UNPROPOSED instance's edges stay in the denominator (skipping a resource never shrinks the ref debt) |
| coverage | is everything accounted for? | claimed / declared-unmapped / **unaccounted** (the silent gap) |

Plus one declared name rule: names are never scored (instances match by
kind + claim overlap — a blind proposer cannot know internal names), except
where a kind's import recipe derives its import id `from_metadata_name`
(S3: the manifest name IS the bucket name) — breaking that derivation
breaks the downstream zero-typing import and is flagged.

## Honesty rules (what keeps scores trustworthy on a shared account)

- **The scored universe is the ground truth only.** A real region also
  contains AWS-implicit resources and unrelated infrastructure; proposals
  about those are reported for information, never rewarded or penalized.
- **Structural invisibility is declared, not discovered.** A resource type
  with no `cloud_control_type_name` in the provider catalog cannot be seen
  by any scan, so it never enters the grouping denominator (the S3
  satellites are properties of the bucket's model, not scannable
  resources).
- **Config-only spec fields leave the spec axis** (declared
  `config_only_attributes`): they exist only in IaC configuration, so
  expecting a scan-driven proposer to produce them would penalize physics.
- **Proposer-extra spec fields are informational**: materializing a
  cloud-observed default explicitly is honest, not wrong.

## Proofs

- **Offline** (`harness_test.go` + `messy_floor_test.go`, run in
  `make test`, creds-free): recorded-shape scan fixtures drive the
  baseline, whose proposal must score perfect on the staples fixture and
  exactly the pinned floor on the messy fixture — and hand-mutated
  proposals prove every axis **discriminates** (mis-grouping,
  literal-instead-of-ref, wrong spec value, duplicate claims, silent gaps,
  dropped ref-carrying instances, name-derivability breaks each produce
  their specific finding). An exam nothing can fail is not an exam.
- **Live** (`TestMappingEval_NetworkStaples` +
  `TestMappingEval_MessyAccount`, opt-in via
  `PLANTON_E2E_MAPPING_EVAL=1`): the same chain against a real account —
  deploy the suite, scan blind, redact fingerprints, propose, score, assert
  the suite's pin (perfect / the floor), destroy everything.
  Create-and-destroy with ambient credentials; artifacts (scan, proposal,
  report per suite) recorded via `PLANTON_E2E_MAPPING_EVAL_ARTIFACTS`.

## Plugging in a proposer

A proposer is anything that turns a `Scan` into an `ImportMappingProposal`.
The contract is enforced at `ParseProposal` (typed-kind parsing, non-empty
claims, no dangling references); the grader consumes only the validated
form. An AI mapper slots into exactly the baseline's seat and is graded by
exactly the same machinery — that symmetry is the point.
