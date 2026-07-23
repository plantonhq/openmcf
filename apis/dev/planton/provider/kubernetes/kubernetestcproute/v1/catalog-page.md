# Kubernetes TCP Route

Provision a Kubernetes Gateway API `TCPRoute` -- namespaced rules that forward raw
TCP connections arriving on a Gateway listener to backend Services. A TCP route
has no matching: connections on the listener's port are forwarded to the rule's
backends. Use it to expose non-HTTP TCP services (databases, brokers, custom
protocols) through a Gateway.

TCPRoute is a GA, standard-channel resource served as
`gateway.networking.k8s.io/v1` (Gateway API v1.6.1); the default
standard-channel CRD install includes it.

## What Gets Created

- A namespaced `gateway.networking.k8s.io/v1` `TCPRoute` custom resource.
- One or more rules (max 16), each forwarding to one or more weighted backend
  refs.

## Prerequisites

- Gateway API standard-channel CRDs installed (`KubernetesGatewayApiCrds`).
- A `Gateway` to attach to via `parentRefs` (`KubernetesGateway`) with a `TCP`
  listener.
- The target namespace (`KubernetesNamespace`).
- The backend Services the route forwards to.

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesTcpRoute
metadata:
  name: postgres-route
spec:
  namespace:
    value: app-ns
  parentRefs:
    - name:
        value: my-gateway
      sectionName: tcp
  rules:
    - backendRefs:
        - name:
            value: postgres
          port: 5432
```

```bash
planton apply -f tcproute.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `namespace` | reference | Namespace to create the route in. |
| `rules` | list | One to 16 routing rules. |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `parentRefs` | list | Gateways (and optional listener `sectionName`) the route attaches to. Each `name` is a reference (defaults to `KubernetesGateway`). |
| `rules[].name` | string | Optional rule name. |
| `rules[].backendRefs` | list | Weighted backends; each `name` is a reference (defaults to `KubernetesService`). |

## Examples

### Port forwarding

```yaml
spec:
  namespace:
    value: app-ns
  parentRefs:
    - name:
        value: my-gateway
      sectionName: tcp
  rules:
    - backendRefs:
        - name:
            value: postgres
          port: 5432
```

### Weighted backends (canary)

```yaml
spec:
  namespace:
    value: app-ns
  parentRefs:
    - name:
        value: my-gateway
      sectionName: tcp
  rules:
    - backendRefs:
        - name:
            value: broker-stable
          port: 9092
          weight: 90
        - name:
            value: broker-canary
          port: 9092
          weight: 10
```

## Stack Outputs

| Output | Description |
|--------|-------------|
| `routeName` | Name of the created TCPRoute (equals metadata.name). |
| `namespace` | Namespace the TCPRoute was created in. |

## Related Components

- [Kubernetes Gateway](kubernetesgateway)
- [Kubernetes TLS Route](kubernetestlsroute)
- [Kubernetes UDP Route](kubernetesudproute)
- [Kubernetes Gateway Class](kubernetesgatewayclass)
- [Kubernetes Gateway API CRDs](kubernetesgatewayapicrds)
- [Kubernetes Namespace](kubernetesnamespace)
