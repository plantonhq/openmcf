---
title: "RabbitMQ Cluster Operator"
description: "RabbitMQ Cluster Operator deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesrabbitmqoperator"
---

# RabbitMQ Cluster Operator

Install the [RabbitMQ Cluster Operator](https://www.rabbitmq.com/kubernetes/operator/operator-overview) — the MPL-2.0 operator maintained by the RabbitMQ team — from its official single-file release manifest (the operator has no Helm chart). The operator reconciles `RabbitmqCluster` custom resources (declared with **RabbitMQ**) into running RabbitMQ clusters: one StatefulSet per cluster, the client and inter-node Services, generated administrator credentials, and rolling upgrades.

This component installs and configures the **engine**. RabbitMQ clusters themselves are declared with RabbitMQ resources — one per cluster — which this operator reconciles.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions the release manifest's documents:

- **The `rabbitmq-system` namespace** — the manifest's FIXED installation namespace, baked into its own cross-references (the webhook client configuration, the certificate DNS names, the CA-injection annotations, the cluster-role binding subjects); it is not configurable
- **The RabbitmqCluster CRD** (`rabbitmqclusters.rabbitmq.com`) — one document of the applied manifest, so it installs AND deletes with this resource; see the lifecycle warning under Key Configuration
- **The operator Deployment** (`rabbitmq-cluster-operator`, image `ghcr.io/rabbitmq/cluster-operator:2.22.3` at the pinned release, 200m CPU / 500Mi memory for both requests and limits, one replica) with its RBAC and the metrics Service on port 8080
- **Admission webhooks** — mutating and validating webhook configurations for RabbitmqCluster with `failurePolicy: Fail`, their serving certificate issued by a cert-manager `Certificate` (self-signed `Issuer`) with CA-injection annotations

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Cluster Side

- **A running cert-manager** — a HARD prerequisite. The admission webhooks' serving certificate is a cert-manager Certificate; without cert-manager it never issues, and because the webhooks fail closed, every RabbitmqCluster admission fails. Declare cert-manager with **Kubernetes Cert Manager** first.
- **No existing install** — the admission webhooks are cluster-scoped singletons with fixed names, so exactly ONE install per cluster is the upstream contract. A second install cannot coexist.

### Watch Scope

Decide up front where RabbitMQ resources will live. Empty `watch_namespaces` watches **ALL namespaces** — the upstream default, and the opposite default from chart-scoped operators like the Altinity or external-secrets operators. Entries fence the watch; every namespace that will hold RabbitMQ resources must then be covered.

## Deploy

### Console

Open the deployment store, find **RabbitMQ Cluster Operator**, and click **Deploy**. The creation wizard walks you through the installation contract (the fixed namespace, the pinned release, the CRD lifecycle), the watch scope, the operator image and pull secrets, the fleet-wide default images, resources, and scheduling. Start from the **Standard** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesRabbitMqOperator
metadata:
  name: rabbitmq-operator
  org: acme-corp
  env: prod
spec: {}
```

```shell
planton apply -f rabbitmq-operator.yaml
```

An empty spec is the production-standard posture: the release manifest's own defaults, with the operator watching ALL namespaces. From that point, RabbitMQ resources in any namespace reconcile into running clusters.

## Key Configuration

**There is no version field — deliberately** — the installed operator (and its RabbitmqCluster CRD schema) is pinned to the release this catalog's typed SDK was generated against; the pinned release manifest ships `ghcr.io/rabbitmq/cluster-operator:2.22.3`. A user-selectable version would silently drift the CRD schema away from the typed resources built on it. Operator upgrades arrive with catalog releases, not spec edits.

**The CRD deletes with this resource** — the RabbitmqCluster CRD is one document of the applied manifest, so destroying the operator deletes the CRD, which CASCADE-DELETES every RabbitmqCluster on the cluster — brokers, queues, and the PVC-backed data behind them. Never destroy the operator while RabbitMQ resources exist; destroy every cluster first, the operator last.

**Watch scope defaults to everything** — empty `watch_namespaces` watches ALL namespaces. Entries fence the watch (rendered as the operator's `OPERATOR_SCOPE_NAMESPACE` environment variable, comma-separated); the fence is silent on the outside — a RabbitMQ cluster declared beyond it is never reconciled, with no event pointing at the fence.

**Fleet-wide image defaults are the air-gap seam** — `default_rabbitmq_image` sets the server image for every RabbitmqCluster that does not pin its own (empty = the compiled-in `rabbitmq:4.2.6-management`; the `-management` variant is REQUIRED — the operator's generated configuration expects the management plugin). `default_user_updater_image` sets the credential-updater sidecar image, consulted only by clusters using the Vault secret backend (empty = the compiled-in `ghcr.io/rabbitmq/default-user-credential-updater:1.0.14`).

**Operator pod knobs** — `operator_image` overrides the operator image itself (keep the mirror's tag at the pinned release), `image_pull_secrets` names pull secrets in the `rabbitmq-system` namespace, `resources` overrides the manifest's 200m CPU / 500Mi memory defaults, and `node_selector` / `tolerations` steer scheduling.

## Outputs and Dependencies

### What This Component Consumes

This component's spec is self-contained — no fields reference other resources' outputs. Its one dependency is environmental: a running cert-manager on the target cluster (see Before You Deploy).

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace the operator is installed into (always `rabbitmq-system` — the release manifest's fixed namespace) | Composition, debugging |
| `deployment_name` | The operator Deployment (`rabbitmq-cluster-operator` — the release manifest's fixed name) | Monitoring, log collection |
| `metrics_endpoint` | In-cluster metrics endpoint of the operator (Prometheus format, port 8080) | Prometheus scrape configuration |
| `crd_name` | The RabbitmqCluster CRD the operator serves (`rabbitmqclusters.rabbitmq.com`) — deleted with this resource | Operational tooling, destroy-order checks |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard** — the one-per-cluster engine install with an empty spec: the release manifest's own defaults ARE the production-standard posture, with the operator watching all namespaces. Start from the **Standard** preset.

**Namespace Fenced** — the multi-tenant posture: the operator watches only the namespaces that are allowed to hold RabbitMQ clusters. Start from the **Namespace Fenced** preset.

**Air-Gapped Mirror** — every image this install (and the clusters it will create) pulls, re-pointed at a private registry with a pull secret. Start from the **Air-Gapped Mirror** preset.

## Works With

- **RabbitMQ** — the RabbitMQ clusters this operator reconciles; deploy the operator FIRST, keep clusters inside the watched namespaces, and destroy every cluster before ever destroying the operator.
- **Kubernetes Cert Manager** — the hard prerequisite; the admission webhooks' serving certificate is a cert-manager Certificate with CA injection.
- **Kubernetes Prometheus** — scrapes the operator's `metrics_endpoint` for reconcile health.
