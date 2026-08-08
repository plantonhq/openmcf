# GcpGcsBucket Guide

The judgment this guide protects: bucket names are globally unique across ALL
of GCP — not your project, not your org — and deleted names can be held by
soft-deleted predecessors. Name buckets like DNS records, and never design a
flow that recreates a bucket under the same name in quick succession.

## Destroy semantics are three independent dials

- `forceDestroy` — may a permitted deletion erase contained objects?
  (default false: a non-empty bucket fails the destroy).
- `deletionPolicy` — is deletion permitted at all? `PREVENT` fails every
  destroy; `ABANDON` unmanages the bucket while it keeps serving.
- `retentionPolicy.isLocked` — the irreversible WORM lock: once locked, not
  even a permitted, forced destroy can remove the bucket until every object
  ages out. Validate retention against real workloads BEFORE locking.

They compose: a data-lake bucket might run `deletionPolicy: PREVENT` +
`forceDestroy: false`; an ephemeral artifacts bucket the opposite.

## Making encryption posture mandatory

`kmsKeyName` sets the default CMEK key — advisory: writers can still choose
another encryption type per object. `encryptionEnforcement` is the mandatory
half: the CMEK-only compliance shape sets `googleManagedRestrictionMode` and
`customerSuppliedRestrictionMode` to `FullyRestricted`. Enforcement applies
to NEW objects only — flipping it on a populated bucket is a policy boundary
in time, not a re-encryption; existing objects keep whatever they had.

## Conventions and gotchas

- Numeric lifecycle conditions carry explicit presence: a set `0` (e.g.
  `ageDays: 0` = every object) is real and distinct from unset — the module
  translates it through the provider's send-zero flags on both engines. The
  size-band conditions (`sizeAboveBytes`/`sizeBelowBytes`) follow the same
  contract.
- `iamMembers[].member` wants the exact IAM member string; a
  `GcpServiceAccount` reference resolves its `member` output, which is
  already that string — no assembly.
- Autoclass and `SetStorageClass` lifecycle rules both manage storage
  classes; the spec rejects combining them. Choose one transition mechanism
  per bucket.

## On the diagram

The bucket is a hub for edges: `iamMembers` references render grant edges
from service-account nodes, `logging.logBucket` renders a bucket-to-bucket
edge, and `ipFilter.vpcNetworkSources[].network` renders edges to VPC nodes.
Grants made outside the manifest (console, gcloud) render as nothing — the
additive `iamMembers` list is what keeps the diagram honest.

## Pairs well with

- `GcpServiceAccount` — the `member` output feeds `iamMembers` directly.
- `GcpKmsKey` — the CMEK key behind `kmsKeyName`, paired with
  `encryptionEnforcement` for the mandatory posture.
- `GcpBackendBucket` — front this bucket with the L7 load-balancer family
  for production HTTPS static sites.
