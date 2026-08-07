---
title: "Kafka MirrorMaker 2"
description: "Kafka MirrorMaker 2 deployment documentation"
icon: "package"
order: 100
componentName: "kuberneteskafkamirrormaker2"
---

# Kafka MirrorMaker 2

Continuously mirrors topics, records, and consumer-group positions from one or more **source** Kafka clusters into a Strimzi-managed **target**. The common story is migration and disaster recovery: point the workers at a running MSK, Confluent Cloud, or datacenter cluster, replicate into your Kubernetes Kafka, then cut consumers over with offsets intact. Each `mirrors` entry is one source; the target connection and the engine's group identity live once on the resource.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **KafkaMirrorMaker2** — the Strimzi custom resource that runs Connect-protocol mirror workers against the declared target and sources
- **Worker pods** — scaled by `replicas`, sized by optional resources/JVM, placed by node selector / tolerations / rack awareness
- **Internal Connect storage topics** on the target (config / status / offset) — default names derive from `metadata.name` and must stay unique among Connect-protocol workloads sharing the target
- **Namespace** (optional) — created when `create_namespace` is true; the natural home is the target Kafka cluster's own namespace

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Kafka Family Side

- A **Kubernetes Kafka** target cluster (and a Strimzi operator watching the workers' namespace)
- Network reachability from the workers to every source bootstrap address
- Matching TLS trust and authentication for each connection (target and per-source)
- A clear **replication-policy** choice: alias-prefixed topic names (default) versus `IdentityReplicationPolicy` when you want identical names on the target for a clean cutover

## Deploy

### Console

Open the deployment store, find **Kafka MirrorMaker 2**, and click **Deploy**. The creation wizard walks you through where the workers run, the target connection, the engine identity, the mirrors repeater (the star step), worker count, sizing, scheduling, and metrics. Start from the **Migrate from MSK** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKafkaMirrorMaker2
metadata:
  name: msk-migration
  org: acme-corp
  env: prod
spec:
  namespace:
    value: kafka
  create_namespace: false
  target:
    bootstrap_servers:
      valueFrom:
        kind: KubernetesKafka
        name: event-bus
        fieldPath: status.outputs.internal_bootstrap_endpoint
  mirrors:
    - source:
        alias: msk
        bootstrap_servers:
          value: b-1.mycluster.abc123.c2.kafka.us-east-1.amazonaws.com:9094
      topics_pattern: ".*"
```

```shell
planton apply -f msk-migration.yaml
```

Watch mirror connector status through the exported REST API endpoint; cut consumers over only after checkpoint lag is acceptable.

### InfraChart

Wire the target bootstrap from the sibling Kafka cluster so the migration graph cannot drift:

```yaml
spec:
  namespace:
    value: kafka
  target:
    bootstrap_servers:
      valueFrom:
        kind: KubernetesKafka
        name: event-bus
        fieldPath: status.outputs.internal_bootstrap_endpoint
    tls:
      trusted_certificates:
        - secret_name:
            valueFrom:
              kind: KubernetesKafka
              name: event-bus
              fieldPath: status.outputs.cluster_ca_cert_secret_name
  mirrors:
    - source:
        alias: msk
        bootstrap_servers:
          value: b-1.mycluster.abc123.c2.kafka.us-east-1.amazonaws.com:9094
```

## Presets

| Preset | Use when |
| --- | --- |
| Migrate from MSK | Lift topics and consumer positions from Amazon MSK into a Strimzi target |
| Migrate from Confluent Cloud | Same migration posture against a Confluent Cloud source |
| Active-passive DR | Continuous mirror into a standby cluster for disaster recovery |

## Outputs

| Output | Purpose |
| --- | --- |
| `namespace` | Namespace the workers run in |
| `mirrormaker_name` | Deployment name (`metadata.name`) |
| `rest_api_endpoint` | In-cluster Connect REST API for read-only mirror status |

## Related Components

- **Kubernetes Kafka** — the usual target cluster
- **Kubernetes Kafka Connect** — general-purpose Connect for non-mirror integrations
- **Kubernetes Kafka UI** — observe topics and consumer lag on both sides during cutover
