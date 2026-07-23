# Kubernetes Ingress

Deploys a Kubernetes `networking/v1` Ingress to a target cluster through a single declarative manifest, covering the complete IngressSpec surface: ingress class selection, a default backend, TLS termination blocks, and host/path rules with all three path types. The IaC module handles label merging, namespace resolution, and reference resolution automatically — and deliberately never blocks waiting for an ingress controller.

## What Gets Created

When you deploy a KubernetesIngress resource, Planton provisions:

- **Ingress** — a `networking/v1` Ingress with the given class, default backend, TLS blocks, and rules
- **Labels** — standard Planton tracking labels merged with any user-provided labels
- **Annotations** — user-provided annotations applied to the Ingress metadata; this is where controller-specific behavior (rewrites, body sizes, cert-manager issuers, external-dns hostnames) is configured

The Ingress is created **without waiting for a controller** (Terraform `wait_for_load_balancer = false`, Pulumi `skipAwait`). An Ingress is inert until an ingress controller claims it via `ingress_class_name` or the cluster's default class — the `load_balancer_ip`/`load_balancer_hostname` outputs stay empty until then and fill in once a controller reconciles the object.

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **Backend Services** in the same namespace as the Ingress — a Kubernetes API constraint; backends can never reference Services in other namespaces
- **An ingress controller** (ingress-nginx, ALB, GCE, ...) for the Ingress to actually serve traffic — not required at creation time
- For TLS: either **cert-manager** installed (issuer annotation + a Secret name for it to create) or an existing `kubernetes.io/tls` **Secret**

## Quick Start

Create a file `ingress.yaml`:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesIngress
metadata:
  name: web-ingress
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesIngress.web-ingress
spec:
  name: web-ingress
  namespace:
    value: web
  ingress_class_name: nginx
  rules:
    - host: app.example.com
      paths:
        - path: /
          path_type: prefix
          backend:
            service_name:
              value: web-svc
            port_number: 8080
```

Deploy:

```shell
planton apply -f ingress.yaml
```

This creates an Ingress named `web-ingress` in the `web` namespace, routing `app.example.com` to port 8080 of the `web-svc` Service.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.name` | `string` | Name of the Ingress (`metadata.name` in the cluster). | 1–253 characters, DNS subdomain |
| `spec.rules` or `spec.default_backend` | — | At least one of the two must be present — an Ingress with neither routes no traffic. | Cross-field rule, rejected at validation |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.namespace` | `StringValueOrRef` | `default` | Namespace of the Ingress — must be the namespace of the backend Services. Literal name (`{ value: web }`) or a `KubernetesNamespace` reference. |
| `spec.ingress_class_name` | `string` | cluster default | Which controller serves this Ingress (`nginx`, `alb`, `gce`). `kubectl get ingressclass` lists what the cluster offers. On clusters without a default class, a class-less Ingress is not served. |
| `spec.default_backend` | `object` | — | Backend for requests no rule matches; set it alone (no rules) to expose one Service on all traffic. |
| `spec.tls` | `list` | `[]` | TLS blocks — each names the `hosts` served under one certificate and the `secret_name` of the `kubernetes.io/tls` Secret (literal or KubernetesSecret reference). Multiplexed on 443 via SNI. |
| `spec.rules` | `list` | `[]` | Host rules: requests match on Host header, then on `paths`. Hosts are precise (`app.example.com`) or single-label wildcards (`*.example.com`); omit to match every host. |
| `spec.labels` / `spec.annotations` | `map<string,string>` | `{}` | Merged with standard Planton labels; annotations carry controller-specific behavior. |

### Backends

Every backend (in `default_backend` and in each path) references a Service in the Ingress's own namespace:

| Field | Type | Description |
|-------|------|-------------|
| `service_name` | `StringValueOrRef` | Literal Service name or a `KubernetesService` reference — where a workload's exported `service` output is wired in. Required. |
| `port_number` | `int32` | Service port by number (1–65535). Exactly one of `port_number`/`port_name`. |
| `port_name` | `string` | Service port by name (e.g. `http`) — survives port-number changes. Exactly one of the two. |

### Path Types

| Value | Semantics |
|-------|-----------|
| `prefix` (default) | Match per path element: `/api` matches `/api` and `/api/users`, not `/apiary`. Longest match wins. Portable across all controllers. |
| `exact` | Exact, case-sensitive match. |
| `implementation_specific` | Delegated to the IngressClass (e.g. regex in ingress-nginx). Non-portable. Only type allowed to omit `path`. |

## Examples

### TLS with cert-manager

The Secret named in `tls` does not exist yet — the issuer annotation instructs cert-manager to create it under exactly that name:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesIngress
metadata:
  name: secure-web
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesIngress.secure-web
spec:
  name: secure-web
  namespace:
    value: web
  ingress_class_name: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  tls:
    - hosts:
        - app.example.com
      secret_name:
        value: app-example-com-tls
  rules:
    - host: app.example.com
      paths:
        - path: /
          path_type: prefix
          backend:
            service_name:
              value: web-svc
            port_number: 8080
```

### Path Fan-Out

One host, multiple backends — frontend on `/`, API on `/api` (by named port):

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesIngress
metadata:
  name: app-fanout
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesIngress.app-fanout
spec:
  name: app-fanout
  namespace:
    value: app
  ingress_class_name: nginx
  rules:
    - host: app.example.com
      paths:
        - path: /
          path_type: prefix
          backend:
            service_name:
              value: frontend-svc
            port_number: 3000
        - path: /api
          path_type: prefix
          backend:
            service_name:
              value: api-svc
            port_name: http
```

### Default Backend Only

Everything the controller routes here goes to one Service — no host or path matching:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesIngress
metadata:
  name: catch-all
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesIngress.catch-all
spec:
  name: catch-all
  namespace:
    value: web
  ingress_class_name: nginx
  default_backend:
    service_name:
      value: web-svc
    port_number: 8080
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `ingressName` | `string` | Name of the Ingress object as created in the cluster |
| `namespace` | `string` | Namespace the Ingress was created in |
| `loadBalancerIp` | `string` | IP the controller's load balancer exposes (IP-based controllers: GCE, ingress-nginx on most clouds); empty until a controller reconciles the Ingress |
| `loadBalancerHostname` | `string` | DNS hostname the controller's load balancer exposes (hostname-based controllers: AWS ALB/ELB); empty until a controller reconciles the Ingress |
| `firstHost` | `string` | First host declared in the rules — the primary public FQDN, ready for DNS records and smoke tests; empty for host-less or default-backend-only Ingresses |

## Related Components

- [KubernetesDeployment](/docs/catalog/kubernetes/kubernetesdeployment) — workloads export a `service` output that Ingress backends route to
- [KubernetesService](/docs/catalog/kubernetes/kubernetesservice) — the default reference kind for backend `service_name`
- [KubernetesSecret](/docs/catalog/kubernetes/kubernetessecret) — the default reference kind for `tls[].secret_name`; holds the certificate when cert-manager is not in play
- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) — provides the target namespace; reference it from `spec.namespace` to deploy both in one run
