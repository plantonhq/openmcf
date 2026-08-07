# KubernetesResourceQuota

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**KubernetesResourceQuotaSpec** governs resource consumption in one
namespace: aggregate caps on what the namespace may use in total (the
ResourceQuota), and per-object defaults and bounds applied to individual
pods, containers, and claims (an optional companion LimitRange). They are
two Kubernetes objects but one governance story — "how much may this
namespace consume, and what does a workload get when it doesn't say?" —
which is why this kind manages both.

The two arms interact: once a quota caps `requests.cpu` or `limits.memory`,
the API REJECTS pods that omit those requests/limits — so a compute quota
without `limit_defaults` (or defaults set elsewhere) breaks naive pod
creation in the namespace. Setting both together is the safe pattern.

For simple T-shirt-size governance on namespaces Planton creates, prefer
KubernetesNamespace's resource profiles (they manage a quota and limit
range internally). This kind is the full-fidelity instrument: scope-filtered
quotas, object-count caps, ratio bounds, and governance for namespaces
Planton does not own.

The spec covers the complete core/v1 ResourceQuotaSpec and LimitRangeSpec
surfaces.

## Example

```yaml
# Full-surface test manifest for the offline plan proofs. Exercises compute
# and object-count caps, a scope, a priority-class scope selector, and all
# three limit-defaults item types.
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesResourceQuota
metadata:
  name: test-resource-quota
spec:
  namespace:
    value: default
  name: test-resource-quota
  labels:
    team: platform-engineering
  hard:
    requests.cpu: "10"
    requests.memory: 20Gi
    limits.cpu: "20"
    limits.memory: 40Gi
    pods: "50"
    persistentvolumeclaims: "20"
    requests.storage: 500Gi
  scopes:
    - not_terminating
  scope_selector:
    - scope_name: priority_class
      operator: NotIn
      values:
        - system-critical
  limit_defaults:
    - type: container
      default_request:
        cpu: 100m
        memory: 128Mi
      default_limit:
        cpu: 500m
        memory: 512Mi
      max:
        cpu: "2"
        memory: 4Gi
      max_limit_request_ratio:
        cpu: "4"
    - type: pod
      max:
        cpu: "4"
        memory: 8Gi
    - type: persistent_volume_claim
      min:
        storage: 1Gi
      max:
        storage: 100Gi
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` |  |  | KubernetesNamespace (`spec.name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.labels` | `map<string, string>` |  |  |  |
| `spec.annotations` | `map<string, string>` |  |  |  |
| `spec.hard` | `map<string, string>` | yes |  |  |
| `spec.scopes` | `[]enum` |  |  |  |
| `spec.scopeSelector` | `[]KubernetesResourceQuotaScopeSelectorRequirement` |  |  |  |
| `spec.scopeSelector[].scopeName` | `enum` |  |  |  |
| `spec.scopeSelector[].operator` | `string` |  |  |  |
| `spec.scopeSelector[].values` | `[]string` |  |  |  |
| `spec.limitDefaults` | `[]KubernetesResourceQuotaLimitDefaults` |  |  |  |
| `spec.limitDefaults[].type` | `enum` |  |  |  |
| `spec.limitDefaults[].max` | `map<string, string>` |  |  |  |
| `spec.limitDefaults[].min` | `map<string, string>` |  |  |  |
| `spec.limitDefaults[].defaultLimit` | `map<string, string>` |  |  |  |
| `spec.limitDefaults[].defaultRequest` | `map<string, string>` |  |  |  |
| `spec.limitDefaults[].maxLimitRequestRatio` | `map<string, string>` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom`

The namespace to govern. Accepts a literal namespace name or a reference
to a KubernetesNamespace resource. When omitted, governs the cluster's
`default` namespace.

- references: KubernetesNamespace (`spec.name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the ResourceQuota object (its `metadata.name` in the
cluster). The companion LimitRange, when `limit_defaults` is set, shares
this name. Must be a valid DNS subdomain: lowercase alphanumeric
characters, hyphens, and dots, at most 253 characters.

- rule: Name must be a valid DNS subdomain (lowercase alphanumeric, hyphens, and dots, no leading/trailing dots or hyphens)
- rule: {"string":{"minLen":"1","maxLen":"253"}}

### spec.labels

`map<string, string>`

Additional labels to apply to the created objects.
Merged with the standard Planton governance labels.

### spec.annotations

`map<string, string>`

Annotations to apply to the created objects.

### spec.hard

`map<string, string>` · required

The aggregate caps: total consumption the namespace may not exceed, as
resource name → quantity. The vocabulary is upstream's:
- Compute: "requests.cpu", "requests.memory", "limits.cpu",
  "limits.memory" (also plain "cpu"/"memory", aliases for requests).
- Storage: "requests.storage" (total claimed),
  "persistentvolumeclaims" (count), and per-class variants
  ("<class>.storageclass.storage.k8s.io/requests.storage").
- Object counts: "pods", "services", "services.loadbalancers",
  "services.nodeports", "secrets", "configmaps", "resourcequotas", and
  the generic "count/<resource>.<group>" form for any countable object.
CAUTION: capping requests/limits makes the API reject pods that omit
them — pair with `limit_defaults` so unspecified workloads inherit sane
values instead of being rejected.

- rule: {"map":{"minPairs":"1","values":{"string":{"minLen":"1"}}}}

### spec.scopes

`[]enum`

Coarse filters on which objects this quota tracks — a pod/claim must
match ALL listed scopes to be counted. When empty, the quota tracks
everything its `hard` entries name. Scoped quotas restrict which
resources they may cap (a best_effort quota may cap only "pods";
pod-behavior scopes cap only pod resources) — the API enforces this.

- rule: {"repeated":{"items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `kubernetes_resource_quota_scope_unspecified` -- Unspecified.
- `terminating` -- Pods with an active deadline (spec.activeDeadlineSeconds set) — i.e. bounded-duration pods.
- `not_terminating` -- Pods WITHOUT an active deadline — long-running workloads.
- `best_effort` -- Pods with BestEffort quality of service (no requests or limits at all). A best_effort-scoped quota can only meter pod counts — "pods" or generic "count/..." object-count entries; the API rejects standard compute resources under this scope (BestEffort pods have nothing to meter).
- `not_best_effort` -- Pods with Burstable or Guaranteed QoS (i.e. anything but BestEffort).
- `priority_class` -- Pods that name ANY priority class. To match specific classes, use scope_selector with the priority_class scope and In/NotIn values.
- `cross_namespace_pod_affinity` -- Pods using cross-namespace pod (anti)affinity terms.
- `volume_attributes_class` -- PersistentVolumeClaims that name any volume attributes class.

### spec.scopeSelector

`[]KubernetesResourceQuotaScopeSelectorRequirement`

Fine-grained scope filters, for scopes that carry values — most usefully
priority_class with In/NotIn (e.g. "quota only pods of priority class
critical"). ANDs with `scopes`.

- rule: values must be non-empty for In/NotIn and empty for Exists/DoesNotExist
- rule: the terminating, not_terminating, best_effort, not_best_effort, and cross_namespace_pod_affinity scopes accept only the Exists operator — only priority_class and volume_attributes_class support In/NotIn

### spec.scopeSelector[].scopeName

`enum`

The scope the filter applies to. Only priority_class and
volume_attributes_class support In/NotIn with values; the pod-behavior
scopes (terminating, best_effort, ...) accept only Exists.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `kubernetes_resource_quota_scope_unspecified` -- Unspecified.
- `terminating` -- Pods with an active deadline (spec.activeDeadlineSeconds set) — i.e. bounded-duration pods.
- `not_terminating` -- Pods WITHOUT an active deadline — long-running workloads.
- `best_effort` -- Pods with BestEffort quality of service (no requests or limits at all). A best_effort-scoped quota can only meter pod counts — "pods" or generic "count/..." object-count entries; the API rejects standard compute resources under this scope (BestEffort pods have nothing to meter).
- `not_best_effort` -- Pods with Burstable or Guaranteed QoS (i.e. anything but BestEffort).
- `priority_class` -- Pods that name ANY priority class. To match specific classes, use scope_selector with the priority_class scope and In/NotIn values.
- `cross_namespace_pod_affinity` -- Pods using cross-namespace pod (anti)affinity terms.
- `volume_attributes_class` -- PersistentVolumeClaims that name any volume attributes class.

### spec.scopeSelector[].operator

`string`

The relationship between the scope and the values: "In", "NotIn",
"Exists", or "DoesNotExist".

- rule: {"string":{"in":["In","NotIn","Exists","DoesNotExist"]}}

### spec.scopeSelector[].values

`[]string`

The values compared against the scope (e.g. priority class names).
Required (non-empty) for In/NotIn; must be empty for Exists/DoesNotExist.

### spec.limitDefaults

`[]KubernetesResourceQuotaLimitDefaults`

Per-object defaults and bounds, managed as a companion LimitRange object
(named after this resource, in the same namespace). This is what makes a
compute quota livable: workloads that omit requests/limits inherit the
defaults instead of being rejected by the quota. Omit entirely to manage
only the quota.

- rule: default_limit and default_request apply only to the container type — pod and persistent_volume_claim items may set only min, max, and max_limit_request_ratio
- rule: a limit_defaults item must set at least one of max, min, default_limit, default_request, or max_limit_request_ratio

### spec.limitDefaults[].type

`enum`

The object type this item governs.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `kubernetes_resource_quota_limit_type_unspecified` -- Unspecified.
- `container` -- Per-container defaults and bounds — the type that carries default_request/default_limit (containers are where requests/limits actually live).
- `pod` -- Per-pod bounds across all its containers summed. Pods cannot carry defaults — only min/max/ratio.
- `persistent_volume_claim` -- Per-PersistentVolumeClaim storage bounds (min/max "storage").

### spec.limitDefaults[].max

`map<string, string>`

Maximum any single object of this type may claim, per resource
(e.g. {"cpu": "2", "memory": "4Gi"}).

- rule: {"map":{"values":{"string":{"minLen":"1"}}}}

### spec.limitDefaults[].min

`map<string, string>`

Minimum any single object of this type must claim, per resource.

- rule: {"map":{"values":{"string":{"minLen":"1"}}}}

### spec.limitDefaults[].defaultLimit

`map<string, string>`

LIMIT applied to containers that omit their own, per resource. Container
type only. When only default_limit is set, the request defaults to the
same value (Kubernetes copies it), producing Guaranteed-QoS containers.

- rule: {"map":{"values":{"string":{"minLen":"1"}}}}

### spec.limitDefaults[].defaultRequest

`map<string, string>`

REQUEST applied to containers that omit their own, per resource.
Container type only. This is the field that keeps a requests-capping
quota from rejecting naive pods.

- rule: {"map":{"values":{"string":{"minLen":"1"}}}}

### spec.limitDefaults[].maxLimitRequestRatio

`map<string, string>`

Maximum limit-to-request ratio per resource (e.g. {"cpu": "4"} allows at
most 4x burst). Both request and limit must then be set and non-zero on
the object.

- rule: {"map":{"values":{"string":{"minLen":"1"}}}}

## Validation Rules

- `scopes.conflicting_best_effort`: best_effort and not_best_effort are conflicting scopes — a quota cannot require both
- `scopes.conflicting_terminating`: terminating and not_terminating are conflicting scopes — a quota cannot require both
- `scopes.best_effort_resources`: a best_effort-scoped quota can only meter pod counts ("pods" or "count/..." entries) — BestEffort pods have no requests or limits to meter

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesResourceQuota, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.resource_quota_name` | `string` | The name of the ResourceQuota object as created in the cluster. |
| `status.outputs.namespace` | `string` | The namespace the quota governs. |
| `status.outputs.limit_range_name` | `string` | The name of the companion LimitRange object. Empty when the spec set no limit_defaults (no LimitRange is created). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |

## See Also

- [Overview](../README.md)
