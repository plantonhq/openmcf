---
title: "Manifest"
description: "Manifest deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesmanifest"
---

# Kubernetes Manifest

Deploys raw Kubernetes YAML manifests to any Kubernetes cluster, acting as a generic escape hatch for resources that do not have a dedicated deployment component. Supports single-document and multi-document manifests (separated by `---`), including Deployments, Services, ConfigMaps, CRDs, Custom Resources, and any other valid Kubernetes resource types. Integrates with Planton's Provider Connections, ValueFromRef, and Stack Job tracking.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kubernetes Namespace** -- created (with the standard governance labels) only when `createNamespace` is `true`, and deleted with the resource; otherwise deploys into an existing namespace
- **Kubernetes Resources from Manifest** -- all resources defined in the `manifestYaml` field are applied to the cluster exactly as written. Multi-document YAML (resources separated by `---`) is supported with automatic CRD ordering: a CRD and its custom resources can ship in the same manifest. The manifest content itself is never mutated -- no injected labels, no rewritten fields.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target Kubernetes cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Appropriate RBAC permissions** for the resource types in your manifest. Cluster-scoped resources (CRDs, ClusterRoles) require cluster-admin access. Namespace-scoped resources require permissions in the target namespace.
- **CRDs pre-installed** if your manifest references custom resources. The manifest component does not handle CRD installation ordering across separate Cloud Resources.

## Deploy

### Console

Open the deployment store, find **Kubernetes Manifest**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Single Resource** preset in the [Presets](#presets) tab to deploy a ConfigMap as a starting point, then replace the YAML with your actual manifest.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesManifest
metadata:
  name: app-config
  org: acme-corp
  env: prod
spec:
  namespace:
    value: "my-app"
  createNamespace: false
  manifestYaml: |
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: app-settings
    data:
      config.yaml: |
        database_pool_size: 10
        log_level: info
```

```shell
planton apply -f manifest.yaml
```

This deploys a ConfigMap into the `my-app` namespace. The `manifestYaml` field accepts any valid Kubernetes YAML, including multi-document manifests with `---` separators. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the manifest deployment to a namespace managed by another Cloud Resource:

```yaml
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: app-namespace
      fieldPath: spec.name
  createNamespace: false
```

The InfraPipeline deploys the namespace first, then applies the manifest resources into it.

## Key Configuration

These are the most important decisions when configuring a Kubernetes Manifest deployment. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Manifest YAML content** -- The `manifestYaml` field is the core of this component. Paste any valid Kubernetes YAML, including multi-document manifests separated by `---`. The IaC module handles resource ordering automatically, applying CRDs before custom resources that depend on them.

**Namespace behavior** -- The `namespace` field is an ANCHOR, not a rewrite: namespaced documents that do not declare their own `metadata.namespace` are applied there, documents with an explicit namespace keep it, and cluster-scoped documents (CRDs, ClusterRoles) are never distorted by it. Set `createNamespace: true` to create the anchor namespace if it does not exist.

**Readiness semantics** -- By default the deploy blocks until every applied resource becomes ready: Deployments, DaemonSets, and StatefulSets complete their rollout, and other kinds pass readiness checks. Set `skipAwait: true` to return as soon as the API server accepts every document -- for manifests whose readiness depends on something deployed later (a webhook configuration waiting on its Service) or that intentionally stay not-ready at install time.

**Use case boundaries** -- A first-class catalog component always wins: typed components validate configuration before deploy, export composable outputs, and document their trade-offs field by field -- raw YAML does none of that. Reach for a manifest only when the catalog has no component for what you need to apply (a vendor's install manifest, a CRD bundle, an exotic custom resource).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **KubernetesNamespace** | `namespace` | `spec.name` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | The anchor namespace: where namespaced documents without an explicit `metadata.namespace` were applied | Application deployment manifests and service discovery |
| `applied_resources` | One entry per applied document, in manifest order, formatted as `<apiVersion>/<Kind>/<name>` (e.g. `apps/v1/Deployment/app`) | Auditing what a manifest shipped without re-parsing its YAML |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Single resource deployment** -- Deploys a single Kubernetes resource (ConfigMap, ServiceAccount, RBAC rule) into an existing namespace. Replace the example ConfigMap with any Kubernetes resource type. Start from the **Single Resource** preset.

## Works With

- [**Kubernetes Namespace**](/cloud-catalog/kubernetes-namespace) -- provides the default namespace for manifest resources