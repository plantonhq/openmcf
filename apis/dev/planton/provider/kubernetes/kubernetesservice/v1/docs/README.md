# Kubernetes Service: Research Documentation

## Introduction

The Service is Kubernetes' answer to a problem as old as the platform itself: pods are ephemeral, but the things that talk to them need a stable address. A Service gives a set of pods one durable virtual IP and DNS name; kube-proxy (or its dataplane successors) keeps traffic flowing to whichever pods currently match the selector. Everything else in the networking stack composes on it — Ingress backends are Services, NetworkPolicies reason about the pods Services select, meshes discover endpoints through Services, and every in-cluster client connects through a Service DNS name.

Planton workload kinds (KubernetesDeployment, KubernetesStatefulSet) already create a Service for their own pods, so the ordinary case needs no separate resource. The standalone **KubernetesService** component exists for everything the workload-owned Service does not cover: exposing pods managed outside Planton, LoadBalancer/NodePort exposure with cloud-provider annotations, ExternalName aliases to endpoints outside the cluster, headless services for custom discovery, dual-stack addressing, and selectorless services fronting manually-managed endpoints.

The spec models the complete core/v1 ServiceSpec surface. The single deliberate omission is the deprecated `loadBalancerIP` field, discussed below.

## Evolution and Historical Context

### The original core (Kubernetes 1.0)

Services shipped in Kubernetes 1.0 with the shape still recognizable today: a selector, a list of ports, a virtual ClusterIP, and the NodePort/LoadBalancer escalation ladder. ExternalName arrived in 1.5 as the odd one out — a pure DNS CNAME with no proxying at all. Headless services (`clusterIP: None`) also date to the earliest days, born from the needs of stateful systems that must address each replica individually.

### Traffic policies (1.7+)

`externalTrafficPolicy: Local` answered a persistent complaint: NodePort/LoadBalancer traffic proxied through a random node masquerades the client source IP and adds a hop. Local-policy routing keeps traffic on the receiving node, preserving the source IP — at the price that nodes without endpoints must be health-checked out by the external load balancer, which is what `healthCheckNodePort` exists for. `internalTrafficPolicy` (1.22+) brought the same node-local option to ClusterIP traffic for agent patterns like node-level DNS caches.

### Dual-stack (stable in 1.23)

IPv4/IPv6 dual-stack introduced `ipFamilies` (which families, in what order) and `ipFamilyPolicy` (SingleStack / PreferDualStack / RequireDualStack). PreferDualStack is the portable opt-in: two families on dual-stack clusters, graceful single-family fallback elsewhere.

### The deprecation of `loadBalancerIP`

The original `spec.loadBalancerIP` promised a pinned load-balancer address but never specified semantics — each cloud interpreted it differently and several ignored it. Upstream deprecated it in 1.24 in favor of provider-specific annotations (`networking.gke.io/load-balancer-ip-addresses`, AWS EIP allocations, Azure PIP names). This spec follows upstream: no `load_balancer_ip` field; pinned addresses are expressed in `annotations`, which is where every cloud actually reads them. In the same era, `loadBalancerClass` (1.24) made multiple LB implementations per cluster first-class, and `allocateLoadBalancerNodePorts` (1.24) let VIP-mode implementations skip the NodePort hop entirely.

### Traffic distribution (stable in 1.33)

`trafficDistribution` superseded the older topology-keys experiments with two modest, honorable hints: `PreferSameZone` (cut cross-zone data transfer cost and latency) and `PreferSameNode` (DaemonSet-style node-local routing). It is a preference, not a guarantee — implementations honor it when safe.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

```bash
kubectl expose deployment my-app --port=80 --target-port=8080
kubectl create service loadbalancer my-app --tcp=443:8443
```

**Pros:**
- Immediate; `kubectl expose` reads the selector off the workload

**Cons:**
- Imperative — no drift detection, no reproducibility
- The interesting surface (traffic policies, LB annotations, dual-stack) is unreachable without follow-up edits

**Verdict:** Fine for a demo. Not for production infrastructure.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: v1
kind: Service
metadata:
  name: public-web
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-type: "nlb"
spec:
  type: LoadBalancer
  selector:
    app: web
  ports:
    - name: https
      port: 443
      targetPort: 8443
  externalTrafficPolicy: Local
```

**Pros:**
- Declarative, version-controllable, full surface available

**Cons:**
- The Service API is a minefield of cross-field rules enforced only at the API server: headless vs. NodePort, ExternalName vs. selector, LoadBalancer-only knobs, affinity timeouts, dual-stack consistency — all fail at apply time
- No state management, no composition with surrounding resources

**Verdict:** The baseline. Everything above it adds lifecycle.

### Level 2: Terraform

```hcl
resource "kubernetes_service_v1" "public_web" {
  metadata {
    name = "public-web"
    annotations = {
      "service.beta.kubernetes.io/aws-load-balancer-type" = "nlb"
    }
  }
  spec {
    type = "LoadBalancer"
    selector = { app = "web" }
    port {
      name        = "https"
      port        = 443
      target_port = 8443
    }
    external_traffic_policy = "Local"
  }
}
```

**Pros:**
- Full IaC lifecycle (plan, apply, destroy, import), drift detection
- The provider waits for the load balancer to provision, so the LB address is a usable output

**Cons:**
- Untyped strings — every cross-field rule still surfaces at apply
- The provider lags the API: `trafficDistribution` is not exposed at all (v3.2.x)

**Verdict:** Production-grade lifecycle, thin validation, incomplete surface.

### Level 3: Pulumi

```go
service, err := corev1.NewService(ctx, "public-web", &corev1.ServiceArgs{
    Spec: &corev1.ServiceSpecArgs{
        Type:     pulumi.String("LoadBalancer"),
        Selector: pulumi.StringMap{"app": pulumi.String("web")},
        Ports: corev1.ServicePortArray{&corev1.ServicePortArgs{
            Port: pulumi.Int(443), TargetPort: pulumi.Int(8443),
        }},
        TrafficDistribution: pulumi.String("PreferSameZone"),
    },
})
```

**Pros:**
- Full programming language, preview before apply
- Tracks the upstream API closely — the full surface including `trafficDistribution` is available
- Await logic blocks on LoadBalancer ingress, so the LB address resolves reliably

**Cons:**
- Constraint violations still surface at apply
- Requires Pulumi runtime and SDK

**Verdict:** Excellent IaC choice with the most complete surface.

### Other Methods

**Helm:** Services templated inside charts are ubiquitous, but the chart owns the shape — a standalone Service for someone else's pods is not what charts are for.

**Operators/meshes:** Meshes (Istio, Linkerd) consume Services rather than replace them; the Service remains the discovery unit underneath.

## Comparative Analysis

| Aspect | kubectl | YAML | Terraform | Pulumi | Planton |
|--------|---------|------|-----------|--------|---------|
| Validation | At creation | API server | Plan time (basic) | Preview time (basic) | Schema + CEL |
| Cross-field rules checked early | No | No | No | No | Yes (11 rules) |
| trafficDistribution | Via edit | Yes | **No (provider gap)** | Yes | Yes (Pulumi engine) |
| State Management | None | None | Full | Full | Full (via IaC) |
| LB address as output | Manual | Manual | Yes | Yes | Yes (ip + hostname) |
| Namespace as reference | No | No | Manual wiring | Manual wiring | First-class |
| Dual IaC | N/A | N/A | TF only | Pulumi only | Both |

## The Planton Approach

### Full surface, validated early

The spec models the entire ServiceSpec surface and moves the API server's cross-field rules to validation time. Eleven CEL rules each mirror a live kube-apiserver rejection:

- ExternalName requires `external_dns_name`, and forbids selector, ports, and cluster IP (it is a DNS alias, not a proxy)
- `external_dns_name` is rejected on every other type
- Headless is incompatible with NodePort/LoadBalancer (no virtual IP to build on) and with a static `cluster_ip_address` (headless IS clusterIP "None")
- Every proxying type requires at least one port
- `health_check_node_port` only exists for Local-policy LoadBalancers
- The affinity timeout only applies to ClientIP affinity
- LoadBalancer-only knobs (`load_balancer_source_ranges`, `load_balancer_class`, `allocate_load_balancer_node_ports`) are rejected on other types
- Dual-stack: family entries must be distinct, and SingleStack allows at most one explicit family

Field-level rules do the same for formats: the Service name is validated as a DNS-1035 label (Services reject leading digits that other kinds accept, because DNS SRV labels cannot start with a digit), port names as IANA service names, `target_port` as number-or-name, IPs and CIDRs as such, and `external_dns_name` as an RFC-1123 hostname.

### Semantic booleans over magic strings

Upstream encodes "headless" as the literal string `"None"` in `clusterIP` — a string-typed field carrying two unrelated meanings. The spec separates them: `headless: true` for the discovery mode, `cluster_ip_address` strictly for a static IP (validated as an IP, so `"None"` cannot sneak in). The modules translate `headless: true` back to `clusterIP: "None"` on the wire.

### Omission-correct field handling

For several Service fields, an empty value is not the same as an absent one: `clusterIP`, `healthCheckNodePort`, and `loadBalancerClass` are immutable or type-gated, and sending them empty (or on the wrong type) is an API rejection. Both modules therefore send each optional field only when the user actually set it, guarded by the same conditions in both engines.

### The one parity exception: `traffic_distribution`

The Terraform kubernetes provider (v3.2.x) does not expose `spec.trafficDistribution`, so only the Pulumi engine can apply it. The Terraform module refuses to pretend otherwise: a lifecycle precondition fails the plan loudly when the field is set, with an error message directing the user to the Pulumi engine (or to unset the field). Silently dropping a set field would be worse than failing — a topology hint the user asked for would just not exist, invisibly.

### Namespace by value or reference

`spec.namespace` is a `StringValueOrRef`: a literal namespace name, or a reference to a `KubernetesNamespace` resource. The reference form lets an infra chart create the namespace and the Service in one run, with ordering handled by the resource graph. When omitted, the Service lands in `default` — the same behavior as kubectl without a namespace flag.

### Composition with Planton workloads

Planton workloads stamp a stable selector identity on their pods — the `app` label set to the workload's `metadata.name` — and export the full set as their `selector_labels` output. Selecting a Planton workload's pods from this kind is therefore one line: `selector: {app: <workload-metadata-name>}`. The common reasons to do so are an additional differently-shaped exposure (a LoadBalancer in front of a workload whose built-in Service is internal, a headless companion for direct pod addressing) rather than the ordinary case, which the workload's own Service already covers.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates provider init, Service creation, and output export
- **`locals.go`**: Computes merged labels, annotations, the resolved namespace, and — centrally — the enum translations from proto value names to Kubernetes API strings (`load_balancer` → `LoadBalancer`, `client_ip` → `ClientIP`, `prefer_same_zone` → `PreferSameZone`, ...), resolved once so the resource and the outputs agree on wire values
- **`service.go`**: Builds the ServiceSpec with the omission-correct guards described above
- **`outputs.go`**: Exports the eight outputs; Pulumi's await logic waits for LoadBalancer ingress, so the LB address handles resolve before export

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic guard-for-guard:

- **`variables.tf`**: Mirrors `spec.proto`; enums arrive as proto value names, `StringValueOrRef` fields arrive flattened to plain strings
- **`locals.tf`**: The same enum maps and label merging as the Pulumi locals
- **`main.tf`**: Creates the `kubernetes_service_v1` resource with identical conditional guards, plus the `traffic_distribution` precondition
- **`outputs.tf`**: Exports the identical eight-output set (empty strings where not applicable, so both engines flatten the same field shape onto the outputs proto)

### Outputs contract

Both engines export: `service_name`, `namespace`, `type`, `cluster_ip` (empty for headless and ExternalName — the API's literal `"None"` is normalized away), `load_balancer_ip` and `load_balancer_hostname` (a provider populates one or the other, never reliably both — exported independently so DNS automation picks whichever is present), `kube_endpoint` (the in-cluster DNS name, which for ExternalName is the very alias the service exists to provide), and `port_forward_command` (empty for ExternalName — there is nothing to forward to).

### Resource Count

This is a lean component — it creates exactly **one Kubernetes resource**: the Service itself. For `type: load_balancer` the cloud provider provisions a load balancer as a side effect, driven entirely by the Service object and its annotations. The complexity is in the spec validation and enum translation, not in resource orchestration.

## Production Best Practices

### Exposure discipline

1. **Default to ClusterIP**: Every external exposure is attack surface and (for LoadBalancer) a billed cloud resource. HTTP workloads usually want one Ingress in front of many ClusterIP services, not many LoadBalancers
2. **Tune load balancers through annotations**: internal vs. external, NLB vs. ELB, pinned addresses, external-dns records — the annotation set is the reproducible record of LB configuration
3. **Restrict sources where supported**: `load_balancer_source_ranges` on LoadBalancers that only specific CIDRs should reach

### Traffic correctness

1. **Use `external_traffic_policy: local` when the client IP matters**: and run enough replicas that every schedulable node is likely to hold one — nodes without endpoints are health-checked out and receive nothing
2. **Reserve `publish_not_ready_addresses` for bootstrap discovery**: on the headless governing Service of a quorum system it is essential; on a traffic-serving service it routes real requests to pods that cannot handle them
3. **Treat `traffic_distribution` as a hint**: set `prefer_same_zone` only when endpoints are spread evenly enough that zone-local preference cannot overload one zone

### Naming and ports

1. **Name every port on multi-port services**: required by the API; named ports keep Ingress and mesh references readable
2. **Prefer named `target_port` over numbers**: `target_port: "http"` keeps working when the container port number changes
3. **Remember Service names are DNS names**: DNS-1035 rules apply — a name that starts with a digit is rejected

## Conclusion

KubernetesService is the escape hatch and the composition point: workload kinds cover their own pods, and this kind covers every other shape a Service can take — external exposure, DNS aliases, headless discovery, dual-stack, selectorless endpoints. The component's value is the full upstream surface with the API server's own rules enforced before apply, semantic fields where upstream uses magic strings, and a dual-IaC lifecycle whose one parity gap fails loudly instead of silently.

## References

- [Kubernetes Service Documentation](https://kubernetes.io/docs/concepts/services-networking/service/)
- [Service API Reference](https://kubernetes.io/docs/reference/kubernetes-api/service-resources/service-v1/)
- [DNS for Services and Pods](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/)
- [Dual-stack Services](https://kubernetes.io/docs/concepts/services-networking/dual-stack/)
- [Traffic Distribution](https://kubernetes.io/docs/concepts/services-networking/service/#traffic-distribution)
- [Source IP and External Traffic Policy](https://kubernetes.io/docs/tutorials/services/source-ip/)
- [Pulumi Kubernetes Service](https://www.pulumi.com/registry/packages/kubernetes/api-docs/core/v1/service/)
- [Terraform kubernetes_service_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service_v1)
