# gRPC Weighted Canary

Split gRPC traffic for a service across two backends by weight -- the standard
progressive-delivery pattern. Here 90% of calls go to the stable backend and 10%
to the canary; adjust the weights to shift traffic.

## When to Use

- You are rolling out a new version of a gRPC service and want to send a small
  percentage of calls to it first.
- You want weighted traffic splitting without an external progressive-delivery
  controller.

## Key Configuration Choices

- **`backendRefs[].weight`** -- relative weight; each backend receives
  `weight / (sum of weights)` of the traffic. `90` and `10` yield a 90/10 split.
- **`method.service`** -- scopes the split to one gRPC service; omit `matches` to
  split all traffic on the route.
- **`backendRefs[].port`** -- required when the backend is a core Service.
- **`parentRefs[].name` / `backendRefs[].name` are foreign keys** -- write `value: <literal>` for existing resources, or `valueFrom:` to reference a Planton-managed `KubernetesGateway` / `KubernetesService` and deploy in dependency order.

## Prerequisites

- The Gateway API CRDs are installed (`KubernetesGatewayApiCrds`).
- The `Gateway` referenced in `parentRefs` exists (`KubernetesGateway`) and its
  listener accepts HTTP/2.
- The target namespace exists (`KubernetesNamespace`).
- Both backend gRPC Services exist in the route's namespace.

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `<gateway-name>` | Name of the `KubernetesGateway` this route attaches to (inside `name.value`, or switch to `valueFrom`). |
| `api.example.com` | Public hostname this route serves (a literal example value -- replace with your real host). |
| `helloworld.Greeter` | Fully-qualified gRPC service (a literal example value -- replace with yours). |
| `<stable-service-name>` | Name of the stable backend Kubernetes Service. |
| `<canary-service-name>` | Name of the canary backend Kubernetes Service. |

Set `spec.namespace.value` to your namespace, or replace it with a `valueFrom`
reference to a `KubernetesNamespace`.
