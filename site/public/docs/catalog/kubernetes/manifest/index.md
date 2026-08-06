---
title: "Manifest"
description: "Manifest deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesmanifest"
---

# Kubernetes Manifest

**Check the catalog before using this component: a first-class component always wins.** Typed components validate configuration before deploy, export composable outputs, and document their trade-offs — raw YAML does none of that. KubernetesManifest is the escape hatch for resources no component covers: a vendor's install manifest, a CRD bundle, an exotic custom resource.

What it does: applies raw Kubernetes YAML — single or multi-document, core kinds or custom resources, even a CRD and its custom resources together — to a target cluster **exactly as written** (no injected labels, no rewritten fields), wrapped in Planton's declarative apply/update/destroy lifecycle. Both IaC engines handle CRD-before-custom-resource ordering in one pass and await readiness by default.

## What Gets Created

When you deploy a KubernetesManifest resource, Planton provisions:

- **Namespace** — the anchor namespace, created only when `create_namespace` is `true` (with standard Planton governance labels on the namespace object; your manifest's documents are never labeled)
- **Every document in `manifest_yaml`** — applied server-side, exactly as written

**Namespace anchoring** (identical on both engines): documents that declare their own `metadata.namespace` keep it; namespaced documents that declare none land in `spec.namespace`; cluster-scoped documents (CRDs, ClusterRoles, ...) are applied as-is — the anchor never distorts them.

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **The anchor namespace** must already exist, or set `create_namespace` to `true`
- **Valid Kubernetes YAML** — content is applied as-is; the API server is what validates your kinds

## Quick Start

Create a file `manifest.yaml` — a ConfigMap and Secret anchored in one namespace:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesManifest
metadata:
  name: app-config
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesManifest.app-config
spec:
  namespace:
    value: my-app
  create_namespace: true
  manifest_yaml: |
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: app-settings
    data:
      LOG_LEVEL: info
    ---
    apiVersion: v1
    kind: Secret
    metadata:
      name: app-credentials
    type: Opaque
    stringData:
      api-key: replace-me
```

Deploy:

```shell
planton apply -f manifest.yaml
```

Neither document declares a namespace, so both land in `my-app`, which is created first and deleted with the resource.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.namespace` | `StringValueOrRef` | The anchor namespace: where namespaced documents without an explicit `metadata.namespace` are applied. A literal name (`value: my-ns`) or a reference to a KubernetesNamespace resource. | required, non-empty |
| `spec.manifest_yaml` | `string` | The raw Kubernetes manifest YAML — a single document or multiple separated by `---`. Every valid Kubernetes resource is accepted, including a CRD and its custom resources in the same manifest. | required, at least one non-blank document |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.create_namespace` | `bool` | `false` | When `true`, the anchor namespace is created (with standard Planton governance labels) before the manifest applies and deleted with the resource. When `false`, it must already exist. |
| `spec.skip_await` | `bool` | `false` | When `true`, the deploy returns as soon as the API server accepts every document. When `false`, both engines block until readiness (workload rollouts complete; other kinds pass the engine's readiness checks). Skip for manifests whose readiness depends on something deployed later — e.g. a webhook configuration waiting on its service. |

## Examples

### CRD and Custom Resource in One Manifest

Both engines order the CRD install before the custom resource that uses it — one apply, no two-pass workaround. The CRD is cluster-scoped and untouched by the anchor; the custom resource declares no namespace, so it lands in `crontab-demo`:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesManifest
metadata:
  name: crontab-operator-config
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesManifest.crontab-operator-config
spec:
  namespace:
    value: crontab-demo
  create_namespace: true
  manifest_yaml: |
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      name: crontabs.stable.example.com
    spec:
      group: stable.example.com
      scope: Namespaced
      names:
        plural: crontabs
        singular: crontab
        kind: CronTab
      versions:
        - name: v1
          served: true
          storage: true
          schema:
            openAPIV3Schema:
              type: object
              properties:
                spec:
                  type: object
                  properties:
                    cronSpec:
                      type: string
    ---
    apiVersion: stable.example.com/v1
    kind: CronTab
    metadata:
      name: sample-crontab
    spec:
      cronSpec: "*/5 * * * *"
```

### Vendor Install Manifest

The "paste the vendor's install YAML" pattern: the manifest is applied byte-for-byte as published. `skip_await: true` because install bundles often contain resources (webhook configurations, custom resources awaiting their operator) that are not ready until the whole bundle settles:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesManifest
metadata:
  name: vendor-install
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesManifest.vendor-install
spec:
  namespace:
    value: vendor-system
  create_namespace: true
  skip_await: true
  manifest_yaml: |
    # Paste the vendor's published install.yaml here, exactly as downloaded.
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: replace-with-vendor-manifest
```

Documents in a vendor bundle that hard-code their namespace keep it; only unanchored namespaced documents land in `vendor-system`.

### Anchoring Into an Existing Namespace by Reference

The anchor can reference a KubernetesNamespace resource instead of naming a literal, so the namespace and the manifest compose in one deployment graph:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesManifest
metadata:
  name: team-quotas
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesManifest.team-quotas
spec:
  namespace:
    value_from:
      kind: KubernetesNamespace
      name: team-a-namespace
  manifest_yaml: |
    apiVersion: v1
    kind: ResourceQuota
    metadata:
      name: team-quota
    spec:
      hard:
        requests.cpu: "10"
        requests.memory: 20Gi
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `namespace` | `string` | The anchor namespace — where namespaced documents without an explicit `metadata.namespace` were applied |
| `appliedResources` | `list(string)` | One entry per applied document, in manifest order, formatted as `<apiVersion>/<Kind>/<name>` (e.g. `apps/v1/Deployment/app`). Derived from the input YAML, identical on both engines |

## Related Components

- [KubernetesHelmRelease](/docs/catalog/kubernetes/helm-release) — when the vendor ships a Helm chart, use it instead of pasting rendered YAML here
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) — a governed namespace this component's anchor can reference
- [KubernetesDeployment](/docs/catalog/kubernetes/deployment) — the first-class path for workloads; always prefer it over a raw Deployment document
