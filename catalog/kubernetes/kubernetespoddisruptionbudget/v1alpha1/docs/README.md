# Kubernetes Pod Disruption Budget: Research Documentation

## Introduction

Kubernetes moves pods on purpose. Node drains for maintenance, cluster upgrades rolling through the fleet, the cluster autoscaler consolidating underused nodes — all of these are *voluntary disruptions*, and without a guardrail nothing stops them from taking down every replica of a service at the same moment. PodDisruptionBudget is the standard, portable answer: a namespaced `policy/v1` resource that selects pods and declares how many of them may be down at once. The eviction API refuses any eviction that would breach the budget; the drain waits and retries until replacements are ready.

The resource's mental model has two load-bearing facts that are the source of most misuse, so they are worth stating precisely:

- **Budgets govern only voluntary disruptions.** Anything going through the eviction API — `kubectl drain`, upgrades, autoscaler scale-down, descheduler moves — consults the budget. Involuntary failures (node crashes, kernel panics, OOM kills) and direct pod deletion never do. A budget is not a substitute for replicas.
- **A budget can only refuse; it cannot create availability.** `min_available: "1"` over a single-replica workload blocks every drain touching that pod. The budget expresses the floor; the workload's replica count must provide headroom above it, or drains wedge.

Planton's **KubernetesPodDisruptionBudget** component brings the full `policy/v1` surface to the platform with schema-level validation that catches the classic mistakes before apply, namespace composition, and dual-IaC support.

## Evolution and Historical Context

### Origins (policy/v1beta1)

PodDisruptionBudget entered as beta in Kubernetes 1.5 (2016) alongside the eviction subresource, and stayed in `policy/v1beta1` for an unusually long run. The original surface was `minAvailable` plus a selector; `maxUnavailable` arrived in 1.7, giving budgets a form that tracks replica count instead of a fixed floor.

### Graduation to policy/v1 (1.21) and the empty-selector flip

`policy/v1` GA'd in Kubernetes 1.21 (2021) with one semantically significant change: in v1beta1, a budget with a **null selector matched ALL pods** in the namespace; in `policy/v1`, a null selector matches **no pods**, and the *empty* selector (`{}`) is the "all pods" form. This flip is why declaring the selector explicitly matters — the same "I didn't write a selector" manifest means opposite things across the two API versions.

### Unhealthy pod eviction policy (1.26–1.31)

For most of the resource's life, a budget counted only *ready* pods as available — which meant a crash-looping application (running, never ready) could hold its budget below the floor forever and block node drains indefinitely. `unhealthyPodEvictionPolicy` was added to solve exactly this: alpha in 1.26, beta in 1.27, stable in 1.31. `IfHealthyBudget` preserves the historical conservative behavior; `AlwaysAllow` lets running-but-not-ready pods be evicted regardless, unwedging drains at the cost of restarting pods that might have been about to become ready.

### What PodDisruptionBudget never became

Upstream deliberately kept the resource small: one selector, one bound, one unhealthy-pod knob. There are no per-disruption-type budgets, no schedules ("no disruptions during business hours"), no priorities between budgets. Overlap handling is equally blunt: a pod covered by more than one budget makes evictions **fail** rather than resolving to the strictest budget — so "one budget per set of pods" is a hard rule, not a style preference.

## The Semantics in Detail

### The eviction handshake

`kubectl drain` and its programmatic equivalents do not delete pods; they POST to the pod's `eviction` subresource. The API server checks every budget selecting that pod: if the eviction would leave fewer than `minAvailable` (or more than `maxUnavailable`) available pods, the request fails with HTTP 429 and the drainer retries. This is why a correctly sized budget turns "upgrade takes the service down" into "upgrade proceeds one pod at a time."

### min_available vs max_unavailable

Both are IntOrString: an absolute count (`"2"`) or a percentage (`"50%"`) of the workload's desired replicas.

- `min_available` is the floor form. Percentages round UP when computing the floor, making them stricter. `"100%"` blocks all voluntary evictions — including node drains — and should be reserved for pods that must never move.
- `max_unavailable` is the ceiling form. Percentages round UP, giving evictions more room. `"0"` blocks all voluntary evictions. For workloads that scale, the ceiling form tracks replica count where an absolute floor goes stale: `min_available: "2"` written for 3 replicas is far too loose at 10 and fully blocking at 2.

Percentage forms require the controller to discover the *desired* replica count, which it reads through the owning controller's scale subresource — one reason budgets compose best with standard controllers (Deployments, StatefulSets, ReplicaSets) and need integer forms for bare pods.

### The selector and the workload label contract

The selector is a standard Kubernetes label selector: exact-match labels AND set-based expressions. Every Planton workload kind stamps the `app` label — set to the workload's `metadata.name` — on its pods as immutable selection identity and exports the full set as its `selector_labels` output, so `match_labels: {app: <name>}` targets exactly one workload's pods, and `match_expressions` with `tier In [web, api]` spans several.

### Where the standalone kind sits

Planton workload kinds carry a built-in `availability.pod_disruption_budget` block that derives the selector automatically from the workload's own labels — for a Planton Deployment's or StatefulSet's own pods, that block is the right tool and cannot drift. The standalone kind exists for everything else: operator-managed pods (a database operator's replicas carry the operator's labels, not a Planton workload's), non-Planton deployments, and selector-level coverage across several workloads. Pointing both at the same pods creates exactly the multi-budget overlap that makes evictions fail.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

`kubectl create poddisruptionbudget my-pdb --selector=app=checkout --min-available=1` exists and covers the basic form.

**Pros:**
- One-liner for the common case

**Cons:**
- No set-based expressions, no unhealthy-pod policy, imperative and unrecorded

**Verdict:** Fine for experiments; nothing to build on.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: checkout-pdb
  namespace: backend
spec:
  selector:
    matchLabels:
      app: checkout
  minAvailable: 1
```

**Pros:**
- Declarative, version-controllable, full surface

**Cons:**
- Nothing stops both bounds being set (rejected at admission, after the fact), a missing selector (matches nothing, silently), or overlapping budgets
- No plan/preview, no state management

**Verdict:** The baseline; the failure modes are quiet ones.

### Level 2: Terraform

```hcl
resource "kubernetes_pod_disruption_budget_v1" "checkout" {
  metadata {
    name      = "checkout-pdb"
    namespace = "backend"
  }
  spec {
    min_available = "1"
    selector {
      match_labels = { app = "checkout" }
    }
  }
}
```

**Pros:**
- Full IaC lifecycle (plan, apply, destroy, import) with drift detection

**Cons:**
- The provider does not expose `unhealthyPodEvictionPolicy` at all — the stable upstream field is simply unexpressible in HCL
- Bound format and selector contracts surface only at apply

**Verdict:** Production-grade lifecycle with a genuine surface gap.

### Level 3: Pulumi

```go
pdb, err := policyv1.NewPodDisruptionBudget(ctx, "checkout-pdb", &policyv1.PodDisruptionBudgetArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name:      pulumi.String("checkout-pdb"),
        Namespace: pulumi.String("backend"),
    },
    Spec: &policyv1.PodDisruptionBudgetSpecArgs{
        MinAvailable: pulumi.Int(1),
        Selector: &metav1.LabelSelectorArgs{
            MatchLabels: pulumi.StringMap{"app": pulumi.String("checkout")},
        },
        UnhealthyPodEvictionPolicy: pulumi.String("AlwaysAllow"),
    },
})
```

**Pros:**
- Full programming language, preview before apply, complete API surface including the unhealthy-pod policy

**Cons:**
- Types describe the wire shape, not the semantics; both-bounds and overlap mistakes pass the compiler

**Verdict:** Excellent IaC choice; the only mainstream IaC path to `unhealthyPodEvictionPolicy`.

### Other Methods

**Helm:** budgets templated per chart — common in operator and application charts; correctness depends entirely on chart authorship, and many charts still template `policy/v1beta1` semantics.

**Operators:** some operators create budgets for the pods they manage. Check before adding one — a second budget over the same pods fails evictions.

## Comparative Analysis

| Aspect | YAML | Terraform | Pulumi | Planton |
|--------|------|-----------|--------|---------|
| Validation | API server | Plan time (shape only) | Preview time (shape only) | Schema + CEL |
| Exactly-one-bound enforced early | No (admission) | No | No | Yes, pre-apply |
| Missing-selector mistake possible | Yes (matches nothing) | Yes | Yes | No — selector required |
| `unhealthyPodEvictionPolicy` | Yes | **No** | Yes | Yes (Pulumi engine; TF plan fails loudly by design) |
| Namespace as reference | No | Manual wiring | Manual wiring | First-class |
| Dual IaC | N/A | TF only | Pulumi only | Both |

## The Planton Approach

### Full surface, validated early

The spec models the entire `policy/v1` PodDisruptionBudgetSpec — selector, both bound forms, unhealthy-pod eviction policy — and moves the API server's rules plus the known footguns to validation time:

- **Exactly one bound**: the API rejects a budget with both `minAvailable` and `maxUnavailable`; the schema additionally rejects a budget with *neither*, because it protects nothing
- **Int-or-percent contracts**: bounds must be an absolute number or a `0%`–`100%` percentage, checked by format before deployment
- **Required selector**: the `policy/v1` null-selector trap (matches no pods) cannot be expressed; the "all pods" intent must be declared as the explicit empty selector
- **Selector operator contracts**: `In`/`NotIn` require values, `Exists`/`DoesNotExist` forbid them — the exact admission rule, surfaced before deployment

### Deterministic unhealthy-pod policy — and the one parity exception

The Pulumi module always submits `unhealthyPodEvictionPolicy` explicitly, applying the server default (`IfHealthyBudget`) module-side, so the deployed object never depends on server-side defaulting.

The Terraform kubernetes provider cannot express the field at all. The Terraform module therefore carries a plan-time **precondition** that fails any spec requesting `always_allow`, with an error directing the user to the Pulumi provisioner (or to dropping the field). Failing loudly is the design: the server default is exactly what a non-default value would override, so silently deploying `IfHealthyBudget` where `always_allow` was requested would deploy the opposite of the user's intent. For every spec the Terraform module accepts, both engines produce identical deployed objects.

### Namespace by value or reference

`spec.namespace` is a `StringValueOrRef`: a literal namespace name, or a reference to a `KubernetesNamespace` resource, letting an infra chart create a namespace, its workloads, and their budgets in one run with ordering handled by the resource graph. When omitted, the budget lands in `default`.

### The workload label contract

Every Planton workload kind stamps the `app` label — set to the workload's `metadata.name` — on its pods as immutable selection identity. `match_labels: {app: <workload-name>}` is therefore the standard way standalone budgets compose with Planton workloads, and `selector_labels` outputs expose the full set for exact composition.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates provider init, resource creation, and output export
- **`locals.go`**: Computes merged labels, annotations, the resolved namespace, and the resolved unhealthy-pod policy (server default applied)
- **`poddisruptionbudget.go`**: Creates the `policy/v1` PodDisruptionBudget, rendering the selector always (empty block = all pods) and the bounds as IntOrString
- **`outputs.go`**: Exports `pod_disruption_budget_name` and `namespace`

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto` fields as Terraform variables
- **`locals.tf`**: Computes merged labels, resolved namespace, and the resolved unhealthy-pod policy string
- **`main.tf`**: Creates the `kubernetes_pod_disruption_budget_v1` resource, with the parity-exception precondition rejecting `always_allow`
- **`outputs.tf`**: Exports the same two outputs

### Resource Count

This is a lean component — it creates exactly **one Kubernetes resource**: the PodDisruptionBudget itself. The complexity is in the spec validation and the semantic guardrails, not in resource orchestration.

## Production Best Practices

### Sizing discipline

1. **Leave every drain a way to succeed**: at steady state, the bound must permit at least one eviction (`min_available` at most replicas−1, or `max_unavailable` at least 1). A zero-disruption budget turns node maintenance into a paging incident
2. **Prefer `max_unavailable` for anything that scales**: ceilings track replica count; absolute floors go stale in both directions
3. **Reserve `"100%"` / `"0"` for pods that genuinely must never move** — and expect drains touching them to require manual intervention

### Placement discipline

1. **One budget per set of pods**: overlapping budgets fail evictions instead of combining. Check whether an operator already creates one before adding yours
2. **Use the built-in block for Planton workloads' own pods**: it derives the selector automatically; the standalone kind is for operator-managed pods, non-Planton workloads, and cross-workload selectors
3. **Budgets live with their pods**: a budget governs only its own namespace

### Unhealthy-pod discipline

1. **Set `always_allow` on crash-loop-prone workloads**: a running-but-never-ready pod under the default policy holds its budget below the floor forever and wedges drains
2. **Remember the engine boundary**: `always_allow` deploys via the Pulumi provisioner; the Terraform module rejects it at plan time by design rather than silently deploying the default

## Conclusion

KubernetesPodDisruptionBudget is a deliberately complete, deliberately lean component: the full `policy/v1` surface, with the resource's quiet failure modes — the null-selector trap, the both-bounds rejection, the budget-with-no-bound that protects nothing, the crash-loop drain wedge — documented in the schema and guarded by validation before anything reaches a cluster. The one place the two IaC engines genuinely differ (`unhealthyPodEvictionPolicy` is inexpressible in the Terraform provider) fails loudly at plan time instead of deploying the wrong object. Combined with the workload label contract and namespace references, it makes availability floors a pattern that can be stamped out per workload rather than rediscovered per stuck drain.

## References

- [Kubernetes Disruptions Documentation](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/)
- [Specifying a Disruption Budget for your Application](https://kubernetes.io/docs/tasks/run-application/configure-pdb/)
- [PodDisruptionBudget API Reference](https://kubernetes.io/docs/reference/kubernetes-api/policy-resources/pod-disruption-budget-v1/)
- [API-initiated Eviction](https://kubernetes.io/docs/concepts/scheduling-eviction/api-eviction/)
- [Pulumi Kubernetes PodDisruptionBudget](https://www.pulumi.com/registry/packages/kubernetes/api-docs/policy/v1/poddisruptionbudget/)
- [Terraform kubernetes_pod_disruption_budget_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/pod_disruption_budget_v1)
