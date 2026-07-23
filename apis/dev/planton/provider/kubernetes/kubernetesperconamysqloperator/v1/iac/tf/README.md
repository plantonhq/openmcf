# KubernetesPerconaMysqlOperator Terraform Module

Installs the Percona Operator for MySQL (based on Percona XtraDB
Cluster) from the official `pxc-operator` Helm chart
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
  PerconaXtraDBCluster CRDs in its Helm-native `crds/` directory:
  installed on first install, never upgraded or deleted by Helm, so
  uninstalling the release leaves the database clusters intact (the
  upstream safety posture); no CRD knob is needed or modeled.
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout. An operator that never becomes
  ready (an unpullable image from a private mirror is the classic case)
  fails THIS apply with a readiness timeout instead of surfacing later
  as PerconaXtraDBCluster resources that never reconcile.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a `kubernetes_namespace_v1` resource carrying the standard
  governance labels; `helm_release.create_namespace` is always false.

## Rendering Quirks

- **Chart-default-matching values render only on divergence** —
  `watchAllNamespaces`, `logStructured`, `disableTelemetry`, and the
  `featureGates.xtrabackupSidecar` gate all default off upstream and
  render only when on.
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
  custom TAG requires the full `image` override (`<repository>:<tag>`),
  with the repository half falling back to the chart's own default
  (`percona/percona-xtradb-cluster-operator`) when the spec leaves it
  empty.
- **Leader election flattens to four top-level chart values** —
  `leaderElectionEnabled`, `leaseDuration`, `renewDeadline`,
  `retryPeriod` — each rendered on presence.
- **Pull secrets render as the raw Kubernetes object list**
  (`[{name: ...}]`) — the chart pipes `imagePullSecrets` straight into
  the pod spec with `toYaml`.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered YAML.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.percona_mysql_operator` | `spec.create_namespace` |
| `helm_release.percona_mysql_operator` | always |

## Usage

```bash
planton tofu apply --manifest percona-mysql-operator.yaml
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
chart identity and pinned default version (1.20.0), same
`metadata.name` release name, same values rendering (the
divergence-only rendering of chart defaults, the two-value image
override with its precedence rule, the render-only-when-set resources
key), same atomic/wait posture, same outputs.
