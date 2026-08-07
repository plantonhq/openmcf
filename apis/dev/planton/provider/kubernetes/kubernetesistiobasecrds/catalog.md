# Istio Base CRDs on Kubernetes

Installs the Istio base Custom Resource Definitions (the `istio/base` CRD bundle) on any Kubernetes cluster. These are the API types only -- no control plane, gateways, or sidecars. Installing them is the lightweight prerequisite that lets a cluster accept and server-side-validate Istio's networking, security, and telemetry resources.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Istio CRDs** -- cluster-scoped Custom Resource Definitions applied from the official Istio release manifest. They register the Istio API types so the cluster can accept Istio resources; no pods, services, or workloads are created.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

## What It Unlocks

The base CRDs are the prerequisite for the typed Istio resources:

- **Traffic management** -- DestinationRule, ServiceEntry
- **Security** -- PeerAuthentication, RequestAuthentication, AuthorizationPolicy
- **Observability** -- Telemetry
- **Extensibility** -- EnvoyFilter

Install this once per cluster before creating any of the above.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Cluster admin access** required. CRDs are cluster-scoped resources that require elevated permissions to install.
- **No conflicting Istio installation** -- if a different Istio version's CRDs are already present, reconcile that first. The installed CRD set is pinned to a specific Istio release for schema compatibility (see below).

## Deploy

### Console

Open the deployment store, find **Istio Base CRDs on Kubernetes**, and click **Deploy**. The creation wizard walks you through environment and connection selection and a short confirmation of what gets installed. There are no settings to configure -- the CRD set and its version are fixed.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesIstioBaseCrds
metadata:
  name: istio-base-crds
  org: acme-corp
  env: prod
spec: {}
```

```shell
planton apply -f istio-base-crds.yaml
```

This installs the Istio base CRDs cluster-wide. No namespace is required since CRDs are cluster-scoped, and there are no spec fields to set.

## Key Configuration

There is nothing to configure. The Istio CRD set is pinned to the Istio release that the typed Istio resources were built against, so the API schema on the cluster always matches the resources you create. A user-selectable version would risk a mismatch that silently drops or rejects fields, so it is intentionally fixed.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources and operators can reference:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `installed_release` | The Istio release the CRDs were installed from (the pinned SDK tag, e.g. `1.30.3`) | Verifying the CRD schema version present on the cluster |
| `installed_manifest_url` | The exact CRD bundle URL that was applied | Auditing precisely what landed on the cluster |

## Common Patterns

Browse the [Presets](#presets) tab for a ready-to-deploy configuration.

**Standard installation** -- Installs the pinned Istio base CRDs with no additional configuration. Suitable for any cluster that will run typed Istio traffic, security, or telemetry resources. Start from the **standard** preset.

## Works With

This component is a prerequisite for the typed Istio resources (DestinationRule, ServiceEntry, PeerAuthentication, RequestAuthentication, AuthorizationPolicy, Telemetry, EnvoyFilter). It installs the API types only; to run a full Istio mesh (control plane and data plane), use the dedicated Istio mesh component instead.
