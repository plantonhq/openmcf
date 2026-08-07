---
title: "gRPC Route"
description: "gRPC Route deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesgrpcroute"
---

# gRPC Route on Kubernetes

Creates a namespaced Kubernetes Gateway API `GRPCRoute` -- a route that matches **gRPC requests** by hostname (`:authority`), service/method, or header, optionally transforms them with filters, and forwards them to one or more backend Services through a Gateway. GRPCRoute is part of the Gateway API **standard channel** (served as `gateway.networking.k8s.io/v1`). This is the first-class way to expose a gRPC API behind a Gateway -- weighted canaries, header-based routing, and request mirroring included. This component mirrors the upstream Gateway API `GRPCRoute` spec with full fidelity while adding proto validation, typed SDKs, and InfraChart composability.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **A namespaced GRPCRoute** named after `metadata.name` in `spec.namespace`, attached to the Gateway listener(s) in `spec.parentRefs`, matching the `spec.hostnames` and the per-rule `matches`, and forwarding to the backends declared in its `spec.rules`.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

The Gateway controller reconciles the route asynchronously: it reports per-parent `Accepted` / `ResolvedRefs` conditions in the route's status, which you observe with `kubectl` (these controller-managed values are intentionally not stored as stack outputs).

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Gateway API CRDs installed** -- deploy the `KubernetesGatewayApiCrds` component first. GRPCRoute is part of the standard channel, so the standard CRDs are sufficient (no experimental channel required).
- **A parent Gateway with an HTTP/2 listener** -- each `parentRefs` entry should resolve to a `KubernetesGateway` whose listener speaks HTTP/2 (h2c over `protocol: HTTP`, or HTTP/2 over `protocol: HTTPS`), with an `allowedRoutes` policy that admits this route. An HTTP/1-only listener cannot serve gRPC.
- **The target namespace exists** -- `spec.namespace` should resolve to a real `KubernetesNamespace`.
- **The backend gRPC Services exist** -- the `backendRefs` name in-cluster Services in the route's namespace (or in another namespace authorized by a `KubernetesReferenceGrant`).

## Deploy

### Console

Open the deployment store, find **gRPC Route on Kubernetes**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and three spec steps: **Namespace** (immutable), then **Routing** (the `hostnames` to match and the `parentRefs` Gateways to attach to), then **Rules** (per-rule matches, filters, and destination Services with weights). Start from the **gRPC Service Routing** or **gRPC Weighted Canary** preset in the [Presets](#presets) tab for a directly deployable configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesGrpcRoute
metadata:
  name: my-grpc-route
  org: acme-corp
  env: prod
spec:
  namespace:
    value: prod-apps
  parentRefs:
    - name: prod-gateway
      sectionName: grpc
  hostnames:
    - api.example.com
  rules:
    - matches:
        - method:
            service: helloworld.Greeter
      backendRefs:
        - name: greeter
          port: 9000
```

```shell
planton apply -f grpc-route.yaml
```

This creates a GRPCRoute in `prod-apps` that attaches to the `grpc` listener of `prod-gateway`, matches calls to `helloworld.Greeter` on `api.example.com`, and forwards them to the `greeter` Service on port 9000.

## Key Configuration

These are the most important decisions when configuring a GRPCRoute. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Namespace** -- The `namespace` field is where the route lives. It is **immutable**, the default namespace for its backends, and the anchor for cross-namespace rules: parent Gateways and backend Services in this namespace attach without a `ReferenceGrant`; those elsewhere require one. Reference an existing `KubernetesNamespace` or type the name directly.

**Hostnames** -- The `hostnames` (0-16) match the request's `:authority` (Host) pseudo-header. A leading `*.` is a suffix match; a bare IP is never valid. Leave empty to accept every hostname the parent listeners permit.

**Parent Gateways** -- The `parentRefs` (0-32) attach this route to Gateway listeners. Use `sectionName` to target one named listener and/or `port` to pin a port. A parent in another namespace needs a `ReferenceGrant` there and a matching listener `allowedRoutes` policy. Each parent's `name` is a foreign key to `KubernetesGateway`: reference a Planton-managed Gateway (the route then deploys after it), or pass a literal name for a Gateway or ListenerSet created outside Planton.

**Rules** -- A GRPCRoute has 1-16 `rules`, each combining:
- **Matches** (0-64) -- select requests by `method` (`service` and/or `method`, `Exact` or `RegularExpression`) and/or request `headers`. Matches within a rule are ORed; conditions within a match are ANDed. No matches means the rule applies to every request.
- **Filters** -- transform matching traffic: `RequestHeaderModifier` / `ResponseHeaderModifier` (set/add/remove headers), `RequestMirror` (shadow traffic to another backend), or `ExtensionRef` (an implementation-specific filter). Rule-level filters apply to every backend; per-backend filters exist but are implementation-specific, so prefer rule-level for portability.
- **Backends** (1-16) -- the destination Services, each with a `name`, `port`, and optional `weight` for traffic splitting (e.g. a stable/canary split at 90/10; weight `0` drains a backend).

## Outputs and Dependencies

### What This Component Consumes

This component takes foreign-key references to a `KubernetesNamespace` (via `spec.namespace`), to `KubernetesGateway` (each `parentRefs` entry's `name`), and to `KubernetesService` (each backend's `name`), so an InfraChart deploys those targets before the route and the resource graph carries the edges. Literal names cover targets created outside Planton; cross-namespace references require a `KubernetesReferenceGrant`.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `route_name` | Name of the created GRPCRoute (equals `metadata.name`) | Orders the route after its Gateway and backends in an InfraChart |
| `namespace` | The resolved namespace the route was created in | Same-namespace / ReferenceGrant rules for its parent and backend references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**gRPC service routing** -- Match a public hostname and a gRPC service (and optionally a method), then forward to a backend gRPC Service. Start from the **gRPC Service Routing** preset.

**Weighted canary** -- Split gRPC traffic for a service across a stable and a canary backend by weight, the standard progressive-delivery pattern. Start from the **gRPC Weighted Canary** preset.

## Works With

- **KubernetesGatewayApiCrds** -- installs the Gateway API CRDs (standard channel is sufficient); deploy first (prerequisite).
- **KubernetesGateway** -- the Gateway whose HTTP/2 listener this route attaches to (`parentRefs`); install first.
- **KubernetesNamespace** -- the namespace (`spec.namespace`) the route runs in.
- **KubernetesReferenceGrant** -- authorizes cross-namespace parent or backend references from this route.
- **KubernetesService** -- the backend gRPC workloads (`backendRefs`) that receive forwarded requests.
