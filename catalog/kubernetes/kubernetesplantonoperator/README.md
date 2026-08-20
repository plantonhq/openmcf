# Kubernetes Planton Operator

## When NOT to Use This

**One installation per cluster.** The Planton operator enforces this
itself at startup — it scans for sibling operator Deployments (matched by
the chart's `app.kubernetes.io/name: planton-operator` +
`control-plane: controller-manager` labels) and refuses to start beside
one, naming the remedy in its log. The Helm release name is therefore
fixed to `planton-operator` and never derives from `metadata.name`.

Also not the right component when:

- **You want a Planton platform** — this component installs and
  configures the MANAGER. The platforms themselves are declared with
  KubernetesPlantonPlatform resources — one per platform — which the
  operator reconciles. "Install the operator" and "declare a platform"
  are different lifecycles: cluster owners install the first once,
  platform owners declare the second per team or per environment.
- **You want a guided first install** — the Planton desktop app's
  self-hosted deploy journey walks a cluster from kubeconfig to a
  signed-in console with preflight checks and a managed tunnel. This
  component is the declarative/GitOps motion for clusters managed as
  code, and the day-2 home for installs however they began.

## Overview

**KubernetesPlantonOperator** installs the Planton operator from the
official `planton-operator` Helm chart (OCI,
`ghcr.io/plantonhq/charts`). The operator reconciles `PlantonPlatform`
custom resources into running self-hosted Planton platforms — control
plane, web console, identity server (Keycloak), PostgreSQL
(CloudNativePG), cache, workflow engine (Temporal), secrets manager
(OpenBAO), and an in-cluster deployment runner — each platform in its own
namespace, all served by this one manager.

The typed spec covers the chart's meaningful configuration surface, with
a `helm_values` escape hatch (merged last, Helm `-f` semantics, identical
on both engines) for anything beyond it.

**Key design points:**

- **Installing the operator deploys NO platform** — the operator never
  auto-creates a PlantonPlatform; every platform is a deliberate
  KubernetesPlantonPlatform declaration. One operator watches all
  namespaces and serves every platform on the cluster.
- **The CRD is module-owned and survives uninstall** — the
  `plantonplatforms.planton.ai` CRD is applied by the modules from a copy
  staged at the pinned chart version (the chart's own install-once
  `crds/` posture is bypassed with `skip_crds`), so a `chart_version`
  upgrade carries the matching CRD update, and destroying this resource
  never cascade-deletes platform declarations or the platforms behind
  them.
- **Versions are floored, loudly** — charts below `0.7.0` ship operators
  whose reconcilers predate the PlantonPlatform schema the staged CRD
  advertises; both engines refuse them instead of letting the API server
  accept fields a running operator would silently ignore.
- **Destroying the operator strands nothing** — platforms keep running
  (unmanaged — spec edits wait for an operator's return), declarations
  survive on the kept CRD, and platform deletion still completes
  (teardown is Kubernetes garbage collection of owner-referenced
  objects, not operator work).

## Spec Fields

| Field | Required | Default | Description |
|---|---|---|---|
| `namespace` | yes | — | Installation namespace (literal or KubernetesNamespace reference); `planton-operator` is the convention |
| `create_namespace` | no | `false` | Create (and own) the namespace with governance labels |
| `chart_version` | no | `0.7.1` | Exact-semver chart pin; floored at `0.7.0` |
| `skip_crds` | no | `false` | Skip the module-owned CRD (something else manages it) |
| `replicas` | no | `1` | Leader-elected warm standbys — operator failover, not throughput |
| `leader_election` | no | `true` | Required when `replicas > 1` |
| `resources` | no | chart defaults | Requests 10m/256Mi, limits 500m/512Mi when unset |
| `service_account` | no | chart-created | Bring-your-own name or annotations |
| `common_labels` | no | — | Extra labels on every chart-rendered resource |
| `pod_annotations` | no | — | Annotations on the operator pod |
| `node_selector`, `tolerations` | no | — | Operator pod scheduling |
| `image_pull_secrets` | no | — | Secret NAMES for private mirrors |
| `image` | no | chart default | Air-gap mirror override (repository/tag) |
| `helm_values` | no | — | Raw values YAML merged last (never for secrets) |

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesPlantonOperator
metadata:
  name: planton-operator
  org: acme-corp
  env: platform
spec:
  namespace:
    value: planton-operator
  create_namespace: true
  chart_version: "0.7.1"
```

Then declare platforms with KubernetesPlantonPlatform resources — the
zero-config platform needs only a version, and its first visitor becomes
the admin through the console's setup page.

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name (fixed `planton-operator`) |
