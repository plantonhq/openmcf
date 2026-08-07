# KubernetesCloudNativePgOperator Terraform Module

Installs CloudNativePG from the official Helm charts
(`https://cloudnative-pg.github.io/charts` — one repository serves both
charts) as up to TWO real Helm releases in the same namespace. The typed
spec renders into operator-chart values in `locals.tf`
(`local.typed_values`); the `helm_values` escape hatch is passed as a
SECOND values document the provider merges over the first with Helm `-f`
semantics — the exact semantic twin of the Pulumi module's
`buildHelmValues` + `mergeMaps`.

## Module Behavior

- **The operator release name is FIXED to `cnpg`** — the operator
  registers cluster-scoped CRDs and mutating/validating webhooks whose
  service name is baked into the chart (`cnpg-webhook-service` —
  embedded in the webhook certificate and not configurable); one
  installation per cluster is an upstream constraint and the name never
  derives from `metadata.name`.
- **The plugin is its own release** (`plugin-barman-cloud`, when
  `barman_cloud_plugin.enabled`) — upstream forbids folding it into the
  operator's release (Helm ownership of shared resources would
  conflict). It installs AFTER the operator (`depends_on`) so its CNPG-I
  registration lands on a running operator; destroy unwinds in reverse.
  Its name is fixed for the same singleton reason (the plugin's gRPC
  service name `barman-cloud` is baked into its TLS certificate), and it
  carries its OWN chart pin (0.7.0 = plugin v0.13.0) — the plugin chart
  versions independently of the operator chart (0.29.0 = operator
  1.30.0).
- **CRDs keep the databases safe by upstream policy** — the chart stamps
  `helm.sh/resource-policy: keep` on every CRD unconditionally, so
  uninstalling never cascade-deletes the Cluster resources; `crds.create`
  renders only on explicit opt-out.
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout on both releases. A PodMonitor
  rendered without the Prometheus operator CRDs fails the operator
  release; a plugin installed without cert-manager (its
  Issuer/Certificate resources render unconditionally) fails the plugin
  release — both roll back cleanly instead of surfacing later as Cluster
  resources that never reconcile.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a `kubernetes_namespace_v1` resource carrying the standard
  governance labels; `helm_release.create_namespace` is always false.

## Rendering Quirks

- **The typed watch field OWNS `WATCH_NAMESPACE`** — a user entry under
  that key in `operator_config` is always stripped; the key renders only
  from `watch.namespaces` (comma-joined) when `cluster_wide` is false
  (spec CEL guarantees the pairing).
- **One `config` block for three concerns** — `clusterWide` (rendered
  only when fencing), `data`, and `maxConcurrentReconciles`, matching
  the chart's own folding.
- **`helm_values` scopes to the OPERATOR chart only** — the plugin's
  values document renders from its typed resources alone; the two charts
  share value keys (`resources`, `image`), so forwarding one document to
  both would misconfigure the plugin.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered YAML.
- **No `fullnameOverride`** — the chart hard-codes the names that matter;
  there is nothing for an override to pin.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.cloudnative_pg` | `spec.create_namespace` |
| `helm_release.cloudnative_pg` | always |
| `helm_release.barman_cloud_plugin` | `spec.barman_cloud_plugin.enabled` |

## Usage

```bash
planton tofu apply --manifest cloudnative-pg-operator.yaml
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
| `namespace` | Namespace the operator (and the plugin, when enabled) runs in |
| `release_name` | Helm release name of the operator (fixed `cnpg` — one installation per cluster) |
| `barman_plugin_release_name` | Helm release name of the plugin when enabled (`plugin-barman-cloud`); empty otherwise |

## Parity

Kept in lockstep with the Pulumi module (`../pulumi/module/`): same
chart identities and pinned default versions (operator 0.29.0 = 1.30.0;
plugin 0.7.0 = v0.13.0), same fixed release names, same values rendering
(the config folding, the WATCH_NAMESPACE precedence, the
divergence-only rendering of chart defaults), same atomic/wait posture,
same outputs.
