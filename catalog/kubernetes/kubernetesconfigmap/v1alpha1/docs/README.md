# Kubernetes ConfigMap: Research Documentation

## Introduction

ConfigMaps are the standard Kubernetes mechanism for decoupling non-confidential configuration from container images. They carry file-like text values, property settings, and binary payloads that pods consume as environment variables, command-line arguments, or mounted files. Alongside Secrets — their confidential mirror — they are one of the two foundational configuration primitives in the platform: virtually every non-trivial workload consumes at least one.

The primitive itself is simple: two maps (`data` for UTF-8 text, `binaryData` for base64-encoded bytes), an `immutable` flag, and a 1MiB combined size cap. The engineering value in managing them well is not in the object but in the lifecycle: validation before apply, drift detection, versioned rollout patterns, and composition with the namespaces and workloads that surround them.

Planton's **KubernetesConfigMap** component brings that lifecycle to the primitive with full surface coverage — there is nothing an upstream ConfigMap can express that this spec cannot — plus schema-level validation and dual-IaC support.

## Evolution and Historical Context

### Origins (Kubernetes 1.2)

ConfigMaps were introduced in Kubernetes 1.2 (2016) to answer a simple question: where does configuration live if not baked into the image? Before ConfigMaps, teams rebuilt images per environment or injected configuration through bespoke init containers and wrapper scripts. ConfigMaps made "same image, different config" the platform default.

### binaryData (1.10+)

The original `data` map was UTF-8 only. Kubernetes 1.10 added `binaryData` for payloads that are not valid text — compiled files, keystores in non-PEM formats, images, serialized blobs. Values are base64-encoded on the wire, keys share the same character rules as `data` keys, and the API server rejects keys that appear in both maps.

### Immutable ConfigMaps (1.21+)

Kubernetes 1.19 introduced (and 1.21 graduated to stable) the `immutable` field. An immutable ConfigMap's data cannot be changed after creation — only delete-and-recreate. Two motivations:

- **Safety**: accidental edits to shared configuration can no longer propagate to running workloads
- **Scale**: the kubelet stops watching immutable objects, which materially reduces kube-apiserver load in clusters with tens of thousands of ConfigMap consumers

Immutability enabled the now-standard versioned-config pattern: `app-config-v1`, `app-config-v2`, with rollouts happening as workload updates pointing at new names.

### What ConfigMaps never became

ConfigMaps deliberately stayed dumb key-value storage. Proposals for typed schemas, size increases, and cross-namespace references were all rejected upstream. Anything larger than 1MiB or requiring structure belongs in a volume, an object store, or a CRD.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

```bash
# From literals
kubectl create configmap app-config \
  --from-literal=LOG_LEVEL=info \
  --from-literal=FEATURE_X=enabled

# From a file (key defaults to the file name)
kubectl create configmap app-config --from-file=application.properties
```

**Pros:**
- Immediate and intuitive
- `--from-file` handles file content and key naming automatically

**Cons:**
- Imperative — no drift detection, no reproducibility
- No record of what is deployed where
- No validation until the API server sees the object

**Verdict:** Fine for debugging. Not for production infrastructure.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: production
data:
  LOG_LEVEL: info
  application.properties: |
    server.port=8080
    cache.ttl.seconds=300
binaryData:
  logo.png: iVBORw0KGgoAAAANSUhEUg...
immutable: true
```

**Pros:**
- Declarative and version-controllable
- Full surface available (data, binaryData, immutable)

**Cons:**
- No plan/preview — `kubectl apply` is immediate
- Key rules, base64 validity, and data/binaryData overlap only checked at the API server
- No state management, no composition with surrounding resources

**Verdict:** The baseline. Everything above it adds lifecycle.

### Level 2: Terraform

```hcl
resource "kubernetes_config_map_v1" "app_config" {
  metadata {
    name      = "app-config"
    namespace = "production"
  }

  data = {
    LOG_LEVEL = "info"
  }

  binary_data = {
    "logo.png" = filebase64("logo.png")
  }

  immutable = true
}
```

**Pros:**
- Full IaC lifecycle (plan, apply, destroy, import)
- State management with drift detection
- Can reference other Terraform resources

**Cons:**
- Untyped maps — key character rules and overlap constraints surface only at apply
- Namespace is a plain string; no first-class reference to a managed namespace

**Verdict:** Production-grade lifecycle, thin validation.

### Level 3: Pulumi

```go
configMap, err := corev1.NewConfigMap(ctx, "app-config", &corev1.ConfigMapArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name:      pulumi.String("app-config"),
        Namespace: pulumi.String("production"),
    },
    Data: pulumi.StringMap{
        "LOG_LEVEL": pulumi.String("info"),
    },
    Immutable: pulumi.Bool(true),
})
```

**Pros:**
- Full programming language (type safety, IDE support, testing)
- Plan/preview before apply

**Cons:**
- Data maps are still untyped strings; constraint violations surface at apply
- Requires Pulumi runtime and SDK

**Verdict:** Excellent IaC choice; validation gap same as Terraform.

### Other Methods

**Helm:** `templates/configmap.yaml` with `{{ .Values }}` substitution — ubiquitous inside charts, but inherits raw-YAML validation gaps plus template complexity.

**Kustomize:** `configMapGenerator` shines at one thing — content-hashed names (`app-config-7c9fk4t2m8`) that force rollouts on config change. The same outcome is achievable with explicit versioned immutable names, which are also auditable.

## Comparative Analysis

| Aspect | kubectl | YAML | Terraform | Pulumi | Planton |
|--------|---------|------|-----------|--------|---------|
| Validation | At creation | API server | Plan time (basic) | Preview time (basic) | Schema + CEL |
| Key/overlap rules checked early | No | No | No | No | Yes |
| State Management | None | None | Full | Full | Full (via IaC) |
| Drift Detection | No | No | Yes | Yes | Yes |
| Namespace as reference | No | No | Manual wiring | Manual wiring | First-class |
| Dual IaC | N/A | N/A | TF only | Pulumi only | Both |
| Reproducible | No | Partially | Yes | Yes | Yes |

## The Planton Approach

### Full surface, validated early

The spec models the entire ConfigMap surface — `data`, `binary_data`, `immutable` — and moves the API server's own rules to validation time:

- Key character rules (`^[-._a-zA-Z0-9]+$`, max 253 chars) on both maps
- `binary_data` values must be valid base64 (checked with a pattern, not at apply)
- A CEL rule rejects any key present in both `data` and `binary_data` — the exact overlap rule the API server enforces, surfaced before deployment

An empty ConfigMap (no data at all) is deliberately valid, matching Kubernetes: empty ConfigMaps are used as name reservations and coordination markers.

### Namespace by value or reference

`spec.namespace` is a `StringValueOrRef`: a literal namespace name, or a reference to a `KubernetesNamespace` resource. The reference form lets an infra chart create the namespace and the ConfigMap in one run, with ordering handled by the resource graph rather than by apply-order discipline. When omitted, the ConfigMap lands in `default` — the same behavior as kubectl without a namespace flag.

### The Secret mirror

KubernetesConfigMap and KubernetesSecret are deliberate mirrors: same name/namespace handling, same immutability semantics, same key rules. The decision between them is exactly one question — is the data confidential? — and switching a value from one to the other is a mechanical move, not a remodeling.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates resource creation — Kubernetes provider and the ConfigMap resource
- **`locals.go`**: Computes merged labels, annotations, and the resolved namespace
- **`outputs.go`**: Exports `configmap_name` and `namespace`

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto` fields as Terraform variables
- **`locals.tf`**: Computes merged labels and resolved namespace
- **`main.tf`**: Creates the `kubernetes_config_map_v1` resource
- **`outputs.tf`**: Exports name and namespace

### Resource Count

This is a lean component — it creates exactly **one Kubernetes resource**: the ConfigMap itself. The complexity is in the spec validation and composition, not in resource orchestration.

## Production Best Practices

### Content discipline

1. **Never store confidential data**: ConfigMaps are readable by anyone with `get` on ConfigMaps in the namespace and are not covered by secrets encryption at rest. Use KubernetesSecret
2. **Respect the 1MiB cap**: The API server rejects oversized ConfigMaps. Large artifacts belong in volumes, object stores, or image layers
3. **Shape keys after their consumption**: environment-variable-style keys (`LOG_LEVEL`) for `envFrom`, file-style keys (`application.properties`) for volume mounts. Binary keys are only consumable as mounted files, never as environment variables

### Immutability and rollout

1. **Use immutable + versioned names in production**: `app-config-v42` names make config changes atomic, auditable, and revertible — roll out by pointing the workload at the new name
2. **Plan for delete-and-recreate**: Immutable ConfigMaps cannot be edited; the deployment pipeline must create new versions rather than mutate
3. **Know the propagation model for mutable ConfigMaps**: volume-mounted data refreshes eventually (kubelet sync period); environment variables never refresh without a pod restart. Immutable versioned config sidesteps both surprises

### Placement

1. **Same namespace as consumers**: ConfigMap references never cross namespaces. Duplicate per namespace or restructure
2. **Create namespace and config together**: use the namespace reference form in charts so ordering is guaranteed

## Conclusion

KubernetesConfigMap is a deliberately complete, deliberately lean component: the full upstream surface, validated at schema level, with namespace composition and dual-IaC lifecycle. Together with KubernetesSecret it closes the configuration story — one decision (confidential or not) selects the kind, and everything else about the two mirrors is the same.

## References

- [Kubernetes ConfigMaps Documentation](https://kubernetes.io/docs/concepts/configuration/configmap/)
- [Configure a Pod to Use a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/)
- [Immutable ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/#configmap-immutable)
- [ConfigMap API Reference](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/config-map-v1/)
- [Pulumi Kubernetes ConfigMap](https://www.pulumi.com/registry/packages/kubernetes/api-docs/core/v1/configmap/)
- [Terraform kubernetes_config_map_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/config_map_v1)
