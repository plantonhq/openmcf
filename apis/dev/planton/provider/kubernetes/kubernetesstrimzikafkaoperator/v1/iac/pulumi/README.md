# KubernetesStrimziKafkaOperator Pulumi Module

Installs the Strimzi cluster operator from the official
`strimzi-kafka-operator` Helm chart (`https://strimzi.io/charts/`) as a
single Helm release named after `metadata.name`. The typed spec renders
into chart values in `module/values.go`; the `helm_values` escape hatch
merges LAST over them with Helm `-f` semantics (maps deep-merge, later
document wins, lists replace) — the exact semantic twin of the
Terraform module's `helm_release` with `values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true; otherwise the namespace must
   already exist
2. **Helm Release `<metadata.name>`** — the `strimzi-kafka-operator`
   chart (pinned default 1.1.0; chart and operator versions move
   together, and the SERVED index at https://strimzi.io/charts/
   governs the version — the Chart.yaml in the Strimzi source tree is
   a build-time placeholder). The chart's operand-facing resources use
   fixed Strimzi names (`strimzi-cluster-operator`), so a SECOND
   install in one cluster additionally needs
   `create_global_resources: false`

## Rendering Notes

- **Chart-default-matching values render only on divergence** — the
  watch scope, the true-defaulted toggles (`leaderElection`,
  `generateNetworkPolicy`, `generatePodDisruptionBudget`,
  `createGlobalResources`), and the image source render only when the
  spec carries a value — the rendered values stay minimal on both
  engines and an empty spec installs the chart exactly as upstream
  ships it.
- **Watch scope maps to two independent chart values** —
  `watch.any_namespace` renders `watchAnyNamespace: true`;
  `watch.namespaces` renders the `watchNamespaces` LIST (the
  installation namespace is always watched in addition). Spec CEL
  rules make the two arms mutually exclusive, so at most one renders;
  both unset means the operator watches its own namespace (the chart
  default).
- **The `resources` key renders only when the spec sets it** — the
  chart SHIPS default requests/limits (requests 200m/384Mi, limits
  1000m/384Mi), and an empty spec must leave them intact. Helm
  deep-merges per key, so a partial spec block overrides only the
  halves it carries.
- **The image override is three flat chart values** —
  `defaultImageRegistry` / `defaultImageRepository` /
  `defaultImageTag`, each rendered on presence. They steer EVERY
  Strimzi image (the operator and all operand images it deploys) — the
  air-gap path. Pull secrets ride the chart's `image.imagePullSecrets`
  list (raw Kubernetes object list, piped into the pod spec).
- **CRDs never cascade-delete Kafka clusters** — the chart ships the
  Strimzi CRDs in its Helm-native `crds/` directory: installed on
  first install, never upgraded or deleted by Helm, so uninstalling
  the release leaves the Kafka clusters intact (the upstream safety
  posture). A `chart_version` upgrade runs new operator code against
  the EXISTING CRDs — apply the new release's CRDs manually when its
  release notes call for it; no CRD knob is needed or modeled.

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and a 600s
timeout, waiting for readiness. An operator that never becomes ready
(an unpullable image from a private mirror is the classic case) fails
THIS deploy with a readiness timeout instead of surfacing later as
Kafka resources that mysteriously never reconcile.

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
  watch scope, reconciliation timing, logging, feature gates, DNS
  domain, leader election, policy toggles, RBAC ownership, scheduling,
  image) and the escape-hatch merge
- `module/locals.go`: resolved namespace, release name
  (`metadata.name`), chart version — kept in lockstep with the
  Terraform module's `locals.tf`
- `module/vars.go`: chart identity (repository,
  `strimzi-kafka-operator`), the pinned default version (1.1.0), the
  600s timeout
- `module/helpers.go`: shared shape renderers (resources, tolerations,
  the Helm `-f` merge)
