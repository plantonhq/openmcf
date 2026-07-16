# GCP Cloud Composer Environment

Deploys a Google Cloud Composer environment — managed Apache Airflow that provisions and operates the underlying GKE cluster, Cloud SQL metadata database, web server, and DAG bucket. Covers Composer 2.x and 3 with configurable workload sizing, private networking, CMEK encryption, retention policies, and recovery snapshots.

## What Gets Created

When you deploy a GcpCloudComposerEnvironment resource, Planton provisions:

- **Cloud Composer Environment** — a `google_composer_environment` managing the full Airflow stack (scheduler, workers, web server, triggerer, DAG processor)
- **GKE Cluster** — created and managed by Composer to run Airflow workloads
- **Cloud SQL Instance** — Airflow metadata (DAG runs, task instances, connections)
- **Cloud Storage Bucket** — DAG files, plugins, and data (or an existing bucket you supply via `storageBucket`)

The Cloud Composer API is enabled automatically. Environment creation takes 25-45 minutes.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A VPC network and subnetwork** for VPC peering networking (Composer 2.x), or a **PSC network attachment** for Composer 3
- **A service account** holding `roles/composer.worker` if specifying a custom node identity
- **A KMS key** if enabling CMEK encryption

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudComposerEnvironment
metadata:
  name: dev-airflow
spec:
  region: us-central1
  environmentSize: ENVIRONMENT_SIZE_SMALL
  softwareConfig:
    imageVersion: composer-2.9.7-airflow-2.9.3
```

```shell
planton apply -f composer.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `region` | `string` | — (required) | GCP region (e.g., `us-central1`). Immutable. |
| `projectId` | `StringValueOrRef` | provider default | GCP project for the environment. |
| `environmentName` | `string` | `metadata.name` | Explicit GCP resource name. Immutable. |
| `environmentSize` | `string` | GCP default | `ENVIRONMENT_SIZE_SMALL` / `MEDIUM` / `LARGE` / `EXTRA_LARGE`. |
| `resilienceMode` | `string` | standard | `HIGH_RESILIENCE` for multi-zone redundancy. |
| `nodeConfig` | `object` | — | Network, subnetwork, service account, tags, Composer 3 network attachment, ip-masq-agent, and `ipAllocationPolicy` (named range XOR CIDR per range). Immutable. |
| `softwareConfig` | `object` | — | Image version, Airflow overrides, PyPI packages, env vars, plugins mode, data lineage integration. Mutable. |
| `privateEnvironmentConfig` | `object` | — | Composer 2.x private networking (private endpoint, connection type, CIDRs). Immutable. |
| `workloadsConfig` | `object` | — | Per-component sizing: scheduler, web server, workers (min/max autoscaling), triggerer, DAG processor. Mutable. |
| `kmsKeyName` | `StringValueOrRef` | Google-managed | CMEK key for all Composer-managed resources. Immutable. |
| `maintenanceWindow` | `object` | — | RFC3339 window + RRULE recurrence (min 12 hours). Mutable. |
| `recoveryConfig` | `object` | — | Scheduled environment snapshots for disaster recovery. |
| `webServerNetworkAccessControl` | `object` | open | IP allowlist for the Airflow UI. |
| `masterAuthorizedNetworksConfig` | `object` | — | CIDR allowlist for the GKE control plane. |
| `dataRetentionConfig` | `object` | — | Task log storage mode; Airflow metadata retention 30-730 days (Composer 3). |
| `storageBucket` | `StringValueOrRef` | auto-created | Existing bucket for DAGs/plugins/data. Immutable. |
| `enablePrivateEnvironment` / `enablePrivateBuildsOnly` | `bool` | `false` | Composer 3 private environment flags. |
| `labels` | `map<string,string>` | `{}` | User labels merged beneath platform attribution labels. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `environment_id` | Fully qualified resource ID |
| `environment_name` | Short environment name |
| `airflow_uri` | URL of the Airflow web UI |
| `dag_gcs_prefix` | Cloud Storage prefix for DAG uploads |
| `gke_cluster` | Underlying GKE cluster name |

## Related Components

- [GcpVpcNetwork](/docs/catalog/gcp/gcpvpcnetwork) — network for VPC peering deployments
- [GcpSubnetwork](/docs/catalog/gcp/gcpsubnetwork) — subnetwork for node placement
- [GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount) — node identity
- [GcpKmsKey](/docs/catalog/gcp/gcpkmskey) — CMEK encryption key
- [GcpGcsBucket](/docs/catalog/gcp/gcpgcsbucket) — governed DAG or snapshot bucket
- [GcpCloudComposerUserWorkloadsSecret](/docs/catalog/gcp/gcpcloudcomposeruserworkloadssecret) — credentials for DAGs
- [GcpCloudComposerUserWorkloadsConfigMap](/docs/catalog/gcp/gcpcloudcomposeruserworkloadsconfigmap) — non-secret configuration for DAGs
