# KubernetesPerconaMongoOperator Pulumi Module

Installs the Percona Operator for MongoDB from the official
`psmdb-operator` Helm chart
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
2. **Helm Release `<metadata.name>`** — the `psmdb-operator` chart
   (pinned default 1.22.0; chart and operator versions move together
   for this chart). The chart derives every resource name from the
   release through its fullname helper, so distinct release names keep
   multiple namespace-scoped installations from colliding

## Rendering Notes

- **Chart-default-matching values render only on divergence** —
  `watchAllNamespaces`, `logStructured`, and `disableTelemetry` all
  default false upstream and render only when on — the rendered values
  stay minimal on both engines.
- **Watch scope maps to two independent chart values** —
  `watch.cluster_wide` renders `watchAllNamespaces: true`;
  `watch.namespaces` renders the comma-joined `watchNamespace` string.
  Spec CEL rules make the two arms mutually exclusive, so at most one
  renders; both unset means the operator watches its own namespace (the
  upstream default). The chart's own `createNamespace` value (create
  the WATCHED namespaces) is never rendered — the module owns only the
  installation namespace, and watched namespaces must already exist.
- **`maxConcurrentReconciles` renders as a STRING** — the chart's own
  default is the string `"1"` (the deployment template quotes the value
  into an environment variable either way); rendering the string keeps
  both engines byte-identical with the chart's values file.
- **Pull secrets render as the raw Kubernetes object list**
  (`[{name: ...}]`) — the chart pipes `imagePullSecrets` straight into
  the pod spec with `toYaml`.
- **The image override renders only the halves that are set** — the
  chart composes `<repository>:<tag>` itself, so an unset half keeps
  the chart's default for it (repository
  `percona/percona-server-mongodb-operator`; tag = the chart version).
- **CRDs never cascade-delete databases** — the chart ships the
  PerconaServerMongoDB CRDs in its Helm-native `crds/` directory:
  installed on first install, never upgraded or deleted by Helm, so
  uninstalling the release leaves the database clusters intact (the
  upstream safety posture); no CRD knob is needed or modeled.

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and a 600s
timeout, waiting for readiness. An operator that never becomes ready
(an unpullable image from a private mirror is the classic case) fails
THIS deploy with a readiness timeout instead of surfacing later as
PerconaServerMongoDB resources that mysteriously never reconcile.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
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
  watch scope, reconciliation throughput, logging, telemetry,
  scheduling, image) and the escape-hatch merge
- `module/locals.go`: resolved namespace, release name
  (`metadata.name`), chart version — kept in lockstep with the
  Terraform module's `locals.tf`
- `module/vars.go`: chart identity (repository, `psmdb-operator`), the
  pinned default version (1.22.0), the 600s timeout
- `module/helpers.go`: shared shape renderers (resources, tolerations,
  the Helm `-f` merge)
