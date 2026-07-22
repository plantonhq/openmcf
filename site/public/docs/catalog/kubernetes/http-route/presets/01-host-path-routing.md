---
title: "Host + Path Routing"
description: "The most common HTTPRoute: match a public hostname and a path prefix, then forward to a backend Service. This is the standard pattern for exposing a web application behind a Gateway."
type: "preset"
rank: "01"
presetSlug: "01-host-path-routing"
componentSlug: "http-route"
componentTitle: "HTTP Route"
provider: "kubernetes"
icon: "package"
order: 1
---

# Host + Path Routing

The most common HTTPRoute: match a public hostname and a path prefix, then
forward to a backend Service. This is the standard pattern for exposing a web
application behind a Gateway.

## When to Use

- You expose a single application at a hostname (and optionally a path prefix).
- You want all traffic for a host routed to one Service.
- You are wiring an app behind an existing Gateway (Istio, Envoy Gateway, ...).

## Key Configuration Choices

- **`parentRefs[].name`** -- attaches the route to the Gateway; add `sectionName` to target a specific listener. The name is a foreign key: write `value: <literal>` for an existing Gateway, or `valueFrom:` (kind `KubernetesGateway`, fieldPath `status.outputs.gateway_name`) so the route deploys after its Planton-managed Gateway.
- **`hostnames`** -- the Host header values that select this route; a leading `*.` is a suffix match.
- **`path` (PathPrefix `/`)** -- matches every path under the host; narrow it (for example `/api`) to split traffic by path.
- **`backendRefs[].name`** -- a foreign key like the parent ref: `value:` for a literal Service name, `valueFrom:` (kind `KubernetesService`, fieldPath `status.outputs.service_name`) for a Planton-managed backend.
- **`backendRefs[].port`** -- required when the backend is a core Service.

## Prerequisites

- The Gateway API CRDs are installed (`KubernetesGatewayApiCrds`).
- The `Gateway` referenced in `parentRefs` exists (`KubernetesGateway`).
- The target namespace exists (`KubernetesNamespace`).
- The backend Service exists in the route's namespace.

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `<gateway-name>` | Name of the `KubernetesGateway` this route attaches to (inside `name.value`, or switch to `valueFrom`). |
| `app.example.com` | Public hostname this route serves (a literal example value -- replace with your real host). |
| `<service-name>` | Name of the backend Kubernetes Service (inside `name.value`, or switch to `valueFrom`). |

Set `spec.namespace.value` to your namespace, or replace it with a `valueFrom`
reference to a `KubernetesNamespace`.
