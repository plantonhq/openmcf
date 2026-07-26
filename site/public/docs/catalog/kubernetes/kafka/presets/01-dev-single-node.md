---
title: "Dev single node preset"
description: "The smallest declarable Kafka that actually serves: one dual-role node (KRaft controller + broker in a single pod), one plaintext internal listener, and single-node replication settings (RF 1..."
type: "preset"
rank: "01"
presetSlug: "01-dev-single-node"
componentSlug: "kafka"
componentTitle: "Kafka"
provider: "kubernetes"
icon: "package"
order: 1
---

# Dev single node preset

The smallest declarable Kafka that actually serves: one dual-role node
(KRaft controller + broker in a single pod), one plaintext internal
listener, and single-node replication settings (RF 1 everywhere — the
cluster cannot place a second replica anywhere).

For developers and CI who need a real Kafka wire protocol without
production ceremony. The trade-offs are total: no authentication, no
TLS, no durability — every broker restart is a data-availability
event, and losing the volume loses the data. Nothing here transfers to
production except the shape of the manifest; graduate to the
production preset's separate pools and RF-3 settings before anything
depends on the data.

Topics and users still compose the normal way — declare
KubernetesKafkaTopic resources in this cluster's namespace (the entity
operators are enabled by default).

See [01-dev-single-node.yaml](./01-dev-single-node.yaml) for the
manifest.
