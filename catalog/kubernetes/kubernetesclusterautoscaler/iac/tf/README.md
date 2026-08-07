# KubernetesClusterAutoscaler Terraform Module

Installs the Kubernetes Cluster Autoscaler from the official Helm chart
(`cluster-autoscaler` at `https://kubernetes.github.io/autoscaler`) as a
real Helm release. The typed spec renders into chart values in `locals.tf`
(`local.typed_values`); the `helm_values` escape hatch is passed as a
SECOND values document the provider merges last (Helm `-f` semantics) —
the exact semantic twin of the Pulumi module's `buildHelmValues` +
`mergeMaps`.

## Module Behavior

- **The release name is FIXED to `cluster-autoscaler`** — the autoscaler
  leader-elects and owns the cluster-wide scaling decision; a second
  installation would fight the first over every scale-up, so one
  installation per cluster is the operating model.
- **Exactly one provider arm renders** — the proto oneof guarantees one of
  aws/azure/gce/cluster_api/civo/kwok; `cloudProvider` selects the chart's
  per-provider command/env blocks and which keys its credential Secret
  carries.
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout. An autoscaler that never becomes
  ready (bad cloud credentials crash-looping the pod is THE classic
  failure; a ServiceMonitor rendered without the Prometheus operator CRDs
  the other) fails THIS apply with a readiness timeout instead of
  surfacing later as node groups that never scale.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a `kubernetes_namespace_v1` resource carrying the standard
  governance labels; `helm_release.create_namespace` is always false.
  kube-system installs (the upstream convention) leave it false.

## Rendering Quirks

- **The Deployment has a render gate** (`templates/deployment.yaml` line
  1): it only renders when `autoDiscovery.clusterName` /
  `autoDiscovery.namespace` / `autoDiscovery.labels` or
  `autoscalingGroups` is set. `autoscalingGroupsnamePrefix` — the GCE
  contract — does NOT satisfy the gate, and kwok sets no gate key at all.
  The gce and kwok arms therefore render
  `autoDiscovery.clusterName = metadata.name` (a benign, deterministic
  value no gce/kwok template consumes beyond the gate; the chart README
  documents "any-name" as required for GCE). Without it the release would
  "succeed" while installing NO autoscaler pod.
- **Civo node pools ride `helm_values`** — the chart requires
  `autoscalingGroups` for civo (its README), but the typed spec does not
  model civo pools; supplying them through `helm_values` also satisfies
  the render gate.
- **Cluster API discovery needs at least one dimension** — the typed
  `namespace` field satisfies the gate when set; with it empty, ship
  `autoDiscovery.clusterName` or `autoDiscovery.labels` via `helm_values`
  per upstream's Cluster API guidance.
- **Credentials become the chart's own Secret** — `awsSecretAccessKey`,
  `azureClientSecret`, `civoApiKey` (and their non-secret siblings) land
  in chart values which the chart materializes into
  `templates/secret.yaml` and wires to the pod as env vars via
  `secretKeyRef`. They flow ONLY through the values documents — never
  into logs or outputs.
- **`extraArgs` precedence** — the typed `scaling` block renders first,
  `extra_args` merges over it (later `merge()` argument wins per key),
  and the chart's own extraArgs defaults
  (`logtostderr`/`stderrthreshold`/`v`) survive chart-side because Helm
  coalesces the provided map over the chart default per key. Values
  render as strings on both engines — they are CLI flag text.
- **Azure workload identity sets only the extension flag** — VERIFIED: no
  chart template adds the `azure.workload.identity/use` pod label;
  clusters relying on the azure-workload-identity webhook add `podLabels`
  via `helm_values`.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered values;
  cross-arm string picks use `try(coalesce(...), null)` instead of
  chained per-arm ternaries.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.cluster_autoscaler` | `spec.create_namespace` |
| `helm_release.cluster_autoscaler` | always |

## Usage

```bash
planton tofu apply --manifest cluster-autoscaler.yaml
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
| `namespace` | Kubernetes namespace the autoscaler was installed into |
| `release_name` | Helm release name (fixed `cluster-autoscaler` — one installation per cluster) |
| `service_account_name` | The autoscaler's service account, derived from the chart's fullname template: `cluster-autoscaler-<cloudProvider>-cluster-autoscaler` (verified in `_helpers.tpl` — the default name is `<cloudProvider>-<chartName>`, which never equals the release name, so fullname prefixes the release). The subject cloud-side keyless bindings are written against |

## Parity

Kept in lockstep with the Pulumi module (`iac/pulumi/module/`): same chart
identity and pinned default version (9.59.0 — ships autoscaler 1.35.0;
chart and app versions move separately), same values rendering (including
the deployment render gate and the extraArgs precedence), same fixed
release name, same derived service-account output.
