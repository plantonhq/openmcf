# KubernetesPerconaMysqlOperator Pulumi Module

Installs the Percona Operator for MySQL (based on Percona XtraDB
Cluster) from the official `pxc-operator` Helm chart
(`https://percona.github.io/percona-helm-charts`) as a single Helm
release named after `metadata.name`. The typed spec renders into chart
values in `module/values.go`; the `helm_values` escape hatch merges
LAST over them with Helm `-f` semantics (maps deep-merge, later
document wins, lists replace) — the exact semantic twin of the
Terraform module's `helm_release` with `values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true; otherwise the namespace must
   already exist
2. **Helm Release `<metadata.name>`** — the `pxc-operator` chart
   (pinned default 1.20.0; chart and operator versions move together
   for this chart). The chart derives every resource name from the
   release through its fullname helper, so distinct release names keep
   multiple namespace-scoped installations from colliding

## Rendering Notes

- **Chart-default-matching values render only on divergence** —
  `watchAllNamespaces`, `logStructured`, `disableTelemetry`, and the
  `featureGates.xtrabackupSidecar` gate all default off upstream and
  render only when on — the rendered values stay minimal on both
  engines.
- **Watch scope maps to two independent chart values** —
  `watch.cluster_wide` renders `watchAllNamespaces: true`;
  `watch.namespaces` renders the comma-joined `watchNamespace` string.
  Spec CEL rules make the two arms mutually exclusive, so at most one
  renders; both unset means the operator watches its own namespace (the
  upstream default). The chart's own `createNamespace` value (create
  the WATCHED namespaces) is never rendered — the module owns only the
  installation namespace, and watched namespaces must already exist.
- **The `resources` key renders only when the spec sets it** — the
  chart SHIPS default requests/limits (requests 100m/20Mi, limits
  200m/500Mi), and an empty spec must leave them intact. Helm
  deep-merges per key, so a partial spec block overrides only the
  halves it carries.
- **The image override is TWO chart values with a precedence rule** —
  the chart's image helper uses `image` verbatim when non-empty, else
  `<operatorImageRepository>:<chart app version>`. A repository alone
  maps to `operatorImageRepository` (tag stays the chart version); any
  custom TAG requires the full `image` override
  (`<repository>:<tag>`), with the repository half falling back to the
  chart's own default (`percona/percona-xtradb-cluster-operator`) when
  the spec leaves it empty.
- **Leader election flattens to four top-level chart values** —
  `leaderElectionEnabled`, `leaseDuration`, `renewDeadline`,
  `retryPeriod` — each rendered on presence.
- **Pull secrets render as the raw Kubernetes object list**
  (`[{name: ...}]`) — the chart pipes `imagePullSecrets` straight into
  the pod spec with `toYaml`.
- **CRDs never cascade-delete databases** — the chart ships the
  PerconaXtraDBCluster CRDs in its Helm-native `crds/` directory:
  installed on first install, never upgraded or deleted by Helm, so
  uninstalling the release leaves the database clusters intact (the
  upstream safety posture); no CRD knob is needed or modeled.

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and a 600s
timeout, waiting for readiness. An operator that never becomes ready
(an unpullable image from a private mirror is the classic case) fails
THIS deploy with a readiness timeout instead of surfacing later as
PerconaXtraDBCluster resources that mysteriously never reconcile.

## Usage

```shell
planton pulumi up --manifest hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (`metadata.name`) |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → operator release → output exports
- `module/values.go`: typed-spec → chart values rendering (sizing,
  watch scope, reconciliation throughput, logging, telemetry, leader
  election, feature gates, scheduling, image) and the escape-hatch
  merge
- `module/locals.go`: resolved namespace, release name
  (`metadata.name`), chart version — kept in lockstep with the
  Terraform module's `locals.tf`
- `module/vars.go`: chart identity (repository, `pxc-operator`), the
  pinned default version (1.20.0), the chart's default image
  repository (for the tag-only override), the 600s timeout
- `module/helpers.go`: shared shape renderers (resources, tolerations,
  the Helm `-f` merge)
