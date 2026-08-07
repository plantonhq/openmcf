# Kubernetes Listener Set

Provision a Kubernetes Gateway API `ListenerSet` -- a namespaced set of
additional listeners merged into an existing Gateway. ListenerSets are the
per-team delegation model for shared gateways: a platform team runs one central
Gateway, and each tenant team attaches its own listeners (ports, hostnames, TLS
certificates, route-attachment policy) from its own namespace, without touching
the Gateway itself.

ListenerSet is a standard-channel resource served as
`gateway.networking.k8s.io/v1` (standard since Gateway API v1.5). The parent
Gateway must opt in via `allowed_listeners`; Gateways allow no ListenerSet
attachment by default.

## What Gets Created

- A namespaced `gateway.networking.k8s.io/v1` `ListenerSet` custom resource.
- One to 64 listeners (HTTP, HTTPS, TLS, TCP, or UDP) merged into the parent
  Gateway, each with optional per-listener TLS termination/passthrough and
  route-attachment policy.

## Prerequisites

- Gateway API CRDs v1.5.0+ installed on the cluster
  (`KubernetesGatewayApiCrds`).
- A parent `Gateway` (`KubernetesGateway`) whose `allowed_listeners` permits
  attachment from this namespace.
- The target namespace (`KubernetesNamespace`).
- For HTTPS listeners, a TLS Secret (e.g. from `KubernetesCertificate`) in the
  ListenerSet's own namespace.

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesListenerSet
metadata:
  name: tenant-listeners
spec:
  namespace:
    value: tenant-ns
  parentRef:
    name:
      value: shared-gateway
    namespace: gateway-ns
  listeners:
    - name: tenant-https
      hostname: tenant.example.com
      port: 443
      protocol: HTTPS
      tls:
        mode: Terminate
        certificateRefs:
          - name:
              value: tenant-tls
```

```bash
planton apply -f listenerset.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `namespace` | reference | Namespace to create the ListenerSet in. |
| `parentRef` | object | The Gateway the listeners attach to; `name` is a reference (defaults to `KubernetesGateway`). |
| `listeners` | list | One to 64 listeners (name, port, protocol). |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `parentRef.namespace` | string | Parent Gateway's namespace; defaults to the ListenerSet's own. |
| `listeners[].hostname` | string | Virtual host to match (HTTP/HTTPS/TLS); wildcard prefix allowed. |
| `listeners[].tls` | object | Per-listener TLS; `certificateRefs[].name` is a reference (defaults to `KubernetesSecret`, typically issued via `KubernetesCertificate`). |
| `listeners[].allowedRoutes` | object | Which Route kinds/namespaces may attach (defaults to the ListenerSet's own namespace). |

## Examples

### Tenant HTTPS listener with label-selected routes

```yaml
spec:
  namespace:
    value: tenant-ns
  parentRef:
    name:
      value: shared-gateway
    namespace: gateway-ns
  listeners:
    - name: tenant-https
      hostname: tenant.example.com
      port: 443
      protocol: HTTPS
      tls:
        mode: Terminate
        certificateRefs:
          - name:
              value: tenant-tls
      allowedRoutes:
        namespaces:
          from: Selector
          selector:
            matchLabels:
              team: tenant
        kinds:
          - kind: HTTPRoute
```

### TLS passthrough listener

```yaml
spec:
  namespace:
    value: tenant-ns
  parentRef:
    name:
      value: shared-gateway
    namespace: gateway-ns
  listeners:
    - name: tenant-tls-passthrough
      hostname: db.tenant.example.com
      port: 8443
      protocol: TLS
      tls:
        mode: Passthrough
```

## Stack Outputs

| Output | Description |
|--------|-------------|
| `listenerSetName` | Name of the created ListenerSet (target of Route parentRefs with `kind: ListenerSet`). |
| `namespace` | Namespace the ListenerSet was created in. |
| `gatewayName` | Name of the parent Gateway the listeners attach to. |

## Related Components

- [Kubernetes Gateway](kubernetesgateway)
- [Kubernetes Gateway API CRDs](kubernetesgatewayapicrds)
- [Kubernetes Certificate](kubernetescertificate)
- [Kubernetes HTTP Route](kuberneteshttproute)
- [Kubernetes Namespace](kubernetesnamespace)
