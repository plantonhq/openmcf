# GCP Cloud Composer User Workloads ConfigMap

Deploys a Kubernetes ConfigMap into a Cloud Composer environment's workloads namespace — non-secret configuration for Airflow DAGs (feature flags, endpoints, tuning parameters) without baking values into DAG code.

## What Gets Created

A Kubernetes ConfigMap that Composer manages in the environment's GKE cluster. DAGs and `KubernetesPodOperator` tasks consume it by name.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing Composer environment** — referenced via `environment` (a `GcpCloudComposerEnvironment` resource or a literal name)
- **IAM permissions** — Composer environment admin access (e.g. `roles/composer.admin`)

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudComposerUserWorkloadsConfigMap
metadata:
  name: orders-dag-config
spec:
  region: us-central1
  environment:
    valueFrom:
      kind: GcpCloudComposerEnvironment
      name: prod-airflow
      fieldPath: status.outputs.environment_name
  configMapName: orders-dag-config
  data:
    api_endpoint: https://api.example.com/v2
    batch_size: "500"
```

```shell
planton apply -f configmap.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `region` | `string` | — (required) | Region of the Composer environment. Immutable. |
| `environment` | `StringValueOrRef` | — (required) | The Composer environment (its `environment_name`). Immutable. |
| `configMapName` | `string` | — (required) | Kubernetes ConfigMap name DAGs reference. Immutable. |
| `data` | `map<string,string>` | — (required, min 1) | Plain configuration entries. Mutable. |
| `projectId` | `StringValueOrRef` | provider default | Project of the Composer environment. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `name` | Fully qualified resource name |
| `config_map_name` | The Kubernetes ConfigMap name DAGs reference |

## Related Components

- [GcpCloudComposerEnvironment](/docs/catalog/gcp/gcpcloudcomposerenvironment) — the environment the ConfigMap is delivered into
- [GcpCloudComposerUserWorkloadsSecret](/docs/catalog/gcp/gcpcloudcomposeruserworkloadssecret) — the secret-bearing sibling for credentials
