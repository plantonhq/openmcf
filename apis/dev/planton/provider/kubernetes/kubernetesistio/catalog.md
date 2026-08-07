# Istio on Kubernetes

Deploys the Istio service mesh control plane on a Kubernetes cluster: istiod, the Istio CRDs, and -- in ambient mode -- the istio-cni node agent and ztunnel node proxies. Control-plane-only by design: no ingress gateway ships with it; gateways are provisioned on demand from Gateway API Gateways that use the `istio` GatewayClass. One version pin drives every chart, and per-release Helm value escape hatches cover the long tail.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Istio CRDs** -- applied by the module itself via server-side apply, outside any Helm release. Module-owned CRDs are co-ownable with the CRDs-only Istio Base CRDs component, so a cluster running just the CRDs upgrades to the full mesh with a plain redeploy.
- **base Helm Release** -- the validation-webhook plumbing and cluster-wide resources istiod requires.
- **istiod Helm Release** -- the control plane, with your capacity, autoscaling, mesh configuration, proxy defaults, and scheduling choices rendered as chart values.
- **cni Helm Release** -- the istio-cni node agent. Always installed in ambient mode (it is how traffic reaches ztunnel); opt-in in sidecar mode, where it replaces the injected privileged init-container.
- **ztunnel Helm Release** -- the ambient L4 node proxy. Installed in ambient mode only.
- **Kubernetes Namespace** -- created only when `createNamespace` is `true`; hosts everything above.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- One KubernetesIstio per cluster: the CRDs and the validation-webhook plumbing are cluster singletons.
- In sidecar mode, namespaces that should participate in the mesh need the `istio-injection: enabled` label (or workloads opt in individually when the injection policy is `disabled`).
- Auto-provisioned gateways default to LoadBalancer Services -- on clusters without a cloud load-balancer controller (bare metal, kind, k3s), set the gateway service type to NodePort or ClusterIP.

## Deploy

### Console

Open the deployment store, find **Istio on Kubernetes**, and click **Deploy**. The wizard walks you through placement, the release pin, the sidecar-vs-ambient decision, control-plane capacity, mesh identity, proxy defaults, the node agents, gateway defaults, and the Helm escape hatches -- with guidance at each step.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesIstio
metadata:
  name: mesh
  org: acme-corp
  env: prod
spec:
  namespace:
    value: istio-system
  createNamespace: true
  version: 1.30.3
  dataplaneMode: sidecar
  meshConfig:
    trustDomain: prod.acme.internal
  istiod:
    priorityClassName: system-cluster-critical
```

```shell
planton apply -f istio.yaml
```

This installs the control plane into `istio-system` with a production trust domain and eviction protection for istiod; everything else runs on the chart defaults.

## Key Configuration

- **Version** -- one chart version for every release (Istio versions its charts in lockstep). Upgrades step ONE minor at a time -- skipping minors in place is unsupported upstream.
- **Data Plane Mode** -- `sidecar` (a proxy per pod, the classic mode) or `ambient` (sidecar-less: per-node ztunnel for L4/mTLS, optional waypoints for L7). An install-time decision: switching later is a workload-by-workload migration.
- **Control Plane (istiod)** -- fixed replicas or the chart-managed HPA (mutually exclusive), resources, log level, PodDisruptionBudget, priority class, and node placement.
- **Mesh Configuration** -- the trust domain (set a stable, organization-unique value BEFORE production -- changing it re-identifies every workload), the outbound traffic policy (`REGISTRY_ONLY` is the egress-lockdown posture), access logging, and the multi-cluster naming trio.
- **Proxies & Injection** -- per-proxy resources (they multiply across every injected pod), log level, the injection policy, and probe rewriting for strict mTLS.
- **CNI Node Agent / ztunnel** -- node-level traffic redirection and the ambient node proxy, each with chart-default tuning.
- **Gateways & Images** -- the Service type for auto-provisioned gateways, plus the image hub/variant/pull-secret knobs for air-gapped and hardened environments.
- **Helm Values** -- per-release YAML documents merged LAST over the typed fields, for chart surface beyond them.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Description |
|------------|-------------|
| `namespace` | The namespace the control plane installs into. Reference an existing Namespace on Kubernetes or supply the name directly. |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources and operators can reference:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Where the control plane runs | Placing mesh-adjacent tooling beside istiod |
| `istiod_service_name` | The discovery Service proxies connect to | Debugging and custom bootstrap configuration |
| `revision` | The control-plane revision (empty = default) | Labeling workloads `istio.io/rev` under revisioned operation |
| `gateway_class_name` | The GatewayClass this control plane serves | Gateway API Gateways reference it to be provisioned by this mesh |
| `trust_domain` | The identity root (`spiffe://<trust_domain>/...`) | Writing authorization policies that match workload identities |
| `dataplane_mode` | `sidecar` or `ambient` | Runbooks and tooling that branch on the mesh's mode |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

- **Classic sidecar mesh** -- the widely documented default: label namespaces for injection and get a proxy per pod.
- **Ambient mesh** -- sidecar-less L4/mTLS everywhere with waypoints added only where L7 policy is needed -- lower overhead, no workload restarts on proxy upgrades.
- **Egress lockdown** -- `REGISTRY_ONLY` outbound policy with every external service declared as a Service Entry.

## Works With

Istio is the control plane the typed Istio resources configure: **DestinationRule** and **Service Entry** (traffic management), **Peer Authentication**, **Request Authentication**, and **Authorization Policy** (security), **Telemetry** (observability), and **EnvoyFilter** (extensibility). Gateway API **Gateways** using the `istio` GatewayClass are provisioned by this control plane, with **HTTPRoutes** and their siblings attaching to them. A cluster that only needs the API types (no mesh runtime) can install the **Istio Base CRDs** component instead and upgrade to this full mesh later with a plain redeploy.
