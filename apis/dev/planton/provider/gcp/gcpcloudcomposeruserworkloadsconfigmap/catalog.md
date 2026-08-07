# GCP Cloud Composer User Workloads ConfigMap

Deploys a Kubernetes ConfigMap inside a Cloud Composer environment's GKE cluster, managed as a first-class resource: DAGs and workers read its entries at runtime, and configuration changes flow through the same declarative review process as everything else — no `kubectl` against the environment's cluster and no drift. The component integrates with Planton's Provider Connections for GCP credential management and composes against the parent environment via ValueFromRef.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **User Workloads ConfigMap** -- a Kubernetes ConfigMap in the Composer environment's user-workloads namespace, holding the configured key-value data
- **Environment Attachment** -- the ConfigMap is created through the Cloud Composer API against the referenced environment, so its lifecycle follows the environment's

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A Cloud Composer environment** (Composer 3, or Composer 2 with user workloads support) in the target project and region — reference a GcpCloudComposerEnvironment Cloud Resource via ValueFromRef.
- **Cloud Composer API** enabled in the target project.

## Deploy

### Console

Open the deployment store, find **GCP Cloud Composer User Workloads ConfigMap**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **DAG Configuration** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudComposerUserWorkloadsConfigMap
metadata:
  name: dag-config
  org: acme-corp
  env: prod
spec:
  region: us-central1
  environment:
    valueFrom:
      kind: GcpCloudComposerEnvironment
      name: prod-airflow
      fieldPath: status.outputs.environment_name
  configMapName: dag-config
  data:
    batch_size: "100"
    retry_count: "3"
```

```shell
planton apply -f config-map.yaml
```

A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, the environment reference orders the deploy graph automatically: the InfraPipeline provisions the Composer environment first, then creates the ConfigMap inside it.

## Key Configuration

**The attachment pair is the identity** -- `environment` + `configMapName` locate the object; both are immutable (pointing the same data at a different environment is a new ConfigMap there). The name follows Kubernetes object-name rules: lowercase letters, numbers, and hyphens, up to 63 characters.

**Data is the day-2 lever** -- `data` requires at least one entry and updates the Kubernetes object in place; running tasks keep the values they already read, new tasks see the new values.

**Plain text only** -- ConfigMap values are visible to anyone with cluster read access. Credentials belong in a GcpCloudComposerUserWorkloadsSecret.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |
| **GcpCloudComposerEnvironment** | `environment` | `status.outputs.environment_name` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `name` | Fully qualified ConfigMap resource name | Composer API handles, audit |
| `config_map_name` | The short Kubernetes object name | The name DAG code mounts and reads |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**DAG configuration** -- per-environment tuning values (batch sizes, retry counts, schedule toggles) DAGs read at runtime. Start from the **DAG Configuration** preset.

**Feature flags** -- boolean-style flags that gate pipeline behavior between dev and prod without code changes. Start from the **Feature Flags** preset.

## Works With

- [**GCP Cloud Composer Environment**](/cloud-catalog/gcp-cloud-composer-environment) -- the environment this ConfigMap lives in
- [**GCP Cloud Composer User Workloads Secret**](/cloud-catalog/gcp-cloud-composer-user-workloads-secret) -- the sensitive twin for credentials
- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project the environment runs in
