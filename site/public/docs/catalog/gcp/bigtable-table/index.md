---
title: "Bigtable Table"
description: "Bigtable Table deployment documentation"
icon: "package"
order: 100
componentName: "gcpbigtabletable"
---

# GCP Bigtable Table

Creates a Cloud Bigtable table with its column families and per-family garbage-collection policies — the schema-bearing unit inside a Bigtable instance: families with retention, pre-splits, change streams, automated backups, and deletion protection.

## What Gets Created

One table plus one GC-policy object per column family that declares one. Data (rows and cells) is application territory; this resource owns the structure applications write into.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A Bigtable instance** — referenced via `instance` (a `GcpBigtableInstance` resource or a literal name)
- **IAM permissions** — `bigtable.tables.create` (e.g. `roles/bigtable.admin`)

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpBigtableTable
metadata:
  name: sensor-readings
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
```

```shell
planton apply -f table.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `instance` | `StringValueOrRef` | — (required) | Parent instance's short name. Immutable. |
| `tableName` | `string` | `metadata.name` | Table name (1-50 chars). Immutable. |
| `columnFamilies` | `object[]` | `[]` | Families with optional per-family GC policies. Mutable. |
| `splitKeys` | `string[]` | — | Pre-split row keys. Immutable — changing REPLACES the table. |
| `changeStreamRetention` | `string` | off | CDC feed retention (24h-168h; `0` disables). |
| `automatedBackupPolicy` | object | off | `retentionPeriod` + `frequency`. |
| `deletionProtection` | `string` | `PROTECTED` | API-side deletion guard. |
| `rowKeySchema` | `string` | — | Structured row-key schema (Type JSON). |
| `projectId` | `StringValueOrRef` | provider default | Project owning the instance. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `table_id` | Fully qualified table resource path |
| `table_name` | Short table name (what clients open) |
| `instance_name` | The parent instance |

## Related Components

- [GcpBigtableInstance](/docs/catalog/gcp/bigtable-instance) — the parent instance
- [GcpProject](/docs/catalog/gcp/project) — provides the GCP project
