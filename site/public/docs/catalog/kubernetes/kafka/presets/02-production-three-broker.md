---
title: "Production three broker preset"
description: "The standard production posture: a 3-node controller pool (the KRaft quorum tolerates one loss) plus a 3-node broker pool on JBOD storage, RF-3 / min-ISR-2 replication (one broker can die without..."
type: "preset"
rank: "02"
presetSlug: "02-production-three-broker"
componentSlug: "kafka"
componentTitle: "Kafka"
provider: "kubernetes"
icon: "package"
order: 2
---

# Production three broker preset

The standard production posture: a 3-node controller pool (the KRaft
quorum tolerates one loss) plus a 3-node broker pool on JBOD storage,
RF-3 / min-ISR-2 replication (one broker can die without losing
acknowledged writes from acks=all producers), SCRAM-over-TLS and
mutual-TLS listeners, simple ACL authorization with a real admin super
user, zone-spread rack awareness, JMX metrics plus the consumer-lag
exporter, explicit CA validity/renewal policy, a Sunday-night
maintenance window for renewal-triggered rolls, and a fixed 4g heap.

For teams running Kafka as shared infrastructure with multiple
applications and real durability requirements. The trade-offs are
cost and ceremony: six pods plus per-broker JBOD volumes, and — the
one to plan for — authorization is ON, so every client needs a
KubernetesKafkaUser with matching ACLs before it can produce or
consume (the `User:CN=platform-admin` super user is the mTLS admin
identity that can always operate).

Declare topics and users as KubernetesKafkaTopic /
KubernetesKafkaUser resources in this cluster's namespace; workloads
connect at the exported `internal_bootstrap_endpoint` and trust the
exported cluster CA Secret.

See [02-production-three-broker.yaml](./02-production-three-broker.yaml)
for the manifest.
