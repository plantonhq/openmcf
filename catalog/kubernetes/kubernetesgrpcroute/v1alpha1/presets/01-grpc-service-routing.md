# gRPC Service Routing

The most common GRPCRoute: match a public hostname and a gRPC service (and
optionally a method), then forward to a backend gRPC Service. This is the
standard pattern for exposing a gRPC API behind a Gateway.

## When to Use

- You expose a gRPC service at a hostname.
- You want all calls for a service (or a specific service/method) routed to one
  backend.
- You are wiring a gRPC app behind an existing Gateway (Istio, Envoy Gateway, ...).

## Key Configuration Choices

- **`parentRefs[].name`** -- attaches the route to the Gateway; add `sectionName` to target a specific listener. The name is a foreign key: write `value: <literal>` for an existing Gateway, or `valueFrom:` (kind `KubernetesGateway`, fieldPath `status.outputs.gateway_name`) so the route deploys after its Planton-managed Gateway.
- **`hostnames`** -- the authority (Host) values that select this route; a leading `*.` is a suffix match.
- **`method.service`** -- the fully-qualified gRPC service (for example `helloworld.Greeter`); add `method` to match a single RPC, or omit `method` entirely to match all services.
- **`backendRefs[].name`** -- a foreign key like the parent ref: `value:` for a literal Service name, `valueFrom:` (kind `KubernetesService`, fieldPath `status.outputs.service_name`) for a Planton-managed backend.
- **`backendRefs[].port`** -- required when the backend is a core Service.

## Prerequisites

- The Gateway API CRDs are installed (`KubernetesGatewayApiCrds`).
- The `Gateway` referenced in `parentRefs` exists (`KubernetesGateway`) and its
  listener accepts HTTP/2 (h2c over `HTTP`, or HTTP/2 over `HTTPS`).
- The target namespace exists (`KubernetesNamespace`).
- The backend gRPC Service exists in the route's namespace.

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `<gateway-name>` | Name of the `KubernetesGateway` this route attaches to (inside `name.value`, or switch to `valueFrom`). |
| `api.example.com` | Public hostname this route serves (a literal example value -- replace with your real host). |
| `helloworld.Greeter` | Fully-qualified gRPC service (a literal example value -- replace with yours). |
| `<service-name>` | Name of the backend Kubernetes Service (inside `name.value`, or switch to `valueFrom`). |

Set `spec.namespace.value` to your namespace, or replace it with a `valueFrom`
reference to a `KubernetesNamespace`.
