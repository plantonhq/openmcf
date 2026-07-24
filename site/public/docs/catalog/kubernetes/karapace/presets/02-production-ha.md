---
title: "Production HA preset"
description: "The production posture: two replicas with automatic leader election (followers forward writes to the leader at its pod-IP identity — no external coordination), a SASL_SSL connection to the Kafka..."
type: "preset"
rank: "02"
presetSlug: "02-production-ha"
componentSlug: "karapace"
componentTitle: "Karapace"
provider: "kubernetes"
icon: "package"
order: 2
---

# Production HA preset

The production posture: two replicas with automatic leader election
(followers forward writes to the leader at its pod-IP identity — no
external coordination), a SASL_SSL connection to the Kafka cluster,
and the schemas topic created at replication factor 3.

The three decisions that matter:

- **`replication_factor: 3` is a one-shot door.** It applies when the
  registry CREATES the schemas topic; on an existing topic the field
  does nothing and the fix is a Kafka partition reassignment. Since
  the topic is the registry's entire state, start production
  registries here, never by graduating the minimal preset in place.
- **The explicit `SASL_SSL` protocol is spec-enforced** — declaring
  `sasl` under the default (PLAINTEXT) protocol would silently ignore
  the credentials, so the spec rejects the combination outright.
- **Two replicas, not three.** Replicas are availability, not write
  scaling — every write funnels through the leader. Two survives a
  pod loss; more mostly adds read capacity.

The API itself is still open at the Service — add
`http_authentication` (see the rest-proxy-and-tls preset for the
basic-auth shape) before any shared or external exposure. Serving TLS
directly is deliberately NOT combined with 2 replicas: follower
forwarding targets pod IPs, which DNS-name certificates do not
cover — terminate TLS at an Ingress/Gateway instead.

See [02-production-ha.yaml](./02-production-ha.yaml) for the
manifest.
