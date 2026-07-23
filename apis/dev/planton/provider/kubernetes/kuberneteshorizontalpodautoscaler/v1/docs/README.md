# Kubernetes Horizontal Pod Autoscaler: Research Documentation

## Introduction

A fixed replica count encodes a guess about load, and the guess is always wrong twice a day: too few pods at peak, too many at trough. The HorizontalPodAutoscaler is the standard Kubernetes answer — a control loop in the controller manager that periodically compares observed metrics against declared targets and resizes the workload between a floor and a ceiling. It works on anything exposing the `scale` subresource: Deployments, StatefulSets, ReplicaSets, and custom resources that opt in.

The resource's mental model in three facts:

- **The HPA owns the replica count.** Once an autoscaler governs a workload, the workload's own `replicas` field is advisory — writing to it just gets overwritten on the next reconciliation. Exactly one controller may own a target's count.
- **Metrics OR toward scale-out.** Each configured metric independently proposes a replica count and the controller takes the highest. A workload scaled on CPU and queue depth scales up when either is hot, down only when both are cold.
- **The HPA reads metrics; it does not produce them.** Resource metrics need metrics-server; custom and external metrics need an adapter (prometheus-adapter, KEDA, cloud adapters). An HPA pointed at a metric nobody serves reports it unavailable and holds the current count.

Planton's **KubernetesHorizontalPodAutoscaler** component brings the full `autoscaling/v2` surface to the platform with schema-level validation that catches the classic mistakes before apply, namespace and target composition, and dual-IaC support.

## Evolution and Historical Context

### autoscaling/v1: CPU only (1.2)

The HPA GA'd in Kubernetes 1.2 (2016) as `autoscaling/v1` with exactly one signal: target CPU utilization as a percentage of requests. That single-metric shape is still what many teams mean by "an HPA," and it remains a valid API version served by conversion.

### The v2 beta era: metrics grow up (1.6–1.12)

`autoscaling/v2beta1` (1.6) introduced the metric source families — pods, object, and resource metrics — atop the new metrics API pipeline that replaced Heapster. `v2beta2` (1.12) restructured targets into the type/value triple (`Utilization`, `Value`, `AverageValue`) and added external metrics, opening scaling to signals from outside the cluster entirely: queue depths, cloud load balancer QPS, anything an external-metrics adapter serves.

### Behavior tuning (1.18) and v2 GA (1.23)

Configurable scaling behavior arrived in 1.18, answering the two chronic operational complaints — flapping and thundering herds — with per-direction stabilization windows and velocity policies (at most N pods or N percent per period, with Min/Max/Disabled selection between policies). `autoscaling/v2` graduated in Kubernetes 1.23 (2021) with the full surface: five metric source families, three target value forms, per-container resource metrics (`ContainerResource`, stable in 1.30), and behavior.

### What remains feature-gated upstream

Two capabilities exist upstream only behind feature gates, and are deliberately not modeled in this component:

- **Scale-to-zero** (`HPAScaleToZero` gate): allows `minReplicas: 0` when at least one object or external metric is configured. Ungated clusters reject it, and the event-driven ecosystem (KEDA) provides the same outcome portably by managing the HPA itself. The spec's floor is therefore ≥ 1.
- **Configurable tolerance** (`HPAConfigurableTolerance` gate): the controller acts only when the observed/target ratio deviates more than a tolerance, cluster-wide 10% by default; per-HPA tolerance is gated. The spec does not model it.

Both omissions are the conservative reading of the stable API: what every conformant cluster accepts, nothing more.

## The Semantics in Detail

### The control loop

Every reconciliation (15s by default), the controller computes for each metric: `desiredReplicas = ceil(currentReplicas × currentValue / targetValue)`, clamps the highest proposal to `[minReplicas, maxReplicas]`, filters it through behavior (stabilization, velocity policies), and writes the scale subresource. Utilization targets divide observed usage by the pods' **requests** — not limits, not node capacity — which is why utilization metrics on pods without requests are meaningless and rejected upstream.

### The five metric families and where each earns its place

- **`resource`** — cpu/memory averaged across the target's pods. The workhorse. CPU is the reliable signal: it rises and falls with load. Memory rarely falls after scale-out (heaps, caches), which makes it a poor scale-in driver — scale on CPU, alert on memory.
- **`container_resource`** — one named container's resource, isolating the app from sidecars. A pod-level average blends the app with its proxies and log shippers; a hot sidecar can mask an idle app and vice versa.
- **`pods`** — a custom per-pod metric (requests-per-second, active sessions), averaged. Needs a custom-metrics adapter.
- **`object`** — a metric describing ONE other object, e.g. an Ingress's request rate. Needs a custom-metrics adapter.
- **`external`** — a metric with no in-cluster object at all: queue depth, cloud LB QPS. Needs an external-metrics adapter. The natural fit for worker fleets — "one pod per 30 ready messages" is an `average_value` target on a queue-depth metric.

### The three target forms

`utilization` (percent of requests; resource families only), `average_value` (quantity divided across pods — the usual form for pods and external metrics), and the raw form that compares the metric directly against a quantity (object and external metrics). Upstream calls the raw form `Value`; this spec names the enum constant `raw_value` because generated code cannot use the bare word "value" as an enum constant — the IaC modules map it back to the API's `Value` string.

### Behavior: the flap damper and the velocity caps

Defaults: scale UP fast — no stabilization, the higher of doubling or +4 pods per 15 seconds; scale DOWN conservatively — a 300-second stabilization window (the controller scales down only to the highest recommendation of the past 5 minutes), all pods removable in one step once the window agrees. Tuning levers per direction: the stabilization window (0–3600s), velocity policies (pods or percent per period), and `select_policy` — `max_change` (default), `min_change`, or `disabled`, which turns a direction off entirely (freeze scale-down during an incident while leaving scale-up live).

### One controller per target

The HPA assumes sole ownership of the target's replica count. Pointing two autoscalers — or an autoscaler and a human-managed `replicas` field, or an HPA and KEDA (which itself manages an HPA) — at one workload produces oscillation, not consensus. This is also the boundary with Planton's built-in workload autoscaling: KubernetesDeployment's `availability.horizontal_pod_autoscaling` block manages an HPA for the workload's own replicas (CPU/memory targets, replicas as the floor). Use the built-in block for that case; use this standalone kind for operator-managed and non-Planton targets and for the advanced v2 surface (custom/object/external metrics, per-container metrics, behavior). Never both on one target.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

```bash
kubectl autoscale deployment checkout --min=2 --max=10 --cpu-percent=60
```

**Pros:**
- One-liner for the CPU case

**Cons:**
- CPU utilization only — none of the v2 surface; imperative and unrecorded

**Verdict:** Fine for experiments; nothing to build on.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: checkout-hpa
  namespace: backend
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: checkout
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 60
```

**Pros:**
- Declarative, version-controllable, full surface

**Cons:**
- The type/source mismatch trap: a metric whose `type` says one family while a different source block is filled deploys cleanly and is silently ignored by the controller
- Target-form mismatches (a `Utilization` type with an `averageValue`) and floor/ceiling inversions surface only at admission or, worse, as an autoscaler that never acts
- No plan/preview, no state management

**Verdict:** The baseline; the failure modes are the quiet kind that read as "the HPA isn't working."

### Level 2: Terraform

```hcl
resource "kubernetes_horizontal_pod_autoscaler_v2" "checkout" {
  metadata {
    name      = "checkout-hpa"
    namespace = "backend"
  }
  spec {
    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = "checkout"
    }
    min_replicas = 2
    max_replicas = 10
    metric {
      type = "Resource"
      resource {
        name = "cpu"
        target {
          type                = "Utilization"
          average_utilization = 60
        }
      }
    }
  }
}
```

**Pros:**
- Full IaC lifecycle (plan, apply, destroy, import) with drift detection; the v2 resource covers the whole surface

**Cons:**
- HCL blocks reproduce the type/source and type/value-form traps without additional guardrails

**Verdict:** Production-grade lifecycle, thin validation.

### Level 3: Pulumi

```go
hpa, err := autoscalingv2.NewHorizontalPodAutoscaler(ctx, "checkout-hpa", &autoscalingv2.HorizontalPodAutoscalerArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name:      pulumi.String("checkout-hpa"),
        Namespace: pulumi.String("backend"),
    },
    Spec: &autoscalingv2.HorizontalPodAutoscalerSpecArgs{
        ScaleTargetRef: &autoscalingv2.CrossVersionObjectReferenceArgs{
            ApiVersion: pulumi.String("apps/v1"),
            Kind:       pulumi.String("Deployment"),
            Name:       pulumi.String("checkout"),
        },
        MinReplicas: pulumi.Int(2),
        MaxReplicas: pulumi.Int(10),
    },
})
```

**Pros:**
- Full programming language, preview before apply

**Cons:**
- Types describe the wire shape, not the semantics; the same traps pass the compiler

**Verdict:** Excellent IaC choice; validation gap same as Terraform.

### Other Methods

**Helm:** HPAs templated per chart — ubiquitous (`autoscaling.enabled: true`), usually CPU-only, and chart templates rarely expose behavior or the advanced metric families.

**KEDA:** the event-driven layer — ScaledObjects with 60+ scalers that manage an HPA underneath and add scale-to-zero. The right tool when the requirement is event-driven scale-to-zero; for everything else it is an extra controller between you and the API. Never point KEDA and a hand-written HPA at the same target.

**VPA:** the perpendicular axis — resizes pods instead of adding them. Combining VPA (on CPU/memory) with an HPA scaling on the same resource creates a feedback loop; combine only across disjoint signals.

## Comparative Analysis

| Aspect | YAML | Terraform | Pulumi | Planton |
|--------|------|-----------|--------|---------|
| Validation | API server | Plan time (shape only) | Preview time (shape only) | Schema + CEL |
| Type/source mismatch caught | No (silently ignored) | No | No | Yes, rejected pre-apply |
| Target-form/type mismatch caught | Admission | No | No | Yes, rejected pre-apply |
| DaemonSet target rejected early | No | No | No | Yes |
| Floor/ceiling and quantity contracts checked early | Admission | No | No | Yes |
| Target as reference | No | Manual wiring | Manual wiring | First-class |
| Dual IaC | N/A | TF only | Pulumi only | Both |

## The Planton Approach

### Full surface, validated early

The spec models the entire `autoscaling/v2` HorizontalPodAutoscalerSpec — five metric source families, three target value forms, per-direction behavior — and moves the API server's rules plus the known footguns to validation time:

- **Type/source consistency**: each metric must set exactly the source field matching its declared type — the upstream failure mode (a metric the controller silently ignores) cannot be expressed
- **Type/value-form consistency**: each target must set exactly the value field matching its target type (`average_utilization` for utilization, `value` for the raw form, `average_value` for per-pod averages)
- **Floor and ceiling**: `min_replicas ≥ 1`, `max_replicas ≥ min_replicas`
- **Target contracts**: DaemonSets rejected as scale targets (one pod per node by definition); quantities checked against the Kubernetes quantity grammar; utilization percentages bounded 1–100
- **Behavior contracts**: stabilization windows bounded 0–3600s, policy periods 1–1800s, positive policy values, and a `disabled` direction cannot list policies (contradictory intent)

### Deterministic defaults

Both IaC modules resolve the spec defaults module-side — `apps/v1` `Deployment` target, `min_replicas` 1, behavior `select_policy` Max — and always send them explicitly, so both engines submit byte-identical objects for the same manifest. One deliberate asymmetry: when the spec lists **no metrics**, the metrics field is OMITTED entirely, because the API server then applies its own default (80% average CPU utilization) — sending an empty list would instead disable metric-driven scaling.

### What is deliberately unmodeled

- **Scale-to-zero** — feature-gated upstream (`HPAScaleToZero`) and requires at least one object/external metric; ungated clusters reject `minReplicas: 0`. The floor is ≥ 1; KEDA is the portable path to event-driven scale-to-zero.
- **Per-HPA tolerance** — feature-gated upstream (`HPAConfigurableTolerance`); the cluster-wide 10% tolerance applies.
- The `raw_value` enum constant maps to the upstream API's `Value` target type — the name differs only because generated code cannot use the bare word "value" as an enum constant.

### Namespace and target by value or reference

`spec.namespace` is a `StringValueOrRef` to a `KubernetesNamespace`; `scale_target.name` is a `StringValueOrRef` defaulting to a `KubernetesDeployment`'s exported `deployment_name` output — so an infra chart deploys the namespace, the workload, and its autoscaler in one run with ordering handled by the resource graph. The target's *existence* is not required at HPA creation: the controller reports it unresolved until the workload appears, which is the correct steady state for infra deployed ahead of the app.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates provider init, resource creation, and output export
- **`locals.go`**: Computes merged labels, annotations, the resolved namespace, the resolved scale target (defaults applied), the resolved replica floor, and the enum→API-string mappings
- **`horizontalpodautoscaler.go`**: Creates the `autoscaling/v2` HorizontalPodAutoscaler, converting metrics, targets, and behavior one-to-one
- **`outputs.go`**: Exports `horizontal_pod_autoscaler_name`, `namespace`, `scale_target`, `min_replicas`, `max_replicas`

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto` fields as Terraform variables
- **`locals.tf`**: Computes merged labels, resolved namespace, resolved target and floor, and the same enum→API-string maps
- **`main.tf`**: Creates the `kubernetes_horizontal_pod_autoscaler_v2` resource with dynamic blocks rendering at most one source per metric
- **`outputs.tf`**: Exports the same five outputs

### Resource Count

This is a lean component — it creates exactly **one Kubernetes resource**: the HorizontalPodAutoscaler itself. The complexity is in the spec validation and the semantic guardrails, not in resource orchestration.

## Production Best Practices

### Signal discipline

1. **Scale on CPU, alert on memory**: memory rarely falls after scale-out, which makes it a scale-in trap
2. **Declare requests before utilization targets**: utilization is measured against requests; unset or wrong requests make the signal meaningless
3. **Use `container_resource` when sidecars skew the average**: a hot proxy can mask an idle app container and vice versa
4. **Verify the metrics pipeline before trusting the autoscaler**: `kubectl get --raw /apis/metrics.k8s.io/v1beta1` (and the custom/external APIs for adapters) — an HPA on an unserved metric holds the count and reports unavailable

### Bounds discipline

1. **Treat `max_replicas` as a budget decision**: it answers "what is the most this workload may cost," not "a big number so scaling never clips"
2. **Set the floor to the availability minimum, not the average**: the floor is what you run at 4 a.m.; pair it with a PodDisruptionBudget so drains respect it

### Behavior discipline

1. **Lengthen scale-down stabilization for spiky traffic**: 600s plus a percent-per-minute policy bleeds capacity gradually instead of cliff-dropping after each spike
2. **Use `disabled` as the incident lever**: freezing scale-down while leaving scale-up live is a one-line change that prevents flap-induced churn during an investigation
3. **Leave scale-up fast unless proven otherwise**: the default (double or +4 per 15s) exists because under-provisioning during a surge is usually worse than over-provisioning after one

### Ownership discipline

1. **One controller per target**: never two HPAs, never an HPA plus KEDA, never the standalone kind plus the workload's built-in autoscaling block on the same target
2. **Stop setting `replicas` once the HPA owns it**: the field is advisory from then on; fighting the controller causes churn

## Conclusion

KubernetesHorizontalPodAutoscaler is a deliberately complete, deliberately lean component: the full `autoscaling/v2` surface — five metric families, three target forms, behavior tuning — with the resource's quiet failure modes (type/source mismatches the controller silently ignores, target forms that never match, DaemonSet targets, two controllers on one count) documented in the schema and guarded by validation before anything reaches a cluster. What upstream keeps behind feature gates stays unmodeled, keeping every accepted manifest deployable on every conformant cluster. Combined with target references and the built-in-block boundary, it makes metric-driven capacity a pattern that composes into infra charts rather than a hand-tuned artifact per service.

## References

- [Horizontal Pod Autoscaling Documentation](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
- [HorizontalPodAutoscaler Walkthrough](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/)
- [HorizontalPodAutoscaler API Reference (autoscaling/v2)](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/horizontal-pod-autoscaler-v2/)
- [Resource Metrics Pipeline (metrics-server)](https://kubernetes.io/docs/tasks/debug/debug-cluster/resource-metrics-pipeline/)
- [KEDA — event-driven autoscaling built on the external metrics API](https://keda.sh/)
- [Pulumi Kubernetes HorizontalPodAutoscaler](https://www.pulumi.com/registry/packages/kubernetes/api-docs/autoscaling/v2/horizontalpodautoscaler/)
- [Terraform kubernetes_horizontal_pod_autoscaler_v2](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/horizontal_pod_autoscaler_v2)
