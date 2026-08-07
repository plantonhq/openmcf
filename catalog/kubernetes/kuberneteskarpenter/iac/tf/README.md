# KubernetesKarpenter Terraform Module

Installs Karpenter from the official OCI-served Helm charts
(`oci://public.ecr.aws/karpenter/karpenter` and the companion
`karpenter-crd` chart) as TWO real Helm releases in the same namespace. The
typed spec renders into controller-chart values in `locals.tf`
(`local.typed_values`); the `helm_values` escape hatch is passed as a
SECOND values document the provider merges last (Helm `-f` semantics) — the
exact semantic twin of the Pulumi module's `buildHelmValues` + `mergeMaps`.

## Module Behavior

- **Two releases, fixed names** — `karpenter-crd` (when `crds.install`,
  default true) then `karpenter`, which `depends_on` it. Karpenter owns
  the cluster-wide `karpenter.sh` label domain, its CRDs, and node
  lifecycle — one installation per cluster is an upstream constraint, so
  the names never derive from `metadata.name`.
- **OCI chart resolution (engine asymmetry)** — this provider takes
  `repository = "oci://public.ecr.aws/karpenter"` plus the bare chart name
  and joins them internally; Pulumi's `helm.v3.Release` needs the JOINED
  `oci://public.ecr.aws/karpenter/<chart>` reference with no repository
  opts. Same chart bytes, different wiring — both modules carry the split
  with comments.
- **CRDs live in their own release** — upstream's supported mechanism for
  keeping CRDs upgradable: Helm installs the copies bundled inside the
  main chart once and NEVER upgrades them. The controller release
  therefore sets `skip_crds = true` UNCONDITIONALLY, keeping its shape
  deterministic whether or not `crds.install` is on.
- **CRDs are kept on uninstall by default** — `keep_on_uninstall` (default
  true) stamps `helm.sh/resource-policy: keep` onto every CRD through the
  CRD chart's `additionalAnnotations`. Without it, a plain uninstall
  cascade-deletes every NodePool/EC2NodeClass/NodeClaim in the cluster.
- **Readiness is verified at install time** — both releases use `wait` +
  `atomic` + `cleanup_on_fail` with a 600s timeout. A Karpenter that never
  becomes ready (a ServiceMonitor rendered without the Prometheus operator
  CRDs, a bad IRSA trust policy) fails THIS apply instead of surfacing
  later as pods that stay Pending forever because no nodes ever appear.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a `kubernetes_namespace_v1` resource carrying the standard
  governance labels; `helm_release.create_namespace` is always false.
  Leave it false when installing into `kube-system`, upstream's
  recommended home.

## Rendering Notes

- **Defaults are applied module-side**: fields whose spec default mirrors
  the chart default (replicas, logLevel, batching windows, scheduling
  policies, featureGates, priorityClassName, reservedENIs,
  vmMemoryOverheadPercent) render with the default APPLIED, so both
  engines produce byte-identical values whether or not the platform's
  defaulting middleware ran.
- **Type fidelity**: the chart's `settings.reservedENIs` default is the
  STRING `"0"` — the spec's number renders through `tostring()`. The
  chart's `settings.vmMemoryOverheadPercent` default is the NUMBER `0.075`
  — the spec's string renders through `tonumber()`. Rendering the other
  type would change the values document the chart sees.
- **The whole `featureGates` map renders** with defaults applied — the
  deployment template composes the `FEATURE_GATES` env var from all six
  keys unconditionally, and `reservedCapacity` defaults TRUE unlike the
  other five gates.
- **`nodeSelector` merges, `tolerations` replace**: Helm deep-merges maps,
  so spec entries narrow the chart's `kubernetes.io/os: linux` default —
  but lists replace wholesale, so tolerations render only when the spec
  provides them (an empty list would DROP the chart's CriticalAddonsOnly
  default).
- **IRSA rides the service-account annotation**
  (`eks.amazonaws.com/role-arn`), not a settings entry; empty means EKS
  Pod Identity, configured cloud-side with no annotation needed.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered values.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.karpenter` | `spec.create_namespace` |
| `helm_release.crds` | `spec.crds.install` (default true) |
| `helm_release.controller` | always |

## Usage

```bash
planton tofu apply --manifest karpenter.yaml
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
| `namespace` | Kubernetes namespace Karpenter was installed into |
| `release_name` | Controller Helm release name (fixed `karpenter`) |
| `crd_release_name` | CRD Helm release name (fixed `karpenter-crd`; empty when `crds.install` is false) |
| `service_account_name` | The controller's service account (fixed `karpenter` — SA name defaults to the chart fullname, which resolves to the release name because the release name contains the chart name); the subject IRSA trust policies and EKS Pod Identity associations are written against |

## Parity

Kept in lockstep with the Pulumi module (`iac/pulumi/module/`): same chart
identity and pinned default version (1.14.0 — both charts version with the
controller), same values rendering (including both type-fidelity
conversions and the always-rendered featureGates map), same fixed release
names, same unconditional skip of the controller chart's bundled CRDs,
same outputs.
