---
title: "Presets"
description: "Ready-to-deploy configuration presets for Metrics Server"
type: "preset-list"
componentSlug: "metrics-server"
componentTitle: "Metrics Server"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-managed-cloud"
    rank: "01"
    title: "Managed Cloud (EKS / AKS)"
    excerpt: "This preset installs metrics-server into `kube-system` with chart defaults — the posture for managed clouds whose kubelets serve CA-signed certificates. `kubelet_insecure_tls` stays false: EKS and..."
  - slug: "02-self-signed-kubelets"
    rank: "02"
    title: "Self-Signed Kubelets (kind / k3s / kubeadm / on-prem)"
    excerpt: "This preset installs metrics-server on clusters whose kubelets serve self-signed certificates: kind, k3s, kubeadm without kubelet certificate rotation, and many on-prem setups. `kubeletInsecureTls:..."
  - slug: "03-ha-verified-tls"
    rank: "03"
    title: "HA with Verified TLS (production hardening)"
    excerpt: "This preset hardens both availability and trust: two replicas guarded by a PodDisruptionBudget keep the metrics API serving through node drains, and cert-manager issues (and renews) the serving..."
---

# Metrics Server Presets

Ready-to-deploy configuration presets for Metrics Server. Each preset is a complete manifest you can copy, customize, and deploy.
