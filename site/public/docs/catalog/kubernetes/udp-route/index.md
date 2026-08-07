---
title: "UDP Route"
description: "UDP Route deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesudproute"
---

# UDP Route on Kubernetes

Creates a namespaced Kubernetes Gateway API `UDPRoute` -- a route that forwards **UDP datagrams** arriving on a Gateway listener to one or more backend Services. UDPRoute is a **layer-4, connectionless** route: it has no hostnames, no matches, and no filters -- there is no connection or request structure to match on. Datagrams on the parent listener's port are simply forwarded to the rule's backends. Typical backends are DNS servers, syslog collectors, game servers, and other datagram protocols. This component mirrors the upstream Gateway API `UDPRoute` (GA, `gateway.networking.k8s.io/v1`) spec with full fidelity while adding proto validation, typed SDKs, and InfraChart composability.

> **Standard channel from v1.6.** `UDPRoute` graduated to GA and is served as `gateway.networking.k8s.io/v1` in the standard channel from Gateway API v1.6.0 (it was an experimental v1alpha2 resource in earlier releases). Deploy `KubernetesGatewayApiCrds` at v1.6.0+ first.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **A namespaced UDPRoute** named after `metadata.name` in `spec.namespace`, attached to the Gateway listener(s) in `spec.parentRefs`, forwarding datagrams to the backends declared in its `spec.rules`.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

The Gateway controller reconciles the route asynchronously: it reports per-parent `Accepted` / `ResolvedRefs` conditions in the route's status, which you observe with `kubectl` (these controller-managed values are intentionally not stored as stack outputs).

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Gateway API CRDs at v1.6.0+ installed** -- deploy the `KubernetesGatewayApiCrds` component first. UDPRoute is standard-channel only from v1.6.0; older CRD releases will not register it.
- **A parent Gateway with a UDP listener** -- each `parentRefs` entry should resolve to a `KubernetesGateway` that has a listener of `protocol: UDP`, and an `allowedRoutes` policy that admits this route.
- **The target namespace exists** -- `spec.namespace` should resolve to a real `KubernetesNamespace`.
- **The backend Services exist** -- the `backendRefs` name in-cluster Services in the route's namespace (or in another namespace authorized by a `KubernetesReferenceGrant`).

## Deploy

### Console

Open the deployment store, find **UDP Route on Kubernetes**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and three spec steps: **Namespace** (immutable), then **Gateways** (the `parentRefs` Gateways to attach to), then **Rules** (the route rules and their destination Services and weights). Start from the **DNS Forwarding** or **Game Server** preset in the [Presets](#presets) tab for a directly deployable configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesUdpRoute
metadata:
  name: my-udp-route
  org: acme-corp
  env: prod
spec:
  namespace:
    value: prod-apps
  parentRefs:
    - name: prod-gateway
      sectionName: udp-dns
  rules:
    - backendRefs:
        - name: coredns
          port: 53
```

```shell
planton apply -f udp-route.yaml
```

This creates a UDPRoute in `prod-apps` that attaches to the `udp-dns` listener of `prod-gateway` and forwards every datagram arriving on that listener to the `coredns` Service on port 53.

## Key Configuration

These are the most important decisions when configuring a UDPRoute. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Namespace** -- The `namespace` field is where the route lives. It is **immutable**, the default namespace for its backends, and the anchor for cross-namespace rules: parent Gateways and backend Services in this namespace attach without a `ReferenceGrant`; those elsewhere require one. Reference an existing `KubernetesNamespace` or type the name directly.

**Parent Gateways** -- The `parentRefs` (0-32) attach this route to Gateway listeners. Use `sectionName` to target one named UDP listener and/or `port` to pin a port -- the listener's port decides which datagrams arrive, since UDP does no matching. A parent in another namespace needs a `ReferenceGrant` there and a matching listener `allowedRoutes` policy. Each parent's `name` is a foreign key to `KubernetesGateway`: reference a Planton-managed Gateway (the route then deploys after it), or pass a literal name for a Gateway or ListenerSet created outside Planton.

**Rules** -- A UDPRoute has 1-16 `rules`, each with an optional `name` and `backendRefs` (1-16). Each backend's `name` is a foreign key to `KubernetesService` (reference a Planton-managed Service or pass a literal in-cluster name), with a `port` and an optional `weight` for traffic splitting (weight `0` drains a backend). There is no matching -- every datagram on the listener is forwarded to the rule's backends, and there is no session or connection to preserve.

## Outputs and Dependencies

### What This Component Consumes

This component takes foreign-key references to a `KubernetesNamespace` (via `spec.namespace`), to `KubernetesGateway` (each `parentRefs` entry's `name`), and to `KubernetesService` (each backend's `name`), so an InfraChart deploys those targets before the route and the resource graph carries the edges. Literal names cover targets created outside Planton; cross-namespace references require a `KubernetesReferenceGrant`.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `route_name` | Name of the created UDPRoute (equals `metadata.name`) | Orders the route after its Gateway and backends in an InfraChart |
| `namespace` | The resolved namespace the route was created in | Same-namespace / ReferenceGrant rules for its parent and backend references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**DNS forwarding** -- Forward datagrams arriving on a Gateway's UDP listener (port 53) to an in-cluster DNS Service. Start from the **DNS Forwarding** preset.

**Game server** -- Expose a UDP-speaking game server through the Gateway, optionally splitting load across instances by weight. Start from the **Game Server** preset.

## Works With

- **KubernetesGatewayApiCrds** -- installs the Gateway API CRDs (v1.6.0+ standard channel carries UDPRoute); deploy first (prerequisite).
- **KubernetesGateway** -- the Gateway whose UDP listener this route attaches to (`parentRefs`); install first.
- **KubernetesNamespace** -- the namespace (`spec.namespace`) the route runs in.
- **KubernetesReferenceGrant** -- authorizes cross-namespace parent or backend references from this route.
- **KubernetesService** -- the backend workloads (`backendRefs`) that receive forwarded datagrams.
