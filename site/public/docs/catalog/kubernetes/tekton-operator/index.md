---
title: "Tekton Operator"
description: "Tekton Operator deployment documentation"
icon: "package"
order: 100
componentName: "kubernetestektonoperator"
---

# Tekton Operator on Kubernetes

Install the [Tekton Operator](https://tekton.dev/docs/operator/) — the lifecycle manager maintained by the Tekton project — from its official single-file release manifest (the in-repo Helm chart is unpublished and is not a distribution channel). The operator reconciles a `TektonConfig` declaration (declared with **Tekton on Kubernetes**) into running Tekton components — Pipelines, Triggers, Dashboard, Chains — managing their installation, upgrades, and removal through `TektonInstallerSet` resources.

This component installs the **manager only**. Installing it deploys NO pipeline runtime: automatic component installation is disabled by design, so the KubernetesTekton declaration is the single owner of what Tekton actually runs on the cluster.

## What Gets Created

When you deploy this Cloud Resource, the IaC module applies the release manifest's documents:

- **The `tekton-operator` namespace** — the manifest's FIXED installation namespace, baked into its own cross-references (the webhook Service, the RBAC subjects); it is not configurable
- **14 `operator.tekton.dev` CRDs** (including `tektonconfigs.operator.tekton.dev`) — documents of the applied manifest, so they install AND delete with this resource; see the destroy ordering under Key Configuration
- **The operator Deployment** (`tekton-operator`, running `ghcr.io/tektoncd/operator/operator-*` at the pinned `v0.80.0` release; its two containers — lifecycle and cluster-operations — share one image) with its RBAC
- **The admission webhook Deployment** (`tekton-operator-webhook`) validating TektonConfig declarations — including the one-per-cluster singleton rule and the immutable target namespace
- **Auto-installation disabled** — upstream, the operator auto-creates a default TektonConfig (profile `all`) at startup; this install always disables that, because two managers writing one object fight through server-side apply and the cluster's Tekton shape would depend on install order

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Cluster Side

- **No existing install** — exactly ONE operator install per cluster is the upstream contract: its webhooks and CRDs are cluster-scoped singletons with fixed names, and a second install cannot coexist.
- **Registry reachability** — the operator images pull from GitHub Container Registry (`ghcr.io`). Air-gapped clusters set the image overrides and pull secrets (see Key Configuration).

## Deploy

### Console

Open the deployment store, find **Tekton Operator on Kubernetes**, and click **Deploy**. The creation wizard walks you through the installation contract (the fixed namespace, the pinned release, the one-per-cluster rule, the CRD lifecycle), the air-gap image overrides, sizing, and scheduling. Start from the **Operator** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesTektonOperator
metadata:
  name: tekton-operator
  org: acme-corp
  env: prod
spec: {}
```

```shell
planton apply -f tekton-operator.yaml
```

An empty spec is the complete install: the release manifest's own defaults, in its fixed namespace, with automatic component installation disabled. Declare a **Tekton on Kubernetes** resource next to choose what actually runs.

## Key Configuration

**There is no version field — deliberately** — the installed operator (and the TektonConfig schema the KubernetesTekton kind renders against) is pinned to the release this catalog was designed against (`v0.80.0`). A user-selectable version would silently drift the TektonConfig surface away from what the catalog models. Operator upgrades arrive with catalog releases, not spec edits.

**There is no namespace field either** — the release manifest installs into `tekton-operator`, and that name is baked into the manifest's own cross-references. Exactly one install per cluster is the upstream contract.

**The CRDs delete with this resource** — the 14 `operator.tekton.dev` CRDs are documents of the applied manifest, so destroying the operator removes them, which CASCADE-DELETES any TektonConfig on the cluster. Always destroy the KubernetesTekton resource FIRST: its teardown blocks until the operator finishes removing the components, and the `TektonInstallerSet` finalizers are processed only by a RUNNING operator — removing the operator first strands them.

**Image overrides are the air-gap seam** — `operator_image` overrides the image for BOTH containers of the operator Deployment, `webhook_image` the admission webhook's; empty means the release manifest's digest-pinned `ghcr.io/tektoncd/operator/*` images at the pinned release. `image_pull_secrets` names existing `kubernetes.io/dockerconfigjson` Secrets in the fixed `tekton-operator` namespace — references, never credentials.

**Sizing and placement** — the manifest sets NO resource requests or limits (the operator runs unbounded); set `operator_resources` / `webhook_resources` on production clusters with quotas. `node_selector` and `tolerations` steer the operator and webhook pods — scheduling for the Tekton COMPONENT pods lives on the KubernetesTekton resource's placement instead.

## Outputs and Dependencies

### What This Component Consumes

This component's spec is self-contained — no fields reference other resources' outputs, and it has no cluster-side prerequisites beyond registry reachability.

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace the operator runs in (always `tekton-operator` — fixed by the release manifest) | Composition, debugging |

The operator exports no component handles of its own: the Tekton namespace, profile, and dashboard endpoints are the KubernetesTekton resource's outputs — this installation is only the manager that reconciles it.

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Operator** — the complete install with an empty spec, which is deliberately tiny: the operator is a lifecycle manager, not the product. Set the image overrides on air-gapped clusters and resource requests on quota-governed ones. Start from the **Operator** preset.

## Works With

- **Tekton on Kubernetes** — the cluster's TektonConfig declaration this operator reconciles; deploy the operator FIRST, and destroy the declaration FIRST on the way out.
- **Kubernetes Manifest** — Tasks, Pipelines, and their runs are plain custom resources once the Tekton installation converges.
