# Kubernetes Kafka Connect

Declares a Kafka Connect cluster reconciled by the Strimzi cluster
operator — the pluggable integration engine that streams data between
Kafka and external systems (databases via Debezium CDC, object
stores, search indexes, SaaS APIs). Connector plugins arrive four
ways: the stock image, a prebuilt image, OCI image-volume mounts, or
an operator-driven image build pushed to your registry. Individual
pipes compose as KubernetesKafkaConnector resources against this
cluster.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **KafkaConnect** (`kafka.strimzi.io/v1`, named `metadata.name`) —
  the worker fleet: replicas, the Kafka connection (TLS trust +
  authentication), worker config, the plugin arm (image / OCI plugins
  / build), sizing, rack awareness; always annotated
  `strimzi.io/use-connector-resources: "true"` so connectors are
  managed declaratively
- **JMX metrics ConfigMap** (optional, `metrics.enabled`) — the
  canonical Strimzi Prometheus rules (`<name>-connect-metrics`),
  wired as the cluster's `metricsConfig`

The Strimzi operator reconciles these into Connect worker pods and
the `<name>-connect-api` REST Service — and, when `build` is set,
runs the Kaniko/Buildah image build first.

## Prerequisites

- The Strimzi cluster operator on the cluster
  (KubernetesStrimziKafkaOperator), watching this namespace
- A reachable Kafka cluster (`bootstrap_servers` — a KubernetesKafka
  reference or an external literal address)
- For the `plugins` arm: the Kubernetes ImageVolume feature enabled
  AND supported by the container runtime (workers fail to schedule
  without it)
- For the `build` arm: a registry the build pod can push to
  (`push_secret` names a docker-registry Secret in this namespace)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesKafkaConnect
metadata:
  name: dev-connect
spec:
  namespace:
    value: dev-kafka
  bootstrap_servers:
    value: dev-kafka-kafka-bootstrap.dev-kafka.svc.cluster.local:9092
  replicas: 1
  config:
    config.storage.replication.factor: "-1"
    offset.storage.replication.factor: "-1"
    status.storage.replication.factor: "-1"
```

The operator forms the worker group on the stock image; declare pipes
as KubernetesKafkaConnector resources in the same namespace. The
`"-1"` storage entries ride the broker default — Connect's own
default of 3 wedges on a single-broker dev cluster.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the Connect cluster runs in |
| `connect_name` | Connect cluster name (`metadata.name`) — the `strimzi.io/cluster` binding value for KubernetesKafkaConnector resources |
| `rest_api_service_name` | Connect REST API Service (`<name>-connect-api`) |
| `rest_api_endpoint` | In-cluster REST endpoint (port 8083) — read-only inspection; connector management stays declarative |

## Next Steps

Declare connectors as KubernetesKafkaConnector resources in this
cluster's namespace — the placement contract the operator enforces.
Graduate the plugin story deliberately: start on the stock image,
move to a prebuilt `image` when a vendor ships one, and reach for
`build` when you need an exact artifact set baked and pushed to your
own registry. In production, set the three storage replication
factors to `"3"`, size `resources` and `jvm` explicitly, and keep the
group identity fields unique per Connect cluster sharing a Kafka
cluster.
