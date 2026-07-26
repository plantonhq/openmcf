---
title: "RabbitMQ Operator"
description: "RabbitMQ Operator deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesrabbitmqoperator"
---

# Kubernetes RabbitMQ Operator

Installs the RabbitMQ Cluster Operator — the MPL-2.0 operator
maintained by the RabbitMQ team — from its released single-file
manifest (the operator has no Helm chart), pinned to release v2.22.3.
The operator reconciles `RabbitmqCluster` custom resources (declared
with KubernetesRabbitMq) into running RabbitMQ clusters with
StatefulSets, Services, generated credentials, and rolling upgrades.
This component installs and configures the engine; RabbitMQ clusters
are declared separately, one KubernetesRabbitMq resource per cluster.

## What Gets Created

- **The `rabbitmq-system` namespace** — the manifest's fixed
  installation namespace, baked into its own cross-references (not
  configurable; exactly one install per cluster)
- **The RabbitmqCluster CRD** (`rabbitmqclusters.rabbitmq.com`) — one
  document of the applied manifest. **It DELETES with this resource**:
  destroying the operator cascade-deletes every RabbitmqCluster on
  the cluster — never destroy the operator while KubernetesRabbitMq
  resources exist
- **The operator Deployment** (`rabbitmq-cluster-operator`, image
  `ghcr.io/rabbitmq/cluster-operator:2.22.3`, 200m CPU / 500Mi memory
  requests and limits, one replica) with its RBAC and the metrics
  Service on port 8080
- **Admission webhooks** — mutating and validating webhook
  configurations for RabbitmqCluster (`failurePolicy: Fail`), with
  their cert-manager Issuer and Certificates

## Prerequisites

- **A running cert-manager** (compose with KubernetesCertManager) —
  the admission webhooks' serving certificate is a cert-manager
  Certificate with CA injection; without cert-manager the certificate
  never issues and every RabbitmqCluster admission fails
- **No existing install** — the webhooks are cluster-scoped
  singletons with fixed names, so exactly one install per cluster

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesRabbitMqOperator
metadata:
  name: rabbitmq-operator
spec: {}
```

An empty spec is the production-standard posture: the release
manifest's own defaults, with the operator watching ALL namespaces
(the upstream default). From that point, KubernetesRabbitMq resources
in any namespace reconcile into running RabbitMQ clusters.

## Configuration

### Version

There is no version field — deliberately. The installed operator and
its RabbitmqCluster CRD schema are pinned to the release this
catalog's typed SDK was generated against; a user-selectable version
would silently drift the CRD schema away from the typed resources
built on it.

### Watch scope

Empty `watch_namespaces` watches ALL namespaces (the upstream
default — the opposite default from chart-scoped operators like the
Altinity operator). Entries fence the watch, rendered as the
operator's `OPERATOR_SCOPE_NAMESPACE` environment variable
(comma-separated) — every namespace that will hold
KubernetesRabbitMq resources must then be covered.

### Fleet-wide image defaults

For air-gapped clusters that mirror images: `default_rabbitmq_image`
sets the server image for every RabbitmqCluster that does not pin its
own (empty = the compiled-in `rabbitmq:4.2.6-management` — the
`-management` variant is required), and `default_user_updater_image`
sets the credential-updater sidecar image used only by Vault-backed
clusters (empty = the compiled-in
`ghcr.io/rabbitmq/default-user-credential-updater:1.0.14`).

### Operator pod

`operator_image` overrides the operator image itself (the air-gap /
private-mirror path); `image_pull_secrets` names pull secrets in the
`rabbitmq-system` namespace. `resources` overrides the container
sizing (empty = the manifest's 200m CPU / 500Mi memory for both
requests and limits); `node_selector` and `tolerations` steer
scheduling.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator is installed into (always `rabbitmq-system`) |
| `deployment_name` | The operator Deployment (`rabbitmq-cluster-operator`) |
| `metrics_endpoint` | In-cluster Prometheus metrics endpoint (port 8080) |
| `crd_name` | The RabbitmqCluster CRD (`rabbitmqclusters.rabbitmq.com`) — deleted with this resource |

## Related Components

- [KubernetesRabbitMq](/docs/catalog/kubernetes/rabbitmq) —
  declares the RabbitmqCluster resources this operator reconciles
- [KubernetesCertManager](/docs/catalog/kubernetes/kubernetescertmanager)
  — the hard prerequisite; the admission webhooks' serving
  certificate is a cert-manager Certificate

## Next Steps

Declare a RabbitMQ cluster with KubernetesRabbitMq and the operator
reconciles it. If the operator was installed with a fenced
`watch_namespaces` list, keep clusters inside the watched namespaces
— they are silently ignored anywhere else. And remember the destroy
ordering: destroy all KubernetesRabbitMq resources before ever
destroying the operator, because the CRD deletes with it.
