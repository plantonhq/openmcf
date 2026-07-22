---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cilium"
type: "preset-list"
componentSlug: "cilium"
componentTitle: "Cilium"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-kind-dev-cluster"
    rank: "01"
    title: "Kind / Local Dev Cluster"
    excerpt: "This preset installs Cilium as the primary CNI on a local development cluster (kind, or any single-node cluster created without a default CNI). It uses `kubernetes` IPAM — the control plane already..."
  - slug: "02-eks-chaining"
    rank: "02"
    title: "EKS Chaining (on top of the AWS VPC CNI)"
    excerpt: "This preset runs Cilium ON TOP of the AWS VPC CNI on an EKS cluster — CNI chaining, the no-rip-and-replace path. The AWS VPC CNI keeps doing what EKS depends on (IPAM from the VPC, ENI wiring,..."
  - slug: "03-self-managed-primary-kpr"
    rank: "03"
    title: "Self-Managed Primary CNI with Kube-Proxy Replacement"
    excerpt: "This preset makes Cilium the cluster's primary CNI on a self-managed (kubeadm-style) cluster AND replaces kube-proxy entirely with Cilium's eBPF service load-balancing. The operator hands out pod IPs..."
  - slug: "04-production-observability"
    rank: "04"
    title: "Production Observability (Hubble + Prometheus + WireGuard)"
    excerpt: "This preset is the production hardening layer for a Cilium installation: the full Hubble stack (relay, UI, and the core metric families), Prometheus telemetry with ServiceMonitors, transparent..."
---

# Cilium Presets

Ready-to-deploy configuration presets for Cilium. Each preset is a complete manifest you can copy, customize, and deploy.
