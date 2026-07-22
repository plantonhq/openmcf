# KubernetesMetricsServer

Installs metrics-server — the cluster's resource-metrics pipeline — from the
official Helm chart, with a typed spec over the chart's meaningful
configuration surface. It scrapes kubelets and serves live CPU/memory usage
through the `metrics.k8s.io` API: what `kubectl top` reads and what
HorizontalPodAutoscalers need to scale. One installation per cluster.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and owned
  when `create_namespace` is set (`kube-system` installs leave it false)
- **Helm Release** — the metrics-server Deployment, Service, RBAC, and the
  cluster-wide `v1beta1.metrics.k8s.io` APIService registration; the release
  waits for readiness, and readiness requires the first kubelet scrape to
  succeed

## Prerequisites

- A Kubernetes cluster WITHOUT a metrics API already in place — EKS, AKS,
  kind, k3s, kubeadm, or self-managed (GKE ships metrics-server built-in;
  do not install this component there)
- On clusters whose kubelets serve self-signed certificates (kind, k3s,
  kubeadm without kubelet certificate rotation, many on-prem setups): set
  `kubelet_insecure_tls: true`, or the install fails its readiness wait

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesMetricsServer
metadata:
  name: metrics-server
spec:
  namespace:
    value: kube-system
```

On a self-signed-kubelet cluster add `kubeletInsecureTls: true`.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Installation namespace |
| `release_name` | Helm release name (always `metrics-server`) |
| `service_name` | The Service the APIService routes to (always `metrics-server`) |
| `api_service_name` | `v1beta1.metrics.k8s.io`; empty when `api_service.create` is false |

## Next Steps

Create **KubernetesHorizontalPodAutoscaler** resources — they consume the
metrics API this component registers, with no explicit reference needed. For
custom or external metrics (queue depth, requests per second), or for metric
history and dashboards, reach for a monitoring stack instead — metrics-server
serves only instantaneous CPU/memory.
