# KubernetesKeda Terraform Module

Installs KEDA from the official Helm chart (`keda` at
`https://kedacore.github.io/charts`) as a real Helm release. The typed spec
renders into chart values in `locals.tf` (`local.typed_values`); the
`helm_values` escape hatch is passed as a SECOND values document the
provider merges last (Helm `-f` semantics) — the exact semantic twin of the
Pulumi module's `buildHelmValues` + `mergeMaps`.

## Module Behavior

- **The release name is FIXED to `keda`** — the component registers the
  cluster-wide `v1beta1.external.metrics.k8s.io` APIService, a singleton;
  one installation per cluster is an upstream constraint. No
  `fullnameOverride` either: the chart names its components
  `keda-operator` / `keda-operator-metrics-apiserver` /
  `keda-admission-webhooks` independent of the release name — there is
  nothing for an override to pin.
- **CRDs are kept on uninstall by default** — the chart has no native keep
  knob, so the module stamps `helm.sh/resource-policy: keep` onto the CRDs
  through `crds.additionalAnnotations`. Without it, a plain uninstall
  cascade-deletes every ScaledObject/ScaledJob/TriggerAuthentication in
  the cluster. The annotation renders only when `install && keep` — the
  keep only makes sense when this release owns the CRDs.
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 300s timeout. A KEDA that never becomes ready
  (a ServiceMonitor rendered without the Prometheus operator CRDs is THE
  classic install failure; broken internal TLS wiring the other) fails
  THIS apply with a readiness timeout instead of surfacing later as
  ScaledObjects that never scale.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a `kubernetes_namespace_v1` resource carrying the standard
  governance labels; `helm_release.create_namespace` is always false.

## Rendering Quirks

- **The chart's sizing layout is ASYMMETRIC**: replica counts live under
  each component (`operator.replicaCount`, `metricsServer.replicaCount`,
  `webhooks.replicaCount`) while container resources are grouped under ONE
  shared top-level `resources` block keyed per component — and the metrics
  server's key there is `metricServer` (SINGULAR), unlike the
  `metricsServer` component block. The module renders both halves and
  keeps the trap contained here.
- **`certificates` renders nothing for the default** — type `operator`
  (the KEDA operator self-generates certificates and patches the
  APIService caBundle) is the chart's own default and needs no values.
  With `cert_manager` and no issuer reference, the chart generates its own
  self-signed CA + Issuer chain — the issuer block stays absent; an
  explicit issuer renders `generate: false`, the name/kind, and group
  `cert-manager.io`.
- **Pod-identity arms render side by side** — they configure independent
  chart blocks (`podIdentity.aws.irsa`, `podIdentity.azureWorkload`,
  `podIdentity.gcp`).
- **One `spec.prometheus` flag fans out per component** — the chart
  mirrors its per-component layout, so the operator and metrics-server
  telemetry blocks render identically.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered values.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.keda` | `spec.create_namespace` |
| `helm_release.keda` | always |

## Usage

```bash
planton tofu apply --manifest keda.yaml
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
| `namespace` | Kubernetes namespace KEDA was installed into |
| `release_name` | Helm release name (fixed `keda` — one installation per cluster) |
| `operator_service_account_name` | The chart's fixed `keda-operator` service account — the subject cloud-side keyless bindings (IRSA trust policies, GCP WI bindings, Entra federated credentials) are written against |

## Parity

Kept in lockstep with the Pulumi module (`iac/pulumi/module/`): same chart
identity and pinned default version (2.20.1 — chart and app versions move
together), same values rendering (including the asymmetric
replicas/resources layout and the CRD keep annotation), same fixed release
name, same outputs.
