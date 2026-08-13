# Migration Without Certificate

The zero-downtime cutover shape: create the mapping with NO managed
certificate, publish the DNS records while the old host still serves,
then flip to `AUTOMATIC` once traffic has moved.

## What it configures

- `certificateMode: NONE` — the mapping exists and emits its
  `resource_records` output, but Cloud Run does not attempt certificate
  issuance (which would fail while DNS still points at the old host).
- `deletionPolicy: PREVENT` — a destroy fails instead of silently
  un-mapping the domain mid-migration.

## The migration sequence

1. Deploy this preset; the domain keeps serving from its old host.
2. Read the `resource_records` output and stage the records in the
   domain's zone (lower the TTL ahead of the cutover).
3. Cut DNS over to the emitted records.
4. Change `certificateMode` to `AUTOMATIC` and re-apply — the mapping is
   replaced (seconds) and Cloud Run issues the certificate against the
   now-pointing DNS.
5. Relax `deletionPolicy` when the migration is done, if you prefer
   DELETE semantics.

## Adjust before deploying

- **domain** — your real domain, verified FIRST by the deploying
  identity.
- **region** / **route** — must name the target service's region and
  resource.

## When to choose something else

A fresh domain with no existing traffic has nothing to migrate — start
from the **Custom Domain** preset and let `AUTOMATIC` issue the
certificate as soon as DNS lands.
