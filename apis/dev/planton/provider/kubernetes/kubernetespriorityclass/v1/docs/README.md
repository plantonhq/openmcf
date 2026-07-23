# Kubernetes Priority Class: Research Documentation

## Introduction

A Kubernetes scheduler with no priority information treats every pod as equally important: when the cluster runs out of room, whatever arrives next simply waits, regardless of whether it is a revenue-path API server or a nightly report job. PriorityClass is the standard answer — a cluster-scoped `scheduling.k8s.io/v1` resource that names an integer. Pods reference the class by name (`priorityClassName`, surfaced in Planton's shared workload pod spec as `priority_class_name`) and receive its value, and the value drives two scheduler behaviors:

- **Queue ordering**: among Pending pods, higher priority schedules first.
- **Preemption**: when a higher-priority pod cannot schedule anywhere, the scheduler EVICTS lower-priority pods to make room — unless the class disables it with `preemptionPolicy: Never`, in which case its pods jump the queue but never displace anything already running.

The resource itself is tiny — a name, an integer, two flags, and a description — but it is the kind whose misconfiguration is felt cluster-wide: a duplicate global default silently changes what every unmarked pod gets, and an over-eager preempting tier evicts production pods. Planton's **KubernetesPriorityClass** component brings the full surface to the platform with schema-level validation for the value ceiling and the reserved name prefix, deterministic defaults across both IaC engines, and safe replacement semantics for the immutable value.

## Evolution and Historical Context

### Origins (Kubernetes 1.8–1.14)

Pod priority and preemption entered as alpha in Kubernetes 1.8 (2017), motivated by cluster-autoscaling economics: without eviction, guaranteeing capacity for critical services means permanently over-provisioning. The feature went beta in 1.11 — enabled by default — and GA in 1.14 as `scheduling.k8s.io/v1`. The surface has been essentially frozen since: name, integer value, `globalDefault`, `description`.

### System classes and the reserved range

Kubernetes ships two built-in classes — `system-cluster-critical` and `system-node-critical` — whose values (2000000000 and 2000001000) sit ABOVE the user-definable ceiling of one billion. The reservation is structural: user workloads must never outrank the control-plane and node components that keep the cluster alive, so the API rejects user classes above 1,000,000,000, and the `system-` name prefix is reserved for Kubernetes itself.

### Non-preempting classes (1.19)

`preemptionPolicy: Never` graduated in Kubernetes 1.24 (beta and enabled by default since 1.19), completing the model. It decouples the two things priority does — queue position and eviction rights — and exists for exactly one pattern: high-priority BATCH work. A data-science job may deserve to schedule ahead of other pending jobs, but should never evict a running service to do so. Before non-preempting classes, that trade-off was inexpressible.

### What PriorityClass never became

Upstream deliberately kept priority a scheduling concept, not a runtime one. Priority does not change CPU shares, OOM ordering on a healthy node, or network weight — those belong to QoS classes and requests/limits. Priority also never became a placement mechanism: it says WHEN a pod schedules and WHAT may be evicted for it, never WHERE it lands. And per-namespace priority ceilings were left to admission policy rather than the resource itself — quota's `priority_class` scope selector (see KubernetesResourceQuota) is the budgeting tool.

## The Semantics in Detail

### The value drives everything

The integer is the entire mechanism. Higher value = earlier scheduling and eviction rights over lower values. Negative values are legal and idiomatic for "always-preemptable" tiers: a class at -100 yields to every unmarked pod (priority 0), which is precisely right for opportunistic batch work. The value is **immutable after creation** — changing it means replacing the object, and because PriorityClass names are cluster-unique, replacement must delete before it creates.

### Preemption mechanics

When a pod cannot schedule on any node, the scheduler looks for a node where evicting lower-priority pods would make it fit, evicts them (gracefully — termination grace periods apply, and PodDisruptionBudgets are respected on a best-effort basis, meaning they can be violated if no PDB-respecting choice exists), and binds the pending pod. Consequences worth internalizing:

- Preemption is triggered by scheduling failure, not by priority comparison alone — a cluster with headroom never preempts
- The victims are the lowest-priority pods whose removal suffices, not everything below the preemptor
- `preemptionPolicy: Never` pods still benefit from queue ordering AND can still be preempted BY higher classes — the policy governs what the class's own pods may do to others

### The global default

`globalDefault: true` assigns the class to every pod that names no class; such pods would otherwise get priority 0. The field has two sharp edges:

- **The API does not prevent multiple global defaults.** When several classes claim it, Kubernetes uses the SMALLEST such value — a silent, surprising resolution rule
- The default applies at admission: changing it never re-prioritizes existing pods, only pods created afterwards

### The description is load-bearing

`description` is surfaced by `kubectl describe priorityclass` and is the only in-cluster documentation of the ladder's intent. A class named `tier2` with no description forces every workload author to guess; the field deserves a real sentence.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

There is no `kubectl create priorityclass` generator; manual means `kubectl apply -f` of hand-written YAML.

**Verdict:** No shortcut exists; even ad-hoc use is YAML authoring.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: critical-services
value: 1000000
globalDefault: false
description: Revenue-path services; preempts lower tiers under pressure.
preemptionPolicy: PreemptLowerPriority
```

**Pros:**
- Declarative, version-controllable, trivially small

**Cons:**
- Nothing catches a second `globalDefault: true` before it changes every unmarked pod's priority
- The value ceiling and `system-` prefix reservation surface only at apply
- A value edit fails (immutable field) instead of replacing; the delete-then-recreate dance is manual

**Verdict:** The baseline; fine until the ladder needs governance.

### Level 2: Terraform

```hcl
resource "kubernetes_priority_class_v1" "critical_services" {
  metadata {
    name = "critical-services"
  }
  value             = 1000000
  description       = "Revenue-path services; preempts lower tiers under pressure."
  global_default    = false
  preemption_policy = "PreemptLowerPriority"
}
```

**Pros:**
- Full IaC lifecycle (plan, apply, destroy, import) with drift detection; value changes force replacement automatically

**Cons:**
- No validation of the value ceiling or reserved prefix before apply
- Fields left unset follow provider defaults, so two configurations of the "same" class can submit different objects

**Verdict:** Production-grade lifecycle, thin validation.

### Level 3: Pulumi

```go
priorityClass, err := schedulingv1.NewPriorityClass(ctx, "critical-services", &schedulingv1.PriorityClassArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name: pulumi.String("critical-services"),
    },
    Value:            pulumi.Int(1000000),
    Description:      pulumi.String("Revenue-path services; preempts lower tiers under pressure."),
    PreemptionPolicy: pulumi.String("PreemptLowerPriority"),
})
```

**Pros:**
- Full programming language, preview before apply

**Cons:**
- Types describe the wire shape, not the semantics; replacement-on-value-change needs explicit delete-before-replace handling because the name is cluster-unique

**Verdict:** Excellent IaC choice; same validation gap as Terraform.

### Other Methods

**Helm:** ladders templated in platform charts — common, and where duplicate global defaults most often ship, because each chart believes it owns "the" default.

**GitOps policy engines (Kyverno, OPA):** the usual home for "at most one global default" and "namespace X may only use classes Y, Z" rules — cluster governance layered on top of the resource, not a way to create it.

## Comparative Analysis

| Aspect | YAML | Terraform | Pulumi | Planton |
|--------|------|-----------|--------|---------|
| Validation | API server | Plan time (shape only) | Preview time (shape only) | Schema + CEL |
| Value ceiling / `system-` prefix caught early | No | No | No | Yes, rejected pre-apply |
| Deterministic preemptionPolicy across engines | N/A | Provider-dependent | Provider-dependent | Always explicit |
| Value change handled safely | Manual delete/create | Forced replacement | Needs explicit config | Delete-before-replace in both engines |
| Dual IaC | N/A | TF only | Pulumi only | Both |

## The Planton Approach

### Full surface, validated early

The spec models the entire `scheduling.k8s.io/v1` PriorityClass — name, value, global default, description, preemption policy — and moves the API server's rules to validation time:

- **The user-value ceiling**: values above 1,000,000,000 are rejected — the range above belongs to Kubernetes system classes
- **The reserved prefix**: names beginning with `system-` are rejected — those names belong to the built-in classes
- **Enum-checked preemption policy**: only `preempt_lower_priority` and `never` exist

### Deterministic wire objects

Both IaC modules resolve `preemption_policy` module-side — the explicit value, or the server default `PreemptLowerPriority` when omitted — and always send it explicitly. `global_default` is likewise always sent (the Kubernetes default `false` applied). Both engines submit byte-identical objects for the same manifest, and a spec that omits the optional fields deploys exactly what the API server would have defaulted, with no engine drift.

### Safe replacement for the immutable value

The priority value is immutable upstream, and PriorityClass names are cluster-unique — so an in-place update is impossible and a create-before-delete replacement would collide on the name. Both modules therefore force **delete-before-replace** on value change: the old class is removed, then the new one created under the same name. (Pods already running keep the priority integer they were admitted with; the class object is only read at admission.)

### The composition handle

The `priority_class_name` output is the handle workload pod specs reference — every Planton workload kind's shared pod spec carries a `priority_class_name` field. The `value` output supports downstream logic that needs the integer (e.g. documentation generators or ladder audits).

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates provider init, resource creation, and output export
- **`locals.go`**: Computes merged labels, annotations, and the resolved preemption policy (server default applied)
- **`priorityclass.go`**: Creates the `scheduling.k8s.io/v1` PriorityClass with `DeleteBeforeReplace` set
- **`outputs.go`**: Exports `priority_class_name` and `value`

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto` fields as Terraform variables
- **`locals.tf`**: Computes merged labels and the same resolved preemption policy
- **`main.tf`**: Creates the `kubernetes_priority_class_v1` resource (the provider forces replacement on value change)
- **`outputs.tf`**: Exports the same two outputs

### Resource Count

This is the leanest possible component — it creates exactly **one Kubernetes resource**: the PriorityClass itself. The complexity is in the cluster-wide semantics (global default, preemption), not in resource orchestration.

## Production Best Practices

### Ladder design

1. **Small and deliberate**: three to five rungs (e.g. critical / standard / batch) with real descriptions; every additional class multiplies the preemption interactions to reason about
2. **Leave numeric gaps**: space values by 1000 or more so future tiers slot in without renumbering — renumbering is a replace, and every workload referencing the class rides through it
3. **Negative values for opportunistic work**: a batch tier below 0 yields even to unmarked pods — the honest expression of "run only on spare capacity"

### Default discipline

1. **Exactly one global default per cluster**: audit `kubectl get priorityclass` before setting `global_default: true`; multiple defaults resolve to the SMALLEST value, silently
2. **Make the default modest**: the global default is what every unlabeled pod gets; it should sit low in the ladder, not high

### Preemption discipline

1. **`never` for high-priority batch**: queue-jumping without eviction rights is almost always what batch tiers actually want
2. **Pair preempting tiers with PodDisruptionBudgets on their victims**: preemption respects PDBs only best-effort, but a PDB still shapes which victims are chosen first
3. **Expect preemption only under pressure**: a cluster with headroom never preempts; test eviction behavior under induced scarcity, not on an idle cluster

## Conclusion

KubernetesPriorityClass is the smallest resource with the largest blast radius: one integer that reorders the scheduler's queue and licenses eviction cluster-wide. The component keeps the full upstream surface while moving the API's guardrails — the user-value ceiling, the reserved `system-` prefix — to validation time, submitting deterministic objects from both IaC engines, and handling the immutable value's replacement safely. Combined with the shared workload pod spec's `priority_class_name` field and ResourceQuota's priority-class budgeting, it makes the deliberate three-rung ladder — critical, standard, batch — a set of manifests rather than tribal knowledge.

## References

- [Pod Priority and Preemption Documentation](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/)
- [PriorityClass API Reference](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/priority-class-v1/)
- [Non-preempting Priority Classes](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/#non-preempting-priority-class)
- [Guaranteed Scheduling For Critical Add-On Pods](https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/)
- [Pod Disruption Budgets](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/#pod-disruption-budgets)
- [Pulumi Kubernetes PriorityClass](https://www.pulumi.com/registry/packages/kubernetes/api-docs/scheduling/v1/priorityclass/)
- [Terraform kubernetes_priority_class_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/priority_class_v1)
