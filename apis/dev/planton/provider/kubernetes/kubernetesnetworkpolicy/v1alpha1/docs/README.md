# Kubernetes Network Policy: Research Documentation

## Introduction

A Kubernetes cluster ships with a deliberately open network: every pod can reach every other pod, in any namespace, on any port. That default is what makes the platform easy to start with — and what makes it unacceptable for multi-tenant or security-conscious production use. NetworkPolicy is the standard, portable answer: a namespaced `networking/v1` resource that selects pods and declares which traffic is allowed to and from them. It is the in-cluster firewall.

The resource's mental model is unusual and is the source of nearly every authoring mistake, so it is worth stating precisely:

- **Policies are additive allows. There is no deny rule.** A pod not selected by any policy accepts and sends all traffic. Once a pod is selected by any policy in a direction, only traffic allowed by *some* policy in that direction is admitted; everything else is dropped.
- **Multiple policies union.** You can widen what is allowed by adding policies; you can never narrow another policy's allows.
- **Isolation is therefore achieved by selecting pods while allowing little** — in the limit, selecting everything and allowing nothing (default-deny).

Planton's **KubernetesNetworkPolicy** component brings the full `networking/v1` surface to the platform with schema-level validation that catches the classic mistakes before apply, namespace composition, and dual-IaC support.

## Evolution and Historical Context

### Origins (Kubernetes 1.3–1.7)

NetworkPolicy entered as beta in Kubernetes 1.3 (2016) and graduated to `networking.k8s.io/v1` in 1.7 (2017). The original surface governed ingress only: a pod selector plus allow rules for inbound traffic.

### Egress and ipBlock (1.8+)

Kubernetes 1.8 added egress rules and the `ipBlock` peer, completing the direction pair and extending policies beyond cluster-internal selectors to external CIDRs. With egress came the `policyTypes` field and its inference rule — ingress is always inferred, egress only when egress rules are present — which exists for backward compatibility with pre-egress policies and remains a live footgun: an egress-only or deny-all-egress policy that omits `policyTypes` silently governs the wrong directions.

### Port ranges (1.25) and namespace labels (1.22)

`endPort` graduated in Kubernetes 1.25, allowing contiguous port ranges instead of one rule per port. Separately, Kubernetes 1.22 began automatically labelling every namespace with `kubernetes.io/metadata.name: <name>` — finally giving namespace selectors a guaranteed, by-name handle without requiring teams to maintain their own namespace-labelling conventions.

### What NetworkPolicy never became

Upstream deliberately kept NetworkPolicy L3/L4, namespaced, and allow-only. Deny rules, priorities, cluster-scoped policies, and L7 matching were all left to CNI-specific CRDs (Calico's policies, CiliumNetworkPolicy) and to newer efforts like the AdminNetworkPolicy API. The portable core is small on purpose — and its enforcement was *also* deliberately left out of the core: the API server stores NetworkPolicy objects, but only the CNI plugin enforces them.

### The enforcement gap

This is the most operationally important fact about the resource. On a cluster whose CNI does not implement NetworkPolicies — including default kind clusters with kindnet, and cloud clusters without policy enforcement enabled — the object is accepted, stored, and visible in `kubectl get networkpolicy`, while all traffic continues to flow. Enforcement requires Calico, Cilium, or a cloud CNI with policy enforcement enabled (GKE Dataplane V2, AKS with Azure/Calico network policy, EKS with a policy-capable CNI configuration). Any team relying on NetworkPolicies for security must verify enforcement per cluster, not per manifest.

## The Semantics in Detail

### The selector grid

A policy has one top-level `podSelector` choosing the governed pods, and rules containing peers choosing the allowed traffic. Both use standard label selectors, and the empty selector is meaningful everywhere:

- Empty top-level `podSelector` → ALL pods in the namespace (the default-deny building block)
- Empty `podSelector` in a peer → all pods (in the policy's namespace, or in the peer's selected namespaces)
- Empty `namespaceSelector` in a peer → all namespaces (the cluster-wide-allow building block)

### AND vs OR — the classic mistake

Within one peer, `podSelector` and `namespaceSelector` AND: pods matching the selector in namespaces matching the selector. As two separate peers they OR: any pod in the selected namespaces, plus matching pods in the policy's own namespace. In raw YAML the difference is a two-character indentation change (`- namespaceSelector:` + `podSelector:` in one list item vs two list items) and it inverts the security meaning of the policy. Rules themselves OR with each other; within a rule, peers AND ports.

### Direction inference

When `policyTypes` is omitted, Kubernetes infers ingress always and egress only when egress rules exist. Three consequences:

1. A pure-allow ingress policy can safely omit it
2. An egress-only policy MUST set `[egress]` or it also isolates ingress
3. A deny-all-egress policy MUST set `[egress]` explicitly — with no egress rules present, there is nothing to infer egress from

### Egress and DNS

A deny-all-egress policy blocks DNS resolution too, which breaks nearly everything in ways that do not obviously point at the network policy (timeouts, `no such host`). The standard companion is an egress allow for UDP+TCP port 53 to the cluster DNS pods in `kube-system`.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

There is no `kubectl create networkpolicy` generator — the resource is too structural. Manual means `kubectl apply -f` of hand-written YAML.

**Verdict:** No shortcut exists; even ad-hoc use is YAML authoring.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: backend-api-ingress
  namespace: backend
spec:
  podSelector:
    matchLabels:
      app: backend-api
  policyTypes: [Ingress]
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: frontend
      ports:
        - protocol: TCP
          port: 8080
```

**Pros:**
- Declarative, version-controllable, full surface

**Cons:**
- The AND/OR indentation trap and the `policyTypes` inference trap are invisible until traffic breaks (or worse, doesn't)
- Rules in an ungoverned direction are silently ignored by the API server
- No plan/preview, no state management

**Verdict:** The baseline, with the sharpest edges of any core Kubernetes resource.

### Level 2: Terraform

```hcl
resource "kubernetes_network_policy_v1" "backend_api_ingress" {
  metadata {
    name      = "backend-api-ingress"
    namespace = "backend"
  }
  spec {
    pod_selector {
      match_labels = { app = "backend-api" }
    }
    policy_types = ["Ingress"]
    ingress {
      from {
        pod_selector {
          match_labels = { app = "frontend" }
        }
      }
      ports {
        protocol = "TCP"
        port     = "8080"
      }
    }
  }
}
```

**Pros:**
- Full IaC lifecycle (plan, apply, destroy, import) with drift detection

**Cons:**
- HCL blocks reproduce the same AND/OR and inference semantics without additional guardrails
- CIDR formats, port ranges, and selector operator contracts surface only at apply

**Verdict:** Production-grade lifecycle, thin validation.

### Level 3: Pulumi

```go
networkPolicy, err := networkingv1.NewNetworkPolicy(ctx, "backend-api-ingress", &networkingv1.NetworkPolicyArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name:      pulumi.String("backend-api-ingress"),
        Namespace: pulumi.String("backend"),
    },
    Spec: &networkingv1.NetworkPolicySpecArgs{
        PodSelector: &metav1.LabelSelectorArgs{
            MatchLabels: pulumi.StringMap{"app": pulumi.String("backend-api")},
        },
        PolicyTypes: pulumi.ToStringArray([]string{"Ingress"}),
    },
})
```

**Pros:**
- Full programming language, preview before apply

**Cons:**
- Types describe the wire shape, not the semantics; the same traps pass the compiler

**Verdict:** Excellent IaC choice; validation gap same as Terraform.

### Other Methods

**Helm:** policies templated per chart — common, but template conditionals layered on already-subtle semantics multiply the review burden.

**CNI-specific CRDs (Calico, Cilium):** strictly more expressive (deny rules, priorities, L7, cluster scope) at the cost of portability. The right tool when requirements exceed core NetworkPolicy; overkill for the standard segmentation patterns.

## Comparative Analysis

| Aspect | YAML | Terraform | Pulumi | Planton |
|--------|------|-----------|--------|---------|
| Validation | API server | Plan time (shape only) | Preview time (shape only) | Schema + CEL |
| Rules-vs-direction mismatch caught | No (silently ignored) | No | No | Yes, rejected pre-apply |
| CIDR/port/operator contracts checked early | No | No | No | Yes |
| Deterministic policyTypes across engines | N/A | Provider-dependent | Provider-dependent | Always explicit |
| Namespace as reference | No | Manual wiring | Manual wiring | First-class |
| Dual IaC | N/A | TF only | Pulumi only | Both |

## The Planton Approach

### Full surface, validated early

The spec models the entire `networking/v1` NetworkPolicySpec — pod selector, policy types, ingress/egress rules, all three peer forms, ports with names and ranges — and moves the API server's own rules plus the known footguns to validation time:

- **Direction/rule consistency**: ingress rules with `policy_types: [egress]` (or vice versa) are rejected outright — upstream Kubernetes silently ignores such rules
- **Peer integrity**: an empty peer (matches nothing, always a mistake) is rejected; `ip_block` cannot be combined with selectors in the same peer, mirroring the API's own constraint
- **CIDR and port contracts**: CIDR format, numeric port ranges (1–65535), `end_port` requiring a numeric anchor and correct ordering, named-port character rules
- **Selector operator contracts**: `In`/`NotIn` require values, `Exists`/`DoesNotExist` forbid them — the exact admission rule, surfaced before deployment
- **Distinct policy types**: duplicate directions are rejected

### Deterministic policy types

Both IaC modules resolve the governed directions before submission — the explicit `policy_types` when set, otherwise the Kubernetes inference rule (ingress always, egress only with egress rules) — and always send the result explicitly. Both engines submit byte-identical direction sets for the same manifest, and the exported `policy_types` output states the deployed truth even when the spec omitted the field.

### Peers map one-to-one

Each spec peer maps to exactly one API peer and is never split or merged, preserving the AND semantics of a combined pod+namespace selector. The AND/OR distinction the user wrote is the AND/OR distinction that deploys.

### Namespace by value or reference

`spec.namespace` is a `StringValueOrRef`: a literal namespace name, or a reference to a `KubernetesNamespace` resource, letting an infra chart create a namespace, its default-deny policy, and its targeted allows in one run with ordering handled by the resource graph. When omitted, the policy lands in `default`.

### The workload label contract

Every Planton workload kind stamps the `app` label — set to the workload's `metadata.name` — on its pods as immutable selection identity. `match_labels: {app: <workload-name>}` is therefore the standard way policies compose with Planton workloads, both as the governed pods and as peers.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates provider init, resource creation, and output export
- **`locals.go`**: Computes merged labels, annotations, the resolved namespace, and the resolved policy types (explicit or inferred)
- **`networkpolicy.go`**: Creates the `networking/v1` NetworkPolicy, converting peers, ports, and selectors one-to-one
- **`outputs.go`**: Exports `network_policy_name`, `namespace`, and `policy_types`

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto` fields as Terraform variables
- **`locals.tf`**: Computes merged labels, resolved namespace, and the same explicit-or-inferred policy types
- **`main.tf`**: Creates the `kubernetes_network_policy_v1` resource
- **`outputs.tf`**: Exports the same three outputs

### Resource Count

This is a lean component — it creates exactly **one Kubernetes resource**: the NetworkPolicy itself. The complexity is in the spec validation and the semantic guardrails, not in resource orchestration.

## Production Best Practices

### Rollout discipline

1. **Verify enforcement before relying on policies**: confirm the cluster's CNI actually enforces NetworkPolicies; on a non-enforcing CNI every policy is a no-op
2. **Default-deny first, then targeted allows**: apply the namespace-wide deny, then one small policy per workload or flow — additive semantics make each addition safe
3. **Test connectivity, not just apply success**: a policy that applies cleanly proves nothing; probe the actual flows (`kubectl exec ... -- nc -z <svc> <port>`) from allowed and denied peers

### Authoring discipline

1. **Be explicit with `policy_types` whenever the intent is isolation**: egress-only and deny-all-egress policies must declare their directions
2. **Never deny egress without allowing DNS**: pair deny-all-egress with UDP+TCP 53 to the cluster DNS pods
3. **Watch the AND/OR boundary**: one peer with both selectors restricts; two peers widen — review every multi-selector rule with this question
4. **Use selectors for cluster traffic, CIDRs only for external**: pod IPs are ephemeral; `ip_block` is for VPC ranges, external services, and the internet
5. **Prefer the automatic `kubernetes.io/metadata.name` label** for by-name namespace selection — it is guaranteed on every namespace

### Placement

1. **One policy per intent**: small, named, single-purpose policies (`default-deny-all`, `allow-frontend-to-api`, `allow-dns-egress`) audit and compose better than monoliths
2. **Policies live with their pods**: a NetworkPolicy governs only its own namespace; per-namespace duplication is by design, not an anti-pattern

## Conclusion

KubernetesNetworkPolicy is a deliberately complete, deliberately lean component: the full upstream surface, with the resource's famously sharp semantics — allow-only additivity, empty-selector meanings, AND'd vs OR'd peers, direction inference — documented in the schema and guarded by validation that rejects the silent-failure cases before they reach a cluster. Combined with the workload label contract and namespace references, it makes default-deny-plus-targeted-allows a pattern that can be stamped out per namespace rather than hand-crafted per incident.

## References

- [Kubernetes Network Policies Documentation](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [Declare Network Policy Walkthrough](https://kubernetes.io/docs/tasks/administer-cluster/declare-network-policy/)
- [NetworkPolicy API Reference](https://kubernetes.io/docs/reference/kubernetes-api/policy-resources/network-policy-v1/)
- [Automatic namespace labelling](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#automatic-labelling)
- [Pulumi Kubernetes NetworkPolicy](https://www.pulumi.com/registry/packages/kubernetes/api-docs/networking/v1/networkpolicy/)
- [Terraform kubernetes_network_policy_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/network_policy_v1)
