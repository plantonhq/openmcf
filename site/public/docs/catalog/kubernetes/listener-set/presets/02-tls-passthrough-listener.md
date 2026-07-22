---
title: "TLS Passthrough Listener"
description: "Attach a TLS passthrough listener to a shared Gateway through a ListenerSet. The listener forwards encrypted connections untouched -- the backend terminates TLS itself -- so the team exposes an..."
type: "preset"
rank: "02"
presetSlug: "02-tls-passthrough-listener"
componentSlug: "listener-set"
componentTitle: "Listener Set"
provider: "kubernetes"
icon: "package"
order: 2
---

# TLS Passthrough Listener

Attach a TLS passthrough listener to a shared Gateway through a ListenerSet.
The listener forwards encrypted connections untouched -- the backend terminates
TLS itself -- so the team exposes an end-to-end-encrypted service (a database, an
mTLS API) through the platform Gateway without handing its certificate to the
edge.

## When to Use

- A team needs a passthrough entry point on a centrally managed Gateway for a
  backend that holds its own certificate.
- You route by SNI hostname via `TLSRoute`s attached to this listener.
- The platform team does not want to add per-team listeners to the Gateway by
  hand.

## Key Configuration Choices

- **`parentRef.name`** -- the Gateway these listeners merge into. It is a foreign key: this preset uses `value:` with a literal name, which fits a Gateway created outside Planton; use `valueFrom:` (kind `KubernetesGateway`, fieldPath `status.outputs.gateway_name`) to reference a Planton-managed Gateway and deploy in dependency order.
- **`protocol` (`TLS`) + `tls.mode` (`Passthrough`)** -- a TLS listener must declare its mode; Passthrough forwards the encrypted stream and needs no `certificateRefs`.
- **`hostname`** -- the SNI hostname the listener accepts; TLSRoutes narrow it further with their own `hostnames`.
- **`allowedRoutes`** (`Same` + `TLSRoute`) -- only TLSRoutes from the ListenerSet's own namespace may attach.

## Prerequisites

- The Gateway API CRDs are installed (`KubernetesGatewayApiCrds`, standard
  channel as of v1.6 -- ListenerSet joined the standard channel in v1.5).
- The parent Gateway exists and opts in to ListenerSet attachment
  (`allowedListeners.namespaces.from` set to `All`, `Selector`, or `Same` --
  Gateways allow none by default).
- The team namespace exists (`KubernetesNamespace`).
- The backend that terminates TLS is reachable via a `KubernetesTlsRoute`
  attached to this listener.

## Placeholders to Replace

| Placeholder | Description |
|-------------|-------------|
| `<gateway-name>` | Literal name of the Gateway this ListenerSet attaches to (inside `name.value`, or switch to `valueFrom`). |
| `db.team-a.example.com` | The SNI hostname the listener accepts (a literal example value -- replace with your real host). |

Set `spec.namespace.value` and `parentRef.namespace` to your team and Gateway
namespaces, or replace the namespace with a `valueFrom` reference to a
`KubernetesNamespace`.
