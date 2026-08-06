---
title: "KubernetesIngressNginx"
description: "KubernetesIngressNginx deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesingressnginx"
---

# KubernetesIngressNginx

Installs the ingress-nginx controller — the cluster's HTTP(S) entry point —
from the official Helm chart, with a typed spec over the chart's meaningful
configuration surface. Multiple controllers per cluster (public + internal
traffic splits) are first-class: each instance owns its own IngressClass,
and everything from release naming to leader election derives from
`metadata.name`, so instances never collide.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and owned
  when `create_namespace` is set
- **Helm Release** — the ingress-nginx controller (Deployment or DaemonSet,
  Services, IngressClass, RBAC, admission webhook, and optionally the
  default backend), named after `metadata.name` so multiple instances
  coexist per cluster
- **Cloud Load Balancer** (indirect) — with the default `load_balancer`
  service type, the host cloud provisions the entry LB, shaped entirely by
  `service.annotations`

## Prerequisites

- A Kubernetes cluster (EKS, GKE, AKS, or any conformant cluster; on
  clusters without a cloud LB controller — kind, bare metal — use
  `node_port` or host access instead of `load_balancer`)
- For `autoscaling`: metrics-server on the cluster
- For `metrics.service_monitor`: the Prometheus operator CRDs (the release
  fails to install without them)

## Common Postures

- **AWS NLB public entry** — `load_balancer` + the NLB annotations
  (`aws-load-balancer-type: external`, `nlb-target-type: ip`),
  `external_traffic_policy: local` to preserve client source IPs
- **Internal-only controller** — a second instance with its own class
  (e.g. `nginx-internal`) and the cloud's internal-LB annotation
- **Single controller, dual LB** — `service.internal.enabled` adds the
  chart's second, internal Service in front of the same pods (requires at
  least one annotation — the internal-LB annotation is what makes it
  internal)
- **Bare metal / edge** — `controller_kind: daemon_set` with `host_ports`
  (or `host_network`) and a `node_port` service
- **Cluster-wide default TLS** — `default_tls_certificate` referencing a
  KubernetesCertificate's secret output (the cert-manager seam)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesIngressNginx
metadata:
  name: ingress-nginx
spec:
  namespace:
    value: ingress-nginx
  createNamespace: true
  ingressClass:
    name: nginx
    isDefaultClass: true
  replicas: 2
  service:
    externalTrafficPolicy: local
```

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Installation namespace |
| `release_name` | Helm release name (equals `metadata.name`) |
| `ingress_class_name` | The IngressClass this controller owns — what KubernetesIngress resources reference |
| `controller_service_name` | The controller's external Service (`<name>-controller`) |
| `internal_service_name` | The internal Service (when `service.internal.enabled`) |
| `load_balancer_ip` | External IP of the cloud LB (GCP/Azure populate an IP) |
| `load_balancer_hostname` | External hostname of the cloud LB (AWS populates a DNS name) |

## Next Steps

Create **KubernetesIngress** resources referencing this controller's
`ingress_class_name` output to route traffic. Pair with
**KubernetesExternalDns** to publish DNS records for the load balancer's
address, and with the **cert-manager** kinds (KubernetesCertManager,
issuers, KubernetesCertificate) for TLS — including the
`default_tls_certificate` wildcard default.
