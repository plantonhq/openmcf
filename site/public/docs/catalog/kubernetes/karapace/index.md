---
title: "Karapace"
description: "Karapace deployment documentation"
icon: "package"
order: 100
componentName: "kuberneteskarapace"
---

# Kubernetes Karapace

Declares a Karapace schema registry — Aiven's Apache-2.0,
Confluent-API-compatible registry. Producers and consumers register
and fetch Avro, JSON Schema and Protobuf schemas through the standard
Schema Registry REST API (existing Confluent SR clients work
unchanged), with compatibility enforcement between versions. Schemas
live in a compacted Kafka topic on the connected cluster — no
database. An optional REST-proxy role serves produce/consume over
HTTP from the same engine.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **Registry Deployment + Service** (`<name>`) — the schema-registry
  role on upstream's image (`ghcr.io/aiven-open/karapace`, pinned
  6.2.1), configured through `KARAPACE_*` environment variables; each
  pod advertises its own POD IP for leader/follower forwarding
- **REST-proxy Deployment + Service** (`<name>-rest`, optional) — the
  same engine with the role flags flipped, wired to the registry and
  the same Kafka cluster
- **SASL password Secret** (`<name>-sasl`, only when a literal
  `password` is declared) — the module materializes it; the
  credential never rides the pod spec as plaintext

## Prerequisites

- A reachable Kafka cluster (`kafka.bootstrap_servers` — a
  KubernetesKafka reference or an external literal address); the
  registry creates its schemas topic on first start
- For SSL/SASL_SSL: the CA Secret to trust (a KubernetesKafka's
  cluster CA by reference); for SASL: credentials (a
  KubernetesKafkaUser Secret by reference, or a declared password)
- For `server_tls`: a certificate Secret (a KubernetesCertificate's
  output — the cert-manager seam), with `replicas: 1`

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKarapace
metadata:
  name: dev-registry
spec:
  namespace:
    value: dev-kafka
  kafka:
    bootstrap_servers:
      value: dev-kafka-kafka-bootstrap.dev-kafka.svc.cluster.local:9092
```

The registry starts, creates `_schemas` on the cluster, and serves
the SR API at the exported endpoint
(`http://dev-registry.dev-kafka.svc.cluster.local:8081`).

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the registry runs in |
| `service_name` | The registry Service (`<name>`) |
| `endpoint` | In-cluster endpoint — the `schema.registry.url` value for clients |
| `rest_proxy_endpoint` | REST-proxy endpoint; empty when the role is not enabled |
| `schemas_topic` | The Kafka topic storing the schemas |

## Next Steps

Set `registry.replication_factor: 3` before production traffic — the
upstream default of 1 makes the schemas topic a single-broker
data-loss risk, and changing it later is a Kafka topic reassignment,
not a field edit. Move to 2 replicas for availability (leader
election is automatic), add `http_authentication` before any shared
exposure, and wire the exported `endpoint` into producers, consumers,
Connect converters, and a KubernetesKafkaUi console's
`schema_registry.url`.
