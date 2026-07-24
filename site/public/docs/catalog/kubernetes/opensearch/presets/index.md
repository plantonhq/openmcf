---
title: "Presets"
description: "Ready-to-deploy configuration presets for OpenSearch"
type: "preset-list"
componentSlug: "opensearch"
componentTitle: "OpenSearch"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-minimal"
    rank: "01"
    title: "Dev minimal preset"
    excerpt: "The smallest declarable OpenSearch that actually serves: one pool, two all-roles nodes (the manager floor — a single manager-eligible replica cannot survive the operator's bootstrap handoff; the spec..."
  - slug: "02-production-cluster"
    rank: "02"
    title: "Production cluster preset"
    excerpt: "The production topology: three dedicated cluster-manager nodes (the coordination quorum, isolated from query load), three data/ingest nodes on fast persistent storage, PodDisruptionBudgets limiting..."
  - slug: "03-s3-snapshots"
    rank: "03"
    title: "S3 snapshots preset"
    excerpt: "A three-node cluster with the backup story wired end to end: the `repository-s3` plugin installed at node startup, S3 credentials loaded into the OpenSearch keystore from an existing Secret, and an..."
---

# OpenSearch Presets

Ready-to-deploy configuration presets for OpenSearch. Each preset is a complete manifest you can copy, customize, and deploy.
