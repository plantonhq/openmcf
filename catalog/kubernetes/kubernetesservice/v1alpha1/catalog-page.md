# Kubernetes Service

Deploys a standalone Kubernetes Service to a target cluster through a single declarative manifest, covering the complete core/v1 ServiceSpec surface: all four service types (ClusterIP, NodePort, LoadBalancer, ExternalName), headless services, traffic policies, session affinity, LoadBalancer tuning, and dual-stack addressing. The IaC module handles label merging, namespace resolution, and enum translation automatically.

Planton workload kinds (KubernetesDeployment, KubernetesStatefulSet) already create a Service for their own pods — use this standalone kind for everything else: exposing pods managed outside Planton, LoadBalancer/NodePort exposure with cloud annotations, ExternalName aliases, headless discovery, dual-stack, and selectorless services.

## What Gets Created

When you deploy a KubernetesService resource, Planton provisions:

- **Service** — a Kubernetes Service of the requested type with ports, selector, traffic policies, and dual-stack settings
- **Cloud load balancer** — for `type: load_balancer`, the cluster's cloud provider provisions an external load balancer configured through the manifest's annotations
- **Labels** — standard Planton tracking labels (resource name, kind, id, organization, environment) merged with any user-provided labels
- **Annotations** — user-provided annotations applied to the Service metadata (the portable way to tune cloud load balancers)

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that already exists, or a `KubernetesNamespace` resource referenced from `spec.namespace` so both deploy in one run
- **Backend pods** carrying the labels your `selector` lists (not needed for ExternalName or selectorless services)

## Quick Start

Create a file `service.yaml`:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesService
metadata:
  name: my-app
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesService.my-app
spec:
  name: my-app
  selector:
    app: my-app
  ports:
    - name: http
      port: 80
      target_port: "8080"
```

Deploy:

```shell
planton apply -f service.yaml
```

This creates a ClusterIP Service named `my-app` in the `default` namespace, routing port 80 to container port 8080 on pods labeled `app: my-app`. To select a Planton-managed workload's pods, set `app` to the workload's `metadata.name`.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.name` | `string` | Name of the Service (`metadata.name` in the cluster); becomes its DNS name. | DNS-1035 label: 1–63 chars, matches `^[a-z]([-a-z0-9]*[a-z0-9])?$` (must start with a letter) |
| `spec.ports` | `list` | Ports the Service exposes. | At least one port for every type except `external_name` |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.namespace` | `StringValueOrRef` | `default` | Namespace where the Service is created. Accepts a literal name (`{ value: my-namespace }`) or a reference to a `KubernetesNamespace` resource. |
| `spec.type` | `enum` | `cluster_ip` | One of `cluster_ip`, `node_port`, `load_balancer`, `external_name`. |
| `spec.selector` | `map<string, string>` | `{}` | Pods whose labels match ALL listed pairs receive the traffic. Empty for ExternalName or selectorless services. |
| `spec.headless` | `bool` | `false` | Creates a headless service (`clusterIP: None`) — DNS returns pod IPs directly. Incompatible with `node_port`/`load_balancer` and with `cluster_ip_address`. |
| `spec.cluster_ip_address` | `string` | allocated | A specific virtual IP inside the cluster's service CIDR. Immutable after creation; rarely needed. |
| `spec.external_dns_name` | `string` | — | CNAME target for `external_name` services (e.g. `db.prod.example.com`). Required for, and only allowed on, that type. |
| `spec.external_ips` | `list<string>` | `[]` | IPs outside Kubernetes' management for which nodes also accept traffic. |
| `spec.external_traffic_policy` | `enum` | `cluster` | `cluster` or `local` — routing for NodePort/LoadBalancer traffic. `local` preserves the client source IP. |
| `spec.internal_traffic_policy` | `enum` | `internal_cluster` | `internal_cluster` or `internal_local` — routing for ClusterIP traffic. |
| `spec.traffic_distribution` | `enum` | unset | `prefer_same_zone` or `prefer_same_node` topology hint. **Deploys only through the Pulumi engine** — the Terraform module fails the plan loudly when set. |
| `spec.session_affinity` | `enum` | `none` | `none` or `client_ip` (sticky sessions by client IP). |
| `spec.session_affinity_timeout_seconds` | `int32` | 10800 | ClientIP pin duration, 1–86400. Only with `session_affinity: client_ip`. |
| `spec.load_balancer_source_ranges` | `list<string>` | `[]` | Client CIDRs allowed to reach the load balancer. LoadBalancer type only. |
| `spec.load_balancer_class` | `string` | cluster default | Selects the LB implementation when the cluster runs several (e.g. MetalLB alongside the cloud default). Immutable. |
| `spec.allocate_load_balancer_node_ports` | `bool` | `true` (Kubernetes) | Set `false` when the LB implementation routes to pods directly and the NodePort hop is dead weight. |
| `spec.health_check_node_port` | `int32` | allocated | Specific health-check NodePort for `external_traffic_policy: local` LoadBalancers. 30000–32767; immutable. |
| `spec.publish_not_ready_addresses` | `bool` | `false` | Publish pod addresses before Ready — for StatefulSet bootstrap discovery. |
| `spec.ip_families` | `list<enum>` | cluster-assigned | `ipv4` / `ipv6` in preference order; at most two distinct entries. |
| `spec.ip_family_policy` | `enum` | SingleStack | `single_stack`, `prefer_dual_stack`, or `require_dual_stack`. |
| `spec.labels` / `spec.annotations` | `map<string, string>` | `{}` | Labels merge with standard Planton labels; annotations carry cloud LB configuration. |

**Ports:** each entry carries `port` (1–65535, what clients connect to), optional `target_port` (a number `"8080"` or a named container port `"http"`; defaults to `port`), `protocol` (`TCP`/`UDP`/`SCTP`, default TCP), optional `app_protocol` (L7 hint like `http` or `kubernetes.io/h2c`), and optional `node_port` (30000–32767, NodePort/LoadBalancer only). Port `name` is required and unique when there is more than one port.

**Cross-field rules** (each mirrors a live kube-apiserver rejection, caught at validation time): ExternalName requires `external_dns_name` and forbids selector/ports/cluster IP; headless forbids NodePort/LoadBalancer and a static cluster IP; LoadBalancer-only knobs are rejected on other types; the affinity timeout requires ClientIP affinity; dual-stack family lists must be consistent with the policy.

## Examples

### Expose Non-Planton Pods Internally

A ClusterIP service in front of pods deployed by a Helm chart:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesService
metadata:
  name: legacy-api
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesService.legacy-api
spec:
  name: legacy-api
  namespace:
    value: legacy
  selector:
    app.kubernetes.io/name: legacy-api
  ports:
    - name: http
      port: 80
      target_port: "8080"
```

### Public Load Balancer with Cloud Annotations

An AWS NLB preserving the client source IP:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesService
metadata:
  name: public-web
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesService.public-web
spec:
  name: public-web
  namespace:
    value: production
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-type: "nlb"
  type: load_balancer
  selector:
    app: web
  ports:
    - name: https
      port: 443
      target_port: "8443"
  external_traffic_policy: local
  load_balancer_source_ranges:
    - 203.0.113.0/24
```

### Headless Service for StatefulSet Peers

Per-pod DNS with addresses published during bootstrap:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesService
metadata:
  name: db-peers
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesService.db-peers
spec:
  name: db-peers
  namespace:
    value: databases
  headless: true
  publish_not_ready_addresses: true
  selector:
    app: postgres
  ports:
    - name: peer
      port: 5432
```

### ExternalName Alias

One stable in-cluster name for an external database:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesService
metadata:
  name: prod-db
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesService.prod-db
spec:
  name: prod-db
  namespace:
    value: production
  type: external_name
  external_dns_name: db.prod.example.com
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `serviceName` | `string` | Name of the Service object as created in the cluster |
| `namespace` | `string` | Namespace the Service was created in |
| `type` | `string` | Service type as deployed (ClusterIP, NodePort, LoadBalancer, ExternalName) |
| `clusterIp` | `string` | Cluster-internal virtual IP — empty for headless and ExternalName services |
| `loadBalancerIp` | `string` | Provisioned load balancer IP — populated only on providers that expose an IP (GCP, Azure, MetalLB) |
| `loadBalancerHostname` | `string` | Provisioned load balancer hostname — populated only on providers that expose a hostname (AWS ELB/NLB) |
| `kubeEndpoint` | `string` | In-cluster DNS endpoint, e.g. `my-app.my-ns.svc.cluster.local` |
| `portForwardCommand` | `string` | Ready-to-run `kubectl port-forward` command; empty for ExternalName services |

## Related Components

- [KubernetesDeployment](/docs/catalog/kubernetes/kubernetesdeployment) — workload kinds already create a Service for their own pods; select their pods here with `app: <workload-metadata-name>` when you need additional exposure
- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) — provides the target namespace; reference it from `spec.namespace` to deploy both in one run
- [KubernetesIngress](/docs/catalog/kubernetes/kubernetesingress) — L7 routing (hostnames, paths, TLS) in front of a Service backend
