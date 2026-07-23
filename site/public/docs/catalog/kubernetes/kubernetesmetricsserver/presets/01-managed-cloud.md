---
title: "Managed Cloud (EKS / AKS)"
description: "This preset installs metrics-server into `kube-system` with chart defaults — the posture for managed clouds whose kubelets serve CA-signed certificates. `kubelet_insecure_tls` stays false: EKS and..."
type: "preset"
rank: "01"
presetSlug: "01-managed-cloud"
componentSlug: "kubernetesmetricsserver"
componentTitle: "KubernetesMetricsServer"
provider: "kubernetes"
icon: "package"
order: 1
---

# Managed Cloud (EKS / AKS)

This preset installs metrics-server into `kube-system` with chart defaults —
the posture for managed clouds whose kubelets serve CA-signed certificates.
`kubelet_insecure_tls` stays false: EKS and AKS kubelet certificates verify
against the cluster CA, so no trust shortcut is needed on the scrape side.

Do NOT use this (or any) metrics-server installation on GKE — GKE ships
metrics-server built-in, and the `v1beta1.metrics.k8s.io` APIService is a
cluster-wide singleton.

## When to Use

- You run EKS or AKS and HPAs never scale / `kubectl top` reports the
  metrics API unavailable
- You want the upstream-conventional install: `kube-system`, one replica,
  self-signed serving certificate with APIService verification skipped

## Key Configuration Choices

- **`namespace: kube-system`** — the upstream convention; the APIService is
  cluster infrastructure
- **`createNamespace: false`** — kube-system always exists; the module never
  tries to own it
- **Defaults otherwise** — chart 3.13.1, one replica, 15s metric
  resolution, `system-cluster-critical` priority

## Placeholders to Replace

No placeholders — this preset is directly deployable.

## Related Presets

- **02-self-signed-kubelets** — kind / k3s / kubeadm / on-prem clusters
  whose kubelets serve self-signed certificates
- **03-ha-verified-tls** — production hardening: HA replicas and a verified
  serving-certificate chain via cert-manager
