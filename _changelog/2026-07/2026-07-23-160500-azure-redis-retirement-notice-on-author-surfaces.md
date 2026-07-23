# Azure catalog: the classic Cache for Redis retirement notice, surfaced where authors and browsers read

**Date**: 2026-07-23
**Scope**: `azurerediscache` (spec header + regenerated Go stub + catalog page + the propagated site page), `azureredislinkedserver` (spec header + regenerated Go stub). Comments and documentation only — zero behavior change.

## What changed

Azure has announced the retirement of classic Azure Cache for Redis in
favor of Azure Managed Redis, and ARM has begun rejecting NEW cache
creations region by region ("Azure Cache for Redis is retiring, create
Azure Managed Redis instance instead" — observed live on new Premium
creations in some regions while Basic/Standard creations elsewhere still
succeed). That fact lived only in the kind's deep technical reference —
the spec header a manifest author reads and the catalog page a browser
compares kinds on said nothing, and the spec's Premium-tier bullet even
advertised geo-replication via the exact path ARM now rejects.

- **`azurerediscache`**: the spec header and the catalog page now carry
  the retirement notice with its real nuance — existing caches keep
  running and this kind remains the right surface for MANAGING them;
  AzureManagedRedis is the kind for NEW deployments. The catalog-page
  change is propagated to the public site.
- **`azureredislinkedserver`**: geo-replication links require two PREMIUM
  caches — the class ARM now rejects creating in some regions. The spec
  header now routes NEW geo-replicated deployments to AzureManagedRedis's
  native geo-replication (AzureManagedRedisGeoReplication) and states
  that this kind links caches that already exist.

## Why

A user choosing between the classic and Managed Redis kinds today gets no
retirement signal at the decision moment — and a deploy of a new Premium
cache or geo-replication pair can fail with a service-retirement error the
catalog never warned about. The warning belongs on the surfaces where the
choice is made, worded to the observed reality rather than an overclaimed
full-service shutdown.

## Validation

- `buf lint` + `buf format --diff` clean on both edited proto directories;
  stubs regenerated with coverage verified.
- Both kinds' spec tests and Pulumi release-entrypoint builds green;
  repo-wide `make build-go` green.
- Site catalog copy run: exactly one Azure page changed (cache-for-redis).
- No module, CEL, preset, or chart changed; existing E2E results stand by
  construction.
