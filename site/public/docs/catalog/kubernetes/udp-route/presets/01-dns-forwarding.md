---
title: "UDP DNS Forwarding"
description: "The most common UDPRoute: forward DNS queries arriving on a Gateway's UDP listener to a DNS Service on port 53, split across a weighted pair of backends. A UDP route has no matching at all -- the..."
type: "preset"
rank: "01"
presetSlug: "01-dns-forwarding"
componentSlug: "udp-route"
componentTitle: "UDP Route"
provider: "kubernetes"
icon: "package"
order: 1
---

# UDP DNS Forwarding

The most common UDPRoute: forward DNS queries arriving on a Gateway's UDP
listener to a DNS Service on port 53, split across a weighted pair of backends.
A UDP route has no matching at all -- the listener's port selects the traffic,
and the route forwards every datagram to the rule's backends.

## When to Use

- You expose a DNS resolver (CoreDNS, BIND, a split-horizon forwarder) through
  a Gateway.
- You want to shift a fraction of queries to a second resolver (an upgrade
  canary or a standby) using backend weights.

## Key Configuration Choices

- **`parentRefs[].name`** -- attaches the route to the Gateway; add `sectionName` to target the UDP listener (the listener's port determines which datagrams arrive). The name is a foreign key: write `value: <literal>` for an existing Gateway, or `valueFrom:` (kind `KubernetesGateway`, fieldPath `status.outputs.gateway_name`) so the route deploys after its Planton-managed Gateway.
- **`backendRefs[].name`** -- a foreign key like the parent ref: `value:` for a literal Service name, `valueFrom:` (kind `KubernetesService`, fieldPath `status.outputs.service_name`) for a Planton-managed backend.
- **`backendRefs[].weight`** -- datagrams are distributed as `weight / sum(weights)`; 90/10 sends ~10% of queries to the secondary resolver.
- **`backendRefs[].port` (`53`)** -- the standard DNS port on the backend Service.

## Prerequisites

- The Gateway API CRDs are installed (`KubernetesGatewayApiCrds`). UDPRoute is
  part of the standard channel as of Gateway API v1.6.
- The `Gateway` referenced in `parentRefs` exists (`KubernetesGateway`) with a
  listener of protocol `UDP`.
- The target namespace exists (`KubernetesNamespace`).
- Both DNS Services exist in the route's namespace.

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `<gateway-name>` | Name of the `KubernetesGateway` this route attaches to (inside `name.value`, or switch to `valueFrom`). |
| `<primary-dns-service>` | DNS Service receiving the majority of queries. |
| `<secondary-dns-service>` | DNS Service receiving the weighted slice of queries. |

Tune the `weight` values to control the split; set `spec.namespace.value` to
your namespace, or replace it with a `valueFrom` reference to a
`KubernetesNamespace`.
