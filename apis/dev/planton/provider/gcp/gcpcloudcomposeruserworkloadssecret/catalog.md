# GCP Cloud Composer User Workloads Secret

Deploys a Kubernetes Secret inside a Cloud Composer environment's GKE cluster, managed as a first-class resource: DAGs read credentials from it at runtime, and rotation flows through the same declarative review process as everything else — no `kubectl` against the environment's cluster and no plaintext in any manifest. Secret values are org-secret references resolved just-in-time at deploy by the runner. The component integrates with Planton's Provider Connections for GCP credential management and composes against the parent environment via ValueFromRef.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **User Workloads Secret** -- a Kubernetes Secret in the Composer environment's user-workloads namespace, holding the configured entries (base64-encoded values)
- **Environment Attachment** -- the Secret is created through the Cloud Composer API against the referenced environment, so its lifecycle follows the environment's

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project.
- **Org Secrets** -- each data value references an org secret in the Config Manager holding the BASE64-ENCODED value (`echo -n 'value' | base64`). The runner resolves references just-in-time; plaintext never enters the manifest.
- **Planton Runner** -- required when using Runner-based credential delivery.

### GCP Project

- **A Cloud Composer environment** (Composer 3, or Composer 2 with user workloads support) in the target project and region — reference a GcpCloudComposerEnvironment Cloud Resource via ValueFromRef.
- **Cloud Composer API** enabled in the target project.

## Deploy

### Console

Open the deployment store, find **GCP Cloud Composer User Workloads Secret**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields — secret values go through the reference-only picker. Start from the **Airflow Connection** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudComposerUserWorkloadsSecret
metadata:
  name: airflow-connections
  org: acme-corp
  env: prod
spec:
  region: us-central1
  environment:
    valueFrom:
      kind: GcpCloudComposerEnvironment
      name: prod-airflow
      fieldPath: status.outputs.environment_name
  secretName: airflow-connections
  data:
    connection-string: $secret/warehouse-connection-b64
```

```shell
planton apply -f secret.yaml
```

A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, the environment reference orders the deploy graph automatically: the InfraPipeline provisions the Composer environment first, then creates the Secret inside it.

## Key Configuration

**The attachment pair is the identity** -- `environment` + `secretName` locate the object; both are immutable. The name follows Kubernetes object-name rules: lowercase letters, numbers, and hyphens, up to 63 characters.

**Values are references, base64 at the source** -- every `data` value is a `$secret/<slug>` org-secret reference; the referenced org secret must hold the BASE64-ENCODED value (the Composer API rejects raw strings at deploy). The decoded material is never placed in stack outputs.

**Rotation is the day-2 lever** -- update the org secret (or point an entry at a new one) and redeploy; the Kubernetes Secret updates in place and new task runs pick it up.

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
| `name` | Fully qualified Secret resource name | Composer API handles, audit |
| `secret_name` | The short Kubernetes object name | The name DAG code mounts and reads |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Airflow connection** -- a database or warehouse connection URI DAGs read as a mounted secret. Start from the **Airflow Connection** preset.

**API credentials** -- third-party vendor keys tasks read as environment variables. Start from the **API Credentials** preset.

## Works With

- [**GCP Cloud Composer Environment**](/cloud-catalog/gcp-cloud-composer-environment) -- the environment this Secret lives in
- [**GCP Cloud Composer User Workloads ConfigMap**](/cloud-catalog/gcp-cloud-composer-user-workloads-config-map) -- the plain-text twin for non-sensitive settings
- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project the environment runs in
