# Kubernetes OpenSearch Operator

## When NOT to Use This

**This component installs the ENGINE, not a search cluster.** The
OpenSearch Kubernetes Operator reconciles `OpenSearchCluster` custom
resources into running clusters; those clusters are declared with
KubernetesOpenSearch — one resource per cluster. Install the operator
once per Kubernetes cluster (or once per watched namespace), then
declare search clusters against it.

Also not the right component when:

- **You want an OpenSearch cluster** — that is KubernetesOpenSearch;
  this component is the controller that reconciles it.
- **You want a managed cloud search service** — use the host cloud
  provider's managed search kinds; this component is for running
  OpenSearch ON the Kubernetes cluster itself.
- **You need the 3.x operator line today** — the served charts past
  2.8.0 (2.8.3, 2.8.4, 3.0.x) default their manager image to a
  PRERELEASE tag, and the 3.x line migrates the CRDs from the
  `opensearch.opster.io` API group to `opensearch.org`. The pinned
  default (2.8.0) is the newest stable chart/image pairing; pin a 3.x
  chart only after upstream cuts a stable 3.x operator release.

## Overview

**KubernetesOpenSearchOperator** installs the OpenSearch Kubernetes
Operator — the opensearch-project's operator for running OpenSearch
(the Apache-2.0 search and analytics engine) on Kubernetes — from the
official `opensearch-operator` Helm chart
(https://opensearch-project.github.io/opensearch-k8s-operator/). The
operator reconciles `OpenSearchCluster` custom resources into running
search clusters with managed TLS, security bootstrap, safe rolling
upgrades, and OpenSearch Dashboards.

**Key design points:**

- **The module owns the CRD lifecycle.** The chart templates its ten
  OpenSearch CRDs as release-owned resources with NO keep-on-uninstall
  knob upstream — a Helm-owned install would cascade-delete every
  OpenSearchCluster (and its data) on uninstall. The modules therefore
  pin `installCRDs: false` unconditionally and apply the staged CRD
  files themselves: the Terraform module with `kubectl_manifest`
  `apply_only = true` (the provider's Delete is a no-op), the Pulumi
  module with a `retainOnDelete` transformation on each CRD. Destroying
  this resource removes the operator but NEVER the CRDs — an operator
  uninstall never cascade-deletes clusters.
- **Chart/image pairing is deliberate.** The chart's default manager
  image tag is its appVersion; chart 2.8.0 pairs with the stable
  operator 2.8.0 image. A `chart_version` bump is a license and
  API-group re-check, not a routine upgrade.
- **Watch scope defaults to cluster-wide.** The operator watches ALL
  namespaces by default (cluster-wide RBAC). Set `watch_namespace` to
  fence it to one namespace, and pair with `use_role_bindings` to swap
  ClusterRoleBindings for namespace-scoped RoleBindings on shared
  clusters — the spec rejects `use_role_bindings` without a
  `watch_namespace` (a cluster-wide operator cannot run on
  namespace-scoped permissions).
- **The metrics endpoint ships shielded.** The kube-rbac-proxy sidecar
  (chart default: enabled) puts the operator's metrics endpoint behind
  Kubernetes RBAC.
- **`helm_values` is the escape hatch** — additional chart values
  merged LAST over everything the typed fields render (Helm `-f`
  semantics, identical on both engines). Two keys are off limits:
  `installCRDs` is re-pinned by the modules AFTER this merge (the one
  deliberate exception to the escape hatch's last-word contract), and
  `nameOverride`/`fullnameOverride` break the exported
  `deployment_name` output, which derives from the chart's default
  naming.

## Essential Configuration Fields

### Required

- **`spec.namespace`**: namespace to install the operator into —
  literal or a KubernetesNamespace reference (`create_namespace` to
  own it)

### Common

- **`spec.chart_version`**: chart pin (default `2.8.0` — the newest
  served chart whose default manager image is a stable operator
  release)
- **`spec.watch_namespace`**: restrict the operator to one namespace;
  empty = watch all namespaces (the chart default)
- **`spec.use_role_bindings`**: namespace-scoped RoleBindings instead
  of ClusterRoleBindings — only valid together with `watch_namespace`
- **`spec.log_level`**: `debug`, `info` (default), `warn`, or `error`
- **`spec.dns_base`**: the cluster DNS domain baked into generated
  certificates and discovery addresses (default `cluster.local`) — a
  mismatch produces TLS certificates whose SANs do not match the
  service DNS names nodes advertise
- **`spec.parallel_recovery_enabled`**: recover cluster pods in
  parallel after failures (chart default true; upstream marks it
  experimental)
- **`spec.kube_rbac_proxy_enabled`**: the RBAC shield on the metrics
  endpoint (chart default true)
- **`spec.resources`**: manager container resources — empty = the
  chart defaults (requests 100m/350Mi, limits 200m/500Mi)
- **`spec.node_selector` / `spec.tolerations`**: operator pod
  scheduling
- **`spec.image` / `spec.image_pull_secrets`**: air-gap and
  private-mirror path (empty = `opensearchproject/opensearch-operator`
  at the chart's appVersion)
- **`spec.helm_values`**: the escape hatch (see above for the two
  off-limits keys)

## Stack Outputs

| Output | Purpose |
|---|---|
| `namespace` | Namespace the operator is installed into |
| `release_name` | Helm release name of the operator install (= `metadata.name`) |
| `deployment_name` | The operator controller-manager Deployment — the chart's fullname helper plus `-controller-manager` |

## Composing in Infra Charts

- **`spec.namespace`** is a foreign key (default kind
  KubernetesNamespace, field path `spec.name`).
- **KubernetesOpenSearch resources depend on this component**: the
  operator must be running and watching their namespace before their
  `OpenSearchCluster` resources reconcile. With the default
  cluster-wide watch, one operator install serves every namespace;
  with `watch_namespace` set, clusters outside that namespace are
  silently ignored by this install.
- **The install is deliberately blocking**: the Helm release waits for
  the operator Deployment to become Available (atomic, 600s timeout),
  so an unpullable image fails THIS apply instead of surfacing later
  as clusters that mysteriously never reconcile.

## Examples

### Standard cluster-wide install

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesOpenSearchOperator
metadata:
  name: opensearch-operator
spec:
  namespace:
    value: opensearch-operator
  create_namespace: true
```

### Namespace-fenced install on a shared cluster

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesOpenSearchOperator
metadata:
  name: search-team-operator
spec:
  namespace:
    value: search
  create_namespace: true
  watch_namespace: search
  use_role_bindings: true
```

### Private-mirror image

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesOpenSearchOperator
metadata:
  name: opensearch-operator
spec:
  namespace:
    value: opensearch-operator
  create_namespace: true
  image:
    repository: my-mirror.example.com/opensearchproject/opensearch-operator
  image_pull_secrets:
    - mirror-pull-secret
```

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
