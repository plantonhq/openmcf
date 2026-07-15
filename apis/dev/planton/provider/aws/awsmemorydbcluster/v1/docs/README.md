# AwsMemorydbCluster — Design Notes

## Provider mapping

One spec folds three provider resources:

| Spec surface | Provider resource |
|---|---|
| The cluster fields | `aws_memorydb_cluster` |
| `subnet_ids` (folded arm) | `aws_memorydb_subnet_group` (module-managed) |
| `parameters` + `parameter_group_family` (folded arm) | `aws_memorydb_parameter_group` (module-managed) |

Both groups fold with bring-your-own name arms (`subnet_group_name` /
`parameter_group_name`, CEL-exclusive with the folded lists) — the settled
data-family shape: a named list of subnets or parameters owned by exactly
one cluster is configuration, not a composable node.

## Design decisions

- **`acl_name` is a required reference.** MemoryDB's only authentication
  model is ACL attachment, and the provider marks the argument Required —
  so the spec requires it rather than defaulting invisibly to
  "open-access": an unauthenticated-access grant belongs in the manifest,
  stated explicitly. The reference targets `AwsMemorydbAcl`; the built-in
  "open-access" ACL is a literal value.
- **The tls↔open-access coupling is not CEL.** AWS rejects
  `tls_enabled: false` with any ACL other than "open-access", but
  `acl_name` is a reference field and message-level CEL must not
  dereference reference sub-fields (it breaks the Java-side validator) —
  the coupling is documented on both fields and enforced by AWS at create.
- **The cluster name derives from `metadata.name`** (create-time
  immutable, max 40 characters), and the module-managed groups derive
  their names from it — everything the module owns is discoverable by one
  identifier, identically on both engines.
- **Multi-region membership is a join arm.** `multi_region_cluster_name`
  (ForceNew) joins an externally created multi-region cluster — the
  global-cluster join precedent; the multi-region cluster itself is
  deferred (below).
- **`snapshot_retention_limit` is always sent.** 0 explicitly disables
  automatic snapshots on both engines, so the two never diverge on an AWS
  default.

## Deliberately skipped / deferred (with reasons)

- **`aws_memorydb_multi_region_cluster`** — DEFER: a cross-region
  active-active topology with a server-generated name; the cluster's
  `multi_region_cluster_name` join arm composes with zero rework on
  concrete pull.
- **`aws_memorydb_snapshot`** — DEFER: an operational point-in-time
  snapshot object, not cluster shape; the cluster's own
  snapshot/restore/retention surface is fully modeled.
- **`name_prefix` arms** — the naming basis is `metadata.name` (no AWS
  proto uses random-suffix naming).
- **Standalone subnet/parameter group kinds** — folded (see above).

## Update semantics

ForceNew (replaces the cluster): `port`, `tls_enabled`, `kms_key_arn`,
`data_tiering`, `auto_minor_version_upgrade`, `network_type`,
`multi_region_cluster_name`, the subnet-group choice, and the snapshot
restore sources. In place: ACL, engine (redis→valkey only), engine
version (upgrades only), node type, shard/replica counts, parameter
group, windows, retention, SNS topic, security groups (swap-only once
set), `ip_discovery`.

## Operational notes

- Cluster create/delete are slow: AWS takes roughly 15–25 minutes to
  create even a small cluster, and deletes poll only after a 5-minute
  initial wait (the provider allows up to 120 minutes each way).
- `final_snapshot_name` is consumed only at delete time — it never
  affects the running cluster.
