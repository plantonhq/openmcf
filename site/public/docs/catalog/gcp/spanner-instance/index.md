---
title: "Spanner Instance"
description: "Spanner Instance deployment documentation"
icon: "package"
order: 100
componentName: "gcpspannerinstance"
---

# GCP Spanner Instance

Deploys a Cloud Spanner instance with configurable compute capacity (fixed nodes, processing units, or autoscaling), edition selection (Standard, Enterprise, Enterprise Plus), regional or multi-region placement, and automatic backup scheduling. Integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Spanner Instance** -- a managed instance in the specified GCP project with the chosen instance configuration (regional or multi-region), edition, and compute capacity allocation
- **Compute Capacity** -- allocated via one of three modes: fixed nodes (each node = ~10,000 QPS reads), processing units (finer-grained, multiples of 100), or autoscaling with CPU and storage utilization targets
- **Automatic Backup Scheduling** -- created only when `defaultBackupScheduleType` is set to `AUTOMATIC`; GCP creates backup schedules for new databases in this instance
- **GCP Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the Spanner instance will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **Cloud Spanner API** (`spanner.googleapis.com`) enabled in the target project.

## Deploy

### Console

Open the deployment store, find **GCP Spanner Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Free Instance** preset in the [Presets](#presets) tab for a zero-cost development instance, or the **Regional Production** preset for production workloads.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpSpannerInstance
metadata:
  name: app-spanner
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  instanceName: app-spanner-prod
  config: regional-us-central1
  displayName: App Spanner Production
  numNodes: 1
  edition: ENTERPRISE
```

```shell
planton apply -f spanner-instance.yaml
```

This creates a 1-node Enterprise Spanner instance in `us-central1` with no automatic backup scheduling. Databases must be created separately using GcpSpannerDatabase.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the instance to a GCP project deployed in the same InfraPipeline:

```yaml
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: production-project
      fieldPath: status.outputs.project_id
```

The InfraPipeline resolves the dependency graph, deploys the project first, then provisions the Spanner instance with the resolved project ID.

## Key Configuration

These are the most important decisions when configuring a Spanner instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Instance configuration** -- The `config` field determines geographic placement and replication topology. Regional configs (e.g., `regional-us-central1`) provide lower latency and cost. Multi-region configs (e.g., `nam-eur-asia1`, `nam6`) provide higher availability across geographic failures. Immutable after creation.

**Capacity mode** -- Three mutually exclusive options: `numNodes` for coarse-grained allocation (1 node = ~10,000 QPS reads), `processingUnits` for finer control (multiples of 100 below 1000, multiples of 1000 above), or `autoscalingConfig` with min/max bounds and CPU/storage utilization targets. FREE_INSTANCE instances must not set any capacity fields.

**Edition** -- `STANDARD` is cost-optimized for most workloads. `ENTERPRISE` enables granular instance sizing and advanced features. `ENTERPRISE_PLUS` provides a 99.999% multi-region SLA and advanced compliance controls. Cannot be set for FREE_INSTANCE.

**Instance type** -- `PROVISIONED` (default) requires explicit capacity. `FREE_INSTANCE` provides a zero-cost instance with ~10 GB storage, limited to one per billing account. Free instances cannot configure edition, capacity, or automatic backups.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | Fully qualified instance ID (`projects/{p}/instances/{i}`) | IAM bindings, API calls, monitoring dashboards |
| `instance_name` | Short instance name | GcpSpannerDatabase and GcpSpannerBackupSchedule `instance` fields, client library connections |
| `state` | Instance state (`CREATING` or `READY`) | Deployment validation, health checks |
| `config` | The instance configuration the instance was created with (e.g. `regional-us-central1`, `nam6`) | Choosing the matching CMEK key shape (regional vs one-key-per-region) for databases and backup schedules |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Free instance** -- Zero-cost instance for development and testing with ~10 GB storage and limited throughput. One per billing account. No edition, capacity, or automatic backup configuration. Start from the **Free Instance** preset.

**Regional production** -- Single-node Enterprise instance with automatic backup scheduling. Suitable for production workloads with predictable capacity in a single region. Start from the **Regional Production** preset.

**Autoscaling production** -- Enterprise instance with autoscaling between 1 and 3 nodes, targeting 65% CPU and 80% storage utilization. Ideal for workloads with variable or unpredictable traffic patterns. Start from the **Autoscaling Production** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the Spanner instance is created