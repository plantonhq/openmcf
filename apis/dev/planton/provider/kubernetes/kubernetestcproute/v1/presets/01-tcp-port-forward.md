# TCP Port Forwarding

The most common TCPRoute: forward all connections arriving on a Gateway's TCP
listener to a backend Service. A TCP route has no matching -- the listener's port
selects the traffic, and the route forwards it. This is the standard pattern for
exposing a non-HTTP TCP service (a database, a message broker, a custom protocol)
through a Gateway.

## When to Use

- You expose a raw TCP service (Postgres, Redis, Kafka, a custom protocol) behind
  a Gateway.
- You route purely by listener port -- there is no application-layer matching.

## Key Configuration Choices

- **`parentRefs[].name`** -- attaches the route to the Gateway; add `sectionName` to target the TCP listener (the listener's port determines which connections arrive). The name is a foreign key: write `value: <literal>` for an existing Gateway, or `valueFrom:` (kind `KubernetesGateway`, fieldPath `status.outputs.gateway_name`) so the route deploys after its Planton-managed Gateway.
- **`backendRefs[].name`** -- a foreign key like the parent ref: `value:` for a literal Service name, `valueFrom:` (kind `KubernetesService`, fieldPath `status.outputs.service_name`) for a Planton-managed backend.
- **`backendRefs[].port`** -- the backend Service port that receives the forwarded connection.

## Prerequisites

- The Gateway API CRDs are installed (`KubernetesGatewayApiCrds`). TCPRoute is
  part of the standard channel as of Gateway API v1.6.
- The `Gateway` referenced in `parentRefs` exists (`KubernetesGateway`) with a
  listener of protocol `TCP`.
- The target namespace exists (`KubernetesNamespace`).
- The backend Service exists in the route's namespace.

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `<gateway-name>` | Name of the `KubernetesGateway` this route attaches to. |
| `<service-name>` | Name of the backend Kubernetes Service. |

Adjust `backendRefs[].port` to your backend's port (the preset uses `5432`,
Postgres). Set `spec.namespace.value` to your namespace, or replace it with a
`valueFrom` reference to a `KubernetesNamespace`.
