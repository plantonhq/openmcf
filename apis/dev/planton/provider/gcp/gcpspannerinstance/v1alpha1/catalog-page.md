# GCP Spanner Instance

Deploys a Google Cloud Spanner instance with configurable compute capacity via fixed nodes, processing units, or managed autoscaling (including asymmetric per-replica scaling for multi-region topologies). Supports PROVISIONED and FREE_INSTANCE types, three editions (STANDARD, ENTERPRISE, ENTERPRISE_PLUS), user labels, and automatic backup scheduling for databases created within the instance.

## What Gets Created

When you deploy a GcpSpannerInstance resource, Planton provisions:

- **Spanner Instance** — a `google_spanner_instance` with the specified configuration, capacity, and edition; the Spanner API (`spanner.googleapis.com`) is enabled automatically
- **Labels** — your labels merged with Planton's attribution labels (resource kind, name, organization, environment)

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A GCP project** (omit `projectId` to use the provider's default project)
- **Billing account** attached to the project (required even for FREE_INSTANCE — limited to one per billing account)

## Quick Start

Create a file `spanner-instance.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpSpannerInstance
metadata:
  name: orders-spanner
spec:
  config: regional-us-central1
  displayName: Orders Spanner
  numNodes: 1
```

Deploy:

```shell
planton apply -f spanner-instance.yaml
```

This creates a single-node Spanner instance in `us-central1`, named after `metadata.name`, in the provider's default project.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `config` | `string` | Instance configuration defining replication topology (e.g., `regional-us-central1`, `nam6`). Immutable. | Required |
| `displayName` | `string` | Human-readable display name. | 4-30 chars |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project ID. Can reference a GcpProject resource via `valueFrom`. |
| `instanceName` | `string` | `metadata.name` | Unique instance name. Immutable. 6-30 chars, pattern `^[a-z][-a-z0-9]*[a-z0-9]$`. |
| `labels` | `map<string,string>` | — | User labels for cost attribution; merged with Planton's platform labels. |
| `numNodes` | `int32` | — | Node count (1 node = ~10,000 QPS reads, 10 TB storage). Mutually exclusive with the other capacity fields. Mutable, online. |
| `processingUnits` | `int32` | — | Fine-grained capacity (1 node = 1000 PUs; multiples of 100 below 1000). Mutually exclusive. Mutable, online. |
| `autoscalingConfig` | `object` | — | Managed autoscaling within bounds. Mutually exclusive with fixed capacity. |
| `autoscalingConfig.autoscalingLimits` | `object` | required | Min/max in ONE unit: `minNodes`+`maxNodes` or `minProcessingUnits`+`maxProcessingUnits`. |
| `autoscalingConfig.autoscalingTargets` | `object` | GCP defaults | `highPriorityCpuUtilizationPercent` (65 regional / 45 multi-region guidance) and `storageUtilizationPercent` (80). |
| `autoscalingConfig.asymmetricAutoscalingOptions[]` | `list` | — | Per-replica node bounds for multi-region instances: `replicaLocation` + `overrides.minNodes`/`maxNodes`. Requires ENTERPRISE+. |
| `instanceType` | `string` | `PROVISIONED` | `PROVISIONED` or `FREE_INSTANCE`. Free instances cannot set capacity, edition, or AUTOMATIC backups. |
| `edition` | `string` | GCP default | `STANDARD`, `ENTERPRISE`, or `ENTERPRISE_PLUS`. Upgrades apply in place. |
| `defaultBackupScheduleType` | `string` | `NONE` | `NONE` or `AUTOMATIC` — whether new databases get a default backup schedule. |
| `forceDestroy` | `bool` | `false` | When `true`, destroy deletes all backups on the instance first; when `false`, destroy fails while any backup exists. |

## Examples

### Free Instance for Development

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpSpannerInstance
metadata:
  name: dev-spanner
spec:
  config: regional-us-central1
  displayName: Dev Spanner
  instanceType: FREE_INSTANCE
```

### Regional Production with Autoscaling

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpSpannerInstance
metadata:
  name: prod-spanner
spec:
  config: regional-us-central1
  displayName: Production Spanner
  edition: ENTERPRISE
  defaultBackupScheduleType: AUTOMATIC
  autoscalingConfig:
    autoscalingLimits:
      minNodes: 1
      maxNodes: 5
    autoscalingTargets:
      highPriorityCpuUtilizationPercent: 65
      storageUtilizationPercent: 80
```

### Multi-Region with Asymmetric Read-Region Scaling

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpSpannerInstance
metadata:
  name: global-spanner
spec:
  config: nam6
  displayName: Global Spanner
  edition: ENTERPRISE_PLUS
  defaultBackupScheduleType: AUTOMATIC
  autoscalingConfig:
    autoscalingLimits:
      minNodes: 2
      maxNodes: 8
    asymmetricAutoscalingOptions:
      - replicaLocation: us-east1
        overrides:
          minNodes: 3
          maxNodes: 12
```

### Using Foreign Key References

Reference a GcpProject managed by Planton instead of hardcoding the project ID:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpSpannerInstance
metadata:
  name: ref-spanner
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: my-project
      fieldPath: status.outputs.project_id
  config: regional-us-central1
  displayName: Referenced Spanner
  numNodes: 1
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `instance_id` | `string` | Fully qualified instance ID (`projects/{project}/instances/{name}`) |
| `instance_name` | `string` | Short instance name, referenced by GcpSpannerDatabase and GcpSpannerBackupSchedule |
| `state` | `string` | Instance state: `CREATING` or `READY` |
| `config` | `string` | Instance configuration (geographic topology) |

## Related Components

- [GcpProject](/docs/catalog/gcp/gcpproject) — provides the project where the Spanner instance is created
- [GcpSpannerDatabase](/docs/catalog/gcp/gcpspannerdatabase) — creates databases within this instance (references `instance_name`)
- [GcpSpannerBackupSchedule](/docs/catalog/gcp/gcpspannerbackupschedule) — cron-driven full/incremental backups for a database on this instance
