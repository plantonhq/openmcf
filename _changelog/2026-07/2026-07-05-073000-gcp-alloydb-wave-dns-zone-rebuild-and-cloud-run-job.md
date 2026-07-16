# GCP AlloyDB Wave, DNS Zone Rebuild, and Cloud Run Job

**Date**: July 5, 2026
**Type**: Enhancement
**Components**: API Definitions, GCP Provider, IaC Modules, E2E Framework, Presets, Documentation

## Summary

Three GCP catalog advances in one wave. The AlloyDB family is now honestly
decomposed: the cluster stands at the released-provider floor with PSC as a
first-class alternative to Private Service Access and a cross-region DR arm,
while read-pool instances (`GcpAlloydbInstance`, 639) and database users
(`GcpAlloydbUser`, 640) are new composable kinds parented by the cluster's
fully qualified resource path. `GcpDnsZone` was deep-rebuilt from a 3-field
shell to the full managed-zone surface (private zones, DNSSEC, forwarding,
peering, query logging) — and two hidden hazards were removed: inline
`records[]` that duplicated the `GcpDnsRecord` kind, and an **authoritative
project-level `roles/dns.admin` IAM binding** both engines silently created
(two zones in one project would fight over the project's member list).
`GcpCloudRunJob` (720, opening the 720–729 serverless overflow block) brings
run-to-completion batch workloads with the same container/volume/VPC
vocabulary as the Cloud Run service. All five kinds proven live on both
engines — 12/12 scenario-runs green, zero orphans.

## What Changed

### GcpAlloydbCluster (630) — depth rebuild

- **Connectivity is now PSA XOR PSC**, mirroring the provider's exact
  `ExactlyOneOf` rule: `network` (ref → `GcpVpcNetwork`) is optional under a
  reference-safe CEL, and `psc_config.psc_enabled` selects Private Service
  Connect instead.
- **Network ref format fixed against the live API**: AlloyDB rejects full
  `https://` self-link URLs ("malformed network path") — the FK now resolves
  `status.outputs.network_id` (the relative resource path), verified live.
- **Cross-region DR arm**: `cluster_type` (PRIMARY/SECONDARY) +
  `secondary_config.primary_cluster_name` with coherence CEL.
- Depth: `annotations`, `subscription_type` (TRIAL/STANDARD),
  `skip_await_major_version_upgrade`.
- **Deletion-protection posture made honest**: the released GA provider has no
  AlloyDB deletion-protection attribute, so the spec field was removed
  (recorded skip) — and the Pulumi bridged provider's beta-only flag, which
  defaults TRUE and **blocked destroy on one engine only** (found live), is
  explicitly neutralized with an explanatory comment.
- Defects fixed: no API enablement on either engine (added,
  `disable_on_destroy=false`), stale `binary: ./pulumi` in Pulumi.yaml,
  Terraform label keys diverging from Pulumi's (`planton-ai_*` now on both),
  ambient-project fallback, `query_string_length` zero-means-unset escape.
- Registry: `prerequisites: [GcpServiceNetworkingConnection]`.

### GcpAlloydbInstance (639, `gcpadbinst`) — new kind

- The read-pool/scaling node: `cluster` ref → the cluster's fully qualified
  `cluster_id` output, `instance_type` (READ_POOL default, PRIMARY/SECONDARY
  arms), machine config (cpu_count XOR machine_type), `read_pool_config`,
  availability type, database flags, query insights, client connection
  config (require_connectors + ssl_mode), **`activation_policy`** (the
  stop/start lever: ALWAYS/NEVER), public-IP arms with authorized external
  networks, PSC instance config. Managed connection pooling is deliberately
  unmodeled (not on the released provider — recorded skip).

### GcpAlloydbUser (640, `gcpadbusr`) — new kind

- First-class database users: `cluster` ref, `user_id`, `user_type`
  (ALLOYDB_BUILT_IN with `(sensitive)` password XOR passwordless
  ALLOYDB_IAM_USER, enforced by CEL), `database_roles`.

### GcpDnsZone (605) — depth rebuild

- **Removed**: inline `records[]` (the `GcpDnsRecord` kind owns records) and
  `iam_service_accounts` with both engines' hidden authoritative
  project-level `dns.admin` IAM binding — an anti-composition write mode the
  catalog's grants-as-nodes model (additive `GcpProjectIamMember`) replaces.
- **Added**: explicit optional `dns_name` (metadata.name-derived fallback
  preserved), `visibility` (public/private), `private_visibility_config`
  (VPC + GKE cluster refs), `dnssec_config` (state/key specs/non-existence),
  `forwarding_config` (IPv4 + FQDN targets; IPv6 is not on the released
  provider — recorded skip), `peering_config`, `cloud_logging_config`,
  `force_destroy`, labels. CEL: forwarding/peering/private-visibility
  require private visibility; forwarding XOR peering.
- Defects fixed: hardcoded `6.19.0` pin → `~> 6.0`; stale `object({value})`
  ref typing in variables.tf (plan-time failure with the real tfvars
  converter); `zone_id` output parity (both engines now emit the numeric
  `managed_zone_id`); ambient-project contract on both engines; hack
  manifest moved to canonical `iac/hack/`; label-key parity.
- Registry: `prerequisites: [GcpVpcNetwork]`.
- Outputs preserved extend-only — `GcpDnsRecord` and `GcpCertManagerCert`
  FK contracts (`zone_name`) untouched.

### GcpCloudRunJob (720, `cldrunjob`) — new kind

- Batch semantics on the Cloud Run v2 job resource: task `template`
  (containers with env literal-XOR-Secret-Manager refs, cloud_sql/secret/
  empty_dir/gcs/nfs volumes, vpc access connector XOR direct interfaces,
  service-account ref, CMEK ref, per-task timeout + max_retries, GPU
  node_selector), `task_count`, `parallelism` (≤ task_count CEL),
  launch stage, binary authorization, `deletion_protection` default TRUE.
  Deliberately no traffic/ingress/probes (serving concerns) and no
  `description` (absent from the released provider resource and the Pulumi
  SDK — modeling it would be silent data loss).
- Opens the 720–729 serverless overflow enum block.

### E2E framework and harness

- Five new verifiers (AlloyDB cluster/instance/user via the AlloyDB Admin
  API, DNS zone via the Cloud DNS API, Cloud Run job via Run Admin v2) with
  posture assertions; harness gains `alloydb/v1` + `dns/v1` clients.
- The AlloyDB kinds reuse the proven consumer-scoped run-scoped PSA
  prerequisite chain; `GcpAlloydbCluster` published its first
  `e2e/prerequisite.yaml` so child kinds compose against a live cluster.
- Ten new test entrypoints in `e2e/gcp`.
- **Live results (project `planton-e2e`, dual-engine, zero orphans)**:
  AlloyDB cluster 17m39s/16m25s, instance chain 26m20s/28m45s (full
  VPC → PSA → cluster → read-pool), user chain 17m10s/17m30s, DNS zone
  4/4 (public + private-VPC × both engines), Cloud Run job 28s/32s.

### Workflow-rule hardening (learn-once)

- `forge/flow/010`: a Pulumi.yaml without a `runtime:` block makes
  `pulumi up` succeed as a NO-OP (zero resources, zero outputs) — check
  deploy duration and resource count, not just exit code.
- `forge/flow/022`: YAML 1.1 bare `on`/`off`/`yes`/`no` parse as booleans —
  quote them in presets for string-typed enum fields.
- Forge orchestrator: verify which address format a consuming API accepts
  before wiring `default_kind_field_path` (compute APIs take self-links;
  AlloyDB/Service Networking require relative resource paths).
- `pkg/iac/MODULE_PARITY.md`: bridged-provider-only client-side flags (e.g.
  a beta-only deletion-protection default) can change one engine's destroy
  semantics — neutralize explicitly when the released GA schema lacks them.

## Validation

- Spec tests green for all five kinds; `pkg/outputs` conformance cases added
  for all five; `validate-outputs` dry-runs green on both module dirs.
- `secret-coverage --check` and `validate-refs --check` green.
- `tofu validate` + offline `planton tofu plan` through the real tfvars
  converter for all five kinds; release-equivalent Pulumi builds.
- Live dual-engine E2E 12/12 scenario-runs green, zero orphaned resources
  (post-run sweeps across AlloyDB clusters, networks, global addresses,
  DNS zones, and Cloud Run jobs all empty).
- Parity audits: all five kinds Fully Complete with PARITY ✅
  (`docs/audit/2026-07-04-233000.md` in each kind).
