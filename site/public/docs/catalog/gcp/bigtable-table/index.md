---
title: "Bigtable Table"
description: "Bigtable Table deployment documentation"
icon: "package"
order: 100
componentName: "gcpbigtabletable"
---

# GCP Bigtable Table

Deploys a table inside a Cloud Bigtable instance — the schema-bearing unit: column families with their garbage-collection (retention) policies, pre-split keys for load distribution, change streams for CDC, automated backups, and deletion protection. Integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to the parent GcpBigtableInstance and GCP projects.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Bigtable Table** -- a table inside the referenced instance, named for clients to open with project + instance + table name
- **Column Families** -- the units of data organization and retention; each family carries its own GC policy (max age, max versions, a combined UNION/INTERSECTION policy, or a raw nested rule tree)
- **Pre-splits** -- created only when `splitKeys` is set; distributes initial write load across tablets instead of hammering one server
- **Change Streams** -- created only when `changeStreamRetention` is set; a CDC feed consumable by Dataflow, retained 1-7 days
- **Automated Backups** -- created only when `automatedBackupPolicy` is set; built-in backups on the configured frequency and retention

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project.
- **Planton Runner** -- required when using Runner-based credential delivery.

### GCP Project

- **A Bigtable instance** the table will live in. Reference a GcpBigtableInstance Cloud Resource via ValueFromRef or provide the instance's short name directly.
- **Bigtable API** (`bigtable.googleapis.com`) and **Bigtable Admin API** (`bigtableadmin.googleapis.com`) enabled in the target project.

## Deploy

### Console

Open the deployment store, find **GCP Bigtable Table**, and click **Deploy**. The creation wizard walks you through the parent instance, column families with their GC policies, pre-splits and CDC, and deletion protection. Start from the **Time-Series** preset in the [Presets](#presets) tab for the classic measurements-plus-metadata shape.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpBigtableTable
metadata:
  name: sensor-readings
  org: acme-corp
  env: prod
spec:
  instance:
    valueFrom:
      kind: GcpBigtableInstance
      name: prod-bigtable
      fieldPath: status.outputs.instance_name
  columnFamilies:
    - family: measurements
      gcPolicy:
        maxAge: 2160h
    - family: meta
      gcPolicy:
        maxVersions: 1
```

```shell
planton apply -f bigtable-table.yaml
```

This creates a table with two GC-bounded families and deletion protection on (the PROTECTED default) — retention is the lever that controls storage cost.

### InfraChart

When deploying as part of a multi-resource environment, wire the table to its instance deployed in the same InfraPipeline:

```yaml
spec:
  instance:
    valueFrom:
      kind: GcpBigtableInstance
      name: prod-bigtable
      fieldPath: status.outputs.instance_name
```

The InfraPipeline resolves the dependency graph, deploys the instance first, then provisions the table with the resolved instance name.

## Key Configuration

**Column families and GC** -- families must exist before applications can write; columns are created on write and never declared. Each family's GC policy bounds retention: `maxAge` alone, `maxVersions` alone, both combined with `mode` (UNION collects when either condition is met; INTERSECTION when both), or a raw `gcRules` JSON tree for nested policies (mutually exclusive with the typed fields). **A family without a GC policy accumulates every cell version forever** — the most common source of surprise Bigtable bills.

**Split keys** -- immutable one-way door: changing `splitKeys` later REPLACES the table and its data. Set them at creation, or manage splits operationally.

**Change streams** -- `changeStreamRetention` between 24h and 168h retains a CDC feed for Dataflow. The special value `'0'` disables streams on an existing table (distinct from empty, which never enables them).

**Deletion protection** -- an optional PROTECTED/UNPROTECTED string whose absence means PROTECTED: no client can delete the table until it is explicitly set UNPROTECTED first. Both IaC engines apply the default identically.

**Aggregate families** -- `type: intsum` (or intmin/intmax/inthll) declares server-side aggregate cells: counters incremented atomically at write time, with no read-modify-write races.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpBigtableInstance** | `instance` | `status.outputs.instance_name` |
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `table_id` | Fully qualified table path (`projects/{p}/instances/{i}/tables/{t}`) | Admin API calls, IAM bindings |
| `table_name` | Short table name client libraries open | Application connection configuration |
| `instance_name` | Short name of the parent instance | Confirming the parent without chasing references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Time-series** -- a wide measurements family with age-based retention, a small metadata family capped by versions, and key-prefix pre-splits. Start from the **Time-Series** preset.

**CDC-enabled** -- a change-stream feed for Dataflow pipelines, daily automated backups, and a combined age-plus-versions retention policy. Start from the **CDC Enabled** preset.

**Aggregate counters** -- an `intsum` family incremented atomically at write time for usage metering and leaderboards. Start from the **Aggregate Counters** preset.

## Works With

- [**GCP Bigtable Instance**](/cloud-catalog/gcp-bigtable-instance) -- the parent instance whose `instance_name` output this table references
- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project when not inherited from the provider connection
