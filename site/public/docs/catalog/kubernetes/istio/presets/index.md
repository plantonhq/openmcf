---
title: "Presets"
description: "Ready-to-deploy configuration presets for Istio"
type: "preset-list"
componentSlug: "istio"
componentTitle: "Istio"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard Istio Mesh (Sidecar)"
    excerpt: "This preset installs the Istio control plane with the classic sidecar data plane and chart-default sizing. istiod autoscales via the chart's HPA (min 1 / max 5 at 80% CPU)."
  - slug: "02-ambient"
    rank: "02"
    title: "Ambient Mesh (Sidecar-less)"
    excerpt: "This preset installs the Istio control plane in ambient mode: no sidecars — a per-node ztunnel DaemonSet carries mTLS and L4 policy, and optional waypoint proxies add L7 where needed. The istio-cni..."
  - slug: "03-production-sidecar"
    rank: "03"
    title: "Production Sidecar Mesh"
    excerpt: "This preset installs a hardened sidecar-mode control plane: highly available istiod, an organization-unique trust domain, egress lockdown, mesh-wide access logging, and the node-level CNI agent in..."
---

# Istio Presets

Ready-to-deploy configuration presets for Istio. Each preset is a complete manifest you can copy, customize, and deploy.
