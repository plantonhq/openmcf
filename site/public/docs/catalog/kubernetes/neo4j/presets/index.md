---
title: "Presets"
description: "Ready-to-deploy configuration presets for Neo4j"
type: "preset-list"
componentSlug: "neo4j"
componentTitle: "Neo4j"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-single-instance"
    rank: "01"
    title: "Dev single instance preset"
    excerpt: "The smallest declarable Neo4j: community edition, a declared admin password (materialized as a Kubernetes Secret, never in rendered values), a 10Gi data volume on the cluster's default StorageClass,..."
  - slug: "02-production"
    rank: "02"
    title: "Production preset"
    excerpt: "A production single server: credentials referenced from a pre-existing Secret (never declared in the manifest), an explicit heap/page-cache memory split sized to the container, a 100Gi data volume on..."
  - slug: "03-graphrag-apoc"
    rank: "03"
    title: "GraphRAG APOC preset"
    excerpt: "A Neo4j tuned for knowledge-graph and agent-memory workloads: the APOC procedure library activated at startup (it ships inside the official image — no download), apoc.conf arms enabled for the..."
---

# Neo4j Presets

Ready-to-deploy configuration presets for Neo4j. Each preset is a complete manifest you can copy, customize, and deploy.
