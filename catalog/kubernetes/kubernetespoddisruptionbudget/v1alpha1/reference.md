# KubernetesPodDisruptionBudget

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**KubernetesPodDisruptionBudgetSpec** limits how many of a set of pods may
be taken down VOLUNTARILY at once — node drains, cluster upgrades,
autoscaler consolidation. The eviction API refuses to evict a selected pod
when doing so would breach the budget; involuntary failures (node crashes,
OOM kills) are not governed.

A standalone budget is the right shape for pods a Planton workload kind
does not manage: operator-managed pods (a database operator's replicas),
non-Planton deployments, or selector-level coverage across several
workloads. For a Planton Deployment/StatefulSet's OWN pods, prefer the
workload's built-in `availability.pod_disruption_budget` block — it derives
the selector automatically and cannot drift from the workload's labels.

Selectors compose with Planton workloads through the documented label
contract: every workload kind stamps the `app` label (set to the workload's
`metadata.name`) on its pods and exports the full set as its
`selector_labels` output — so `match_labels: {app: <name>}` targets one
workload's pods.

The spec covers the complete policy/v1 PodDisruptionBudgetSpec surface.

## Example

```yaml
# Full-surface test manifest for the offline plan proofs. Exercises the
# selector (labels + expressions) and a percentage bound.
# (unhealthy_pod_eviction_policy: always_allow is exercised by the
# Pulumi-only negative proof — the terraform module rejects it with a
# precondition by design.)
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesPodDisruptionBudget
metadata:
  name: test-pdb
spec:
  namespace:
    value: default
  name: test-pdb
  labels:
    team: platform-engineering
  selector:
    match_labels:
      app: test-app
    match_expressions:
      - key: tier
        operator: In
        values:
          - web
          - api
  max_unavailable: 25%
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` |  |  | KubernetesNamespace (`spec.name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.labels` | `map<string, string>` |  |  |  |
| `spec.annotations` | `map<string, string>` |  |  |  |
| `spec.selector` | `KubernetesPodDisruptionBudgetLabelSelector` | yes |  |  |
| `spec.selector.matchLabels` | `map<string, string>` |  |  |  |
| `spec.selector.matchExpressions` | `[]KubernetesPodDisruptionBudgetLabelSelectorRequirement` |  |  |  |
| `spec.selector.matchExpressions[].key` | `string` | yes |  |  |
| `spec.selector.matchExpressions[].operator` | `string` |  |  |  |
| `spec.selector.matchExpressions[].values` | `[]string` |  |  |  |
| `spec.minAvailable` | `string` |  |  |  |
| `spec.maxUnavailable` | `string` |  |  |  |
| `spec.unhealthyPodEvictionPolicy` | `enum` |  | `if_healthy_budget` |  |

## Field Details

### spec.namespace

`string | valueFrom`

The namespace the budget lives in — it governs only pods in its own
namespace. Accepts a literal namespace name or a reference to a
KubernetesNamespace resource. When omitted, the budget lands in the
cluster's `default` namespace.

- references: KubernetesNamespace (`spec.name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the PodDisruptionBudget (its `metadata.name` in the
cluster). Must be a valid DNS subdomain: lowercase alphanumeric
characters, hyphens, and dots, at most 253 characters.

- rule: Name must be a valid DNS subdomain (lowercase alphanumeric, hyphens, and dots, no leading/trailing dots or hyphens)
- rule: {"string":{"minLen":"1","maxLen":"253"}}

### spec.labels

`map<string, string>`

Additional labels to apply to the PodDisruptionBudget object.
Merged with the standard Planton governance labels.

### spec.annotations

`map<string, string>`

Annotations to apply to the PodDisruptionBudget object.

### spec.selector

`KubernetesPodDisruptionBudgetLabelSelector` · required

Selects the pods this budget protects, within the budget's namespace.
Required: an EMPTY selector (present but with no match_labels or
match_expressions) protects ALL pods in the namespace — declare that
shape explicitly; a budget with NO selector protects nothing in policy/v1
and is always a mistake. To protect one Planton workload's pods, match on
its `app` label: `match_labels: {app: <workload-name>}`.

- rule: {"required":true}

### spec.selector.matchLabels

`map<string, string>`

Exact-match label requirements: every key/value pair listed must be
present, e.g. {"app": "checkout"}.

### spec.selector.matchExpressions

`[]KubernetesPodDisruptionBudgetLabelSelectorRequirement`

Set-based label requirements, for selections exact-match cannot express
(key existence, value-in-set).

- rule: values must be non-empty for In/NotIn and empty for Exists/DoesNotExist

### spec.selector.matchExpressions[].key

`string` · required

The label key the requirement applies to.

- rule: {"string":{"minLen":"1"}}

### spec.selector.matchExpressions[].operator

`string`

The relationship between the key and the values: "In" (value must be one
of `values`), "NotIn" (must not be), "Exists" (key present, `values` must
be empty), or "DoesNotExist" (key absent, `values` must be empty).

- rule: {"string":{"in":["In","NotIn","Exists","DoesNotExist"]}}

### spec.selector.matchExpressions[].values

`[]string`

The values compared against the label's value. Required (non-empty) for
In/NotIn; must be empty for Exists/DoesNotExist.

### spec.minAvailable

`string`

Minimum pods that must stay available — absolute ("2") or percentage
("50%", of desired replicas). Percentages round UP when computing the
floor. "100%" blocks ALL voluntary evictions — it also blocks node
drains, so use it only for pods that must never move. Set exactly one of
min_available or max_unavailable.

- rule: min_available must be an absolute number ("2") or a percentage between 0% and 100% ("50%")

### spec.maxUnavailable

`string`

Maximum pods that may be down — absolute or percentage (percentages
round UP, giving evictions more room). "0" blocks all voluntary
evictions. Alternative to min_available; do not set both. For workloads
that scale, prefer max_unavailable — it tracks replica count where an
absolute min_available floor goes stale.

- rule: max_unavailable must be an absolute number ("1") or a percentage between 0% and 100% ("25%")

### spec.unhealthyPodEvictionPolicy

`enum` · optional (explicit presence)

How running-but-not-ready pods are treated by the eviction API.
Default: if_healthy_budget

- default: `if_healthy_budget`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `kubernetes_pod_disruption_budget_unhealthy_pod_eviction_policy_unspecified` -- Unspecified. Defaults to if_healthy_budget.
- `if_healthy_budget` -- Not-yet-ready pods may be evicted only while the healthy count meets the budget — the default, most conservative behavior.
- `always_allow` -- Not-yet-ready pods may ALWAYS be evicted. Prevents a crash-looping application from wedging node drains forever — the practical choice for budgets over workloads that can crash-loop.

## Validation Rules

- `availability.exactly_one`: Set exactly one of min_available or max_unavailable

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesPodDisruptionBudget, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.pod_disruption_budget_name` | `string` | The name of the PodDisruptionBudget object as created in the cluster. |
| `status.outputs.namespace` | `string` | The namespace the budget was created in. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |

## See Also

- [Overview](../README.md)
