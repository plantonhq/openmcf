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
  (existing, E2E-proven scenarios — suites compose proven fixtures, never
  invent parallel ones) plus the scan scope. Deploying it yields the
  **ground truth**: per instance, the manifest as authored (references
  unresolved) and the scan-visible cloud resources it owns (read from its
  IaC state, translated to scan coordinates through the provider import
  catalog's `cloud_control_type_name`).
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
  **perfect score** on the seeded suites — any drop is a harness, recipe,
  or scanner regression, never model variance — and the floor an AI mapper
  must beat before it earns the fuzzy cases. Its mappers are deliberately
  bounded to the suites' kinds; generalizing hand-written mapping code is
  exactly the infeasible work the AI proposer exists to replace.
- **`Score`**: the grader. Entirely structural — driven by the kinds' own
  proto schemas and the shared `StringValueOrRef` encoding — so it works
  for any component on any provider with zero per-kind grading code.

## The axes

| Axis | Question | Unit |
|------|----------|------|
| grouping | did each discovered resource land in the right instance? | claims vs ground-truth ownership (precision + recall) |
| spec | do the manifests reconstruct the declared settings? | recall over the leaves the ground-truth spec sets |
| refs | are dependencies wired as references, not frozen literals? | edges at the same spec location, targets translated through instance matching |
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

- **Offline** (`harness_test.go`, runs in `make test`, creds-free): the
  recorded-shape scan fixture drives the baseline, whose proposal must
  score perfect against the ground truth assembled from the real
  network-staples suite manifests — and hand-mutated proposals prove every
  axis **discriminates** (mis-grouping, literal-instead-of-ref, wrong spec
  value, duplicate claims, silent gaps, name-derivability breaks each
  produce their specific finding). An exam nothing can fail is not an exam.
- **Live** (`TestMappingEval_NetworkStaples`, opt-in via
  `PLANTON_E2E_MAPPING_EVAL=1`): the same chain against a real account —
  deploy the suite, scan blind, propose, score, assert perfect, destroy
  everything. Create-and-destroy with ambient credentials; artifacts (scan,
  proposal, report) recorded via `PLANTON_E2E_MAPPING_EVAL_ARTIFACTS`.

## Plugging in a proposer

A proposer is anything that turns a `Scan` into an `ImportMappingProposal`.
The contract is enforced at `ParseProposal` (typed-kind parsing, non-empty
claims, no dangling references); the grader consumes only the validated
form. An AI mapper slots into exactly the baseline's seat and is graded by
exactly the same machinery — that symmetry is the point.
