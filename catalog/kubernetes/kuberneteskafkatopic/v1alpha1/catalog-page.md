# Kubernetes Kafka Topic

Declares one Kafka topic on the Strimzi KafkaTopic custom resource. The target cluster's topic operator reconciles it into a real topic — creating it, growing partitions (never shrinking them), and applying `config` entries declaratively; changes made to the topic outside the resource are reverted. Deleting the resource deletes the topic and its data.

## What Gets Created

- **KafkaTopic** — named after the resource, in the Kafka cluster's own namespace, bound to the cluster through the `strimzi.io/cluster` label
- **The Kafka topic itself** — created by the cluster's topic operator (not by the IaC modules); named `spec.topic_name` when set, otherwise `metadata.name`

## Prerequisites

- A Kafka cluster (**KubernetesKafka**) with its topic operator enabled (the default)
- The Strimzi operator watching the cluster's namespace (**KubernetesStrimziKafkaOperator**)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesKafkaTopic
metadata:
  name: order-events
spec:
  namespace:
    value: kafka # MUST be the Kafka cluster's own namespace
  kafkaCluster:
    value: my-kafka
  partitions: 6
  replicas: 3
  config:
    retention.ms: "604800000" # 7 days
```

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Where the KafkaTopic resource lives (the Kafka cluster's namespace) |
| `topic_name` | The actual Kafka topic name producers and consumers subscribe to |

## Next Steps

Point clients at this resource's `topic_name` output plus the KubernetesKafka cluster's `internal_bootstrap_endpoint` output; declare client identities and ACLs with **KubernetesKafkaUser**.
