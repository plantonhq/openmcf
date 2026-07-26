# Kubernetes Kafka MirrorMaker 2

Declares a MirrorMaker 2 replication engine reconciled by the Strimzi
cluster operator — continuous, offset-aware mirroring of topics and
consumer groups from one or more source Kafka clusters into one
target cluster. The migration on-ramp off Confluent, MSK, or any
self-hosted Kafka: records, topic configuration, and consumer
positions replicate continuously, then consumers cut over with their
offsets intact.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **KafkaMirrorMaker2** (`kafka.strimzi.io/v1`, named
  `metadata.name`) — one `target` cluster connection (plus the
  engine's group identity and storage topics) and one `mirrors` entry
  per source: connection, topic/group scope patterns, and per-mirror
  MirrorSourceConnector / MirrorCheckpointConnector tuning
- **JMX metrics ConfigMap** (optional, `metrics.enabled`) — the
  canonical Strimzi rules (`<name>-mm2-metrics`), wired as the
  deployment's `metricsConfig`

The Strimzi operator reconciles these into MirrorMaker 2 worker pods
running the mirror connectors against the target cluster.

## Prerequisites

- The Strimzi cluster operator on the cluster
  (KubernetesStrimziKafkaOperator), watching this namespace
- A reachable target cluster (usually the KubernetesKafka sibling by
  reference) and source clusters (usually external literal
  bootstraps)
- Source credentials in Secrets — SCRAM passwords, Confluent API
  secrets, or client certificates, referenced by Secret name

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKafkaMirrorMaker2
metadata:
  name: dev-mirror
spec:
  namespace:
    value: dev-kafka
  replicas: 1
  target:
    bootstrap_servers:
      value: dev-kafka-kafka-bootstrap.dev-kafka.svc.cluster.local:9092
    config:
      config.storage.replication.factor: "-1"
      offset.storage.replication.factor: "-1"
      status.storage.replication.factor: "-1"
  mirrors:
    - source:
        alias: legacy
        bootstrap_servers:
          value: legacy-kafka.example.internal:9092
      topics_pattern: "orders.*"
```

Topics arrive prefixed with the source alias (`legacy.orders...`)
under the default replication policy; set IdentityReplicationPolicy
on both connectors of a mirror to keep original names — the usual
migration posture.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the deployment runs in |
| `mirrormaker_name` | The deployment's name (`metadata.name`) |
| `rest_api_endpoint` | In-cluster Connect REST endpoint of the underlying engine (port 8083) — read-only inspection of mirror connector status |

## Next Steps

Enable the checkpoint connector's group-offset sync
(`sync.group.offsets.enabled: "true"`) and `auto_restart` on both
connectors for a migration that survives transient source outages.
Size `source_connector.tasks_max` against the source's partition
volume and `replicas` for the task spread. Keep the target alias
distinct from every source alias (spec-enforced; the target alias
defaults to `target`), and give each MirrorMaker 2 instance sharing a
target cluster a distinct `metadata.name` — the engine's group
identity derives from it.
