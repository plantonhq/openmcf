---
title: "Presets"
description: "Ready-to-deploy configuration presets for Kube Prometheus Stack"
type: "preset-list"
componentSlug: "kube-prometheus-stack"
componentTitle: "Kube Prometheus Stack"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-single-cluster"
    rank: "01"
    title: "Dev single cluster preset"
    excerpt: "The smallest honest kube-prometheus-stack for a kind cluster, a laptop lab, or a private single-node environment: small Prometheus and Alertmanager volumes, short retention, and the control-plane..."
  - slug: "02-managed-cloud"
    rank: "02"
    title: "Managed cloud preset"
    excerpt: "The production-shaped install for EKS, GKE, AKS and peers: a 50Gi Prometheus volume with 15-day retention, a durable Alertmanager volume, persistence on the bundled Grafana, and every control-plane..."
  - slug: "03-production-sized"
    rank: "03"
    title: "Production sized preset"
    excerpt: "The full production posture: an HA Prometheus pair (each replica scrapes and stores the complete target set — duplication, not sharding), a quorum-safe three-replica Alertmanager gossip cluster,..."
---

# Kube Prometheus Stack Presets

Ready-to-deploy configuration presets for Kube Prometheus Stack. Each preset is a complete manifest you can copy, customize, and deploy.
