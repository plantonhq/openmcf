# KubernetesMetricsServer Terraform Module

Installs metrics-server from the official Helm chart (`metrics-server` at
`https://kubernetes-sigs.github.io/metrics-server/`) as a real Helm release.
The typed spec renders into chart values in `locals.tf`
(`local.typed_values`); the `helm_values` escape hatch is passed as a SECOND
values document the provider merges last (Helm `-f` semantics) — the exact
semantic twin of the Pulumi module's `buildHelmValues` + `mergeMaps`.

## Module Behavior

- **The release name is FIXED to `metrics-server`** — the component
  registers the cluster-wide `v1beta1.metrics.k8s.io` APIService, a
  singleton; one installation per cluster is an upstream constraint. The
  chart's `fullnameOverride` is pinned to the release name too, so chart
  objects (the Service the APIService routes to) get deterministic names.
- **The module owns the chart's `defaultArgs` list** — it re-renders the
  chart's default argument list with the typed substitutions
  (`kubelet_preferred_address_types`, `metric_resolution`) applied, keeping
  the pod spec canonical instead of appending duplicate flags.
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 300s timeout. The chart's `/readyz` probe only
  passes once the first kubelet scrape succeeds, so a wrong scrape-side TLS
  posture (self-signed kubelets without `kubelet_insecure_tls`) fails THIS
  apply with a readiness timeout instead of surfacing later as HPAs that
  never scale — and a green apply means metrics are actually flowing.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a `kubernetes_namespace_v1` resource carrying the standard
  governance labels; `helm_release.create_namespace` is always false.
  kube-system installs leave `create_namespace` false.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.metrics_server` | `spec.create_namespace` |
| `helm_release.metrics_server` | always |

## Usage

```bash
planton tofu apply --manifest metrics-server.yaml
```

## Local Development

```bash
terraform init
terraform validate
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Inputs

See `variables.tf` for the full variable specification (generated from the
spec proto).

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Kubernetes namespace metrics-server was installed into |
| `release_name` | Helm release name (always `metrics-server`) |
| `service_name` | The Service the APIService routes to (the chart fullname, pinned to the release name) |
| `api_service_name` | `v1beta1.metrics.k8s.io`; empty when `spec.api_service.create` is false |

## Parity

Kept in lockstep with the Pulumi module (`iac/pulumi/module/`): same chart
identity and pinned default version, same values rendering, same fixed
release name, same outputs. Conditional objects use the null-prune idiom
throughout so numbers and booleans keep their types in the rendered values.
