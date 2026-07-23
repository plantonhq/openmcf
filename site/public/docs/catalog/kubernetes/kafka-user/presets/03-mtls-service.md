---
title: "mTLS Service (Producer + Consumer, Quotas)"
description: "This preset declares a full service identity on mutual TLS: the user operator issues a client certificate from the cluster's clients CA into the user's Secret, the ACLs cover both directions (produce..."
type: "preset"
rank: "03"
presetSlug: "03-mtls-service"
componentSlug: "kafka-user"
componentTitle: "Kafka User"
provider: "kubernetes"
icon: "package"
order: 3
---

# mTLS Service (Producer + Consumer, Quotas)

This preset declares a full service identity on mutual TLS: the user operator issues a client certificate from the cluster's clients CA into the user's Secret, the ACLs cover both directions (produce to one topic, consume another under its own group), and broker-enforced quotas cap how hard the service can push the cluster. The certificate renews with the clients CA; the manifest carries no secret material.

## When to Use

- A service that both produces and consumes, on a cluster with tls-auth listeners
- Certificate-based identity requirements — the principal is the certificate subject (`CN=payment-service`), not a password
- Multi-tenant clusters where per-client throughput caps matter

## Key Configuration Choices

- **`type: tls`** -- must match a tls-auth listener; the generated Secret carries `user.crt` / `user.key` in PEM form plus `user.p12` / `user.password` keystore forms. For certificates issued OUTSIDE the cluster, use `tls-external` instead — no Secret is generated for those
- **Producer + consumer ACLs on one principal** -- `Write`/`Describe` on `payment-events`, `Read`/`Describe` on `order-events`, and `Read` on the `payment-service` group; each direction is its own rule, reviewable independently
- **`producerByteRate` / `consumerByteRate`** -- bytes-per-second ceilings (per broker) before the brokers throttle the client
- **`requestPercentage`** -- caps the share of broker request-handler time this client may consume
- **ACLs enforced only with cluster-side `simple` authorization** -- declared against a cluster without it, the rules are rejected at reconcile and the resource reports NotReady; quotas apply regardless

## Values to Adapt

| Value | Description | Where to Find |
|---|---|---|
| `kafka` (namespace) | The Kafka cluster's own namespace | The KubernetesKafka resource's `namespace` |
| `my-kafka` (kafkaCluster) | The Kafka cluster's name | The KubernetesKafka resource's `metadata.name` or its `cluster_name` output |
| `payment-events` / `order-events` | The topics this service writes and reads | The KubernetesKafkaTopic resources' `topic_name` outputs |
| `payment-service` (group) | The consumer group id | The client application's `group.id` configuration |
| Quota numbers | Throughput and request-time ceilings | Your capacity planning for the cluster |

## Related Presets

- **01-producer** -- Write-only SCRAM identity with a prefix ACL
- **02-consumer** -- Read-only SCRAM identity
