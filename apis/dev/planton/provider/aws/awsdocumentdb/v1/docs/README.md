# AwsDocumentDb — Architecture and Design

## Overview

AwsDocumentDb provisions an Amazon DocumentDB (with MongoDB compatibility) cluster: the shared-storage cluster resource plus its compute instances, the subnet group, and an optional cluster parameter group, in a single declarative resource. DocumentDB speaks the MongoDB 4.0/5.0 wire protocol; existing MongoDB drivers, tools, and application code connect unchanged.

The cluster's AWS identifier is taken from `metadata.name` — create-time immutable, so renaming means replacement.

## The Cluster/Instances Split

DocumentDB separates storage from compute the same way Aurora does: the **cluster** owns the storage volume (replicated six ways across three AZs), the endpoints, credentials, backups, and encryption; the **instances** are stateless compute that serve queries from that shared volume.

The spec folds instances in as a per-name list rather than a separate kind: an instance is a sub-resource of exactly one cluster and is referenced by nothing else. Both IaC modules manage each entry as its own provider resource keyed by `name` (`<cluster>-<name>`), so:

- Adding a reader is appending a list entry — an in-place update that touches nothing else.
- Removing a reader deletes exactly that instance.
- Renaming an entry replaces that one instance.

The instance with the lowest promotion tier that is available is the writer; all others serve reads from the shared volume and double as failover targets (promotion takes seconds because no data copy is needed).

`instances` may only be empty for headless shapes that attach compute later — a snapshot or point-in-time restore, or a global-cluster member (CEL-enforced).

## Provisioned vs Serverless

Two compute models, selected by instance class:

- **Provisioned** — instances with fixed classes (`db.r6g.large`, `db.t4g.medium`, ...). Predictable cost for steady load.
- **DocumentDB Serverless** — a `serverlessV2Scaling` block (DCU bounds, 0.5–256 in half-steps) plus instances of class `db.serverless`. Each serverless instance scales independently within the bounds. The minimum DCU is the idle cost floor — DocumentDB Serverless does not pause to zero.

CEL enforces coherence in both directions: a scaling block requires every instance to be `db.serverless`, and any `db.serverless` instance requires the scaling block.

One asymmetry worth knowing (the module comments carry it too): adding or modifying the serverless block on a live cluster applies in place, but **removing** it replaces the cluster — AWS cannot switch a cluster off serverless.

## Credentials

Two mutually exclusive strategies (CEL-enforced):

1. **Managed master password** (`manageMasterUserPassword: true`, the recommended default) — AWS generates the password, stores it in Secrets Manager, and rotates it on schedule. No secret ever appears in the manifest or IaC state. The secret's ARN is exported as `master_user_secret_arn`.
2. **Direct password** (`masterPassword`, sensitive) — supplied by the operator and stored in IaC state. Supported for migration paths; prefer the managed strategy.

`masterUsername` is required for a brand-new cluster (AWS has no default and rejects a blank value); clusters created from a snapshot, a point-in-time restore, or as global-cluster members inherit credentials from their source and leave it empty (CEL allows exactly these shapes).

## Restore Shapes (Create-Time)

Two mutually exclusive restore sources seed a new cluster at creation:

- **`snapshotIdentifier`** — restore from a manual or automated cluster snapshot.
- **`restoreToPointInTime`** — restore from another cluster's continuous backup: a named source cluster, a timestamp XOR latest-restorable-time, and a restore type of `full-copy` (independent storage) or `copy-on-write` (a fast clone that shares storage with the source and pays only for divergence — the prod-data staging pattern).

Both are create-time only. `globalClusterIdentifier` joins an existing DocumentDB global cluster; the first joiner is the global writer, later joiners are read-only secondaries.

## Observability

CloudWatch log exports come in two halves that must agree: the export list (`enabledCloudwatchLogsExports: [audit, profiler]`) turns on delivery, and the matching cluster parameters (`audit_logs`, `profiler` + thresholds) make the engine produce the events. The presets set both.

Performance Insights is **instance-scoped** on DocumentDB (unlike Aurora there is no cluster-level toggle) — enable it per instance entry, with an optional KMS key per instance.

## Parameter Groups

Inline `parameters` produce a module-managed cluster parameter group whose family is derived from the pinned `engineVersion` ("5.0.0" → `docdb5.0`) — which is why inline parameters require a pinned version (CEL). Alternatively, `dbClusterParameterGroupName` points at an existing group. The two are mutually exclusive.

## Networking and Encryption

- At least two `subnetIds` in distinct AZs (the module manages the subnet group) or an existing `dbSubnetGroupName` (CEL-enforced either/or).
- `securityGroupIds` attach by reference; ingress rules live on the referenced `AwsSecurityGroup` nodes. Empty uses the VPC default group.
- `storageEncrypted` (recommended default true) is a create-time one-way door; `kmsKeyId` optionally replaces the AWS-managed key.
- `networkType: DUAL` enables dual-stack IPv4+IPv6 given IPv6-capable subnets.
- `port` (default 27017) accepts 1150–65535 — DocumentDB rejects lower ports — and is create-time only.

## Infra Chart Composability

### Inputs (StringValueOrRef)

| Field | Default Reference |
|-------|-------------------|
| `subnetIds` | `AwsSubnet.status.outputs.subnet_id` |
| `securityGroupIds` | `AwsSecurityGroup.status.outputs.security_group_id` |
| `kmsKeyId` | `AwsKmsKey.status.outputs.key_arn` |
| per-instance `performanceInsightsKmsKeyId` | `AwsKmsKey.status.outputs.key_arn` |

### Outputs (for downstream)

| Output | Downstream Use |
|--------|---------------|
| `endpoint` + `port` | Application connection config (writes) |
| `reader_endpoint` | Read-only connection pool |
| `master_user_secret_arn` | Runtime credential fetch from Secrets Manager |
| `arn` / `cluster_resource_id` | IAM policies, metric dimensions, PITR sources |
| `hosted_zone_id` | Route53 alias records |

### Typical DAG Position

```
Layer 0: AwsVpc
Layer 1: AwsSubnet, AwsSecurityGroup, AwsKmsKey
Layer 2: AwsDocumentDb  ← this component
Layer 3: Application configs reading endpoint + secret ARN
```

## Deliberately Omitted (v1)

| Feature | Reason |
|---------|--------|
| DocumentDB Elastic Clusters (`aws_docdbelastic_cluster`) | A separate product surface with its own resource, sharding, and auth model; deferred until real demand appears |
| Creating a global cluster | Join-as-member via `globalClusterIdentifier` is supported; provisioning the global resource itself is a separate deferred kind |
| Event subscriptions (`aws_docdb_event_subscription`) | Account-level notification plumbing, not cluster shape; deferred |
| Terraform write-only password (`master_password_wo`) | An engine-asymmetric ephemeral that Pulumi has no equivalent for; the managed-password strategy already keeps secrets out of state |

## References

- [Amazon DocumentDB Developer Guide](https://docs.aws.amazon.com/documentdb/latest/developerguide/what-is.html)
- [DocumentDB Serverless](https://docs.aws.amazon.com/documentdb/latest/developerguide/docdb-serverless.html)
- [Managing DocumentDB users with Secrets Manager](https://docs.aws.amazon.com/documentdb/latest/developerguide/secrets-manager.html)
- [Auditing DocumentDB events](https://docs.aws.amazon.com/documentdb/latest/developerguide/event-auditing.html)
- [DocumentDB global clusters](https://docs.aws.amazon.com/documentdb/latest/developerguide/global-clusters.html)
