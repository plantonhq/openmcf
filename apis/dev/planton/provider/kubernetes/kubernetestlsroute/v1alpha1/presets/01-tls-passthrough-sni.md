# TLS Passthrough by SNI

The most common TLSRoute: match a TLS connection by its SNI hostname and forward
it, unmodified (passthrough), to a single backend Service. The backend terminates
TLS itself -- the Gateway never sees the plaintext. This is the standard pattern
for exposing a service that does its own TLS termination (databases, mTLS
services, or apps that must hold their own certificate).

## When to Use

- You want end-to-end TLS where the backend, not the Gateway, terminates the
  connection.
- You route by SNI hostname only (TLS routes have no path/header matching).
- The parent Gateway has a TLS listener in `Passthrough` mode.

## Key Configuration Choices

- **`parentRefs[].name`** -- attaches the route to the Gateway; add `sectionName` to target the TLS listener. The name is a foreign key: write `value: <literal>` for an existing Gateway, or `valueFrom:` (kind `KubernetesGateway`, fieldPath `status.outputs.gateway_name`) so the route deploys after its Planton-managed Gateway.
- **`hostnames`** -- the SNI hostnames that select this route; a leading `*.` is a suffix match. At least one is required, and IP addresses are not allowed (RFC 6066).
- **`backendRefs[].name`** -- a foreign key like the parent ref: `value:` for a literal Service name, `valueFrom:` (kind `KubernetesService`, fieldPath `status.outputs.service_name`) for a Planton-managed backend.
- **`backendRefs[].port`** -- the backend port that accepts the passthrough TLS connection.

## Prerequisites

- The Gateway API CRDs are installed (`KubernetesGatewayApiCrds`).
- The `Gateway` referenced in `parentRefs` exists (`KubernetesGateway`) and has a
  listener of protocol `TLS` with `tls.mode: Passthrough`.
- The target namespace exists (`KubernetesNamespace`).
- The backend Service exists in the route's namespace.

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `<gateway-name>` | Name of the `KubernetesGateway` this route attaches to (inside `name.value`, or switch to `valueFrom`). |
| `secure.example.com` | SNI hostname this route serves (a literal example value -- replace with your real host). |
| `<service-name>` | Name of the backend Kubernetes Service (inside `name.value`, or switch to `valueFrom`). |

Set `spec.namespace.value` to your namespace, or replace it with a `valueFrom`
reference to a `KubernetesNamespace`.
