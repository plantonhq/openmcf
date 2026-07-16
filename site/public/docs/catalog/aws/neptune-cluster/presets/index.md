---
title: "Presets"
description: "Ready-to-deploy configuration presets for Neptune Cluster"
type: "preset-list"
componentSlug: "neptune-cluster"
componentTitle: "Neptune Cluster"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-production-graph"
    rank: "01"
    title: "Production Neptune (Provisioned)"
    excerpt: "This preset creates a production-shaped Neptune cluster: one writer and one reader instance on shared cluster storage, IAM database authentication, encrypted storage, deletion protection, seven days..."
  - slug: "02-serverless"
    rank: "02"
    title: "Neptune Serverless"
    excerpt: "This preset creates a Neptune Serverless cluster: a single `db.serverless` instance that scales between 1 and 32 NCUs as traversal load moves, IAM database authentication, encrypted storage, and..."
---

# Neptune Cluster Presets

Ready-to-deploy configuration presets for Neptune Cluster. Each preset is a complete manifest you can copy, customize, and deploy.
