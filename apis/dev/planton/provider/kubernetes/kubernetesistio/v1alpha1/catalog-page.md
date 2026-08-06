# Kubernetes Istio

Installs the Istio service-mesh control plane from the official Helm charts, with a typed spec over the charts' meaningful configuration surface and the data plane mode (sidecar or ambient) as a first-class choice. One installation per cluster; gateways and mesh policy compose as separate resources against its outputs.

## What Gets Created

- **Namespace** (optional) — the control-plane namespace (`istio-system` by convention), created and owned when `createNamespace` is set
- **Istio CRDs** — the full pinned CRD bundle, applied by the module via server-side apply (co-ownable with a CRDs-only [KubernetesIstioBaseCrds](/docs/catalog/kubernetes/kubernetesistiobasecrds) install, so upgrading a CRDs-only cluster to a full mesh is a plain redeploy)
- **Helm Release `istio-base`** — the validation-webhook plumbing (CRDs excluded — module-owned above)
- **Helm Release `istiod`** — the control plane (`istiod-<revision>` when a revision is named)
- **Helm Release `istio-cni`** — the node agent, in ambient mode or when `cni.enabled` in sidecar mode
- **Helm Release `ztunnel`** — the per-node L4 proxy, in ambient mode

No gateway is deployed: istiod implements the Kubernetes Gateway API, so gateways compose from [KubernetesGateway](/docs/catalog/kubernetes/kubernetesgateway) resources with `gatewayClassName: istio` and istiod provisions their deployments automatically.

## Prerequisites

- A Kubernetes cluster (EKS, GKE, AKS, kind, or any conformant cluster) — the control plane itself calls no cloud APIs
- No existing Istio control plane on the cluster (the CRDs and validation webhooks are cluster singletons)

## Quick Start

### Minimal (sidecar mode)

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesIstio
metadata:
  name: my-istio
spec:
  namespace:
    value: istio-system
  createNamespace: true
```

Sidecar mode is the default. Enroll a namespace by labeling it `istio-injection=enabled`; pods created after that get a sidecar proxy.

### Ambient (sidecar-less)

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesIstio
metadata:
  name: ambient-mesh
spec:
  namespace:
    value: istio-system
  createNamespace: true
  dataplaneMode: ambient
  # On clusters without a cloud load balancer (kind, k3s, bare metal):
  gatewayDefaults:
    serviceType: NodePort
```

Ambient additionally installs the `istio-cni` node agent and the `ztunnel` per-node proxy. Namespaces enroll with the `istio.io/dataplane-mode=ambient` label — no sidecars, no pod restarts.

### Production (sidecar, hardened)

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesIstio
metadata:
  name: prod-mesh
spec:
  namespace:
    value: istio-system
  createNamespace: true
  dataplaneMode: sidecar
  istiod:
    autoscale:
      enabled: true
      minReplicas: 2
      maxReplicas: 5
      targetCpuUtilizationPercent: 75
    podDisruptionBudget: true
    priorityClassName: system-cluster-critical
  meshConfig:
    # Set a stable, organization-unique trust domain BEFORE production —
    # changing it later re-identifies every workload in the mesh.
    trustDomain: prod.mesh.example.internal
    # Egress lockdown: external destinations must be declared with a
    # KubernetesServiceEntry.
    outboundTrafficPolicyMode: REGISTRY_ONLY
    accessLogFile: /dev/stdout
  cni:
    enabled: true
```

## Composing Against the Mesh

### North-south: a Gateway and a route

istiod serves the `istio` GatewayClass. Creating a Gateway with that class makes istiod provision and program the gateway deployment (named `<gateway>-istio`) — no gateway install step exists or is needed:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesGateway
metadata:
  name: web-gateway
spec:
  namespace:
    value: istio-system
  gatewayClassName:
    value: istio
  listeners:
    - name: http
      port: 80
      protocol: HTTP
      allowedRoutes:
        namespaces:
          from: All
---
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesHttpRoute
metadata:
  name: web-route
spec:
  namespace:
    value: app-ns
  parentRefs:
    - name:
        value: web-gateway
      namespace: istio-system
  hostnames:
    - app.example.com
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /
      backendRefs:
        - name:
            value: web-service
          port: 8080
```

`gatewayClassName` can also be wired as a `valueFrom` reference to the mesh's `status.outputs.gateway_class_name`, so the Gateway deploys after the control plane in one infra chart. `spec.gatewayDefaults.serviceType` on the mesh sets the default Service type for every gateway provisioned this way (LoadBalancer upstream default; NodePort/ClusterIP for clusters without a cloud LB).

### Mesh policy: the typed Istio kinds

Traffic policy composes from the typed kinds, which need only the CRDs this component installs — for example strict mTLS plus an allow-list:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesPeerAuthentication
metadata:
  name: strict-mtls
spec:
  namespace:
    value: app-ns
  mtls:
    mode: STRICT
```

The mesh's `status.outputs.trust_domain` is the prefix of the SPIFFE principals that [KubernetesAuthorizationPolicy](/docs/catalog/kubernetes/kubernetesauthorizationpolicy) rules match on.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Control-plane namespace |
| `istiod_service_name` | istiod Service name (`istiod`, or `istiod-<revision>`) — the discovery address |
| `revision` | Installed control-plane revision (`default` when unnamed) |
| `gateway_class_name` | GatewayClass istiod serves (`istio`) |
| `trust_domain` | Identity root of every workload certificate |
| `dataplane_mode` | `sidecar` or `ambient` |

## Related Components

- [KubernetesGateway](/docs/catalog/kubernetes/kubernetesgateway) / [KubernetesHttpRoute](/docs/catalog/kubernetes/kuberneteshttproute) — north-south exposure through the mesh's `istio` GatewayClass
- [KubernetesDestinationRule](/docs/catalog/kubernetes/kubernetesdestinationrule), [KubernetesPeerAuthentication](/docs/catalog/kubernetes/kubernetespeerauthentication), [KubernetesAuthorizationPolicy](/docs/catalog/kubernetes/kubernetesauthorizationpolicy), [KubernetesRequestAuthentication](/docs/catalog/kubernetes/kubernetesrequestauthentication), [KubernetesServiceEntry](/docs/catalog/kubernetes/kubernetesserviceentry), [KubernetesTelemetry](/docs/catalog/kubernetes/kubernetestelemetry), [KubernetesEnvoyFilter](/docs/catalog/kubernetes/kubernetesenvoyfilter) — mesh traffic policy
- [KubernetesIstioBaseCrds](/docs/catalog/kubernetes/kubernetesistiobasecrds) — CRDs only, for clusters that use the typed Istio kinds without running a mesh
- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) — the control-plane namespace via `valueFrom` reference
