# GcpSecretManagerSecret Guide

The judgment this guide protects: a secret is three decisions wearing
one resource — where payloads LIVE (replication/region, permanent),
how they ROTATE (versions, never edits), and who READS them (scoped
grants, never project-wide roles). Getting any of the three wrong is
invisible until it is expensive.

## Scope is permanent

Global-with-replication versus regional is decided at create and only
changes by destroy-and-recreate — which destroys every version with it.
Data-residency regimes need the regional form from day one; everyone
else wants global with automatic replication (the omitted-replication
default) and should not pin `userManaged` replicas without a residency
reason.

## The payload is a version, and version 1 is a seed

`initialVersion.data` exists so ONE manifest yields a readable secret —
it seeds version 1 and never changes it. Rotation is additive: new
versions from tooling or pipelines, `versionAliases` re-pointing
consumers atomically. Editing the manifest's payload expecting an
update is the classic misread; the field is immutable by design.

## GCP's rotation "feature" rotates nothing

`rotation` + `topics` publishes REMINDERS to Pub/Sub on a schedule. The
subscriber — a Cloud Function, a pipeline — performs the actual
rotation and adds the version. Setting rotation without building the
subscriber creates a calendar, not a rotation.

## Grant on the secret, never the project

`roles/secretmanager.secretAccessor` at project level reads EVERY
secret in the project — the blast-radius mistake this kind's
`iamMembers` exists to prevent. Secret-scoped grants: each workload's
service account on exactly the secrets it consumes. The grants are
additive (iam_member semantics) and compose with grants made elsewhere.

## The destroy surface has three levers

`deletionProtection` (engine-side plan blocker) and `deletionPolicy:
PREVENT` (API-side refusal) guard the whole secret; `versionDestroyTtl`
(≥24h) turns version destruction into disable-then-destroy with a
restore window. Production credentials deserve all three — and note
that whole-secret DELETE destroys every version immediately regardless
of the TTL.

## Chart wiring

Feed `initialVersion.data` by valueFrom from a producing resource's
sensitive output (a generated key, an OAuth client credential) — the
payload then never transits a human. Consumers take `secret_name` for
mounts and `latest_version_name` for version pinning.
