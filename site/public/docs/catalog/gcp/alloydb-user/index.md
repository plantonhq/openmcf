---
title: "AlloyDB User"
description: "AlloyDB User deployment documentation"
icon: "package"
order: 100
componentName: "gcpalloydbuser"
---

# GCP AlloyDB User

Creates a database user on an AlloyDB cluster — classic username/password (`ALLOYDB_BUILT_IN`) or passwordless IAM authentication (`ALLOYDB_IAM_USER`).

## What Gets Created

- **AlloyDB API enablement** on the project
- **AlloyDB user** — `google_alloydb_user` on the referenced cluster

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpAlloydbUser
metadata:
  name: orders-app-user
spec:
  cluster:
    valueFrom:
      kind: GcpAlloydbCluster
      name: my-orders-cluster
      fieldPath: status.outputs.cluster_id
  userId: orders-app
  password: a-strong-generated-password
```

## Configuration Reference

| Field | Default | Description |
|-------|---------|-------------|
| `cluster` | — (required) | Cluster ref → `cluster_id`. |
| `userId` | — (required) | Database role name. |
| `userType` | `ALLOYDB_BUILT_IN` | `ALLOYDB_IAM_USER` for IAM auth. |
| `password` | — | BUILT_IN only; never for IAM users. |
| `databaseRoles` | — | e.g. `alloydbiamuser`. |

## Related Components

- [GcpAlloydbCluster](/docs/catalog/gcp/alloydb-cluster)
- [GcpServiceAccount](/docs/catalog/gcp/service-account) — identity for IAM users
