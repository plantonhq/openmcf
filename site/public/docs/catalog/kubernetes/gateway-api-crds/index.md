---
title: "Gateway API CRDs"
description: "Gateway API CRDs deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesgatewayapicrds"
---

# Kubernetes Gateway API CRDs

Installs the Kubernetes Gateway API Custom Resource Definitions (CRDs) on a target Kubernetes cluster. The Gateway API is the next-generation, role-oriented API for managing ingress and service mesh traffic, replacing the legacy Ingress resource with richer routing primitives such as Gateway, HTTPRoute, GRPCRoute, and ReferenceGrant. This component fetches the official CRD manifests from the upstream `kubernetes-sigs/gateway-api` releases and applies them directly to the cluster.

## What Gets Created

When you deploy a KubernetesGatewayApiCrds resource, Planton provisions:

- **Standard Channel CRDs** — as of Gateway API v1.6: `GatewayClass`, `Gateway`, `ListenerSet`, `HTTPRoute`, `GRPCRoute`, `TLSRoute`, `TCPRoute`, `UDPRoute`, `ReferenceGrant`, and `BackendTLSPolicy` custom resource definitions, enabling the full stable Gateway API surface
- **Experimental Channel CRDs** (when `installChannel` is set to `experimental`) — all standard CRDs (with additional experimental fields) plus experimental resources such as `XBackend`, `XBackendTrafficPolicy`, and `XMesh`

No namespaced workloads are created. The CRDs are cluster-scoped and make the Gateway API resource types available for any namespace in the cluster.

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **Cluster-admin privileges** on the target cluster, because CRD installation requires cluster-wide write access
- **Network access** from the deployment runner to `https://github.com/kubernetes-sigs/gateway-api/releases/download` to fetch CRD manifests

## Quick Start

Create a file `gateway-api-crds.yaml`:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesGatewayApiCrds.gateway-api
spec: {}
```

Deploy:

```shell
planton apply -f gateway-api-crds.yaml
```

This installs the standard-channel Gateway API CRDs at the default version (v1.6.1) on the cluster configured in your environment.

## Configuration Reference

### Required Fields

This component has no strictly required spec fields. An empty `spec: {}` installs the standard-channel CRDs at the default version.

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `version` | `string` | `v1.6.1` | Gateway API release version to install. Must match the pattern `v<major>.<minor>.<patch>` with an optional pre-release suffix (e.g., `v1.6.1`, `v1.6.0-rc1`). Installing an older version narrows what the catalog's Gateway API kinds can deploy (TCPRoute/UDPRoute are standard-channel from v1.6.0; ListenerSet from v1.5.0). |
| `installChannel.channel` | `enum` | `standard` | CRD installation channel. `standard` installs the full stable set (GatewayClass, Gateway, ListenerSet, all five route kinds, ReferenceGrant, BackendTLSPolicy). `experimental` adds experimental fields and resources such as XBackend, XBackendTrafficPolicy, and XMesh. |

## Examples

### Standard CRDs at Default Version

Installs the stable Gateway API CRDs using all defaults:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesGatewayApiCrds.gateway-api
spec: {}
```

### Experimental Channel with Specific Version

Installs all Gateway API CRDs, including the experimental resources, at a pinned version:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api-experimental
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: staging.KubernetesGatewayApiCrds.gateway-api-experimental
spec:
  version: "v1.6.1"
  installChannel:
    channel: experimental
```

### Target a Specific GKE Cluster

Installs the standard CRDs on a named GKE cluster in a production environment:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesGatewayApiCrds
metadata:
  name: gateway-api-prod
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesGatewayApiCrds.gateway-api-prod
spec:
  version: "v1.6.1"
  installChannel:
    channel: standard
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `installedVersion` | `string` | Gateway API version that was installed (e.g., `v1.6.1`) |
| `installedChannel` | `string` | Installation channel that was used (`standard` or `experimental`) |
| `installedManifestUrl` | `string` | Full URL of the Gateway API CRD bundle that was applied (encodes version + channel, e.g., `.../releases/download/v1.6.1/standard-install.yaml`) |

## Related Components

- [KubernetesHelmRelease](/docs/catalog/kubernetes/helm-release) — deploy a Gateway API controller (such as Envoy Gateway or Istio) after the CRDs are in place
- [KubernetesGatewayClass](/docs/catalog/kubernetes/gateway-class) and [KubernetesGateway](/docs/catalog/kubernetes/gateway) — the first-class Gateway API kinds that require these CRDs
- [KubernetesHttpRoute](/docs/catalog/kubernetes/http-route) and its route siblings (gRPC, TLS, TCP, UDP) — attach routing to a Gateway
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) — create namespaces for Gateway API controller workloads
