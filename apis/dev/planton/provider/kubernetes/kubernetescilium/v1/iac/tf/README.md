# KubernetesCilium Terraform Module

Installs Cilium from the official Helm chart (`cilium` at
`https://helm.cilium.io`) as a real Helm release. The typed spec renders
into chart values in `locals.tf` (`local.typed_values`); the `helm_values`
escape hatch is passed as a SECOND values document the provider merges last
(Helm `-f` semantics) — the exact semantic twin of the Pulumi module's
`buildHelmValues` + `mergeMaps`.

## Module Behavior

- **The release name is FIXED to `cilium`** — Cilium is the node dataplane:
  the agent DaemonSet, operator, and generated CNI configuration are
  cluster singletons, so one dataplane per cluster is an upstream
  constraint and the release name never derives from `metadata.name`.
- **No `fullnameOverride`** (unlike sibling Helm modules): the cilium chart
  names its workloads with FIXED names — DaemonSet `cilium`, Deployment
  `cilium-operator` — regardless of the release name, so there is nothing
  to pin.
- **The install waits for the whole dataplane** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout (not the usual 300): the agent
  DaemonSet must roll out on EVERY node plus the operator, and on a fresh
  cluster nodes transition NotReady→Ready only as Cilium wires each one —
  the rollout itself unblocks scheduling. A dataplane that never converges
  fails THIS apply instead of surfacing later as pods stuck in
  ContainerCreating.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a `kubernetes_namespace_v1` resource carrying the standard
  governance labels; `helm_release.create_namespace` is always false.
  kube-system installs leave `create_namespace` false.

## Rendering Quirks

- **`kubeProxyReplacement` is a STRING in the chart's values**
  (historically it took `"strict"`/`"partial"`; today `"true"`/`"false"`,
  declared default `"false"`) — the module renders `"true"` as a string,
  keeping the values document byte-identical with what the chart declares
  and with the Pulumi module. Only rendered when true.
- **`k8sServicePort` is also a string in values.yaml** (default `""`), so
  the spec's number renders through `tostring()` — the Pulumi twin uses
  `strconv.Itoa` for the same reason.
- **`hubble.metrics.enabled` is upstream's LIST of metric families** (null
  disables) — not a boolean despite the name.
- **`cni.exclusive` renders on presence, not truth**: an explicit false is
  exactly the value chaining setups must send (the spec's CEL rule enforces
  it), while unset keeps the chart default (true).
- **One `operator` map, two spec arms**: operator sizing
  (`operator.replicas`/`resources`) and operator telemetry
  (`operator.prometheus`, driven by the same `spec.prometheus` toggle as
  the agent's) are built into ONE map so the later arm never overwrites the
  earlier.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered values (the
  `merge(concat(...))` alternative silently unifies primitive-only sibling
  objects into `map(string)`).

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.cilium` | `spec.create_namespace` |
| `helm_release.cilium` | always |

## Usage

```bash
planton tofu apply --manifest cilium.yaml
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
| `namespace` | Kubernetes namespace Cilium was installed into |
| `release_name` | Helm release name (fixed `cilium` — one dataplane per cluster) |
| `cluster_name` | Cluster identity Cilium runs under (resolved `spec.cluster_name`; `default` when unset) — the name in Hubble flows and any future Cluster Mesh |
| `hubble_relay_service_name` | `hubble-relay` (fixed by the chart) when `hubble.relay` is enabled; empty otherwise |
| `hubble_ui_service_name` | `hubble-ui` (fixed by the chart) when `hubble.ui` is enabled; empty otherwise |
| `gateway_class_name` | `cilium` (fixed by the chart) when `gateway_api` is enabled; empty otherwise |

## Parity

Kept in lockstep with the Pulumi module (`iac/pulumi/module/`): same chart
identity and pinned default version (1.19.6), same values rendering
(including the string-typed `kubeProxyReplacement`/`k8sServicePort`), same
fixed release name, same outputs.
