---
title: "KEDA"
description: "KEDA deployment documentation"
icon: "package"
order: 100
componentName: "kuberneteskeda"
---

# Kubernetes KEDA

Installs KEDA — Kubernetes Event-Driven Autoscaling — from the official
Helm chart, with a typed spec over the chart's meaningful configuration
surface. KEDA scales workloads on real-world signals (queue depth, stream
lag, database rows, cron schedules, cloud metrics — 70+ scalers) instead of
only CPU/memory: its operator watches ScaledObject/ScaledJob resources,
drives the workload's HPA including scale-to-ZERO, and serves the
`external.metrics.k8s.io` API the HPA controller reads. One installation
per cluster.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and owned
  when `create_namespace` is set (`keda` is the upstream convention)
- **Helm Release** — the `keda-operator` Deployment, the
  `keda-operator-metrics-apiserver` Deployment (registers the cluster-wide
  `v1beta1.external.metrics.k8s.io` APIService), the
  `keda-admission-webhooks` Deployment, RBAC, and the KEDA CRDs
  (ScaledObject, ScaledJob, TriggerAuthentication, ...) — annotated to
  survive uninstall by default, so removing the release does not
  cascade-delete every scaling declaration in the cluster

## Prerequisites

- A cluster WITHOUT an existing external-metrics provider — the
  `v1beta1.external.metrics.k8s.io` APIService is a singleton and
  Kubernetes allows only one
- With `certificates.type: cert_manager`: cert-manager on the cluster
  (KubernetesCertManager)
- With `prometheus.service_monitor`: the Prometheus operator CRDs — the
  release fails to install without them
- Cloud identity arms (IRSA, Azure Workload Identity, GCP Workload
  Identity): the cloud-side trust/binding written against the
  `keda-operator` service account

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesKeda
metadata:
  name: keda
spec:
  namespace:
    value: keda
  createNamespace: true
```

The engine then watches ScaledObjects in all namespaces. Deploy
ScaledObject/ScaledJob/TriggerAuthentication resources alongside the
workloads they scale — this component installs the engine, not the
declarations.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Installation namespace |
| `release_name` | Helm release name (always `keda`) |
| `operator_service_account_name` | Always `keda-operator` — the subject cloud-side keyless bindings are written against |

## Next Steps

Create ScaledObject resources next to the workloads they scale — a cron
trigger is the simplest deterministic proof the engine works. For scalers
that read cloud metric sources without stored keys, enable the matching
`pod_identity` arm (IRSA / Azure Workload Identity / GCP Workload
Identity). When autoscaling becomes load-bearing, harden with standby
replicas, `system-cluster-critical` priority, cert-manager internal TLS,
and Prometheus telemetry. For CPU/memory-only scaling, plain
HorizontalPodAutoscaler with metrics-server is the lighter tool.
