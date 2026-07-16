# GCP Redis Instance + Router NAT Depth Rebuild

**Date**: July 4, 2026
**Type**: Feature
**Components**: API Definitions, GCP Provider, IAC Modules, E2E Framework, InfraCharts

## Summary

Deep-rebuilt `GcpRedisInstance` (631) and `GcpRouterNat` (612) to the released v6.50.0 provider floor with full dual-engine parity, live dual-engine E2E on both kinds, RouterNat address-reference fix (no inline address creation), Redis `deletion_protection` skip with 7.x return path documented, and the gke-environment chart RouterNat template fix.

## Problem Statement / Motivation

`GcpRedisInstance` carried stale 80/20 coverage, missing provider-floor fields, Pulumi-only `deletion_protection` (false parity on `~> 6.0`), and no live E2E proof for direct peering or private services access.

`GcpRouterNat` was a four-field stub that created addresses inline (black-box defect), missing the full NAT surface (private NAT, rules, drain IPs, port tuning, timeouts, endpoint types), and the gke-environment chart template omitted required `projectId`/`routerName`/`natName`.

## Solution / What's New

### `GcpRedisInstance` deep-rebuild (631)

- Added `alternative_location_id`, `secondary_ip_range`, `maintenance_version`, `rdb_snapshot_start_time`, maintenance `minute`; removed `deletion_protection` (7.x-only — documented skip).
- Extended outputs: `server_ca_certs`, `persistence_iam_identity`, `effective_reserved_ip_range`, `instance_name`, `region`.
- Both engines: `redis.googleapis.com` API enablement, `google ~> 6.0`, ambient-project fallback, plain-string ref typing.
- Registry `prerequisites: [GcpServiceNetworkingConnection]`.
- Three presets (basic cache, HA production, private services access).
- E2E: `basic-direct-peering` + `private-service-access` scenarios; Redis verifier in `aa_e2e`.

### `GcpRouterNat` deep-rebuild (612)

- Full NAT surface (~25 spec fields): private NAT + rules, drain IPs, port tuning with power-of-two CEL, five timeouts, endpoint types, subnetwork scoping, minimal router arm.
- Address creation removed — references `GcpAddress` via `nat_ip_names` (`default_kind = GcpAddress`).
- Outputs contract: `nat_ips` self links (manual-only; honestly empty for AUTO).
- Registry `prerequisites: [GcpVpc, GcpAddress]`.
- Three presets (all-subnets auto, static-ip allowlisting, private NAT).
- E2E: `auto-allocate` + `manual-static-ip` scenarios; RouterNat verifier.

### Chart fix

- `charts/gcp/gke-environment/templates/network.yaml`: RouterNat node now supplies required `projectId`, `routerName`, `natName`.

## Validation

Offline green: spec tests (51 cases each), `secret-coverage --check`, `validate-refs --check`, `validate-outputs`, `pkg/outputs` conformance (2 new cases), release-equivalent Pulumi builds, `tofu validate` + offline plan, all presets/hacks/scenarios validated.

Parity audits: both kinds **100% Fully Complete — PARITY ✅**.

Live E2E (project `planton-e2e`, zero orphans):

| Kind | Scenario | Pulumi | Terraform |
|------|----------|--------|-----------|
| Redis | basic-direct-peering | ✅ ~11m | ✅ ~12m |
| Redis | private-service-access | ✅ ~12m | ✅ ~11m |
| RouterNat | auto-allocate | ✅ ~3m18s | ✅ ~3m19s |
| RouterNat | manual-static-ip | ✅ ~3m50s | ✅ ~3m45s |

## Learn-once uplifts

- E2E README: GCP batches that include PSA prerequisite chains or Memorystore create/destroy need `-timeout=90m` or higher per run; individual PSA scenarios need ≥35m (Redis TF destroy can exceed 15m).
