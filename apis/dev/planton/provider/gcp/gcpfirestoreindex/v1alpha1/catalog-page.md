# GCP Firestore Index

Creates a composite index on a Cloud Firestore database — the prerequisite for multi-field queries that filter or order on more than one field.

## What Gets Created

A Firestore composite index with the declared field roles (sort order, array-contains, or vector search). The Firestore API is enabled automatically.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing Firestore database** — referenced via `database` (a `GcpFirestoreDatabase` resource or literal name; empty falls back to `(default)`)
- **IAM permissions** — Firestore index admin (e.g. `roles/datastore.indexAdmin`)

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpFirestoreIndex
metadata:
  name: orders-by-customer-date
spec:
  database:
    valueFrom:
      kind: GcpFirestoreDatabase
      name: prod-firestore
      fieldPath: status.outputs.database_name
  collection: orders
  fields:
    - fieldPath: customerId
      order: ASCENDING
    - fieldPath: createdAt
      order: DESCENDING
```

```shell
planton apply -f index.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `collection` | `string` | — (required) | Collection ID. Immutable. |
| `fields` | `object[]` | — (required) | Indexed fields with order, arrayConfig, or vectorConfig. |
| `database` | `StringValueOrRef` | `(default)` | Database name. Immutable. |
| `queryScope` | `string` | `COLLECTION` | Query scope. Immutable. |
| `apiScope` | `string` | `ANY_API` | API surface. Immutable. |
| `density` | `string` | GCP default | Index density. Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project owning the database. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `index_id` | Fully qualified index resource name |
| `collection` | Collection (group) the index applies to |

## Related Components

- [GcpFirestoreDatabase](/docs/catalog/gcp/gcpfirestoredatabase) — the parent database
- [GcpProject](/docs/catalog/gcp/gcpproject) — provides the GCP project
