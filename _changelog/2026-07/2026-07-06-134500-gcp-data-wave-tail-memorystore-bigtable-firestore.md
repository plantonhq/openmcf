# GCP Data-Wave Tail: Memorystore, Bigtable, and Firestore at the Released Floor

**Date**: July 6, 2026
**Type**: Feature
**Components**: GCP Provider, API Definitions, IaC Modules, E2E Framework

## Summary

The GCP data family is complete. Three existing kinds are brought to the
released google 6.50.0 provider floor — `GcpMemorystoreInstance` (636),
`GcpBigtableInstance` (635), and `GcpFirestoreDatabase` (632) — and four new
composable kinds are forged: `GcpServiceConnectionPolicy` (715), the
Private Service Connect authorization that PSC-first managed services
require on a network before any instance can exist; `GcpBigtableTable`
(642), the many-per-instance table with per-family garbage-collection
policies folded in; and `GcpFirestoreBackupSchedule` (643) +
`GcpFirestoreIndex` (644), Firestore's managed-backup protection and
composite-index nodes. All seven kinds are live-proven on both engines,
including the composed chain in which a Memorystore instance's PSC
endpoints materialize through a service connection policy resolved from
first-class VPC and subnetwork nodes.

## Problem Statement / Motivation

Memorystore for Valkey could not actually be deployed through composition:
its PSC connectivity requires a `google_network_connectivity_service_connection_policy`
on the consumer network, and no kind modeled it. Bigtable had no table
kind (structure lived outside IaC entirely), Firestore had no
backup-schedule or index kinds, and all three existing specs trailed the
released provider surface.

### Pain Points

- **A live-blocking FK format defect**: the Memorystore spec pointed its
  network reference at the VPC's `https://` self-link, but the Service
  Connectivity API only accepts relative resource paths — every composed
  deploy would have failed at create.
- **A destroy-parity break**: the Memorystore spec's plain-bool deletion
  protection made Terraform send an explicit FALSE when a manifest omitted
  it while Pulumi left the provider's TRUE default — the same manifest
  destroyed differently per engine.
- **A live label-parity break**: the Bigtable and Memorystore Terraform
  modules stamped legacy unprefixed label keys while their Pulumi modules
  used the shared `planton-ai_*` set.
- **Converter-contract violation**: the Memorystore Terraform module typed
  reference fields `object({ value = string })`, failing at plan against
  the real tfvars converter.

## What Changed

### GcpServiceConnectionPolicy (715, gcpscp) — forged

Location + network + service class (one policy per triple, all
immutable), `psc_config` with subnetworks / connection limit / producer
resource-hierarchy allowlist, labels, `policy_name` → `metadata.name`
fallback. Both engines normalize `https://` self-links to the relative
resource paths the API requires. Registry
`prerequisites: [GcpVpcNetwork, GcpSubnetwork]`; 3 presets
(memorystore-valkey / shared-vpc-guarded / producer-allowlist).

### GcpMemorystoreInstance (636) — depth rebuild

Network FK repointed to `network_id` (relative path); deletion protection
remodeled `optional bool` default TRUE, explicit on both engines; user
labels beneath `planton-ai_*` platform labels; cross-region DR
(`cross_instance_replication_config` with PRIMARY/SECONDARY coherence
CEL); seed-from-GCS XOR seed-from-managed-backup; per-entry PSC endpoint
project optional with ambient resolution (gated data-source on TF,
only-when-needed GetClientConfig on Pulumi); outputs extend-only
(`name`, `backup_collection`); registry
`prerequisites: [GcpServiceConnectionPolicy]`. Recorded skips
(schema-probe verified): `server_ca_mode`/`server_ca_pool`,
`maintenance_version`, `allow_fewer_zones_deployment` (absent from the
pinned Pulumi SDK — one-engine field), expanded node types (main-only).

### GcpBigtableInstance (635) — labels + conformance

User `labels` added (the released resource supports them; the spec had
none); the Terraform module's legacy unprefixed label keys fixed to the
shared `planton-ai_*` set — closing the live label-parity break; Bigtable
Admin API enablement both engines; canonical Pulumi.yaml. Recorded
skips: `edition`, `tags` (main-only), `instance_type` (deprecated).

### GcpBigtableTable (642, gcpbttbl) — forged

Column families with aggregate type shortcuts, **per-family GC policies
folded in** (max-age/max-versions/mode/raw JSON under mutual-exclusion
CEL — the provider's one-per-family `google_bigtable_gc_policy` has no
independent lifecycle, so it fails the split test), split keys,
change-stream retention, automated backup policy, `deletion_protection`
default PROTECTED explicit on both engines (the API-side guard). Registry
`prerequisites: [GcpBigtableInstance]`. Recorded skips:
authorized/logical/materialized views + schema bundle (Tier-2).

### GcpFirestoreDatabase (632) — depth rebuild

`app_engine_integration_mode` added; outputs extended
(`version_retention_period`, `key_prefix`, `update_time`); Firestore API
enablement both engines; `deletion_policy` hardcoded DELETE identically.
Recorded skips (schema-probe verified): the ENTERPRISE access modes and
`user_creds` (absent from released 6.50.0), resource tags, documents and
single-field/TTL overrides.

### GcpFirestoreBackupSchedule (643, gcpfstbs) — forged

Database ref, retention ≤ 14 weeks (the only mutable field), daily XOR
weekly-day recurrence under exactly-one CEL; the daily-plus-weekly
pattern composes as two resources; backups-outlive-the-schedule
semantics taught in spec and docs.

### GcpFirestoreIndex (644, gcpfstidx) — forged

Collection + query/api scope + density + fields with exactly-one-role
CEL per field (order XOR array-contains XOR vector with dimension — the
vector arm enables nearest-neighbor queries); fully immutable with
create-before-destroy rebuild semantics. Recorded skips: `search_config`,
`multikey`, `unique` (absent from released 6.50.0).

### E2E framework and harness

- Network Connectivity, Bigtable Admin, and Firestore Admin typed clients
  added to the GCP harness; an ADC-authenticated plain REST client added
  for services whose typed Go client is not in the pinned
  `google.golang.org/api` line (Memorystore for Valkey) — verifiers probe
  the documented REST GET path.
- Seven verifiers with posture assertions; scenarios and prerequisite
  profiles for all seven kinds; 14 test entrypoints.
- Two live-found lessons folded into `e2e/README.md`: first-ever API
  activation on a project outruns in-module enablement (pre-enable new
  service APIs on the test project once), and reservation-window resource
  classes (Firestore holds deleted database IDs for minutes) need
  consumer-unique prerequisite names, not just run-unique ones.
- Terraform-module rule: ambient-project data sources must be gated with
  `count` on actual need, or offline plans break against sample projects.

## Validation

- Offline: spec tests 28+68+52+32+50+17+25; release-equivalent Pulumi
  builds ×7; `tofu validate` ×7; offline `planton tofu plan` through the
  real tfvars converter ×7 (fields inspected); `secret-coverage --check`;
  `validate-refs --check`; `validate-outputs` on BOTH module dirs ×7;
  seven `pkg/outputs` conformance cases; every preset/hack/scenario/
  prerequisite manifest through `planton validate` on a freshly built CLI
  (run-id tokens substituted); framework tests green.
- Live (project `planton-e2e`, dual-engine, zero orphans): service
  connection policy 2m48s/4m21s (VPC → subnetwork chain); Memorystore
  instance 22m24s/[terraform] through the composed SCP chain; Bigtable
  instance minimal + autoscaling 50s/2m55s; Bigtable table chain
  1m07s/2m12s; Firestore database 25s/51s; backup schedule chain
  41s/1m21s; composite index chain 6m39s/7m12s (the async index build
  dominates).
- Audits: all seven kinds **Fully Complete — PARITY ✅**
  (`docs/audit/2026-07-06-131500.md` per kind); site catalog regenerated
  (4 pages created, 3 refreshed).
- Recorded live exclusions (proven offline): Memorystore DR pair and
  CMEK, seed-from-backup/GCS, Bigtable CMEK, Firestore CMEK.

## Impact

Private in-memory caching (Valkey), wide-column storage (Bigtable with
IaC-owned tables and retention), and serverless documents (Firestore with
reviewable indexes and managed backups) are now fully composable GCP
architectures — every edge a first-class reference, both engines
behaviorally identical.
