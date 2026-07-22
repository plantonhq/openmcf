# UDP Game Server

Forward all datagrams arriving on a Gateway's UDP listener to a single game
server Service on a custom port. This is the pattern for any single-backend
datagram protocol -- game servers, syslog collectors, metrics agents -- where the
listener's port selects the traffic and the route simply forwards it.

## When to Use

- You expose a UDP workload (a game server, syslog, a custom datagram protocol)
  behind a Gateway.
- One backend Service receives all traffic on the listener -- no splitting.

## Key Configuration Choices

- **`parentRefs[].name`** -- attaches the route to the Gateway; add `sectionName` to target the UDP listener. The name is a foreign key: write `value: <literal>` for an existing Gateway, or `valueFrom:` (kind `KubernetesGateway`, fieldPath `status.outputs.gateway_name`) so the route deploys after its Planton-managed Gateway.
- **`backendRefs[].name`** -- a foreign key like the parent ref: `value:` for a literal Service name, `valueFrom:` (kind `KubernetesService`, fieldPath `status.outputs.service_name`) for a Planton-managed backend.
- **`backendRefs[].port` (`7777`)** -- the backend Service port that receives the forwarded datagrams; set it to your server's port.

## Prerequisites

- The Gateway API CRDs are installed (`KubernetesGatewayApiCrds`). UDPRoute is
  part of the standard channel as of Gateway API v1.6.
- The `Gateway` referenced in `parentRefs` exists (`KubernetesGateway`) with a
  listener of protocol `UDP`.
- The target namespace exists (`KubernetesNamespace`).
- The backend Service exists in the route's namespace.

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `<gateway-name>` | Name of the `KubernetesGateway` this route attaches to (inside `name.value`, or switch to `valueFrom`). |
| `<game-server-service>` | Name of the backend Kubernetes Service. |

Adjust `backendRefs[].port` to your server's UDP port (the preset uses `7777`).
Set `spec.namespace.value` to your namespace, or replace it with a `valueFrom`
reference to a `KubernetesNamespace`.
