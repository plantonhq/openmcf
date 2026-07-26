---
title: "Presets"
description: "Ready-to-deploy configuration presets for Qdrant"
type: "preset-list"
componentSlug: "qdrant"
componentTitle: "Qdrant"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-single-node"
    rank: "01"
    title: "Dev single node preset"
    excerpt: "The smallest declarable Qdrant: one node, no authentication (the upstream default), a 5Gi data volume on the cluster's default StorageClass, and the chart's own resource defaults. For developers who..."
  - slug: "02-production-cluster"
    rank: "02"
    title: "Production cluster preset"
    excerpt: "A production Qdrant: 3 nodes (the quorum posture — Raft survives one member loss), generated read-write and read-only API keys living in the chart-owned Secret (never in a manifest), 8Gi of memory..."
  - slug: "03-rag-workload"
    rank: "03"
    title: "RAG workload preset"
    excerpt: "A single-node Qdrant sized for a typical RAG corpus: 8Gi of memory (the real capacity bound — Qdrant serves from RAM-resident segments and HNSW indexes; 8Gi holds an embedding set in the..."
---

# Qdrant Presets

Ready-to-deploy configuration presets for Qdrant. Each preset is a complete manifest you can copy, customize, and deploy.
