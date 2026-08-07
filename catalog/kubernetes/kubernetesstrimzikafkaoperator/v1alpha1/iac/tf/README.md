# KubernetesStrimziKafkaOperator Terraform Module

Installs the Strimzi cluster operator from the official
`strimzi-kafka-operator` Helm chart (`https://strimzi.io/charts/`) as a
single Helm release named after `metadata.name`. The typed spec renders
into chart values in `locals.tf` (`local.typed_values`); the
`helm_values` escape hatch is passed as a SECOND values document the
provider merges over the first with Helm `-f` semantics — the exact
semantic twin of the Pulumi module's `buildHelmValues` + `mergeMaps`.

## Module Behavior

- **The release name is `metadata.name`** — but the chart's
  operand-facing resources use fixed Strimzi names
  (`strimzi-cluster-operator`), so a SECOND install in one cluster
  additionally needs `create_global_resources: false` (the fixed-name
  ClusterRoles are owned by the first release).
- **CRDs never cascade-delete Kafka clusters** — the chart ships the
  Strimzi CRDs in its Helm-native `crds/` directory: installed on
  first install, never upgraded or deleted by Helm, so uninstalling
  the release leaves the Kafka clusters intact (the upstream safety
  posture). The same posture means a `chart_version` upgrade runs new
  operator code against the EXISTING CRDs — apply the new release's
  CRDs manually when its release notes call for it; no CRD knob is
  needed or modeled.
- **The chart version is the SERVED version** — the pinned default
  (1.1.0) must exist in the repository index at
  https://strimzi.io/charts/; the Chart.yaml inside the Strimzi source
  tree carries a build-time placeholder and never reflects the served
  version. Chart and operator versions move together for this chart.
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout. An operator that never
  becomes ready (an unpullable image from a private mirror is the
  classic case) fails THIS apply with a readiness timeout instead of
  surfacing later as Kafka resources that never reconcile.
- **The module (not Helm) owns namespace creation** —
  `create_namespace` drives a `kubernetes_namespace_v1` resource
  carrying the standard governance labels;
  `helm_release.create_namespace` is always false.

## Rendering Quirks

- **Chart-default-matching values render only on divergence** — the
  watch scope, the true-defaulted toggles (`leaderElection`,
  `generateNetworkPolicy`, `generatePodDisruptionBudget`,
  `createGlobalResources`), and the image source all render only when
  the spec carries a value, so an empty spec installs the chart
  exactly as upstream ships it.
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
  air-gap path.
- **Pull secrets render as the raw Kubernetes object list**
  (`image.imagePullSecrets: [{name: ...}]`) — the chart pipes the list
  straight into the pod spec.
- **Leader election nests a single enable flag**
  (`leaderElection.enable`), rendered on presence.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered YAML.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.strimzi_kafka_operator` | `spec.create_namespace` |
| `helm_release.strimzi_kafka_operator` | always |

## Usage

```bash
planton tofu apply --manifest strimzi-kafka-operator.yaml
```

## Local Development

```bash
terraform init
terraform validate
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Inputs

See `variables.tf` for the full variable specification (generated from
the spec proto). The spec arrives from the proto→tfvars converter in
snake_case with the `namespace` foreign key (KubernetesNamespace)
resolved to a literal string before Terraform runs.

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (`metadata.name`) |

## Parity

Kept in lockstep with the Pulumi module (`../pulumi/module/`): same
chart identity and pinned default version (1.1.0), same
`metadata.name` release name, same values rendering (the
divergence-only rendering of chart defaults, the two-arm watch scope,
the render-only-when-set resources key, the flat image override), same
atomic/wait posture, same outputs.
