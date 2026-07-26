# Kubernetes Kafka

Declares a KRaft-mode Apache Kafka cluster reconciled by the Strimzi
cluster operator — node pools for controllers and brokers, listeners
with TLS and SCRAM or mutual-TLS authentication across seven exposure
types, ACL authorization, operator-managed certificates, Cruise
Control rebalancing, and Prometheus metrics. One resource per Kafka
cluster; workloads connect at the exported bootstrap endpoint, and
topics and users compose as KubernetesKafkaTopic / KubernetesKafkaUser
resources.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **KafkaNodePool** (`kafka.strimzi.io/v1`, one per `node_pools`
  entry) — the machines: roles (controller/broker), replica count,
  storage (persistent-claim, ephemeral, or JBOD), sizing, scheduling;
  bound to the cluster by the `strimzi.io/cluster` label
- **Kafka** (`kafka.strimzi.io/v1`, named `metadata.name`) — the
  cluster: listeners, broker config, authorization, entity operators,
  Cruise Control, exporter, CAs, rack awareness, maintenance windows
- **JMX metrics ConfigMap** (optional, `metrics.enabled`) — the
  canonical Strimzi Prometheus rules (`<name>-kafka-metrics`), wired
  as the cluster's `metricsConfig`

The Strimzi operator reconciles these into broker/controller pods,
Services, and certificate Secrets — including the exported
`<name>-kafka-bootstrap` Service and `<name>-cluster-ca-cert` Secret.

## Prerequisites

- The Strimzi cluster operator on the cluster
  (KubernetesStrimziKafkaOperator) — it must watch this cluster's
  namespace (the default posture watches its OWN namespace: install
  the operator there, or widen its watch)
- For `loadbalancer` listeners: a cloud LoadBalancer controller (the
  annotations surface on `configuration.bootstrap` / `brokers`)
- For `ingress` listeners: an NGINX ingress controller with SSL
  passthrough enabled, plus a bootstrap host and one host per broker
- For custom listener certificates: a cert-manager-issued Secret
  referenced by `broker_cert_chain_and_key`

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKafka
metadata:
  name: dev-kafka
spec:
  namespace:
    value: kafka
  node_pools:
    - name: dual
      roles: [controller, broker]
      replicas: 1
      storage:
        size: 10Gi
  listeners:
    - name: plain
      port: 9092
      tls: false
  config:
    default.replication.factor: "1"
    min.insync.replicas: "1"
    offsets.topic.replication.factor: "1"
    transaction.state.log.replication.factor: "1"
    transaction.state.log.min.isr: "1"
```

The operator forms the single-node KRaft cluster; workloads connect at
the exported `internal_bootstrap_endpoint`
(`dev-kafka-kafka-bootstrap.kafka.svc.cluster.local:9092`).

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster runs in |
| `cluster_name` | Kafka resource name (equals `metadata.name`) — the `strimzi.io/cluster` binding value |
| `bootstrap_service_name` | Internal bootstrap Service (`<name>-kafka-bootstrap`) |
| `internal_bootstrap_endpoint` | In-cluster `bootstrap.servers` value for the first internal listener; empty when only external listeners are declared |
| `cluster_ca_cert_secret_name` | Cluster CA Secret (`<name>-cluster-ca-cert`, key `ca.crt`) for TLS client truststores |

## Next Steps

Declare topics and users as KubernetesKafkaTopic /
KubernetesKafkaUser resources in this cluster's namespace — the entity
operators (enabled by default) reconcile them into real topics and
credential Secrets. Move to a 3-node controller pool and a 3-node
broker pool with `default.replication.factor: "3"` and
`min.insync.replicas: "2"` for production — that shape survives one
broker loss without losing acknowledged writes (producers using
acks=all). Enable `authorization: simple` with real admin principals
in `super_users` before the first client depends on open access.
