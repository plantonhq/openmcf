---
title: "Presets"
description: "Ready-to-deploy configuration presets for OpenBao"
type: "preset-list"
componentSlug: "openbao"
componentTitle: "OpenBao"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-mode"
    rank: "01"
    title: "Dev-mode preset"
    excerpt: "OpenBao with zero ceremony: dev mode auto-initializes and auto-unseals at startup, the root token is literally `root`, and all data lives in memory. Port-forward the `openbao-dev` Service and start..."
  - slug: "02-production-ha"
    rank: "02"
    title: "Production HA (integrated Raft) preset"
    excerpt: "Three OpenBao servers with integrated Raft storage: each replica persists to its own 10Gi PVC, the module synthesizes the `retry_join` stanzas for every peer (the chart alone ships none — without..."
  - slug: "03-production-ha-gcp-auto-unseal"
    rank: "03"
    title: "Production HA + GCP Cloud KMS auto-unseal preset"
    excerpt: "The production-ha shape with the restart toil removed: the master key is wrapped by a Cloud KMS crypto key, so every server unseals ITSELF at startup — pod restarts, node replacements and scale..."
---

# OpenBao Presets

Ready-to-deploy configuration presets for OpenBao. Each preset is a complete manifest you can copy, customize, and deploy.
