---
title: "Cilium"
description: "Cilium deployment documentation"
icon: "package"
order: 100
componentName: "kubernetescilium"
---

# Cilium on Kubernetes

Installs Cilium -- the eBPF-based networking, network-security, and observability engine -- from the official Helm chart (`cilium` at https://helm.cilium.io). Cilium is the cluster's CNI: it wires pod networking, enforces NetworkPolicy (standard policies plus Cilium's own L7-aware ones), can replace kube-proxy entirely with eBPF service load-balancing, and streams flow-level observability through Hubble. The typed configuration covers the chart's meaningful surface; a merged-last `helm_values` document remains as the escape hatch for the long tail.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **The Cilium Helm release** (release name fixed to `cilium` -- ONE installation per cluster) in `spec.namespace`: the agent DaemonSet on every node, the cilium-operator Deployment, and the generated CNI configuration.
- **Hubble components** when enabled: the relay (cluster-wide flow aggregation) and the UI (service-map console).
- **The `cilium` GatewayClass** when Gateway API support is enabled.
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- **Primary CNI mode**: the cluster must be created WITHOUT another CNI (kind: `disableDefaultCNI: true`; AKS: BYOCNI via `--network-plugin none`; kubeadm: no CNI addon). Nodes stay NotReady until Cilium installs -- that is by design.
- **Chaining mode**: the incumbent CNI stays; nothing special is required beyond choosing the right chaining mode (e.g. `aws-cni` on EKS).
- **Kube-proxy replacement**: create the cluster without kube-proxy (kind: `kubeProxyMode: none`; kubeadm: `--skip-phases=addon/kube-proxy`) and have the API server address at hand.
- **ServiceMonitors** (Hubble metrics / telemetry): the Prometheus operator CRDs must exist -- the release fails to install without them.
- **Gateway API support**: the Gateway API CRDs (deploy `KubernetesGatewayApiCrds` first) AND kube-proxy replacement.

## Deploy

### Console

Open the deployment store, find **Cilium on Kubernetes**, and click **Deploy**. The creation wizard walks the platform engineer's decision sequence: **Namespace** (kube-system convention), **Installation** (pinned chart version, cluster identity), **CNI Mode** (primary vs chaining -- the load-bearing choice), **Kube-Proxy** (eBPF replacement + the API-server address it requires + Gateway API), **IPAM & Routing**, **Cloud Integration** (at most one arm), **Security** (enforcement posture, transparent encryption), **Observability** (Hubble, Prometheus), **Performance**, **Sizing**, and the **Helm Values** escape hatch. Start from the **Kind / Local Dev Cluster** or **EKS Chaining** preset in the [Presets](#presets) tab for a directly deployable configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesCilium
metadata:
  name: cilium
  org: acme-corp
  env: dev
spec:
  namespace:
    value: kube-system
  chartVersion: 1.19.6
  ipam:
    mode: kubernetes
  kubeProxyReplacement: true
  k8sServiceHost: kind-control-plane
  k8sServicePort: 6443
  hubble:
    relay: true
    ui: true
```

```shell
planton apply -f cilium.yaml
```

This installs Cilium as the primary CNI on a kind cluster (created with `disableDefaultCNI: true` and `kubeProxyMode: none`), with eBPF service load-balancing and the Hubble service map.

## Key Configuration

These are the most important decisions when configuring Cilium. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**CNI mode** -- The load-bearing choice. As the PRIMARY CNI, Cilium owns pod networking. With CHAINING (`cni.chaining_mode`: `aws-cni`, `flannel`, `generic-veth`, `portmap`), the incumbent CNI keeps IPAM and routing while Cilium adds eBPF policy, load-balancing, and observability. Chaining forces `cni.exclusive` to false -- the chained CNI's configuration must survive.

**Kube-proxy replacement** -- eBPF service load-balancing instead of kube-proxy. Requires `k8s_service_host` (the agent must reach the API server before any service load-balancing exists) and unlocks Gateway API support (`gateway_api` creates the `cilium` GatewayClass).

**IPAM & routing** -- `ipam.mode` mirrors the agent's own vocabulary (`cluster-pool` chart default, `kubernetes` for kind/kubeadm, `eni`/`azure` for cloud addressing, `delegated-plugin` for chaining). `routing.mode` is `tunnel` (works anywhere) or `native` (faster, needs a routable fabric plus the native-routing CIDR).

**Cloud integration** -- At most one arm: `aws_eni` (primary CNI on EKS/EC2), `aks_byocni`, or `gke`. None is correct for kind, self-managed, and all chaining setups.

**Security** -- `policy_enforcement_mode` (`default` / `always` / `never`); transparent encryption via WireGuard (automatic keys, node-kernel module required) or IPsec (pre-created key Secret, rotation policies).

**Observability** -- Hubble is enabled by default upstream; the relay aggregates flows cluster-wide, the UI serves the service map, `metrics` exports flow-metric families. Cilium's own agent/operator telemetry is separate (`prometheus`).

**Helm values** -- The escape hatch, merged LAST with `-f` semantics on both engines: BGP control plane, L2 announcements, Cluster Mesh, per-image overrides. Never the substitute for the typed fields, and never a place for secrets.

## Outputs and Dependencies

### What This Component Consumes

`spec.namespace` is a foreign-key reference to a `KubernetesNamespace` (or the literal `kube-system`). Gateway API support additionally expects the Gateway API CRDs on the cluster.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace Cilium was installed into | Co-locating companion resources |
| `release_name` | Helm release name (fixed `cilium`) | Operational tooling and runbooks |
| `cluster_name` | Cluster identity in Hubble flows / Cluster Mesh | Multi-cluster naming |
| `hubble_relay_service_name` | The hubble-relay Service (when relay is enabled) | hubble CLI / UI flow access |
| `hubble_ui_service_name` | The hubble-ui Service (when UI is enabled) | Port-forward to open the service map |
| `gateway_class_name` | The GatewayClass Cilium registers (when Gateway API is on) | `KubernetesGateway.gateway_class_name` |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Local development** -- Primary CNI on kind with kube-proxy replacement and the Hubble UI. Start from the **Kind / Local Dev Cluster** preset.

**EKS without rip-and-replace** -- Chain onto the AWS VPC CNI for policy and observability while AWS keeps addressing. Start from the **EKS Chaining** preset.

**Self-managed production** -- Primary CNI with kube-proxy replacement on kubeadm-class clusters. Start from the **Self-Managed Primary CNI with Kube-Proxy Replacement** preset.

**Production observability** -- Hubble relay/UI/metrics, Prometheus ServiceMonitors, and WireGuard encryption. Start from the **Production Observability** preset.

## Works With

- **KubernetesGatewayApiCrds** -- prerequisite for Gateway API support; Cilium then registers the `cilium` GatewayClass.
- **KubernetesGateway** -- consumes the exported GatewayClass for north-south routing.
- **KubernetesNetworkPolicy** -- the policies Cilium enforces under the chosen enforcement posture.
- **KubernetesNamespace** -- the namespace (`spec.namespace`) the release installs into.
- **KubernetesKubePrometheusStack** (or an existing Prometheus operator) -- required for the ServiceMonitor toggles; scrapes the agent, operator, and Hubble metrics.
