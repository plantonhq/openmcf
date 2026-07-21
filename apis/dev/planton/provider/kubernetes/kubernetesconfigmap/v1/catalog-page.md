# Kubernetes ConfigMap

Deploys a Kubernetes ConfigMap to a target cluster through a single declarative manifest, covering the complete ConfigMap surface: UTF-8 `data` entries, base64-encoded `binaryData` entries, and the `immutable` flag. The IaC module handles label merging and namespace resolution automatically.

## What Gets Created

When you deploy a KubernetesConfigMap resource, Planton provisions:

- **ConfigMap** — a Kubernetes ConfigMap with the given `data`, `binaryData`, and `immutable` settings
- **Labels** — standard Planton tracking labels (`managed-by`, `resource`, `resource-kind`) merged with any user-provided labels
- **Annotations** — user-provided annotations applied to the ConfigMap metadata

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that already exists, or a `KubernetesNamespace` resource referenced from `spec.namespace` so both deploy in one run
- **Configuration values** ready to supply — plain strings for `data`, base64-encoded content for `binaryData`

## Quick Start

Create a file `configmap.yaml`:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesConfigMap
metadata:
  name: app-config
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesConfigMap.app-config
spec:
  name: app-config
  data:
    LOG_LEVEL: info
    FEATURE_X: "enabled"
```

Deploy:

```shell
planton apply -f configmap.yaml
```

This creates a ConfigMap named `app-config` in the `default` namespace with two keys.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.name` | `string` | Name of the ConfigMap (`metadata.name` in the cluster). Pods reference this name in `configMapRef`/`configMapKeyRef` env sources and `configMap` volumes. | 1–253 characters, matches `^[a-z0-9]([a-z0-9.-]{0,251}[a-z0-9])?$` |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.namespace` | `StringValueOrRef` | `default` | Namespace where the ConfigMap is created. Accepts a literal name (`{ value: my-namespace }`) or a reference to a `KubernetesNamespace` resource. |
| `spec.data` | `map<string, string>` | `{}` | UTF-8 configuration entries. Keys become file names (volume mounts) or env var names (`envFrom`); keys match `^[-._a-zA-Z0-9]+$`, max 253 chars. |
| `spec.binaryData` | `map<string, string>` | `{}` | Binary entries with base64-encoded values, for payloads that are not valid UTF-8. Same key rules as `data`; keys must not overlap with `data` keys. Consumable only as mounted files. |
| `spec.immutable` | `bool` | `false` | When `true`, data cannot be updated after creation (delete-and-recreate only). Reduces API server watch load; recommended for versioned production config. |
| `spec.labels` | `map<string, string>` | `{}` | Additional labels merged with standard Planton labels. |
| `spec.annotations` | `map<string, string>` | `{}` | Additional annotations applied to the ConfigMap. |

**Size limit:** the Kubernetes API server caps the combined size of `data` and `binaryData` at 1MiB; oversized ConfigMaps are rejected at apply time.

## Examples

### Application Configuration

Scalar settings plus a properties file, consumed as env vars and a mounted file respectively:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesConfigMap
metadata:
  name: backend-config
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesConfigMap.backend-config
spec:
  name: backend-config
  namespace:
    value: backend
  data:
    LOG_LEVEL: info
    CACHE_TTL_SECONDS: "300"
    application.properties: |
      server.port=8080
      management.endpoints.web.exposure.include=health,metrics
```

### Immutable, Versioned Configuration

A roll-forward config version — changes ship as `app-config-v2` plus a workload update, never as an in-place edit:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesConfigMap
metadata:
  name: app-config-v1
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesConfigMap.app-config-v1
spec:
  name: app-config-v1
  namespace:
    value: production
  immutable: true
  data:
    LOG_LEVEL: warn
    MAX_CONNECTIONS: "100"
```

### Binary Payload

A base64-encoded binary entry alongside a text entry. Binary keys mount as files only:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesConfigMap
metadata:
  name: branding-assets
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesConfigMap.branding-assets
spec:
  name: branding-assets
  namespace:
    value: frontend
  data:
    theme: dark
  binaryData:
    favicon.ico: AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAA=
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `configmapName` | `string` | Name of the created ConfigMap — the handle workloads use in `configMapRef`, `configMapKeyRef`, and volume definitions |
| `namespace` | `string` | Namespace where the ConfigMap was created. Consumers must live in the same namespace — ConfigMap references never cross namespaces |

## Related Components

- [KubernetesSecret](/docs/catalog/kubernetes/kubernetessecret) — the confidential mirror of this component; use it for passwords, tokens, keys, and certificates
- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) — provides the target namespace; reference it from `spec.namespace` to deploy both in one run
- [KubernetesDeployment](/docs/catalog/kubernetes/kubernetesdeployment) — consumes ConfigMaps as environment variables or mounted files
