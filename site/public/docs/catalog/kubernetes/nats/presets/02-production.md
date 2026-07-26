---
title: "Production preset"
description: "A 3-server NATS cluster with JetStream on 20Gi file stores per server, authenticated clients, and Prometheus metrics. Three servers is the smallest count that keeps replicated (R3) streams available..."
type: "preset"
rank: "02"
presetSlug: "02-production"
componentSlug: "nats"
componentTitle: "NATS"
provider: "kubernetes"
icon: "package"
order: 2
---

# Production preset

A 3-server NATS cluster with JetStream on 20Gi file stores per server,
authenticated clients, and Prometheus metrics. Three servers is the
smallest count that keeps replicated (R3) streams available through a
pod loss — create availability-critical streams with `--replicas 3`
and JetStream's RAFT placement does the rest.

The auth model is deliberately least-privilege: `orders-service` may
only publish and subscribe under `orders.>` (plus `_INBOX.>` —
request/reply needs it), while `platform-admin` carries full access
for operations. Every password is module-generated into the
`nats-auth` Secret, one key per username; workloads mount their own
key and the manifest stays credential-free. Unauthenticated
connections are rejected outright.

`memory` deserves a note: the 2Gi limit leaves headroom for JetStream
metadata and route buffers — if you enable a `memory_store_max_size`
tier later, raise the limit comfortably above it.

Change first: switch to `accounts` when tenants must not see each
other's subjects (accounts are isolated namespaces, not just
permissions); put TLS on the client listener (`tls` with a
cert-manager Secret) before any traffic leaves the cluster network;
for external clients, set `service` to LoadBalancer with your cloud's
annotations or enable `websocket` behind first-class exposure kinds.

See [02-production.yaml](./02-production.yaml) for the manifest.
