# KubernetesCertManager: Research and Design

## Introduction

cert-manager is the de-facto certificate machinery for Kubernetes: a
controller that watches Certificate resources and drives issuance, a
validating webhook, and a cainjector that maintains CA bundles in webhook
configurations. This component installs that machinery from the official
Helm chart. Signing authorities (issuers) and certificates are separate
first-class kinds — this component deliberately manages ONLY the controller
machinery.

## Design Authority

Designed from the pinned upstream chart's `values.yaml` (deploy/charts/
cert-manager) and validated against the chart's own defaults. Chart identity
(`cert-manager` at `https://charts.jetstack.io`) is fixed in both engines'
modules — cross-engine chart drift deploys two different products.

## Typed Surface vs Escape Hatch

The typed spec covers the chart's meaningful configuration surface:

- **CRD lifecycle** (`crds`): install with the release (Planton default TRUE
  — the upstream chart defaults false and expects a separate kubectl apply;
  one component owning both halves is strictly simpler), keep on uninstall
  (TRUE, upstream-aligned: deleting the CRDs cascades to every certificate
  object cluster-wide). Kept CRDs pin the install namespace — they retain
  the Helm release's namespace in their ownership metadata, so
  re-installing into a different namespace fails with Helm's
  release-ownership error on the surviving CRDs; treat the namespace as
  permanent
- **Controller**: replicas, resources, log level, leader-election namespace,
  cluster-resource namespace, certificate owner refs, feature gates (typed
  map rendered to the chart's comma-string), max concurrent challenges
- **DNS-01 self-check** (`dns01_self_check`): recursive resolvers + only-flag
  — the split-horizon fix (in-cluster DNS serving a private view makes
  cert-manager's pre-flight TXT check hang issuance forever)
- **Workload identity**: the shared per-cloud oneof, rendered to the
  ServiceAccount annotation each cloud's webhook expects (plus the AKS pod
  label). The chart owns the ServiceAccount; the identity rides
  `serviceAccount.annotations`
- **Sub-components**: webhook (replicas, timeout, host-network + secure-port
  — the EKS-custom-CNI fix), cainjector (enabled/replicas/resources),
  startupapicheck (enabled/timeout), prometheus (metrics + opt-in
  ServiceMonitor), scheduling (node selector, tolerations), image registry
  (air-gapped mirror), pod disruption budget

`helm_values` merges LAST with Helm `-f` semantics on both engines (Terraform
natively via the two-document values list; Pulumi module-side with the same
deep-merge). Deliberately unmodeled as typed fields: per-sub-component
scheduling overrides, extraArgs/extraEnv/volumes (the config-file and
values escape hatches cover them), and the chart's securityContext blocks
(upstream defaults are correct; overriding them is an expert move that
belongs in helm_values).

## Install Semantics

Both engines install a REAL Helm release and wait for full readiness
including the startupapicheck hook Job — the post-install probe that proves
the webhook actually serves. A cert-manager whose webhook is down rejects
every Issuer/Certificate apply, so returning success early would only move
the failure into the next resource. Atomic + cleanup-on-fail: a failed
install never leaves a half-deployed cert-manager.

## The Release Name Is Fixed

The Helm release is always named `cert-manager` (so the chart-derived
ServiceAccount name is stable for cloud-side identity bindings, and because
one installation per cluster is an upstream architectural constraint —
cluster-scoped CRDs and webhook configurations). A manifest-derived release
name would only enable a second, broken install.

## Outputs as Composition Seams

`service_account_name` — what the cloud side must trust for keyless DNS-01
(IRSA trust policy, GKE WI binding, Azure federated credential).
`cluster_resource_namespace` — where cert-manager reads Secrets for
cluster-scoped resources; KubernetesClusterIssuer's namespace FK defaults to
this output, so issuer credentials always land where the controller looks.

## E2E

Chart-default and tuned installs run on the kind cluster, both engines, with
the install verifier requiring the three Deployments Available and the core
CRDs Established. This component is also the registry prerequisite fixture
for every issuer/certificate scenario.
