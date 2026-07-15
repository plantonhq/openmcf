# Identity-and-egress eval suite: every import-recipe kind is now examined

**Date:** 2026-07-12
**Scope:** `apis/dev/planton/provider/aws/aa_eval/suites/` (third suite + suite-owned members), `pkg/iac/mappingeval` (repeated-reference spec leaves, pinned floor test), `e2e/aws` (third live lane)

## Summary

The mapping eval harness gains its third exam — the **coverage** exam.
Where `network-staples` proves the measuring instrument and `messy-account`
measures headroom, the new `identity-and-egress` suite brings the last two
import-recipe kinds — the IAM role and the NAT gateway — under examination.
With it, **every kind that carries an import recipe appears in at least one
exam**, so a future mapping proposer is graded on identity and egress
infrastructure, not just the tiers earlier suites happened to cover. The
baseline's pinned imperfect report is this suite's floor
(`TestBaselineFloorOnIdentityAndEgress` offline + the live lane's report
artifact).

## What this exam adds that the other two do not

- **The first GLOBAL-service member.** Region scopes nothing on IAM:
  scanning `AWS::IAM::Role` lists every role in the account. The live run
  scanned 62 resources of which 47 were IAM roles — 46 of them foreign
  (service-linked roles, other projects) — and all of that noise stayed
  strictly informational, proving the universe-only honesty rules hold on
  a tier where noise is the norm. The offline fixture carries
  service-linked and foreign roles for the same proof, creds-free.
- **The first `nat_gateway` route target.** The routed subnet's route
  targets the NAT gateway; the deterministic baseline wires only
  internet-gateway targets, so both NAT-touching edges (the route target
  and the gateway's own `subnet_id`) are guaranteed headroom.
- **A second name-derived identity beside S3.** The `awsiamrole` recipe
  derives its import id `from_metadata_name`, so a proposed role's
  manifest name must BE the role name and the scorer enforces it — the
  rule stays exercised beyond buckets.

## Composition constraints (verified, recorded in the suite header)

- The NAT gateway is **private connectivity by structural necessity, not
  thrift**: a public gateway requires an Elastic IP allocation id that is
  unknowable before deploy, and `aws_eip` has no catalog scan-side
  identity — an `AwsElasticIp` member would deploy zero scan-visible
  resources, an instance no proposer could ever match. The suite's live
  lane is also the first live deploy of the NAT module's
  private-connectivity arm.
- The role's managed policy attaches by **literal ARN** for the same
  zero-claim reason (`aws_iam_policy` has no catalog entry).
- Exam-fairness probed live before authoring: the role's Cloud Control
  model carries the trust document, managed attachments, inline policies
  (as decoded JSON), path, description, and session duration — every
  member-set field is scan-reconstructable with **no enrichment needed**.

## The floor (pinned offline AND live, both green)

Grouping 5/7 correct (the NAT gateway and the role unclaimed — and
DECLARED unmapped, never silently dropped); spec 12/21 leaves (zero
mismatches; every missing leaf belongs to the two unproposed instances);
refs 2/4 (both subnet `vpc_id` edges correct; the `nat_gateway` route
target and the NAT's `subnet_id` owed); zero
misassigned/duplicate/wrong-target/unaccounted; zero name-derivability
breaks. Live lane `TestMappingEval_IdentityAndEgress` PASS (401s,
create-and-destroy, teardown verified clean by read-only probes).

## Scorer fix this suite forced (and now pins)

The role's `managed_policy_arns` is the first repeated `StringValueOrRef`
any exam manifest sets, and it panicked the spec-axis leaf walker (the
singular-ref path called `.Message()` on a list value). Repeated
references now participate exactly like their singular form at every
point (matched diff, unmatched-instance accounting, proposer-extra): the
literal arms are ONE whole-list spec leaf — mirroring scalar lists, so
element count never inflates the denominator — and the `value_from` arms
are edges, the refs axis's business.
`TestScorerComparesRepeatedRefLiterals` pins the wrong-value arm; the
staples and messy pins were re-run after the fix and held unchanged.

## Verification

- Offline: `go test ./pkg/iac/mappingeval/...` 18/18 (3 new tests);
  `go test ./pkg/iac/importmap/...` (conformance guard); `make gazelle`;
  `go build ./pkg/...`; `make e2e-vet`; `make e2e-build`.
- Live: `PLANTON_E2E_MAPPING_EVAL=1` `TestMappingEval_IdentityAndEgress`
  PASS on the e2e account (ambient credentials), artifacts recorded via
  `PLANTON_E2E_MAPPING_EVAL_ARTIFACTS`, leak check clean.
