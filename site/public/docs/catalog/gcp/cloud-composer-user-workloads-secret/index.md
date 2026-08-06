---
title: "Cloud Composer User Workloads Secret"
description: "Cloud Composer User Workloads Secret deployment documentation"
icon: "package"
order: 100
componentName: "gcpcloudcomposeruserworkloadssecret"
---

# GCP Cloud Composer User Workloads Secret

Deploys a Kubernetes Secret into a Cloud Composer environment's workloads namespace — credentials for Airflow DAGs (connection URIs, passwords, API tokens) without baking them into DAG code.

## What Gets Created

A Kubernetes Secret that Composer manages in the environment's GKE cluster. `KubernetesPodOperator` tasks and Airflow connections consume it by name.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing Composer environment** — referenced via `environment` (a `GcpCloudComposerEnvironment` resource or a literal name)
- **IAM permissions** — Composer environment admin access (e.g. `roles/composer.admin`)

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpCloudComposerUserWorkloadsSecret
metadata:
  name: orders-db-connection
spec:
  region: us-central1
  environment:
    valueFrom:
      kind: GcpCloudComposerEnvironment
      name: prod-airflow
      fieldPath: status.outputs.environment_name
  secretName: orders-db-connection
  data:
    # base64 of "postgresql://appuser:s3cr3t@10.20.0.12:5432/orders"
    connection: cG9zdGdyZXNxbDovL2FwcHVzZXI6czNjcjN0QDEwLjIwLjAuMTI6NTQzMi9vcmRlcnM=
```

```shell
planton apply -f secret.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `region` | `string` | — (required) | Region of the Composer environment. Immutable. |
| `environment` | `StringValueOrRef` | — (required) | The Composer environment (its `environment_name`). Immutable. |
| `secretName` | `string` | — (required) | Kubernetes Secret name DAGs reference. Immutable. |
| `data` | `map<string,string>` | — (required, min 1) | Entries with base64-encoded values (`echo -n 'value' | base64`). Mutable. |
| `projectId` | `StringValueOrRef` | provider default | Project of the Composer environment. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `name` | Fully qualified resource name |
| `secret_name` | The Kubernetes Secret name DAGs reference |

The Secret's data is deliberately never exported; it is held as a secret in IaC state.

## Related Components

- [GcpCloudComposerEnvironment](/docs/catalog/gcp/cloud-composer-environment) — the environment the Secret is delivered into
- [GcpCloudComposerUserWorkloadsConfigMap](/docs/catalog/gcp/cloud-composer-user-workloads-configmap) — the non-secret sibling for plain configuration
