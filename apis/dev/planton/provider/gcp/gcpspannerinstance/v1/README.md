# GcpSpannerInstance

Provisions a [Google Cloud Spanner](https://cloud.google.com/spanner) instance — the unit of compute and storage capacity that hosts one or more Spanner databases.

## What It Does

Cloud Spanner is a fully managed, globally distributed, strongly consistent relational database. This component creates and manages a Spanner **instance**, which pins a geographic topology (the instance configuration) and a capacity envelope shared by every database on it.

A Spanner instance is **not** a database. It is the container of databases — the "server" that hosts them. Databases ([GcpSpannerDatabase](../gcpspannerdatabase/v1)) and backup schedules ([GcpSpannerBackupSchedule](../gcpspannerbackupschedule/v1)) are separate composable resources that reference this instance by name. The Spanner API is enabled automatically in the target project.

## When to Use

- You need a relational database that scales horizontally without sharding your application
- Your workload requires strong consistency across regions (globally distributed)
- You need 99.999% availability for mission-critical data
- You want a managed database that handles replication, failover, and scaling automatically

## Key Configuration

### Instance Configuration (`config`)

The `config` field determines where your data is replicated:

- **Regional** (e.g., `regional-us-central1`) — all replicas in one region. Lower latency, lower cost. 99.99% SLA.
- **Multi-region** (e.g., `nam6`, `nam-eur-asia1`) — replicas across regions. Higher availability (99.999% SLA with ENTERPRISE_PLUS), higher cost, higher write latency.

This field is **immutable** — changing it requires recreating the instance and moving the data. It is the most consequential choice on the instance.

### Capacity

Exactly one of three capacity methods may be chosen (if none is set, GCP defaults to 1 node):

| Method | When to Use |
|---|---|
| `num_nodes` | Simple allocation. 1 node = ~10,000 QPS reads, 10 TB storage. Predictable workloads. |
| `processing_units` | Finer-grained. 1 node = 1000 PUs; multiples of 100 below 1000. Smaller workloads or precise sizing. |
| `autoscaling_config` | Spanner manages capacity within min/max bounds and CPU/storage targets. Variable workloads. |

Capacity changes are online operations — scaling up or down (or switching to/from autoscaling) never recreates the instance.

Autoscaling additionally supports **asymmetric per-replica options** on multi-region instances: a read-heavy replica region gets its own node range instead of sizing every region for the hottest one (requires ENTERPRISE or ENTERPRISE_PLUS edition).

### Editions

| Edition | SLA | Features |
|---|---|---|
| STANDARD | 99.99% regional | Cost-optimized, single-region feature set |
| ENTERPRISE | 99.99% regional | Granular sizing, asymmetric autoscaling, incremental backups |
| ENTERPRISE_PLUS | 99.999% multi-region | Highest availability, advanced compliance |

Edition upgrades apply in place.

### Free Instances

Set `instance_type: FREE_INSTANCE` for the billing account's one zero-cost development instance (~10 GB storage). Free instances cannot set capacity, edition, or automatic backups — the spec validates all of this pre-deploy. Upgrading to PROVISIONED works in place; there is no downgrade.

## Outputs

| Output | Description |
|---|---|
| `instance_id` | Fully qualified path (`projects/{project}/instances/{name}`) — the IAM/API handle |
| `instance_name` | Short name (referenced by GcpSpannerDatabase and GcpSpannerBackupSchedule) |
| `state` | CREATING or READY |
| `config` | The instance configuration (geographic topology) |

## Relationships

- **Depends on**: GcpProject (`project_id`; falls back to the provider's default project when omitted)
- **Referenced by**: GcpSpannerDatabase (`instance`), GcpSpannerBackupSchedule (`instance`)

## Deployment

```shell
planton apply -f spanner-instance.yaml
```

For copy-paste ready manifests, see the [presets](presets/).
