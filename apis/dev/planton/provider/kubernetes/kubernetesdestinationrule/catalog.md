# Destination Rule on Kubernetes

Defines an Istio DestinationRule: a namespaced resource that tunes *how* traffic is sent to a service once routing has chosen it. It controls load balancing, connection-pool sizing, circuit breaking (outlier detection), and the TLS the sidecar originates upstream, and it can carve a service into named subsets (versions) that route rules target. This is how you add resilience, session affinity, and secure egress to a service without changing the application.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **A DestinationRule** -- a namespaced Istio policy that applies a traffic policy to a service host, plus any named subsets and per-port overrides.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

## What It Controls

- **Host** -- the service from the registry this rule applies to (a Kubernetes service or a host declared by a ServiceEntry). A rule for an unknown host is ignored.
- **Traffic Policy** -- load balancing (simple algorithm or consistent-hash session affinity, plus locality settings and warmup), connection-pool limits (TCP and HTTP), outlier detection (circuit breaking), and the client TLS the sidecar originates upstream. Settings apply across all ports unless a per-port override is set.
- **Subsets** -- named versions of the service selected by label, each able to override the traffic policy. Route rules send traffic to subsets by name for canary and blue/green rollouts.

## Important Behavior

The traffic policy is optional -- an empty policy means Envoy's defaults apply. Per-port overrides **fully replace** the destination-level policy for that port (there is no field-level merge), so set everything a port needs inside its override. A subset's policy only takes effect once a route rule actually sends traffic to that subset. For short host names, istiod resolves the host relative to this rule's namespace, so prefer fully-qualified hosts. The `host` is a foreign key defaulting to a Kubernetes Service reference: wiring it with `valueFrom` orders this rule after the Service whose traffic it shapes and can never drift from the Service's actual name, while a literal `value:` covers ServiceEntry hosts, external FQDNs, and wildcards. The `workload_selector` and TLS `credential_name` remain plain runtime references (no dependency edge) -- order this rule after the workloads or secret it uses with `metadata.relationships`.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Istio installed** -- the Istio CRDs (Istio Base CRDs on Kubernetes) must be present and the Istio control plane (istiod) running. The rule is only honored where istiod is active.
- **Target namespace exists** -- the rule is created in a specific namespace; reference an existing one or create it first.

## Deploy

### Console

Open the deployment store, find **Destination Rule on Kubernetes**, and click **Deploy**. The creation wizard walks you through the namespace, the destination host, the traffic policy, and any subsets, with guidance at each step.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesDestinationRule
metadata:
  name: reviews-circuit-breaker
  org: acme-corp
  env: prod
spec:
  namespace:
    value: prod-apps
  host:
    value: reviews.prod-apps.svc.cluster.local
  trafficPolicy:
    loadBalancer:
      simple: LEAST_REQUEST
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        http2MaxRequests: 1000
        maxRequestsPerConnection: 10
    outlierDetection:
      consecutive5xxErrors: 5
      interval: 10s
      baseEjectionTime: 30s
      maxEjectionPercent: 50
```

```shell
planton apply -f destination-rule.yaml
```

This caps how many connections and requests callers push at `reviews`, and ejects any host that returns five consecutive 5xx errors -- classic circuit breaking that keeps a struggling backend from taking down its callers.

## Key Configuration

- **Namespace** -- the namespace the rule is created in. It is fixed once created; short host names resolve relative to it.
- **Host** -- the service registry host the rule applies to: reference a Planton-managed Kubernetes Service (the rule then deploys after it), or pass a literal fully-qualified name (e.g. `reviews.prod.svc.cluster.local`) for ServiceEntry hosts and external services.
- **Traffic Policy** -- **load balancing** (Least Request / Round Robin / Random / Passthrough, or consistent-hash on a header, cookie, source IP, or query parameter), **connection pool** (TCP connection caps and timeouts; HTTP request and retry limits), **outlier detection** (consecutive 5xx, scan interval, ejection time and cap), **TLS** (Disable / Simple / Mutual / Istio Mutual), and a **retry budget** bounding concurrent retries as a share of active traffic.
- **Subsets** -- named versions selected by labels, each with optional traffic-policy overrides.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Description |
|------------|-------------|
| `namespace` | The namespace the rule is created in. Reference an existing Namespace on Kubernetes or supply the name directly. |
| `host` | The service whose traffic this rule shapes. Reference an existing Kubernetes Service (resolves to its in-cluster FQDN and orders the rule after it) or supply a literal host for ServiceEntry hosts, external FQDNs, and wildcards. |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources and operators can reference:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `destination_rule_name` | Name of the created DestinationRule (equals `metadata.name`) | Ordering resources that depend on the rule being in place |
| `namespace` | The namespace the rule was created in | Confirming where the rule applies |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

- **Circuit breaking and outlier detection** -- cap connection-pool size and eject hosts that keep returning errors, so a struggling backend cannot take down its callers and the mesh routes around unhealthy endpoints. Start from the **circuit-breaking-outlier-detection** preset.
- **mTLS origination to an egress host** -- have the sidecar originate mutual TLS to an external service, presenting client certificates loaded from a Kubernetes secret, so in-mesh workloads can reach a partner API or managed database that requires client-cert auth. Start from the **mtls-origination-egress** preset.

## Works With

DestinationRule is part of the Istio traffic-management family. It requires the Istio Base CRDs on Kubernetes and a running Istio control plane. It pairs naturally with a **Service Entry** (to bring an external host into the registry so this rule can configure it), a **VirtualService** (to route traffic across the subsets this rule defines), and a **Peer Authentication** (the shared workload selector scopes the rule to specific sidecars). Referencing a Kubernetes Service for the host gives the rule its dependency edge automatically; for the workloads or TLS secret it uses, express ordering through `metadata.relationships`.
