# Kubernetes Cilium

Installs Cilium — the eBPF-based CNI, network-security, and observability
engine — from the official Helm chart, with a typed spec over the chart's
meaningful configuration surface. Cilium wires pod networking, enforces
NetworkPolicy (including L7-aware Cilium policies), can replace kube-proxy
with eBPF service load-balancing, and streams flow-level observability
through Hubble. Runs as the cluster's primary CNI or chained on top of an
existing one. One installation per cluster.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and owned
  when `create_namespace` is set (`kube-system` installs leave it false)
- **Helm Release** — the `cilium` agent DaemonSet, the `cilium-operator`
  Deployment, generated CNI configuration, RBAC, and — per spec toggles —
  Hubble relay/UI, ServiceMonitors, and the `cilium` GatewayClass. The
  release installs atomically and waits (600s) for the whole dataplane: on
  a fresh cluster, nodes go NotReady→Ready as Cilium wires each one

## Prerequisites

- **Primary CNI mode**: a cluster created WITHOUT another CNI (kind:
  `disableDefaultCNI: true`; AKS: `--network-plugin none` with
  `cloud.aks_byocni`; kubeadm: no CNI addon). Nodes stay NotReady until
  this component installs — by design
- **Chaining mode**: the incumbent CNI stays; set `cni.chaining_mode` and
  `cni.exclusive: false`
- **`kube_proxy_replacement`**: the cluster must be created without
  kube-proxy, and `k8s_service_host` must point at the API server
- **`gateway_api`**: requires `kube_proxy_replacement` and the Gateway API
  CRDs (`KubernetesGatewayApiCrds`)
- **ServiceMonitor toggles**: require the Prometheus operator CRDs — the
  release fails to install without them

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesCilium
metadata:
  name: cilium
spec:
  namespace:
    value: kube-system
  ipam:
    mode: kubernetes
  operator:
    replicas: 1
```

The posture above fits kind and other single-node clusters whose control
plane already assigns node PodCIDRs. On multi-node clusters keep the
2-replica operator default and choose the IPAM mode for your environment.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Installation namespace |
| `release_name` | Helm release name (always `cilium`) |
| `cluster_name` | Cluster identity Cilium runs under (Hubble flows, Cluster Mesh) |
| `hubble_relay_service_name` | `hubble-relay` when relay is enabled; empty otherwise |
| `hubble_ui_service_name` | `hubble-ui` when the UI is enabled; empty otherwise |
| `gateway_class_name` | `cilium` when `gateway_api` is enabled; empty otherwise |

## Next Steps

Apply NetworkPolicy (or Cilium's L7-aware policies) — with Cilium as the
enforcing CNI they take effect immediately. Enable Hubble relay and UI for
the live service map. With `gateway_api: true`, create
**KubernetesGateway** resources against the `cilium` GatewayClass. For the
chart surface beyond the typed fields (BGP control plane, L2 announcements,
Cluster Mesh, Cilium's ingress controller), use `helm_values`.
