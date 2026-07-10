# AWS Aurora PostgreSQL

The PostgreSQL a team points production at, deployed the way it should be
from day one: Aurora Serverless v2 capacity (or provisioned instances) with
the master password managed and rotated in Secrets Manager, storage encrypted
at rest, a week of point-in-time recovery, deletion protection, Performance
Insights, and a dedicated security group that admits exactly one thing — the
PostgreSQL port from your application network.

Databases are where infrastructure shortcuts hurt the most, and where they
are hardest to fix later (encryption is create-time-only; a plaintext
password in a repo is forever). This chart makes the secure, recoverable
posture the starting point instead of the retrofit.

## Architecture

```
  applications (10.0.0.0/16)
        │  tcp/5432 only
        ▼
  ┌──────────────────────┐      ┌────────────────────────────────────┐
  │  AwsSecurityGroup    │      │  AwsRdsCluster (aurora-postgresql) │
  │  5432 from app_cidr  │─────▶│   writer  [db.serverless 0.5–16 ACU│
  │  deny-all egress     │      │            or provisioned class]   │
  └──────────────────────┘      │   reader  [optional, tier-1        │
                                │            failover + read scaling]│
   private subnets (≥2 AZs) ───▶│   storage encrypted • PITR 7d      │
                                │   final snapshot • deletion guard  │
                                └──────────────┬─────────────────────┘
                                               │ master_user_secret_arn
                                               ▼
                                     Secrets Manager (managed +
                                     rotated master password)
```

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| Database security group | `AwsSecurityGroup` | The network contract: PostgreSQL from the app range, nothing else, deny-all egress |
| Aurora cluster | `AwsRdsCluster` | The database — Serverless v2 or provisioned, managed master password, encrypted, backed up, deletion-protected |

## Parameters

| Parameter | Description | Default | Type |
|-----------|-------------|---------|------|
| `aws_region` | Region for the cluster and security group | `us-east-1` | string |
| `cluster_name` | Cluster name and companion-resource prefix | `my-app-db` | string |
| `vpc_id` | VPC for the database security group | placeholder | string |
| `db_subnet_ids` | ≥2 private subnet ids in different AZs (builds the subnet group) | placeholders | list |
| `app_cidr` | The only range allowed to reach port 5432 | `10.0.0.0/16` | string |
| `database_name` | Initial database created in the cluster | `appdb` | string |
| `master_username` | Master login (password is AWS-managed, never an input) | `dbadmin` | string |
| `engine_version` | Pinned engine version (empty = AWS default) | `16.4` | string |
| `min_capacity` | Serverless v2 ACU floor (0 = allow auto-pause) | `0.5` | number |
| `max_capacity` | Serverless v2 ACU ceiling — the spend guard | `16` | number |
| `provisioned_instance_class` | Instance class when serverless is off | `db.r6g.large` | string |
| `serverless_enabled` | Serverless v2 capacity vs provisioned instances | `true` | bool |
| `reader_enabled` | Reader instance: read scaling + fast failover (doubles instance cost) | `false` | bool |
| `deletion_protection` | Refuse cluster deletion while enabled | `true` | bool |

## Composing with network-foundation

Deployed next to the `network-foundation` chart, the placement literals
become references to its outputs — same environment, resolved by name:

```yaml
vpcId:
  valueFrom:
    kind: AwsVpc
    name: core-vpc
    fieldPath: status.outputs.vpc_id
subnetIds:
  - valueFrom:
      kind: AwsSubnet
      name: core-private-us-east-1a
      fieldPath: status.outputs.subnet_id
  - valueFrom:
      kind: AwsSubnet
      name: core-private-us-east-1b
      fieldPath: status.outputs.subnet_id
```

The chart's parameter surface uses literal ids so it also works against any
existing VPC, chart-managed or not.

## After deploying

1. **Read the connection secret**: the cluster exports
   `status.outputs.master_user_secret_arn`. Fetch the credentials with
   `aws secretsmanager get-secret-value --secret-id <arn>` — the secret JSON
   carries `username` and `password`.
2. **Connect**: writer endpoint at `status.outputs.endpoint`, reader endpoint
   (when `reader_enabled`) at `status.outputs.reader_endpoint`, port 5432.
   Applications should read the endpoint and the secret ARN from outputs
   rather than hardcoding either.
3. **Point apps at the reader for reads** where the framework supports a
   read/write split — the reader endpoint load-balances across present and
   future readers automatically.

## Day-2 guidance

- **Tighter network scoping**: `app_cidr` admits a range; the stricter shape
  is SG-to-SG — replace the ingress rule's `ipv4Cidrs` with
  `sourceSecurityGroupIds` referencing the application tier's security
  group, so membership (not addressing) grants access.
- **Scaling**: raise `max_capacity` (in-place) as load grows; move to
  provisioned instances (`serverless_enabled: false`) when utilization is
  high and steady enough that reserved pricing wins. Adding
  `reader_enabled: true` later is also in-place.
- **Watching capacity**: the `ServerlessDatabaseCapacity` CloudWatch metric
  shows the live ACU level — if it sits at `max_capacity`, the ceiling is
  throttling you; if it never leaves the floor, lower the floor.
- **Restoring**: any point in the retention window restores to a NEW cluster
  (AWS's model — never in-place). The cluster spec's restore fields
  (`restoreToPointInTime`, `snapshotIdentifier`) drive it declaratively.
- **Tearing down, deliberately**: disable `deletion_protection` first, then
  destroy — the final snapshot (`<cluster_name>-final`) is still taken, so
  even a completed teardown leaves a restore point.

---

© Planton. Licensed under [Apache-2.0](../../../LICENSE).
