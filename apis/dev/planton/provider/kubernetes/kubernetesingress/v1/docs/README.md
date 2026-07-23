# Kubernetes Ingress: Research Documentation

## Introduction

Ingress is the standard Kubernetes answer to a question every cluster eventually asks: how does HTTP(S) traffic from outside reach Services inside? A `networking/v1` Ingress object declares host rules and path matches routing to Service backends, with optional TLS termination from certificate Secrets. It is deliberately a *declaration*, not an implementation — the object does nothing on its own.

That split is the single most important thing to understand about Ingress: the object is **inert until an ingress controller** runs in the cluster and claims it. The controller (ingress-nginx, AWS ALB, GCE, Traefik, HAProxy, ...) watches Ingress objects, selects the ones addressed to it via `ingress_class_name` or a default IngressClass, and programs a real proxy or cloud load balancer to match. An Ingress with no controller is valid but unserved; its load-balancer status stays empty forever.

Planton's **KubernetesIngress** component brings full lifecycle management to the primitive — the complete `networking/v1` surface, schema-level validation of every constraint the API server would otherwise catch at apply time, reference-based composition with workloads and Secrets, and dual-IaC support.

## Evolution and Historical Context

### Origins (extensions/v1beta1)

Ingress appeared as a beta API in Kubernetes 1.1 (2015) and stayed beta for an unusually long five years. The core idea — a portable, declarative L7 routing object with pluggable implementations — was right, but the surface was underspecified: path matching semantics varied wildly between controllers, and the controller-selection mechanism was a mere annotation (`kubernetes.io/ingress.class`).

### networking/v1 (Kubernetes 1.19, 2020)

The GA API fixed the two big ambiguities:

- **`pathType` became required** with three explicit values: `Prefix` (per-path-element matching, identical on every conforming controller), `Exact`, and `ImplementationSpecific` (explicitly delegated to the IngressClass — e.g. ingress-nginx treats such paths as regex candidates). Portability stopped being accidental
- **IngressClass became a real cluster-scoped object**, replacing the annotation. Classes ship with their controllers; a cluster advertises what it offers via `kubectl get ingressclass`, and one class may be annotated as the default

The GA API also added the `resource` backend variant — an ObjectRef to an arbitrary same-namespace object (e.g. a static-asset bucket CRD) as an alternative to a Service backend. In practice it is controller-specific and rarely implemented.

### What Ingress never became — and Gateway API

Ingress deliberately stayed minimal: no traffic splitting, no header matching, no cross-namespace routing delegation. Everything beyond host/path/backend goes through **annotations**, which became the de-facto extension contract — powerful but non-portable. The Gateway API is the upstream successor for the advanced cases; Ingress remains the stable, universal baseline that every cluster and every controller supports, and it is the right tool for the overwhelming majority of "expose this Service on this hostname" needs.

## The Composition Model

Exposure in Planton is **composed, never embedded**. Three independent kinds cooperate:

1. A **workload kind** (KubernetesDeployment and friends) deploys the app and exports its Service name as a `service` output
2. **KubernetesIngress** routes a hostname to that Service — the backend's `service_name` is a `StringValueOrRef`, so a chart wires the workload's exported output straight in (via valueFrom) without copying names by hand
3. A **certificate** materializes the TLS Secret the `tls` block names — either cert-manager (issuer annotation; it creates the Secret) or a KubernetesSecret resource

Every piece of the exposure path is a visible, independently managed node in the resource graph. The alternative — embedding an ingress sub-spec inside every workload kind — duplicates a half-featured Ingress surface into each workload and hides the routing topology inside opaque modules.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

```bash
kubectl create ingress web --class=nginx \
  --rule="app.example.com/*=web-svc:8080"
```

**Pros:**
- Immediate; fine for a demo cluster

**Cons:**
- The rule mini-language covers a fraction of the surface (no default backend, awkward TLS)
- Imperative — no drift detection, no reproducibility

**Verdict:** Debugging only.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: web
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  ingressClassName: nginx
  tls:
    - hosts: [app.example.com]
      secretName: app-example-com-tls
  rules:
    - host: app.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: web-svc
                port:
                  number: 8080
```

**Pros:**
- Declarative, version-controllable, full surface

**Cons:**
- Constraint violations (both port forms set, missing path on a Prefix path, no rules and no default backend) surface only at the API server — or worse, not at all: an Ingress with neither rules nor a default backend is *accepted* by the API and silently routes nothing
- No state management, no composition with the workload or the certificate

**Verdict:** The baseline. Everything above it adds lifecycle.

### Level 2: Terraform

```hcl
resource "kubernetes_ingress_v1" "web" {
  wait_for_load_balancer = false

  metadata {
    name      = "web"
    namespace = "web"
  }

  spec {
    ingress_class_name = "nginx"
    rule {
      host = "app.example.com"
      http {
        path {
          path      = "/"
          path_type = "Prefix"
          backend {
            service {
              name = "web-svc"
              port { number = 8080 }
            }
          }
        }
      }
    }
  }
}
```

**Pros:**
- Full IaC lifecycle with state and drift detection
- `wait_for_load_balancer` gives explicit control over whether creation blocks on a controller

**Cons:**
- Deeply nested block syntax; easy to mis-structure
- No schema-level validation of host syntax, port exclusivity, or the rules-or-default-backend rule

**Verdict:** Production-grade lifecycle, thin validation.

### Level 3: Pulumi

```go
ingress, err := networkingv1.NewIngress(ctx, "web", &networkingv1.IngressArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name:      pulumi.String("web"),
        Namespace: pulumi.String("web"),
    },
    Spec: &networkingv1.IngressSpecArgs{
        IngressClassName: pulumi.String("nginx"),
        // rules elided
    },
})
```

**Pros:**
- Full programming language; preview before apply
- `pulumi.com/skipAwait` gives the same non-blocking control as Terraform's flag

**Cons:**
- By default Pulumi *waits* for the Ingress to get a load-balancer address — deploys hang forever on clusters where the controller is not up yet, unless skipAwait is set
- Same validation gap as Terraform

**Verdict:** Excellent IaC choice; the await default is a real operational trap.

### Other Methods

**Helm:** Nearly every public chart templates an Ingress behind `ingress.enabled`. Ubiquitous, but each chart re-invents the values shape, and the routing topology is buried inside the release.

**Gateway API:** The successor for advanced L7 (traffic splitting, header routing, cross-namespace delegation). Not yet as universally available; Ingress remains the portable baseline.

## Comparative Analysis

| Aspect | kubectl | YAML | Terraform | Pulumi | Planton |
|--------|---------|------|-----------|--------|---------|
| Validation | At creation | API server | Plan time (basic) | Preview time (basic) | Schema + CEL |
| "Routes nothing" caught | No | No (API accepts it) | No | No | Yes — rejected |
| Port number/name exclusivity | No | API server | Apply time | Apply time | Validation time |
| Blocks on controller | No | No | Configurable | Waits by default | Never (deliberate) |
| Backend as reference to workload output | No | No | Manual wiring | Manual wiring | First-class |
| Dual IaC | N/A | N/A | TF only | Pulumi only | Both |

## The Planton Approach

### Full surface, one deliberate omission

The spec models the complete `networking/v1` IngressSpec: `ingress_class_name`, `default_backend`, `tls[]` (hosts + secret name), and `rules[]` (host + paths, each with path, path type, and backend). The single omission is the **`resource` backend variant** — controller-specific, rarely implemented, and Service backends cover the real exposure paths. Controller-specific behavior (rewrites, timeouts, body sizes, auth) goes through `annotations`, which is the upstream contract.

### Upstream's silent failure modes become validation errors

- **At least one rule or a default backend**: the upstream API accepts an Ingress with neither — it is always a mistake, and it silently routes nothing. A CEL rule rejects it
- **Exactly one of `port_number`/`port_name`** per backend — the API server's own rule, surfaced before deployment
- **`path` required for `prefix` and `exact`** path types; only `implementation_specific` may leave it to the controller
- **Host syntax**: precise DNS names or single-leading-label wildcards (`*.example.com` matches one label — `a.example.com`, not `b.a.example.com`); IPs and ports rejected
- **A rule must carry at least one path**: upstream treats a missing rule value as "controller-defined", which in practice means ignored

### Non-blocking creation, by design

Both modules create the Ingress **without waiting for a controller** — Terraform sets `wait_for_load_balancer = false`, Pulumi sets `pulumi.com/skipAwait` (the exact same choice, mirrored). An Ingress object is valid without a controller; infra charts routinely deploy the workload and its exposure before the ingress controller wave, and blocking every deploy until a controller populates the load-balancer status would couple this kind to cluster addon ordering.

The trade-off is honest outputs: `load_balancer_ip` and `load_balancer_hostname` read the object's live status. On a cluster where a controller reconciles quickly, the values land on the same deploy; on a cluster with no controller yet, they export empty — matching the object's real state — and fill in once one reconciles.

### References where composition happens

- **`namespace`**: literal or KubernetesNamespace reference. It must be the namespace of the backend Services — Ingress backends can only reference Services in their own namespace, a Kubernetes API constraint
- **`backend.service_name`**: literal or KubernetesService reference (default field path `status.outputs.service_name`) — where a workload's exported `service` output lands
- **`tls[].secret_name`**: literal or KubernetesSecret reference. With cert-manager, the Secret deliberately does not exist yet: the issuer annotation instructs cert-manager to create it under exactly the name written here. Omitting it asks the controller to serve those hosts with its default certificate

### Why `ingress_class_name` is a plain string

IngressClass objects ship with their controllers — they are cluster inventory, not user-managed resources. `kubectl get ingressclass` lists what the cluster offers. Omitting the class defers to the cluster's default class (the one annotated `ingressclass.kubernetes.io/is-default-class: "true"`); on clusters without a default, a class-less Ingress is not served by any controller.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates provider init, Ingress creation, output export
- **`locals.go`**: Computes merged labels, annotations, resolved namespace, and the first rule host
- **`ingress.go`**: Builds the `networking/v1` Ingress args (class, default backend, TLS, rules) and adds the `skipAwait` engine annotation on top of user annotations
- **`outputs.go`**: Exports `ingress_name`, `namespace`, `load_balancer_ip`, `load_balancer_hostname`, `first_host` — the load-balancer reads are nil-tolerant because the status may legitimately be empty

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto`; StringValueOrRef fields arrive flattened to plain strings, enums as proto value names (`prefix`, `exact`, `implementation_specific`)
- **`locals.tf`**: Merged labels, namespace default, the proto-enum-to-API path-type map (`prefix` → `Prefix`), and the first rule host
- **`main.tf`**: One `kubernetes_ingress_v1` resource with `wait_for_load_balancer = false`
- **`outputs.tf`**: The same five outputs, with `try()`-guarded status reads

### Resource Count

This is a lean component — it creates exactly **one Kubernetes resource**: the Ingress itself. The complexity is in the validation and the composition, not in resource orchestration. The controller, the backend Services, and the certificates are all separate nodes in the graph, which is the point.

## Production Best Practices

### Routing discipline

1. **Pin `ingress_class_name` explicitly**: relying on the cluster default class works until it doesn't (a second controller, a cluster migration). Explicit classes are portable and reviewable
2. **Prefer `prefix` path types**: the one matching semantics every conforming controller implements identically. Reserve `implementation_specific` for controller features (e.g. nginx regex paths) you consciously depend on — they are non-portable by definition
3. **Prefer named ports**: `port_name: http` survives Service port-number refactors; numbers silently break
4. **Mind prefix semantics**: `/api` matches `/api` and `/api/users` but not `/apiary` — matching is per path element, and the longest matching path wins

### TLS

1. **Let cert-manager own the Secret**: name it in `tls`, add the issuer annotation, and never handle certificate material by hand. The Secret appears under exactly the name written in `secret_name`
2. **Keep TLS hosts and rule hosts aligned**: TLS hosts must appear in the certificate's SANs and should match the hosts in `rules`; the controller multiplexes entries on 443 via SNI
3. **Wildcard entries follow single-label semantics** — `*.example.com` covers `a.example.com` only

### Placement and composition

1. **Backends live in the Ingress's namespace** — absolute API constraint. Plan namespace layout around it; there is no cross-namespace backend
2. **Wire backends by reference in charts**: deploy the workload, reference its `service` output — no hand-copied names to drift
3. **Compose DNS on the outputs**: point records at `load_balancer_ip`/`load_balancer_hostname`, or let external-dns read the same status; `first_host` is the ready-made primary FQDN for smoke tests and dashboards
4. **Don't panic at empty load-balancer outputs**: they mean no controller has reconciled the object yet — check that a controller is installed and that the class matches

## Conclusion

KubernetesIngress is a deliberately complete, deliberately lean component: the full `networking/v1` surface minus one controller-specific corner, upstream's silent failure modes promoted to validation errors, creation decoupled from controller availability, and every point of composition — namespace, backends, TLS Secrets — expressed as a reference the resource graph can order. The Ingress object stays what upstream designed it to be: a small, portable routing declaration, with the controller, the workloads, and the certificates as visible neighbors rather than hidden internals.

## References

- [Kubernetes Ingress Documentation](https://kubernetes.io/docs/concepts/services-networking/ingress/)
- [Ingress Controllers](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/)
- [IngressClass](https://kubernetes.io/docs/concepts/services-networking/ingress/#ingress-class)
- [Ingress API Reference](https://kubernetes.io/docs/reference/kubernetes-api/service-resources/ingress-v1/)
- [cert-manager: Securing Ingress Resources](https://cert-manager.io/docs/usage/ingress/)
- [Gateway API](https://gateway-api.sigs.k8s.io/)
- [Pulumi Kubernetes Ingress](https://www.pulumi.com/registry/packages/kubernetes/api-docs/networking/v1/ingress/)
- [Terraform kubernetes_ingress_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/ingress_v1)
