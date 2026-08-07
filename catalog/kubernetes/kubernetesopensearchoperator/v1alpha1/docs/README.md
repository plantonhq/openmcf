# KubernetesOpenSearchOperator: Research and Design

## Introduction

KubernetesOpenSearchOperator installs the OpenSearch Kubernetes
Operator from the official `opensearch-operator` Helm chart
(https://opensearch-project.github.io/opensearch-k8s-operator/, pinned
2.8.0) as a single Helm release named after `metadata.name`. The
operator is the ENGINE of the OpenSearch story in this catalog:
KubernetesOpenSearch declares `OpenSearchCluster` custom resources, and
this operator reconciles them into node StatefulSets, Services, TLS
Secrets, security bootstrap, rolling upgrades, and Dashboards
deployments.

## The Deployment Landscape

OpenSearch without an operator is the classic stateful anti-pattern:
cluster-manager quorums, per-node certificates, drain ordering during
upgrades, and security-plugin initialization are Day-2 concerns no
plain StatefulSet or Helm chart encodes. The operator carries that
expertise as a reconciler, which is why the catalog splits the concern
in two: this kind installs the engine once, KubernetesOpenSearch
declares each cluster.

### The chart/image pinning truth

The chart's default manager image tag is its appVersion, and the served
repository index is the source of truth for what actually installs:

- **Chart 2.8.0** pairs with the **stable operator 2.8.0 image** — the
  pinned default, and the newest stable pairing.
- **The newer served charts (2.8.3, 2.8.4, 3.0.x)** default their
  manager image to a **PRERELEASE tag** (verified against the served
  charts' own appVersion). They additionally bundle next-generation
  CRDs in the `opensearch.org` API group that the stable operator does
  not serve.
- **The 3.x line migrates the CRDs** from `opensearch.opster.io` to
  `opensearch.org` — a version bump is therefore a license and
  API-group re-check, not a routine upgrade. Pin those lines only
  after upstream cuts a stable 3.x operator release.

## The CRD Lifecycle: Module-Owned, Keep-on-Uninstall

The chart TEMPLATES its ten OpenSearch CRDs (OpenSearchCluster plus the
opensearchusers, opensearchroles, opensearchismpolicies and the rest of
the `opensearch.opster.io` group) as release-owned resources, with no
keep-on-uninstall knob upstream. Left to Helm, uninstalling the
operator would delete the CRDs — and deleting a CRD cascade-deletes
every custom resource of that type cluster-wide, taking every
OpenSearch cluster and its data with it.

The modules therefore own the CRD lifecycle end to end:

- **`installCRDs` pins false unconditionally** in the rendered chart
  values — and is re-pinned AFTER the `helm_values` escape-hatch merge,
  the one deliberate exception to the escape hatch's last-word
  contract. Letting an override hand the CRDs to Helm would arm the
  cascade-delete this design exists to prevent.
- **The staged CRD files (the chart 2.8.0 set) apply as module-owned
  resources**, keyed by each CRD's OWN `metadata.name` (never a
  positional index), so state addresses stay stable across file
  renames and reorderings.
- **Keep-on-uninstall, per engine:**
  - Terraform: `kubectl_manifest` with `apply_only = true` — the
    provider's Delete is a NO-OP (verified in the provider source), so
    destroy removes the CRDs from state without touching the cluster.
  - Pulumi: the classic yaml `ConfigGroup` with `retainOnDelete`
    delivered through a resource TRANSFORMATION — the one mechanism the
    SDK propagates to the ConfigGroup's in-process children (the yaml
    SDK forwards only parent/version options to children, and yaml/v2
    children are created provider-side, beyond the reach of any
    SDK-side option — both verified in the pinned pulumi-kubernetes
    source).
- **Server-side apply** keeps the megabyte-scale CRD schemas free of
  client-side last-applied-configuration annotation bloat and lets a
  re-adopting apply never conflict with itself.
- **The release depends on the CRDs**, so the operator never starts
  against an unregistered API group.

## Watch Scope and RBAC Posture

By default the operator watches ALL namespaces, on cluster-wide RBAC
(ClusterRoleBindings) — one install serves every KubernetesOpenSearch
resource on the cluster. Two knobs tighten that on shared clusters:

- **`watch_namespace`** fences the operator to one namespace;
  OpenSearchCluster resources anywhere else are ignored by this
  install.
- **`use_role_bindings`** swaps ClusterRoleBindings for
  namespace-scoped RoleBindings. It is only valid together with
  `watch_namespace` — the chart grants roles only in the watched
  namespace, so a cluster-wide operator cannot run on namespace-scoped
  permissions. The spec enforces this with a CEL rule rather than
  letting the install proceed into a permission-starved operator.

## Design Decisions

- **The install is blocking.** The Helm release waits for the operator
  Deployment to become Available (atomic, 600s timeout, cleanup on
  fail): an operator that never becomes ready — an unpullable image
  from a private mirror is the classic case — fails THIS apply with a
  readiness timeout instead of surfacing later as OpenSearch clusters
  that mysteriously never reconcile.
- **The module owns namespace creation** (`create_namespace`), never
  the Helm release — pre-existing-namespace installs leave the flag
  false.
- **Chart-default-matching values render only on divergence** (watch
  scope, the true-defaulted toggles, the image source), so the
  rendered values stay minimal on both engines — with the one
  `installCRDs` exception above.
- **`deployment_name` replays the chart's own naming.** The exported
  Deployment name derives exactly as the chart's fullname helper does
  it (the release name, suffixed with the chart name when absent,
  truncated to 63 characters) plus `-controller-manager`. This is why
  `nameOverride`/`fullnameOverride` are off limits in `helm_values`.
- **`dns_base` exists because certificates do.** The operator bakes
  the cluster DNS domain into generated certificates and discovery
  addresses; on clusters with a non-default DNS domain, a mismatch
  produces TLS certificates whose SANs do not match the service names
  nodes advertise.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Chart | `opensearch-operator` at https://opensearch-project.github.io/opensearch-k8s-operator/ | Pinned 2.8.0 (spec default) |
| Manager image | `opensearchproject/opensearch-operator` at the chart's appVersion | 2.8.0 = the stable pairing; 2.8.3+/3.0.x default to prerelease tags |
| CRD API group | `opensearch.opster.io` | The 3.x line migrates to `opensearch.org` |
| CRD files | staged with the module (the chart 2.8.0 set) | Upgrade together with the chart pin |
| Deployment | chart fullname + `-controller-manager` | Exported as `deployment_name` |
| Watch scope | all namespaces (chart default) | `watch_namespace` fences; `use_role_bindings` requires it |

## IaC Twins

Pulumi (`module/crds.go` + `module/values.go`) and Terraform
(`main.tf` + `locals.tf`) render identical chart values, the same
module-owned CRD set, and the same keep-on-uninstall posture
(`retainOnDelete` transformation / `apply_only = true`). Keep the
typed-value rendering and the deployment-name derivation in lockstep.
