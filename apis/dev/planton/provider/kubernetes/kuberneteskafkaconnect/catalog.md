# Kafka Connect

Deploys a Strimzi **Kafka Connect** cluster — the worker pool that runs connector plugins between Kafka and external systems. Connectors themselves are declared separately as KubernetesKafkaConnector resources against this cluster's `connect_name`. Plugin delivery is the star decision: the stock image carries ONLY the MirrorMaker 2 connector classes; every real integration needs a prebuilt `image`, OCI image-volume `plugins`, or an operator `build`.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **KafkaConnect** — the Strimzi custom resource that runs the worker Deployment and the Connect REST API Service (`<name>-connect-api`)
- **Internal Connect storage topics** on the target Kafka (config / status / offset) — default names derive from `metadata.name` and must stay unique among Connect-protocol workloads sharing the cluster
- **Namespace** (optional) — created when `create_namespace` is true; a Strimzi operator must watch that namespace or the resource is accepted and silently never reconciled

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Kafka Family Side

- A **Kubernetes Kafka** cluster the workers talk to (bootstrap, TLS trust, authentication)
- A clear **plugin-delivery** plan — stock image alone cannot run Debezium, JDBC, S3, or other real connectors
- Unique **group identity** when multiple Connect clusters share one Kafka (colliding group IDs / storage topics corrupt each other's state)

## Deploy

### Console

Open the deployment store, find **Kafka Connect**, and click **Deploy**. The creation wizard walks you through placement, the Kafka connection, plugin delivery (the star step), worker count, group identity, worker config, sizing, scheduling, and metrics. Start from the **Debezium prebuilt image** preset in the [Presets](#presets) tab for a real integration posture.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKafkaConnect
metadata:
  name: cdc-connect
  org: acme-corp
  env: prod
spec:
  namespace:
    value: kafka
  create_namespace: false
  bootstrap_servers:
    valueFrom:
      kind: KubernetesKafka
      name: event-bus
      fieldPath: status.outputs.internal_bootstrap_endpoint
  image: quay.io/debezium/connect:2.5
  replicas: 2
```

```shell
planton apply -f cdc-connect.yaml
```

Then declare pipes with **Kubernetes Kafka Connector** against this cluster's `connect_name`.

### InfraChart

Wire bootstrap and trust from the sibling Kafka cluster:

```yaml
spec:
  namespace:
    value: kafka
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
  image: quay.io/debezium/connect:2.5
```

## Presets

| Preset | Use when |
| --- | --- |
| Minimal stock Connect | Explore the Connect surface; stock image = MirrorMaker connectors only |
| Debezium prebuilt image | Run real CDC/integrations from a known Connect image |
| Operator-built image | Let Strimzi build a Connect image from declared Maven/URL artifacts |

## Outputs

| Output | Purpose |
| --- | --- |
| `namespace` | Namespace the workers run in |
| `connect_name` | Value KubernetesKafkaConnector binds via `connect_cluster` |
| `rest_api_service_name` | Connect REST API Service name |
| `rest_api_endpoint` | In-cluster REST endpoint for read-only inspection |

## Related Components

- **Kubernetes Kafka** — the cluster workers connect to
- **Kubernetes Kafka Connector** — individual pipes on this cluster
- **Kubernetes Kafka MirrorMaker 2** — purpose-built mirroring (also Connect-protocol)
- **Kubernetes Kafka UI** — observe Connect status alongside topics
