---
title: "Kafka Connector"
description: "Kafka Connector deployment documentation"
icon: "package"
order: 100
componentName: "kuberneteskafkaconnector"
---

# Kubernetes Kafka Connector

Declares one data pipe on the Strimzi KafkaConnector custom resource,
run by a KubernetesKafkaConnect cluster's workers: a source connector
streams an external system into Kafka (Debezium CDC, file tails, SaaS
APIs), a sink connector streams Kafka topics out. The connector's
lifecycle is fully declarative — desired state, automatic restarts
with back-off, and annotation-triggered offset inspection and
override.

## What Gets Created

- **KafkaConnector** (`kafka.strimzi.io/v1`, named `metadata.name`) —
  in the Connect cluster's own namespace, bound to the cluster
  through the `strimzi.io/cluster` label (rendered from
  `connect_cluster`)
- **The running connector itself** — created by the Connect cluster's
  operator through the Connect group (not by the IaC modules); its
  consumer group is `connect-<name>`

## Prerequisites

- A KubernetesKafkaConnect cluster in the SAME namespace (the
  placement contract: a connector elsewhere is accepted and silently
  never reconciled)
- The connector class on that cluster's workers — via the stock
  image, a prebuilt `image`, OCI `plugins`, or a `build`
- For `${secrets:...}` config references: the
  KubernetesSecretConfigProvider enabled through the Connect
  cluster's `config` (`config.providers` entries)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesKafkaConnector
metadata:
  name: first-pipe
spec:
  namespace:
    value: dev-kafka
  connect_cluster:
    value: dev-connect
  connector_class: org.apache.kafka.connect.mirror.MirrorSourceConnector
  tasks_max: 1
  config:
    source.cluster.alias: src
    target.cluster.alias: dev
    source.cluster.bootstrap.servers: dev-kafka-kafka-bootstrap.dev-kafka.svc.cluster.local:9092
    target.cluster.bootstrap.servers: dev-kafka-kafka-bootstrap.dev-kafka.svc.cluster.local:9092
    topics: orders
    replication.factor: "1"
    offset-syncs.topic.replication.factor: "1"
```

The stock image carries only the MirrorMaker 2 connector classes
(Kafka's FileStream examples are not on the distribution's
classpath), so a MirrorSource self-mirror is the pipe that reaches
RUNNING with zero plugin machinery — records produced to `orders`
appear on `src.orders`. Real integration classes (Debezium, S3,
search) arrive through the Connect cluster's `image`, `plugins`, or
`build` arms.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the KafkaConnector lives in (the Connect cluster's namespace) |
| `connector_name` | Connector name inside the Connect cluster (`metadata.name`) — what the REST API and `connect-<name>` consumer groups key off |

## Next Steps

Enable `auto_restart` on every production pipe — transient
source-system outages otherwise leave the connector FAILED until a
human intervenes. Move credentials out of `config` literals and into
`${secrets:<namespace>/<secret>:<key>}` config-provider references.
For operational offset work, declare `list_offsets` /
`alter_offsets` ConfigMap targets and trigger the verbs with the
`strimzi.io/connector-offsets: list|alter` annotation (`alter`
requires the connector `stopped`) — the replay and
skip-poison-record mechanism.
