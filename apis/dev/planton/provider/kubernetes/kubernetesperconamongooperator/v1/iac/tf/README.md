# KubernetesPerconaMongoOperator Terraform Module

Installs the Percona Operator for MongoDB from the official
`psmdb-operator` Helm chart
(`https://percona.github.io/percona-helm-charts`) as a single Helm
release named after `metadata.name`. The typed spec renders into chart
values in `locals.tf` (`local.typed_values`); the `helm_values` escape
hatch is passed as a SECOND values document the provider merges over
the first with Helm `-f` semantics — the exact semantic twin of the
Pulumi module's `buildHelmValues` + `mergeMaps`.

## Module Behavior

- **The release name is `metadata.name`** — the chart derives every
  resource name from the release (Deployment, ServiceAccount, RBAC
  through its fullname helper), so distinct release names keep multiple
  namespace-scoped installations from colliding.
- **CRDs never cascade-delete databases** — the chart ships the
  PerconaServerMongoDB CRDs in its Helm-native `crds/` directory:
  installed on first install, never upgraded or deleted by Helm, so
  uninstalling the release leaves the database clusters intact (the
  upstream safety posture); no CRD knob is needed or modeled.
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout. An operator that never becomes
  ready (an unpullable image from a private mirror is the classic case)
  fails THIS apply with a readiness timeout instead of surfacing later
  as PerconaServerMongoDB resources that never reconcile.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a `kubernetes_namespace_v1` resource carrying the standard
  governance labels; `helm_release.create_namespace` is always false.

## Rendering Quirks

- **Chart-default-matching values render only on divergence** —
  `watchAllNamespaces`, `logStructured`, and `disableTelemetry` all
  default false upstream and render only when on.
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
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered YAML.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.percona_mongo_operator` | `spec.create_namespace` |
| `helm_release.percona_mongo_operator` | always |

## Usage

```bash
planton tofu apply --manifest percona-mongo-operator.yaml
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
chart identity and pinned default version (1.22.0), same
`metadata.name` release name, same values rendering (the
divergence-only rendering of chart defaults, the string-typed
`maxConcurrentReconciles`, the watch-scope mapping), same atomic/wait
posture, same outputs.
