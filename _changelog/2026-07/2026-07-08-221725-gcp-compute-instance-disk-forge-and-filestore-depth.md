# GCP Compute Instance Depth Rebuild, GcpComputeDisk Forge, and Filestore Floor Completion

**Date**: July 8, 2026
**Type**: Feature
**Components**: GcpComputeInstance, GcpComputeDisk (new), GcpFilestoreInstance, GCP E2E harness, pkg/outputs, MODULE_PARITY.md

## Summary

The general-purpose compute node and the shared-filesystem node now stand
at the released `google ~> 6.0` provider floor, and persistent disks are a
first-class composable kind. `GcpComputeInstance` (617) was deep-rebuilt
from a 17-field spec to the full released surface of
`google_compute_instance`; `GcpComputeDisk` (730, `gcpdisk`, opening the
730–739 compute overflow block) was forged so a VM's data disks are
referenceable nodes with their own lifecycle instead of plain strings;
and `GcpFilestoreInstance` (700) was completed to the floor — gaining the
user-labels surface it entirely lacked, the replication arm, restore-from-
backup, network modes, and honest connect-mode semantics. All three kinds
are live-proven on both engines (8/8 scenario-runs green on the test
project, zero orphans).

## What Changed

### GcpComputeInstance (617) — deep rebuild

- **Spec to the released floor** (contiguous renumbering): boot disk with
  the exactly-one-source contract (image XOR snapshot XOR a bootable
  GcpComputeDisk reference) plus CMEK, disk labels, hyperdisk tuning,
  architecture, and confidential-mode arms; attached disks as
  GcpComputeDisk references (`status.outputs.self_link`); local-SSD
  scratch disks; NIC depth (Shared-VPC `subnetwork_project`, static
  internal `network_ip` → GcpAddress, external `nat_ip` → GcpAddress,
  IPv6 access configs, `stack_type`, `nic_type`, `queue_count`, alias
  ranges); shielded + confidential VM (SEV/SEV_SNP/TDX); advanced machine
  features (SMT control, nested virt, PMU, turbo mode); GPUs; reservation
  affinity; sole-tenant node affinities; `max_run_duration` XOR
  `termination_time`; TIER_1 egress; Resource Manager tags; resource
  policies; `hostname`, `min_cpu_platform`, `can_ip_forward`,
  `enable_display`, `desired_status` (RUNNING/SUSPENDED/TERMINATED
  stop/start lever), `key_revocation_action_type`.
- **The Spot triplication removed.** The old spec said "spot" three ways
  (`preemptible`, `spot`, `scheduling.*`); the single switch is now
  `scheduling.provisioning_model: SPOT`, with both engines deriving the
  API's required legacy preemptible flag and forcing automatic restart
  off. Legacy preemptible-without-Spot is a recorded skip.
- **Three cross-engine parity breaks closed, all user-visible:** the
  Terraform module stamped legacy label keys (`resource`, `organization`,
  ...) while Pulumi stamped `planton-ai_*` — AND merged user labels last,
  letting users clobber platform attribution; Pulumi silently dropped
  `ssh_keys`; Pulumi never exported the `status` output. Both engines now
  share the byte-identical label merge (user-first, platform wins), the
  identical `ssh-keys` metadata fold, and the identical 9-output set.
- **Module hygiene:** hardcoded `google 6.19.0` → `~> 6.0`; API
  enablement on both engines; canonical Pulumi.yaml; the converter
  contract (plain resolved strings) replacing the stale `object({value})`
  tfvars shapes that dropped `valueFrom` references.
- Registry `prerequisites: [GcpVpcNetwork, GcpSubnetwork,
  GcpServiceAccount, GcpComputeDisk]`; 60+-case spec test; three presets
  (Spot dev VM, hardened web server, stateful data VM composing the
  disk); timeless docs with probe-verified recorded skips (dynamic-NIC
  fields absent from released 6.50.0; `partner_metadata` and
  `graceful_shutdown` beta-only; raw CSEK keys — secret material in
  manifests/state — with CMEK as the modeled path; PSC network
  attachments deferred to a future kind).

### GcpComputeDisk (730) — new kind

- Zonal persistent disk at the released floor: the full type surface
  (pd-* + hyperdisk family), grow-only sizing, at-most-one-source
  (image / snapshot / self-kind clone; empty disk requires size — both
  CEL-enforced), CMEK → GcpKmsKey with the confidential-compute-requires-
  CMEK rule, hyperdisk provisioned IOPS/throughput, access modes,
  `create_snapshot_before_destroy` (the last-resort recovery net for
  precious volumes), physical block size, storage pool, labels, Resource
  Manager tags.
- Outputs export the attachment composition keys (`self_link` is what
  instances consume); `type` is normalized to the plain name on BOTH
  engines (the bridged-attribute-format parity class applied at forge
  time). Regional dual-zone disks are a recorded skip: a separate
  provider resource with a materially different surface that fails the
  dual-scope fold test.
- Full anatomy: 4 protos, both modules, spec test, 3 presets, docs,
  registry entry, E2E leaf scenario + published prerequisite profile,
  `pkg/outputs` conformance case.

### GcpFilestoreInstance (700) — floor completion

- **User labels added** (the spec had NO user-labels surface at all) plus
  Resource Manager tags; `file_share.source_backup` restore arm;
  `network_config.modes` (previously hardcoded `MODE_IPV4` in both
  engines); `initial_replication` with peers as self-kind
  GcpFilestoreInstance references; `instance_name` and `project_id`
  conformed to the catalog grain (metadata.name fallback, ambient
  project); outputs +`reserved_ip_range` +`etag` (extend-only).
- **Live-found FK format defect fixed:** the network reference resolved
  the VPC's self-link URL, which the Filestore API rejects
  ("network project is: , mismatch with instance project") — the
  reference now resolves `status.outputs.network_name`, proven live on
  both engines.
- **Parity fixes:** the TF module's legacy label keys → `planton-ai_*`
  with the shared merge order; the off-grain `provider_config` variable
  removed; `file.googleapis.com` enablement on both engines; the
  Pulumi.yaml `binary:` line removed; the bridged provider's client-side
  `deletion_policy` pinned to DELETE (PARITY comment).
- Probe-verified recorded skips: `psc_config`,
  `nfs_export_options.network`, LDAP `directory_services`,
  `source_backupdr_backup`, `deletion_policy` — all absent from released
  6.50.0. Filestore backup/snapshot kinds deliberately deferred (one-shot
  point-in-time lifecycle fits continuously-reconciled IaC poorly; GCP
  has no Filestore schedule resource).

### E2E harness

- Three new verifiers: compute instance + compute disk (typed compute
  client; posture assertions include the live label-parity canary and
  machine-type/disk-type output cross-checks) and filestore (the
  REST-probe pattern — no typed `file/v1` client in the pinned API line).
- New scenarios: instance minimal (VPC → subnetwork chain + ephemeral
  external IP) and features (Spot + shielded + SA identity + a
  pre-created GcpComputeDisk attached by self-link — the disk FK proven
  live); disk leaf; filestore BASIC_HDD on the chained VPC via
  DIRECT_PEERING. Consumer-scoped subnetwork prerequisite for the
  instance; the disk publishes a prerequisite profile.
- **Live dual-engine proof on the test project: 8/8 scenario-runs green,
  zero orphans** (instance 4/4 incl. the composed five-node chain; disk
  2/2; filestore 2/2 at ~3-4.5 min creates — under the plan-only
  threshold).

### Durable guidance

- `pkg/iac/MODULE_PARITY.md` gained the optional-output export class,
  found live twice this work: `ApplyT` callbacks must match the bridged
  field's pointer-ness (a `func(string)` against a `*string` output
  panics only at deploy), and chaining lazy accessors through
  possibly-empty nested lists (`.Index(0)` on a private VM's empty
  access-config list) panics the program — export through one
  struct-slice `ApplyT` with len/nil guards degrading to `""`, mirrored
  by `try(..., "")` on Terraform.

## Impact

Stateful-VM architectures are now honestly composable: a database VM is
a `GcpComputeDisk` (surviving the VM), a `GcpComputeInstance` referencing
it, a reserved `GcpAddress` for its VIP, and a `GcpServiceAccount`
identity — every edge a resolvable reference, deployable identically on
either engine. The compute instance was the catalog's last big
first-generation spec; the filestore was the last kind with no user
labels. Both now match the catalog's converter, labeling, ambient-project,
and name-fallback grain, and the disk kind opens the 730–739 compute
block for the MIG/template family when demand pulls it.
