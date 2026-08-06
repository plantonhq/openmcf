# Kubernetes TLS Route

Provision a Kubernetes Gateway API `TLSRoute` -- namespaced TLS passthrough rules
that attach to a Gateway and forward connections, by SNI hostname, to backend
Services. The backend terminates TLS, so the encrypted stream is forwarded end to
end (the Gateway never sees plaintext).

TLSRoute is a standard-channel resource served as
`gateway.networking.k8s.io/v1` (standard since Gateway API v1.5).

## What Gets Created

- A namespaced `gateway.networking.k8s.io/v1` `TLSRoute` custom resource.
- Exactly one rule that forwards to one or more weighted backend refs.

## Prerequisites

- Gateway API CRDs installed on the cluster (`KubernetesGatewayApiCrds`).
- A `Gateway` to attach to via `parentRefs` (`KubernetesGateway`) with a `TLS`
  listener (typically `tls.mode: Passthrough`).
- The target namespace (`KubernetesNamespace`).
- The backend Services the route forwards to.

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesTlsRoute
metadata:
  name: secure-route
spec:
  namespace:
    value: app-ns
  parentRefs:
    - name:
        value: my-gateway
      sectionName: tls
  hostnames:
    - secure.example.com
  rules:
    - backendRefs:
        - name:
            value: secure-backend
          port: 8443
```

```bash
planton apply -f tlsroute.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `namespace` | reference | Namespace to create the route in. |
| `hostnames` | list | One to 1024 SNI hostnames that select this route (no IPs). |
| `rules` | list | Exactly one routing rule. |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `parentRefs` | list | Gateways (and optional listener `sectionName`) the route attaches to. Each `name` is a reference (defaults to `KubernetesGateway`). |
| `rules[].name` | string | Optional rule name. |
| `rules[].backendRefs` | list | Weighted backends; each `name` is a reference (defaults to `KubernetesService`). |

## Examples

### TLS passthrough by SNI

```yaml
spec:
  namespace:
    value: app-ns
  parentRefs:
    - name:
        value: my-gateway
      sectionName: tls
  hostnames:
    - secure.example.com
  rules:
    - backendRefs:
        - name:
            value: secure-backend
          port: 8443
```

### Weighted backends (canary)

```yaml
spec:
  namespace:
    value: app-ns
  parentRefs:
    - name:
        value: my-gateway
      sectionName: tls
  hostnames:
    - secure.example.com
  rules:
    - backendRefs:
        - name:
            value: secure-stable
          port: 8443
          weight: 90
        - name:
            value: secure-canary
          port: 8443
          weight: 10
```

## Stack Outputs

| Output | Description |
|--------|-------------|
| `routeName` | Name of the created TLSRoute (equals metadata.name). |
| `namespace` | Namespace the TLSRoute was created in. |

## Related Components

- [Kubernetes Gateway](kubernetesgateway)
- [Kubernetes HTTP Route](kuberneteshttproute)
- [Kubernetes Gateway Class](kubernetesgatewayclass)
- [Kubernetes Gateway API CRDs](kubernetesgatewayapicrds)
- [Kubernetes Namespace](kubernetesnamespace)
