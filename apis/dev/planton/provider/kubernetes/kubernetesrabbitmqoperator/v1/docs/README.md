# KubernetesRabbitMqOperator: Research and Design

## Introduction

KubernetesRabbitMqOperator installs the RabbitMQ Cluster Operator —
the MPL-2.0 operator maintained by the RabbitMQ team
(rabbitmq/cluster-operator) — from its released single-file manifest,
pinned to release v2.22.3. The operator is the ENGINE of the RabbitMQ
story in this catalog: KubernetesRabbitMq declares `RabbitmqCluster`
custom resources — one per cluster — and this operator reconciles
them into running RabbitMQ clusters (StatefulSet, Services, generated
credentials, rolling upgrades).

## The Deployment Landscape

RabbitMQ without an operator is the classic stateful anti-pattern:
Erlang-cookie coordination, peer discovery, credential generation,
and safe rolling upgrades are Day-2 concerns no plain StatefulSet
encodes. The RabbitMQ team carries that expertise in the Cluster
Operator, which is why the catalog splits the concern in two: this
kind installs the engine once, KubernetesRabbitMq declares each
cluster.

One boundary is deliberate: this component installs the CLUSTER
operator only. RabbitMQ's messaging-topology-operator — declarative
Queue/Exchange/User/Vhost/Policy custom resources — is a separate
upstream product and is NOT installed here.

## Distribution: the Release Manifest, Not a Chart

The operator has NO Helm chart. Its official distribution is the
single-file release manifest (`cluster-operator.yml` attached to each
GitHub release), which the modules fetch from the pinned release tag
and apply per document: the `rabbitmq-system` namespace, the
RabbitmqCluster CRD, RBAC, the webhook and metrics Services, the
cert-manager Issuer and Certificates, and the mutating/validating
webhook configurations. The released asset pins the operator image
(`ghcr.io/rabbitmq/cluster-operator:2.22.3` at the pinned release).

Every document applies verbatim (faithful distribution) except the
operator Deployment, onto which the spec's typed overrides are
patched: the watch-scope and default-image environment variables, the
operator image, container resources, node selector, tolerations, and
image-pull secrets.

### Why the spec has no version field

The installed operator — and its RabbitmqCluster CRD schema — is
pinned to the release this catalog's typed SDK was generated against
(see `pkg/kubernetes/kubernetestypes/Makefile`
`rabbitmq_cluster_operator_release`). A user-selectable version would
silently drift the CRD schema away from the typed resources built on
it. The pin is always an exact release TAG, never a branch: a branch
ref moves as patches land, so the same deployed resource would
install different operators at different times.

## The Fixed Namespace and the Singleton Contract

The release manifest installs into `rabbitmq-system`, and that name
is baked into the manifest's own cross-references: the webhook client
configuration, the cert-manager Certificate DNS names, the
CA-injection annotations, and the cluster-role binding subjects. It
is not configurable — the spec has no namespace field.

Exactly ONE install per cluster is the upstream contract: the
admission webhooks are cluster-scoped singletons with fixed names, so
a second install cannot coexist.

## cert-manager: a Hard Prerequisite

The manifest ships mutating and validating admission webhooks for
RabbitmqCluster with `failurePolicy: Fail`. Their serving certificate
is a cert-manager `Certificate` (backed by a self-signed `Issuer`)
with CA-injection annotations. Installing without a running
cert-manager leaves the webhook certificate unissued — and every
RabbitmqCluster admission failing. cert-manager is therefore declared
as a registry prerequisite of this kind; compose with
KubernetesCertManager first.

## The CRD Lifecycle: Deletes With the Resource

The RabbitmqCluster CRD is one document of the applied manifest — it
installs and is REMOVED with the resource. This is the opposite of
chart-based sibling operators whose charts keep CRDs on uninstall:
destroying this operator deletes the CRD, which cascade-deletes every
RabbitmqCluster on the cluster. Never destroy the operator while
KubernetesRabbitMq resources exist. The warning rides the spec, the
`crd_name` output, and both modules' headers.

## Watch Scope

Empty `watch_namespaces` (the upstream default) watches ALL
namespaces — note this is the opposite default from chart-scoped
operators like the Altinity or external-secrets operators, whose
defaults fence to the install namespace. Setting entries fences the
watch, rendered as the operator's `OPERATOR_SCOPE_NAMESPACE`
environment variable (comma-separated); every namespace that will
hold KubernetesRabbitMq resources must then be covered — a fenced
operator silently ignores clusters elsewhere.

## Fleet-Wide Image Defaults

Two spec fields exist for air-gapped clusters that mirror images:

- **`default_rabbitmq_image`** (the operator's
  `DEFAULT_RABBITMQ_IMAGE` environment variable) applies to every
  RabbitmqCluster that does not pin its own `image`. Empty = the
  operator's compiled-in default (`rabbitmq:4.2.6-management` at the
  pinned release). The `-management` variant is required — the
  operator's generated configuration expects the management plugin.
- **`default_user_updater_image`** (`DEFAULT_USER_UPDATER_IMAGE`) is
  the credential-updater sidecar image, consulted only for
  KubernetesRabbitMq resources using the Vault secret backend. Empty
  = `ghcr.io/rabbitmq/default-user-credential-updater:1.0.14` at the
  pinned release.

The operator image itself is overridden through `operator_image`
(repo/tag/pull secret), with `image_pull_secrets` naming pull secrets
in the `rabbitmq-system` namespace.

## Design Decisions

- **Faithful distribution, one patched document.** The modules apply
  the release manifest per document exactly as shipped, patching only
  the operator Deployment with the spec's typed overrides. The
  manifest's own documents keep their upstream labels untouched.
- **Server-side apply is REQUIRED, not stylistic.** The
  RabbitmqCluster CRD document (~342 KB) exceeds the client-side
  last-applied-configuration annotation cap (256 KB). Both engines
  apply server-side for this reason.
- **The Deployment defaults stay upstream.** Empty `resources` keeps
  the release manifest's 200m CPU / 500Mi memory (requests AND
  limits), one replica, and the metrics Service on port 8080.
- **Outputs are the manifest's fixed handles.** The namespace,
  Deployment name, metrics Service, and CRD name are baked into the
  release manifest, so the outputs are constants of the pin:
  `rabbitmq-system`, `rabbitmq-cluster-operator`,
  `http://rabbitmq-cluster-operator-metrics-service.rabbitmq-system.svc.cluster.local:8080/metrics`,
  and `rabbitmqclusters.rabbitmq.com`.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Release | rabbitmq/cluster-operator `v2.22.3` | The pinned release tag; the single-file `cluster-operator.yml` asset is the distribution |
| Operator image | `ghcr.io/rabbitmq/cluster-operator:2.22.3` | Pinned inside the released asset |
| Default RabbitMQ image | `rabbitmq:4.2.6-management` | Compiled-in default; the `-management` variant is required |
| Default user-updater image | `ghcr.io/rabbitmq/default-user-credential-updater:1.0.14` | Compiled-in default; Vault-backed clusters only |
| Namespace | `rabbitmq-system` | Fixed — baked into the manifest's cross-references |
| Deployment | `rabbitmq-cluster-operator` | Fixed name; 200m/500Mi requests and limits, one replica |
| Metrics | `rabbitmq-cluster-operator-metrics-service`, port 8080 | Prometheus format |
| CRD | `rabbitmqclusters.rabbitmq.com` | One manifest document — DELETES with the resource |

## IaC Twins

Pulumi (`module/deployment_patch.go`) and Terraform (`locals.tf` +
`main.tf`) fetch the same pinned manifest, perform the identical
Deployment patch, apply every other document verbatim with
server-side apply, and derive the same fixed outputs. Keep the patch
surfaces and the release pin in lockstep — the pin must also match
`rabbitmq_cluster_operator_release` in
`pkg/kubernetes/kubernetestypes/Makefile`, which generates the typed
SDK KubernetesRabbitMq is built against.

## Validation Status

The component is offline-validated: the spec's validation tests pass,
and both engines' modules carry plan/preview proofs against the
pinned release manifest across full and minimal shapes. Live
end-to-end verification on a running cluster is pending.
