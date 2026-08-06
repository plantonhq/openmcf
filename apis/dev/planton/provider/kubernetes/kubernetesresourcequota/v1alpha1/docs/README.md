# Kubernetes Resource Quota: Research Documentation

## Introduction

A Kubernetes namespace is an accounting boundary with, by default, no account: any pod can request any amount of CPU and memory, any team can create unbounded Services and PersistentVolumeClaims, and the first tenant to misbehave takes the cluster down for everyone. ResourceQuota is the standard multi-tenancy answer — a namespaced `core/v1` resource holding aggregate caps that the API server enforces at admission: any object whose creation or update would push the namespace's usage over a cap is rejected.

But ResourceQuota alone is only half the governance story, and the half-measure is actively dangerous:

- **Once a quota caps a compute resource (`requests.cpu`, `limits.memory`, ...), the API rejects pods that omit that request or limit.** A naive `kubectl run nginx` in a compute-governed namespace fails admission with a quota error.
- The upstream fix is a second object — **LimitRange** — that injects per-container defaults so unspecified workloads inherit sane values instead of being rejected, and additionally bounds what any single container, pod, or claim may ask for.
- The two objects answer one question — "how much may this namespace consume, and what does a workload get when it doesn't say?" — and deploying one without considering the other is the classic quota rollout failure.

Planton's **KubernetesResourceQuota** component treats the pair as one kind: the ResourceQuota always, and a companion LimitRange when `limit_defaults` is set — sharing the quota's name and namespace, one governance pair with one identity.

## Evolution and Historical Context

### Origins (Kubernetes 1.0–1.2)

ResourceQuota and LimitRange are among the oldest resources in Kubernetes — both shipped in `core/v1` at 1.0 (2015), designed for the original multi-tenant use case Kubernetes inherited from Borg: many teams sharing one cluster, each confined to a namespace with a budget. The admission-time enforcement model has never changed: quota is checked when objects are created, not while they run.

### Scopes (1.8+) and scope selectors (1.12+)

Quota scopes arrived to let one namespace carry several quotas tracking different slices of its workloads: `Terminating`/`NotTerminating` split by active deadline (bounded jobs vs long-running services), `BestEffort`/`NotBestEffort` by QoS class. Kubernetes 1.12 added `scopeSelector` with the `PriorityClass` scope, letting quota compose with the then-new pod priority feature — "this namespace may run at most 4 CPUs of critical-priority pods" — using `In`/`NotIn`/`Exists` operators. Later releases added the `CrossNamespacePodAffinity` scope (1.22, to contain a scheduling-abuse vector) and `VolumeAttributesClass` (1.29+, tracking claims by volume attributes class).

### Object-count quotas and the generic form (1.9+)

Early quotas counted a fixed list of object types (`pods`, `services`, `secrets`, ...). Kubernetes 1.9 generalized this with the `count/<resource>.<group>` syntax, making any countable API object — including CRDs — quotable without upstream changes.

### What ResourceQuota never became

Upstream deliberately kept quota namespaced and admission-time. Cluster-wide or hierarchical budgets (a parent team quota subdividing into child namespaces) were left to addons like the Hierarchical Namespace Controller and Kueue's cohort model. Runtime throttling was never in scope: a quota rejects creations; it never slows a running pod — that is the kubelet's and the scheduler's territory, driven by the requests/limits that LimitRange defaults inject.

## The Semantics in Detail

### Admission, not runtime

Quota usage is tracked by the quota controller and enforced by the ResourceQuota admission plugin. Consequences:

- Objects created BEFORE the quota existed are counted against usage but never evicted; a namespace can be born over-quota and nothing breaks until the next creation is rejected
- Quota errors surface at creation time, on whatever created the pod — a Deployment's ReplicaSet silently stops scaling, with the rejection visible only in its events

### The compute-caps-require-declarations rule

When a quota caps `requests.cpu`, `requests.memory`, `limits.cpu`, or `limits.memory`, every pod in the namespace must declare that request or limit — omission is rejection. This is the single most operationally important fact about ResourceQuota, and the reason LimitRange defaults exist:

- **`defaultRequest`** — the request injected into containers that omit their own; the field that keeps a requests-capping quota from rejecting naive pods
- **`default`** (limit) — the limit injected into containers that omit their own; when set without `defaultRequest`, Kubernetes copies the limit into the request, producing Guaranteed-QoS containers
- Defaults exist only where requests/limits live — on containers. Pod and PersistentVolumeClaim limit items may carry only `min`/`max`/`maxLimitRequestRatio`; the API rejects defaults on them

### Scope algebra

An object must match ALL of a quota's scopes to be tracked (scopes AND). `scopeSelector` requirements AND with `scopes`. The API rejects contradictory quotas (`BestEffort` + `NotBestEffort`, `Terminating` + `NotTerminating`) and restricts what scoped quotas may cap — a `BestEffort` quota may cap only `pods`, because BestEffort pods by definition have no requests or limits to meter. Only `PriorityClass` and `VolumeAttributesClass` carry values (`In`/`NotIn`); the pod-behavior scopes accept only `Exists`.

### Quantities and the ratio bound

All caps and bounds are Kubernetes quantities (`"10"`, `500m`, `20Gi`). `maxLimitRequestRatio` bounds burstiness: with `cpu: "4"`, a container's CPU limit may be at most 4x its request, and both must be set and non-zero. It is the lever against namespaces full of `request: 1m, limit: 4` containers that overcommit nodes.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

`kubectl create quota my-quota --hard=pods=50,requests.cpu=10` exists and covers simple caps; scopes, scope selectors, and the companion LimitRange do not fit the generator and require YAML anyway.

**Verdict:** Fine for experiments; the pair cannot be expressed.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: team-alpha-quota
  namespace: team-alpha
spec:
  hard:
    requests.cpu: "10"
    requests.memory: 20Gi
    limits.cpu: "20"
    limits.memory: 40Gi
---
apiVersion: v1
kind: LimitRange
metadata:
  name: team-alpha-quota
  namespace: team-alpha
spec:
  limits:
    - type: Container
      defaultRequest:
        cpu: 100m
        memory: 128Mi
      default:
        cpu: 500m
        memory: 512Mi
```

**Pros:**
- Declarative, version-controllable, full surface

**Cons:**
- Nothing ties the two objects together; the quota half routinely ships without the LimitRange half, and the namespace starts rejecting naive pods
- Scope contradictions, scoped-resource restrictions, and the defaults-on-container-only rule surface only at apply
- No plan/preview, no state management

**Verdict:** The baseline; the pairing discipline is entirely on the author.

### Level 2: Terraform

```hcl
resource "kubernetes_resource_quota_v1" "team_alpha" {
  metadata {
    name      = "team-alpha-quota"
    namespace = "team-alpha"
  }
  spec {
    hard = {
      "requests.cpu"    = "10"
      "requests.memory" = "20Gi"
    }
  }
}

resource "kubernetes_limit_range_v1" "team_alpha" {
  metadata {
    name      = "team-alpha-quota"
    namespace = "team-alpha"
  }
  spec {
    limit {
      type = "Container"
      default_request = {
        cpu    = "100m"
        memory = "128Mi"
      }
    }
  }
}
```

**Pros:**
- Full IaC lifecycle (plan, apply, destroy, import) with drift detection

**Cons:**
- Still two independent resources; the pairing is a convention, not a contract
- The admission rules (scope conflicts, BestEffort restrictions, container-only defaults) surface only at apply

**Verdict:** Production-grade lifecycle, thin validation, no pairing.

### Level 3: Pulumi

```go
quota, err := corev1.NewResourceQuota(ctx, "team-alpha-quota", &corev1.ResourceQuotaArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name:      pulumi.String("team-alpha-quota"),
        Namespace: pulumi.String("team-alpha"),
    },
    Spec: &corev1.ResourceQuotaSpecArgs{
        Hard: pulumi.StringMap{
            "requests.cpu":    pulumi.String("10"),
            "requests.memory": pulumi.String("20Gi"),
        },
    },
})
```

**Pros:**
- Full programming language, preview before apply; the pairing CAN be encapsulated in a component

**Cons:**
- Types describe the wire shape, not the admission semantics; every team re-invents the pairing convention

**Verdict:** Excellent IaC choice; the governance-pair abstraction is left as an exercise.

### Other Methods

**Helm:** quotas and limit ranges templated into namespace-onboarding charts — common in platform teams, with the same untyped pairing convention.

**Hierarchical/queue-based quota (HNC, Kueue):** the right tools when budgets must span namespaces or gate batch admission by queue; strictly beyond what core ResourceQuota models.

## Comparative Analysis

| Aspect | YAML | Terraform | Pulumi | Planton |
|--------|------|-----------|--------|---------|
| Quota + LimitRange as one unit | No (two objects, convention) | No | Component DIY | Yes, one kind |
| Scope conflicts caught early | No (apply time) | No | No | Yes, rejected pre-apply |
| BestEffort-caps-only-pods rule | Apply time | Apply time | Apply time | Schema + CEL |
| Defaults-on-container-only rule | Apply time | Apply time | Apply time | Schema + CEL |
| Scope-selector operator contract | Apply time | Apply time | Apply time | Schema + CEL |
| Namespace as reference | No | Manual wiring | Manual wiring | First-class |
| Dual IaC | N/A | TF only | Pulumi only | Both |

## The Planton Approach

### One kind, one governance pair

The spec covers the complete `core/v1` ResourceQuotaSpec and LimitRangeSpec surfaces in one resource. The ResourceQuota is always created; the companion LimitRange is created exactly when `limit_defaults` is non-empty, sharing the quota's name and namespace. The `limit_range_name` output states which case occurred — the LimitRange's name, or empty when none exists.

### The admission rules, moved to validation time

The schema mirrors the API server's own rejections so they surface before anything reaches a cluster:

- **Conflicting scope pairs**: `best_effort` + `not_best_effort` and `terminating` + `not_terminating` are rejected outright
- **The BestEffort restriction**: a `best_effort`-scoped quota may cap only `pods` (and `count/` entries) — BestEffort pods have no requests or limits to meter
- **The operator/values contract**: `In`/`NotIn` require values, `Exists`/`DoesNotExist` forbid them; the pod-behavior scopes accept only `Exists` — only `priority_class` and `volume_attributes_class` support `In`/`NotIn`
- **Container-only defaults**: `default_limit`/`default_request` on a `pod` or `persistent_volume_claim` item is rejected, as is an item that constrains nothing at all

### Deterministic API strings

Both IaC modules map the proto enum names to the Kubernetes API strings identically (`best_effort` → `BestEffort`, `persistent_volume_claim` → `PersistentVolumeClaim`) and both stamp the same resource-identity labels, so the two engines submit identical objects for the same manifest.

### Namespace by value or reference

`spec.namespace` is a `StringValueOrRef`: a literal namespace name, or a reference to a `KubernetesNamespace` resource, letting an infra chart create a namespace and its governance in one run with ordering handled by the resource graph. When omitted, the quota governs `default`.

### Division of labor with KubernetesNamespace

The `KubernetesNamespace` kind's resource profiles manage a quota and limit range internally as T-shirt sizes — the right tool for simple governance on namespaces Planton creates. KubernetesResourceQuota is the full-fidelity instrument: scope-filtered quotas, object-count caps, ratio bounds, priority-class budgets, and governance for namespaces Planton does not own. The two are alternatives on the same namespace, not layers.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates provider init, quota creation, conditional LimitRange creation, and output export
- **`locals.go`**: Computes merged labels, annotations, the resolved namespace, the LimitRange name (empty when no `limit_defaults`), and the enum→API-string mappings
- **`resourcequota.go`**: Creates the `core/v1` ResourceQuota (hard, scopes, scope selector) and, when configured, the companion `core/v1` LimitRange with all item types
- **`outputs.go`**: Exports `resource_quota_name`, `namespace`, and `limit_range_name`

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto` fields as Terraform variables
- **`locals.tf`**: Computes merged labels, resolved namespace, the same enum→API-string maps, and the create-LimitRange condition
- **`main.tf`**: Creates `kubernetes_resource_quota_v1` and a count-gated `kubernetes_limit_range_v1`
- **`outputs.tf`**: Exports the same three outputs

### Resource Count

This component creates **one or two Kubernetes resources**: the ResourceQuota always, and the companion LimitRange exactly when `limit_defaults` is set. The complexity is in the admission-rule validation and the pairing semantics, not in resource orchestration.

## Production Best Practices

### Rollout discipline

1. **Never cap compute without defaults**: a quota on `requests.*`/`limits.*` without `limit_defaults` (or defaults established elsewhere) makes the API reject every pod that omits them — the number-one quota rollout failure
2. **Start with object counts on live namespaces**: `pods`, `services.loadbalancers`, `persistentvolumeclaims` caps constrain nothing pods must declare; compute caps deserve a measured rollout with defaults in place first
3. **Watch `kubectl describe resourcequota` after rollout**: it reports used vs hard per resource; pre-existing usage counts immediately, and a namespace born over-quota rejects the very next creation

### Sizing discipline

1. **Size from observed usage, not headcount**: read current namespace consumption before choosing caps; a too-tight quota surfaces as mysterious scaling stalls in ReplicaSet events
2. **Keep defaults modest**: `default_request` is what every naive container is billed against the quota; generous defaults silently exhaust the budget with idle reservations
3. **Bound burstiness where overcommit hurts**: `max_limit_request_ratio` caps the request-vs-limit gap that lets a namespace look small at scheduling time and huge at runtime

### Targeting discipline

1. **Use `best_effort`-scoped pod caps as the naive-pod guard**: they contain unbounded pods without touching declared workloads
2. **Budget priority tiers separately**: a `priority_class` scope selector gives critical and batch tiers independent caps in the same namespace
3. **Prefer KubernetesNamespace resource profiles for the simple case**: reach for this kind when you need scopes, ratios, object counts, or namespaces Planton does not own

## Conclusion

KubernetesResourceQuota is a deliberately paired component: the full upstream ResourceQuota surface plus the LimitRange that makes compute quotas livable, managed as one kind because they answer one governance question. The admission rules that upstream enforces only at apply time — scope conflicts, the BestEffort restriction, container-only defaults, the operator/values contract — are moved to validation time, and the quota-without-defaults footgun that breaks naive pod creation is documented in the schema and solved by the pairing itself. Combined with namespace references, it makes per-team consumption governance a manifest that can be stamped out per namespace rather than a two-object convention every platform team re-invents.

## References

- [Kubernetes Resource Quotas Documentation](https://kubernetes.io/docs/concepts/policy/resource-quotas/)
- [Kubernetes Limit Ranges Documentation](https://kubernetes.io/docs/concepts/policy/limit-range/)
- [ResourceQuota API Reference](https://kubernetes.io/docs/reference/kubernetes-api/policy-resources/resource-quota-v1/)
- [LimitRange API Reference](https://kubernetes.io/docs/reference/kubernetes-api/policy-resources/limit-range-v1/)
- [Quality of Service for Pods](https://kubernetes.io/docs/concepts/workloads/pods/pod-qos/)
- [Pulumi Kubernetes ResourceQuota](https://www.pulumi.com/registry/packages/kubernetes/api-docs/core/v1/resourcequota/)
- [Terraform kubernetes_resource_quota_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/resource_quota_v1)
- [Terraform kubernetes_limit_range_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/limit_range_v1)
