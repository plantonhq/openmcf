# AwsNeptuneCluster — Architecture and Design

## Overview

AwsNeptuneCluster provisions an Amazon Neptune graph database cluster: the shared-storage cluster resource plus its compute instances, the subnet group, and an optional cluster parameter group, in a single declarative resource. Neptune serves property-graph queries (Apache TinkerPop Gremlin, openCypher) and RDF queries (SPARQL) from the same cluster.

The cluster's AWS identifier is taken from `metadata.name` — create-time immutable, so renaming means replacement.

## The Cluster/Instances Split

Neptune separates storage from compute the same way Aurora does: the **cluster** owns the storage volume (replicated six ways across three AZs), the endpoints, backups, and encryption; the **instances** are stateless compute that serve queries from that shared volume.

The spec folds instances in as a per-name list rather than a separate kind: an instance is a sub-resource of exactly one cluster and is referenced by nothing else. Both IaC modules manage each entry as its own provider resource keyed by `name` (`<cluster>-<name>`), so:

- Adding a reader is appending a list entry — an in-place update that touches nothing else.
- Removing a reader deletes exactly that instance.
- Renaming an entry replaces that one instance.

The instance with the lowest promotion tier that is available is the writer; all others serve reads and double as failover targets (promotion takes seconds because no data copy is needed).

`instances` may only be empty for headless shapes that attach compute later — a snapshot restore, a replica, or a global-cluster member (CEL-enforced).

## No Master Password — By Design

Neptune has no master username or password anywhere in its API; that is AWS's design, not an omission of this spec. Access control is:

1. **Network reachability** — security groups on port 8182, always in effect.
2. **IAM database authentication** (`iamDatabaseAuthenticationEnabled`) — requires every request to be SigV4-signed by an IAM identity. This is the production posture; policies can scope access down to the cluster's `cluster_resource_id`.

Engine-side integrations authenticate the other direction through `iamRoles`: the S3 bulk loader and Neptune ML assume roles the cluster associates (the roles own their policies — this cluster never mutates them).

## Provisioned vs Serverless

Two compute models, selected by instance class:

- **Provisioned** — instances with fixed classes (`db.r6g.large`, ...). Predictable cost for steady traversal load.
- **Neptune Serverless** — a `serverlessV2Scaling` block (NCU bounds, 1–128 on both ends) plus instances of class `db.serverless`. Each serverless instance scales independently within the bounds. The 1-NCU minimum is the idle cost floor — Neptune Serverless does not pause to zero.

CEL enforces coherence in both directions: a scaling block requires every instance to be `db.serverless`, and any `db.serverless` instance requires the scaling block.

## Restore and Replication Shapes (Create-Time)

- **`snapshotIdentifier`** — restore from a manual or automated cluster snapshot (create-time only).
- **`replicationSourceIdentifier`** — make this cluster a read replica of the source cluster ARN; promote by clearing the field.
- **`globalClusterIdentifier`** — join an existing Neptune global database; the first joiner is the global writer, later joiners are read-only secondaries.

## Major Version Upgrades

A major `engineVersion` change needs two fields set together (CEL-enforced): `allowMajorVersionUpgrade` (the deliberate-upgrade gate) and `neptuneInstanceParameterGroupName` (the instance-level parameter group AWS applies to the cluster's instances during the upgrade). Requiring them as a pair means an upgrade can never fail halfway for a missing parameter group.

## Observability

CloudWatch log exports come in two halves that must agree: the export list (`enabledCloudwatchLogsExports: [audit, slowquery]`) turns on delivery, and the matching cluster parameters (`neptune_enable_audit_log`, the slow-query threshold parameters) make the engine produce the events. The presets set both.

## Parameter Groups

Inline `parameters` produce a module-managed cluster parameter group whose family is derived from the pinned `engineVersion` ("1.4.5.1" → `neptune1.4`) — which is why inline parameters require a pinned version (CEL). Alternatively, `neptuneClusterParameterGroupName` points at an existing group. The two are mutually exclusive. Per-instance tunables go on the instance entries' `neptuneParameterGroupName`.

## Networking and Encryption

- At least two `subnetIds` in distinct AZs (the module manages the subnet group) or an existing `neptuneSubnetGroupName` (CEL-enforced either/or).
- `securityGroupIds` attach by reference; ingress rules live on the referenced `AwsSecurityGroup` nodes. Empty uses the VPC default group.
- `storageEncrypted` (recommended default true) is a create-time one-way door; `kmsKeyId` optionally replaces the AWS-managed key.
- `storageType: iopt1` opts into I/O-Optimized storage (engine 1.3+) for I/O-heavy traversal workloads.

## Infra Chart Composability

### Inputs (StringValueOrRef)

| Field | Default Reference |
|-------|-------------------|
| `subnetIds` | `AwsSubnet.status.outputs.subnet_id` |
| `securityGroupIds` | `AwsSecurityGroup.status.outputs.security_group_id` |
| `kmsKeyId` | `AwsKmsKey.status.outputs.key_arn` |
| `iamRoles` | `AwsIamRole.status.outputs.role_arn` |

### Outputs (for downstream)

| Output | Downstream Use |
|--------|---------------|
| `endpoint` + `port` | Application connection config (writes) |
| `reader_endpoint` | Read-only traversal pool |
| `cluster_resource_id` | IAM database-auth policy scoping |
| `arn` | IAM policies, metric dimensions |
| `hosted_zone_id` | Route53 alias records |

### Typical DAG Position

```
Layer 0: AwsVpc
Layer 1: AwsSubnet, AwsSecurityGroup, AwsKmsKey, AwsIamRole
Layer 2: AwsNeptuneCluster  ← this component
Layer 3: Application configs reading endpoint + IAM auth
```

## Deliberately Omitted (v1)

| Feature | Reason |
|---------|--------|
| Creating a global database | Join-as-member via `globalClusterIdentifier` is supported; provisioning the global resource itself is a separate deferred kind |
| Custom cluster endpoints (`aws_neptune_cluster_endpoint`) | Named instance subsets for specialized routing; deferred until real demand appears |
| Event subscriptions (`aws_neptune_event_subscription`) | Account-level notification plumbing, not cluster shape; deferred |
| Neptune Analytics | A separate product (`aws_neptunegraph_graph`) with its own resource model; a candidate kind of its own |

## References

- [Amazon Neptune User Guide](https://docs.aws.amazon.com/neptune/latest/userguide/intro.html)
- [Neptune Serverless capacity scaling](https://docs.aws.amazon.com/neptune/latest/userguide/neptune-serverless-capacity-scaling.html)
- [IAM database authentication](https://docs.aws.amazon.com/neptune/latest/userguide/iam-auth.html)
- [Neptune bulk loader](https://docs.aws.amazon.com/neptune/latest/userguide/bulk-load.html)
- [Neptune global databases](https://docs.aws.amazon.com/neptune/latest/userguide/neptune-global-database.html)
