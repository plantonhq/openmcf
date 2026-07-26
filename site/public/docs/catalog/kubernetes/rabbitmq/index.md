---
title: "RabbitMQ"
description: "RabbitMQ deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesrabbitmq"
---

# Kubernetes RabbitMQ

Declares a RabbitMQ cluster — the queue-messaging broker behind task
queues, work distribution, RPC and per-message routing (AMQP
0-9-1/1.0, MQTT and STOMP via plugins) — as a `RabbitmqCluster`
reconciled by the RabbitMQ Cluster Operator. One resource carries the
node count, per-node persistent storage, TLS from an existing
certificate Secret, RabbitMQ's own configuration vocabulary, and
placement. The operator generates the administrator credentials
itself — consumers read them from the exported Secret, and workloads
connect at the exported AMQP (5672) and management (15672) endpoints.

> **Queue, not log**: RabbitMQ delivers each message and forgets it on
> acknowledgement. For append-only event streaming and replayable
> history at scale, use KubernetesKafka — that boundary is the first
> decision, not a tuning knob.

## What Gets Created

- **Namespace** (optional) — created and owned when
  `create_namespace` is set
- **RabbitmqCluster** (`rabbitmq.com/v1beta1`, named `metadata.name`)
  — topology, storage, TLS, configuration, placement

The operator reconciles these into the StatefulSet
(`<name>-server`, one data PVC per pod), the client Service
(`<name>`, ClusterIP by default), the headless inter-node Service
(`<name>-nodes`), and the generated admin credentials Secret
(`<name>-default-user` — keys: username, password, host, port,
connection_string, ...).

## Prerequisites

- The RabbitMQ Cluster Operator on the cluster
  (KubernetesRabbitMqOperator) — a registry prerequisite, declared
  first; it watches ALL namespaces by default
- A StorageClass for the data volumes (most managed clusters provide
  a default; or reference a KubernetesStorageClass) — unless the
  cluster is `ephemeral`
- For `tls`: an existing kubernetes.io/tls Secret covering the
  `<name>.<namespace>.svc` DNS forms (a KubernetesCertificate
  reference works directly)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesRabbitMq
metadata:
  name: orders-mq
spec:
  namespace:
    value: messaging
  create_namespace: true
```

Every field has a working default: a single node on a 10Gi volume
(the cluster's default StorageClass), the operator's default
`-management` image, operator-default sizing. Workloads connect at
`orders-mq.messaging.svc.cluster.local:5672` with the credentials
from the `orders-mq-default-user` Secret (it even carries a
ready-made `connection_string`); the management UI is on port 15672.

## Sizing and Availability: Read This Before Production

Three facts deserve conscious ownership. **Replicas are odd, and a
one-way door**: quorum queues and the Raft-based metadata store
survive node loss with 3, 5 or 7 nodes (a 2-node cluster loses
availability when either node fails), and the operator does NOT
support scaling down — removed brokers strand their queue replicas,
so sizing down means migrating to a new cluster. **Memory requests
must equal limits**: RabbitMQ derives its memory high watermark from
the container memory LIMIT, so a gap between the two triggers flow
control at the wrong threshold. And **`spread_across_nodes` is the
production posture**: without it, co-located brokers make quorum
queues pointless against node loss — but a cluster with more replicas
than schedulable nodes sits Pending when it is on.

## Configuration

### Topology

`replicas` (default 1 — dev only) is the node count; production runs
an odd count. `spread_across_nodes` forces every broker onto a
different Kubernetes node; `node_selector` (rendered as required node
affinity) and `tolerations` handle placement.

### Storage

`disk_size` (default 10Gi) sizes each node's data volume — queues,
quorum-queue Raft logs and the metadata store; PVCs cannot shrink.
`storage_class` picks the class. `ephemeral: true` swaps all of it
for emptyDir — everything vanishes with each pod restart, throwaway
dev only (the spec rejects it alongside `disk_size` /
`storage_class`).

### Credentials

The operator generates the admin credentials; the spec never carries
them. The `secret_backend` block redirects delivery — `vault` (the
credential-updater sidecar rotates; the Secret output goes empty) or
`external_secret_name` (a pre-existing Secret, e.g. from
KubernetesExternalSecret).

### TLS

`tls.secret_name` mounts an existing certificate Secret (the
cert-manager seam); `ca_secret_name` adds mutual TLS;
`disable_non_tls_listeners` closes every plain port — AMQPS moves to
5671, management to 15671, and the exported endpoints follow.

### Server configuration

`configuration.additional_plugins` extends the always-on essentials
(rabbitmq_management, rabbitmq_prometheus,
rabbitmq_peer_discovery_k8s) — e.g. rabbitmq_mqtt, rabbitmq_stomp,
rabbitmq_shovel. `additional_config` appends rabbitmq.conf lines to
the operator-generated file; `advanced_config`, `env_config` and
`erlang_inet_config` cover the rest of RabbitMQ's own vocabulary.
Changing plugins or config rolls the cluster.

### Exposure

No ingress resources are created. The client Service's `service.type`
(cluster_ip / load_balancer / node_port) plus `service.annotations`
(NLB annotations, external-dns hostname, ...) are the cloud-exposure
surface.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster runs in |
| `cluster_name` | RabbitmqCluster resource name (= `metadata.name`) |
| `service_name` | Client Service (`<name>`) |
| `headless_service_name` | Inter-node Service (`<name>-nodes`) |
| `amqp_endpoint` | In-cluster AMQP endpoint (5672; 5671 when plain listeners are closed) |
| `management_endpoint` | Management API / UI endpoint (http 15672 / https 15671) |
| `default_user_secret_name` | `<name>-default-user` credentials Secret; empty under the Vault backend |
| `port_forward_command` | Management-UI access from a workstation |

## Related Components

- [KubernetesRabbitMqOperator](/docs/catalog/kubernetes/rabbitmq-operator)
  — the engine; a registry prerequisite, declared first
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) —
  provides the target namespace via reference
- [KubernetesStorageClass](/docs/catalog/kubernetes/storage-class)
  — referenced by the data volumes' storage class
- [KubernetesCertificate](/docs/catalog/kubernetes/kubernetescertificate)
  — issues the TLS Secret the `tls` block references
- [KubernetesKafka](/docs/catalog/kubernetes/kafka) — the
  event-streaming sibling for replayable logs at scale

## Next Steps

Move to 3 replicas with `spread_across_nodes` before real traffic
arrives, and keep the count odd — remember scaling down is a
migration, not a field change. Keep memory requests equal to limits
so the high watermark lands where the limit says. Declare `tls` (and
close the plain listeners) before anything crosses a trust boundary,
and compose external exposure from the client Service's type and
annotations — this component never embeds it.
