---
title: "UDP Route"
description: "UDP Route deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesudproute"
---

# Kubernetes UDP Route

Provision a Kubernetes Gateway API `UDPRoute` -- namespaced rules that forward
UDP datagrams arriving on a Gateway listener to backend Services. A UDP route
has no matching: traffic on the listener's port is forwarded to the rule's
backends. Use it to expose datagram workloads (DNS servers, syslog collectors,
game servers, and other UDP protocols) through a Gateway.

UDPRoute is a GA, standard-channel resource served as
`gateway.networking.k8s.io/v1` (Gateway API v1.6.1); the default
standard-channel CRD install includes it.

## What Gets Created

- A namespaced `gateway.networking.k8s.io/v1` `UDPRoute` custom resource.
- One or more rules (max 16), each forwarding to one or more weighted backend
  refs.

## Prerequisites

- Gateway API standard-channel CRDs installed (`KubernetesGatewayApiCrds`).
- A `Gateway` to attach to via `parentRefs` (`KubernetesGateway`) with a `UDP`
  listener.
- The target namespace (`KubernetesNamespace`).
- The backend Services the route forwards to.

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesUdpRoute
metadata:
  name: dns-route
spec:
  namespace:
    value: app-ns
  parentRefs:
    - name:
        value: my-gateway
      sectionName: dns
  rules:
    - backendRefs:
        - name:
            value: coredns
          port: 53
```

```bash
planton apply -f udproute.yaml
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

### DNS forwarding

```yaml
spec:
  namespace:
    value: app-ns
  parentRefs:
    - name:
        value: my-gateway
      sectionName: dns
  rules:
    - backendRefs:
        - name:
            value: coredns
          port: 53
```

### Weighted backends (canary)

```yaml
spec:
  namespace:
    value: app-ns
  parentRefs:
    - name:
        value: my-gateway
      sectionName: dns
  rules:
    - name: dns-forward
      backendRefs:
        - name:
            value: dns-backend
          port: 53
          weight: 90
        - name:
            value: dns-backend-canary
          port: 53
          weight: 10
```

## Stack Outputs

| Output | Description |
|--------|-------------|
| `routeName` | Name of the created UDPRoute (equals metadata.name). |
| `namespace` | Namespace the UDPRoute was created in. |

## Related Components

- [Kubernetes Gateway](kubernetesgateway)
- [Kubernetes TCP Route](kubernetestcproute)
- [Kubernetes TLS Route](kubernetestlsroute)
- [Kubernetes Gateway Class](kubernetesgatewayclass)
- [Kubernetes Gateway API CRDs](kubernetesgatewayapicrds)
- [Kubernetes Namespace](kubernetesnamespace)
