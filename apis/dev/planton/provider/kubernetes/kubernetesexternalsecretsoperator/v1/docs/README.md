# KubernetesExternalSecretsOperator: Research and Design

## Introduction

The External Secrets Operator (ESO) is the de-facto machinery for keeping
Kubernetes Secrets in sync with external secret stores: a controller
reconciles store connections and secret syncs, a validating webhook checks
every ESO resource at admission, and a cert-controller bootstraps and
rotates the webhook's serving certificate. This component installs that
machinery from the official Helm chart. Which stores exist and which
secrets sync are separate first-class kinds
(KubernetesClusterSecretStore / KubernetesSecretStore /
KubernetesExternalSecret) — this component deliberately manages ONLY the
operator machinery, the same split as cert-manager vs. issuers and
certificates.

## Upstream Architecture

ESO's CRD model has two halves. **Stores** (SecretStore, namespace-scoped;
ClusterSecretStore, cluster-scoped) describe a backend connection: which
provider (AWS Secrets Manager, GCP Secret Manager, Azure Key Vault, Vault,
and dozens more) and how to authenticate. **ExternalSecrets** describe a
sync: which store, which remote keys, what shape of Kubernetes Secret to
materialize, and how often to refresh. The controller reconciles both; the
webhook validates them at admission (a disabled webhook moves
misconfiguration failures from apply time to reconcile time); the
cert-controller exists solely so the webhook can serve TLS without
requiring cert-manager.

The CRDs serve `external-secrets.io/v1` — the stable API surface that
store and secret kinds compile against.

## The Release Name Is Fixed

The Helm release is always named `external-secrets`. One installation per
cluster is an upstream architectural constraint — the CRDs and the
validating webhook configuration are cluster-global, so a second
installation fights the first over both. Fixing the release name (and with
it the chart-derived ServiceAccount name) makes the identity subject
deterministic for cloud-side bindings and makes a second, broken install
inexpressible. For genuine multi-operator sharding upstream provides
`controller_class` (stores opt into a class; each installation reconciles
only its class) and `scoped_namespace` + `scoped_rbac` (one namespace, a
Role instead of a ClusterRole) — both typed.

## Identity: Ambient vs Per-Store

Store authentication has two layers. **Per-store identities** live in each
store's auth block (dedicated ServiceAccounts, credential Secrets) — the
finer-grained posture, recommended for multi-team clusters, and entirely
outside this component. **Ambient identity** is this component's
`workload_identity`: the shared per-cloud oneof rendered onto the
CONTROLLER ServiceAccount (EKS IRSA role-arn annotation, GKE Workload
Identity email annotation, AKS client-id annotation plus the
`azure.workload.identity/use` pod label). Stores that leave their auth
block empty fall back to it — the simplest posture when one cloud identity
may read every synced secret. The two layers make cross-cloud combinations
first-class: a GKE cluster syncing AWS Secrets Manager puts AWS credentials
(or an assumable role) on the store, touching nothing here.

## Typed Surface vs Escape Hatch

The typed spec covers: CRD lifecycle (install with the release, keep on
uninstall), controller scaling (replicas + leader election — enforced
together by validation, since replicas without leader election race;
concurrent reconciliation), sharding and scoping (controller_class,
scoped_namespace + scoped_rbac), logging, resources, workload identity,
scheduling (node selector, tolerations, priority class, pod disruption
budget), Prometheus (opt-in ServiceMonitor), per-component webhook and
cert-controller tuning (enabled/replicas/resources), and the image
registry override (air-gapped mirrors).

`helm_values` merges LAST with Helm `-f` semantics on both engines
(Terraform natively via the two-document values list; Pulumi module-side
with the same deep-merge). Deliberately unmodeled as typed fields:

- **PushSecret write-back** (cluster → external store) — a niche inversion
  of ESO's core direction; unmodeled until demand, and its CRDs ship with
  the release regardless
- **cert-manager webhook-cert integration** — the chart's
  `webhook.certManager` values ride `helm_values`; the chart then skips the
  cert-controller Deployment on its own. The typed `cert_controller.enabled`
  knob covers the common case (on by default)
- **Generators** (dynamic credentials: ECR/GCR tokens, STS sessions, ...) —
  workload-side CRD usage, not installation configuration
- **Per-component scheduling overrides** — controller-level scheduling is
  typed; webhook/cert-controller differences belong in `helm_values`

## CRD Lifecycle

`installCRDs` matches the chart's own default (true) — one component owning
both halves is strictly simpler. `keep_on_uninstall` (default TRUE) has NO
chart knob: the chart templates its CRDs and Helm would DELETE them on
uninstall, cascading to every ExternalSecret and SecretStore object
cluster-wide. The modules render the standard
`helm.sh/resource-policy: keep` annotation onto the CRDs (via the chart's
`crds.annotations`), so that destructive act requires an explicit `false`.

Kept CRDs also pin the install namespace: they retain the Helm release's
namespace in their ownership metadata, so re-installing the operator into
a different namespace fails with Helm's release-ownership error on the
surviving CRDs. Treat the namespace as permanent — moving requires first
deleting the kept CRDs, with the cascade above.

## Install Semantics

Both engines install a REAL Helm release and wait (600s timeout) for all
three Deployments — controller, webhook, cert-controller — to become
Available. The webhook validates every ESO resource at admission: an
operator whose webhook is down rejects every SecretStore/ExternalSecret
apply, so returning success early would only move the failure into the next
resource. Atomic + cleanup-on-fail: a failed install never leaves a
half-deployed operator.

## Upgrade Posture

`chart_version` pins the chart (default: the validated pin). Chart and
operator versions are aligned upstream — chart 2.8.0 ships operator v2.8.0
— so one pin decides both, unlike charts whose app versions float
separately. Upgrades re-run the release with the new chart, which also
upgrades the CRDs it installed; the CRDs serve `external-secrets.io/v1`.

## Outputs as Composition Seams

`controller_service_account` (fixed to `external-secrets`) — what the cloud
side must trust for ambient keyless access: the IRSA trust policy subject,
the GKE Workload Identity member, the Azure federated credential subject.
`namespace` — where cluster-scoped stores' credential Secrets convention-
ally live. `release_name` — the handle for verification and Helm-level
operations.

## E2E

Chart-default and tuned installs run on the kind cluster, both engines,
with the install verifier requiring the three Deployments Available and the
core CRDs (ExternalSecret, SecretStore, ClusterSecretStore) Established.
This component is also the prerequisite fixture for every store and
external-secret scenario.
