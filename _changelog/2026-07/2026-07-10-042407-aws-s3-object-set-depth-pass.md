# AWS S3 Object Set at the Full Per-Object Provider Surface

**Date**: July 10, 2026
**Type**: Enhancement
**Components**: AWS Provider, API Definitions, Terraform Modules, Pulumi Modules, Testing Framework

## Summary

`AwsS3ObjectSet` was rebuilt from an 8-field spec to the complete
`aws_s3_object` per-object surface: HTTP presentation headers, user metadata,
storage classes, per-object encryption overrides with a first-class
`AwsKmsKey` reference, upload-integrity checksums, Object Lock retention and
legal holds, and the governance-bypass `force_destroy`. Three cross-engine
behavioral defects were fixed along the way, the hand-written Terraform
contract moved to the generator under the drift guard, and the kind gained
its first own end-to-end suite — live dual-engine runs green on both
scenarios.

## Problem Statement / Motivation

The spec covered 8 of the ~20 configurable arguments on `aws_s3_object`.
Anything beyond inline content and basic caching — a `Content-Disposition`
header, a storage-class placement, an encryption override for a single
sensitive object, a WORM-retained audit artifact — was unrepresentable.
Beyond breadth, the module pair carried real behavioral defects:

### Pain Points

- **Reordering the manifest churned objects (Pulumi)**: object resources
  were named by list index (`object-0`, `object-1`, …), so inserting an
  entry at the top replaced every object below it. Terraform keyed
  `for_each` by object key — the engines disagreed about identity.
- **The same manifest stored different Content-Types per engine**: with
  `content_type` unset, Pulumi sent the spec default
  (`application/octet-stream`) while Terraform sent nothing, letting S3
  store its own `binary/octet-stream`.
- **Legacy identity tags**: both engines still emitted `planton.org/*` keys
  instead of the settled `Name` + `planton.ai/*` convention.
- **No guards**: a hand-written `variables.tf` outside the drift guard, no
  outputs-conformance case (this is a map-output kind — the exact class the
  conformance guard exists for), and no scenarios/profile/entrypoints of its
  own (only a verifier and a fixture other kinds borrow).

## Solution / What's New

### Full per-object surface

Each entry in `objects` now models everything the provider's PutObject path
accepts for declarative content: `content_disposition`, `content_language`,
lowercase-keyed `metadata` (pattern-validated so what reads back always
matches the manifest), `website_redirect`, `storage_class` (the full
13-value set), `server_side_encryption` (AES256 / aws:kms / aws:kms:dsse /
aws:fsx), `kms_key` (a `StringValueOrRef` foreign key to
`AwsKmsKey.status.outputs.key_arn`), `bucket_key_enabled`,
`checksum_algorithm` (CRC32 / CRC32C / CRC64NVME / SHA1 / SHA256), the
Object Lock trio with pairing and RFC 3339 CELs, `force_destroy`, and a
value-set CEL on the existing `acl`.

Deliberately excluded, with reasons in the research doc: the local-file
`source` arms (meaningless in a hosted declarative manifest) and
`override_provider.default_tags` (provider-level tag plumbing).

### Honest composition over invented convenience

No set-level defaults block was added. Uniform posture belongs on the
bucket — `AwsS3Bucket` models default encryption, and S3 applies it to
every uploaded object automatically. The per-object security fields are
overrides, and the field comments teach exactly that division so neither a
human nor an agent reaches for per-object encryption when the bucket
setting is the right home.

### Cross-engine identity and behavior convergence

- Pulumi object resources are now named by their S3 key, matching both the
  object's S3 identity and Terraform's `for_each` keying — manifest
  reordering no longer touches unrelated objects.
- Both engines send the resolved `content_type` explicitly, so the stored
  Content-Type always matches the manifest regardless of engine.
- Identity tags converged to `Name` + `planton.ai/*` with user labels
  merged underneath, in both engines.

### The force_destroy lesson (live-caught)

The first live Pulumi destroy failed with S3's
`x-amz-bypass-governance-retention is only applicable to Object Lock
enabled buckets`. The argument's name suggests "purge versions on delete";
the provider's delete function shows it is ONLY the GOVERNANCE-retention
bypass header — ordinary destroys already remove all versions, and the flag
is invalid on regular buckets. The spec comment, both modules, the
scenarios, and the docs now state the real semantic, and the forge rules
gained a standing guardrail: verify a flag's semantics in the provider's
CRUD functions, never from its name — delete-time flags especially, because
no offline proof exercises the destroy path.

## Implementation Details

- **Contract**: generator-owned `variables.tf` (kind enrolled in the drift
  guard); provider floor `>= 6.1.0`, verified as the first release whose
  vendored S3 SDK (v1.82.0) accepts the FSx-flavored storage-class and
  encryption enum values the spec allows (v6.0.0's SDK v1.80.1 predates
  them).
- **Outputs**: fields renumbered contiguously; new `object_arns` map
  (`arn:aws:s3:::bucket/key` per key) for IAM-policy composition, exported
  by both engines.
- **Guards**: a `pkg/outputs` conformance case with mustPopulate for all
  three per-key maps — the keys carry slashes and dots, exercising the
  verbatim map-key routing path.
- **E2E**: first own artifacts — profile, a minimal scenario and a
  full-surface scenario (header/metadata breadth, SHA256 checksum, SSE-S3
  override, STANDARD_IA, website redirect, tag-merge precedence, and a
  `${...}` template introducer proving the tfvars escape live), dual-engine
  entrypoints. Recorded exclusions: `acl` (the shared fixture bucket's
  BucketOwnerEnforced ownership disables ACLs), Object Lock and
  `force_destroy` (need a lock-enabled bucket), SSE-KMS (proven at
  plan/preview level via the full-surface hack manifest),
  special-infrastructure storage classes.
- **Docs**: README, research doc, and catalog page rewritten to the full
  surface; two new presets (static-website-assets,
  encrypted-compliance-drop); site catalog mirror regenerated.

## Benefits

- Advanced S3 content patterns — compliance drops under WORM retention,
  per-object KMS overrides, archive-tier placements, website redirect
  stubs — are now first-class and validated at manifest time.
- Manifest edits are safe: object identity is the S3 key in both engines,
  so list order is purely cosmetic.
- The stored Content-Type is engine-independent and always matches the
  manifest.
- The kind is under every standing guard (drift, conformance, refs, secret
  coverage) and has a live-proven E2E suite.

## Impact

Purely additive for existing manifests — the Lambda-consumed zip-staging
fixture, the existing preset, and all stored manifests remain valid. The
Pulumi identity change would replace objects in pre-existing stacks; none
exist. Live dual-engine E2E: 4/4 lanes green (Pulumi full-surface 58.7s /
minimal 50.6s; Terraform 59.7s / 56.9s, each including the ~24-second
bucket-fixture chain), zero-orphan sweep clean.

## Related Work

- The S3 bucket depth pass established the bucket-level posture this kind
  now composes with (default encryption, versioning, ownership controls).
- The KMS depth pass provides the `key_arn` output the new `kms_key`
  reference consumes.
- The map-output framework fix (typed proto population for dot-flattened
  map keys) is what makes this kind's per-key outputs — and now its
  conformance case — work.

---

**Status**: ✅ Production Ready
