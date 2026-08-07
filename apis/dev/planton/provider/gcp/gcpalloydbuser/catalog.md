# GCP AlloyDB User

Deploys a database user (`google_alloydb_user`) on an AlloyDB cluster. Users are first-class nodes: one per application with ALLOYDB_BUILT_IN (password) or ALLOYDB_IAM_USER (passwordless IAM database authentication).

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **AlloyDB User** — database role on the cluster with the chosen authentication type
- **Database Roles** — PostgreSQL roles granted to the user (e.g. `alloydbiamuser`)

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** — credentials for the target project
- **Planton Runner** — when using Runner-based credential delivery

### GCP Prerequisites

- **GcpAlloydbCluster** — the cluster this user lives on. Reference `status.outputs.cluster_id`.
- **Password secret** — for ALLOYDB_BUILT_IN, store the password as an org secret reference (`$secret/slug`) or inline for dev only.

## Deploy

### Console

Open the deployment store, find **GCP AlloyDB User**, and click **Deploy**. Start from the **Application User** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpAlloydbUser
metadata:
  name: orders-app
spec:
  cluster:
    valueFrom:
      kind: GcpAlloydbCluster
      name: orders-db
      fieldPath: status.outputs.cluster_id
  userId: orders-app
  password: $secret/orders-app-db-password
  databaseRoles:
    - alloydbiamuser
```

### InfraChart

Wire the user to a cluster deployed in the same InfraPipeline:

```yaml
spec:
  cluster:
    valueFrom:
      kind: GcpAlloydbCluster
      name: orders-db
      fieldPath: status.outputs.cluster_id
  password: $secret/orders-app-db-password
```

## Key Configuration

**Authentication type** — `ALLOYDB_BUILT_IN` (default) uses a password. `ALLOYDB_IAM_USER` enables passwordless IAM database authentication via AlloyDB Auth Proxy or Language Connectors.

**User ID** — Immutable after creation. One user per application service is the recommended pattern.

**Database roles** — Grant only the PostgreSQL roles the application needs. `alloydbiamuser` is the standard application role for built-in users.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |
| **GcpAlloydbCluster** | `cluster` | `status.outputs.cluster_id` |

### What This Component Provides

| Output | Description |
|--------|-------------|
| `user_id` | The AlloyDB user resource name |
| `name` | Short user name within the cluster |

## Common Patterns

**Application user** — password-authenticated built-in user with `alloydbiamuser` role. Start from **Application User**.

**IAM user** — passwordless IAM-authenticated user for services connecting via Auth Proxy. Start from **IAM User**.

## Works With

- [**GCP AlloyDB Cluster**](/cloud-catalog/gcp-alloydb-cluster) — parent cluster this user lives on
- [**GCP AlloyDB Instance**](/cloud-catalog/gcp-alloydb-instance) — read pools on the same cluster
- [**GCP Project**](/cloud-catalog/gcp-project) — optional project override
