---
title: "Karapace Schema Registry"
description: "Karapace Schema Registry deployment documentation"
icon: "package"
order: 100
componentName: "kuberneteskarapace"
---

# Karapace Schema Registry

Deploy a [Karapace](https://github.com/Aiven-Open/karapace) schema registry — the Apache-2.0, Confluent Schema Registry API-compatible engine from Aiven. Producers and consumers register and fetch Avro, JSON Schema, and Protobuf schemas through the standard REST API; existing Confluent SR client libraries, Connect converters, and tooling work unchanged.

Schemas live in a compacted Kafka topic (`_schemas` by convention) on the connected cluster — the same Kafka-native architecture Confluent SR uses. Multiple replicas coordinate leadership through a consumer group; followers forward writes to the leader.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Deployment + Service** — the registry engine (Karapace has no official Helm chart; the module owns the manifests)
- **Optional second Deployment** — the Kafka REST-proxy role (`<name>-rest`) when `rest_proxy.enabled` is true
- **Schemas topic** — created on first start at the replication factor you declare (applies at creation only)

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Kafka Cluster Side

- A **Kafka cluster** reachable at the bootstrap address you declare — reference a **Kubernetes Kafka** resource for in-cluster wiring.
- A **listener whose security posture matches** your declared `security_protocol` — SSL forms need TLS trust material; SASL forms need credentials. The SASL password is exactly one of a Secret reference (recommended — reference a **Kubernetes Kafka User**) or a declared value the module materializes into a Secret.

## Deploy

### Console

Open the deployment store, find **Karapace**, and click **Deploy**. The creation wizard walks you through namespace placement, the Kafka connection (with the SASL password XOR), replica leadership, registry behavior, optional REST proxy, server TLS, HTTP authentication, runtime, and pod sizing. Start from the **Minimal** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKarapace
metadata:
  name: dev-registry
  org: acme-corp
  env: dev
spec:
  namespace:
    value: dev-registry
  create_namespace: true
  kafka:
    bootstrap_servers:
      value: dev-kafka-kafka-bootstrap.dev-kafka.svc.cluster.local:9092
```

```shell
planton apply -f karapace-registry.yaml
```

### InfraChart

Wire the bootstrap endpoint from the Kafka cluster itself:

```yaml
spec:
  namespace:
    value: kafka
  kafka:
    bootstrap_servers:
      valueFrom:
        kind: KubernetesKafka
        name: event-bus
        fieldPath: status.outputs.internal_bootstrap_endpoint
```

## Key Configuration

**Namespace placement** — unlike a KafkaTopic, the registry has no same-namespace contract with the Kafka cluster; it connects over the network. The namespace locks after creation (moving a Deployment is a delete-and-recreate).

**Replication factor is a one-shot door** — it applies when the registry creates the `_schemas` topic. Start production registries at RF 3; graduating later means Kafka topic reassignment, not a field edit.

**Server TLS + multiple replicas** — followers forward writes to the leader at its pod IP, which DNS-name certificates do not cover. Pair server TLS with one replica, or terminate TLS at an Ingress/Gateway in front of plain-HTTP replicas.

**HTTP authentication** — absent means the API is open to anyone who can reach the Service. Configure Basic (authfile Secret) or OIDC before shared exposure.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | Where the registry Deployment runs |
| `spec.kafka.bootstrap_servers` | KubernetesKafka (`status.outputs.internal_bootstrap_endpoint`) | The cluster storing schemas |
| `spec.kafka.tls.ca_secret_name` | KubernetesKafka (`status.outputs.cluster_ca_cert_secret_name`) | TLS trust for SSL listeners |
| `spec.kafka.sasl.password_secret` | KubernetesKafkaUser (`status.outputs.secret_name`) | SASL credentials from an operator-generated Secret |
| `spec.server_tls.secret_name` | KubernetesCertificate (`status.outputs.secret_name`) | HTTPS serving certificate |
| `spec.http_authentication.basic.secret_name` | KubernetesSecret (`status.outputs.secret_name`) | Authfile for HTTP Basic |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `endpoint` | In-cluster registry URL — the `schema.registry.url` for clients | Producers, consumers, Connect converters |
| `rest_proxy_endpoint` | REST-proxy URL (empty when the role is disabled) | HTTP produce/consume clients |
| `schemas_topic` | The Kafka topic storing schemas | Monitoring, backup planning |
| `service_name` / `namespace` | Service identity handles | Ingress/Gateway composition |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Minimal** — one replica, plaintext bootstrap, everything else on defaults. Start from the **Minimal** preset.

**Production HA** — two replicas, SASL_SSL, schemas topic at RF 3. Start from the **Production HA** preset.

**REST proxy and TLS** — single replica with server TLS, HTTP Basic auth, and the REST-proxy role enabled. Start from the **REST Proxy and TLS** preset.

## Works With

- **Kubernetes Kafka** — the cluster storing schemas; wire `bootstrap_servers` from its outputs.
- **Kubernetes Kafka User** — operator-generated SASL credentials via `password_secret`.
- **Kubernetes Certificate** — cert-manager-issued TLS for `server_tls`.
- **Workloads and Connect** — point `schema.registry.url` at the `endpoint` output.
