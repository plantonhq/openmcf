---
title: "Presets"
description: "Ready-to-deploy configuration presets for DocumentDB"
type: "preset-list"
componentSlug: "documentdb"
componentTitle: "DocumentDB"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-production-managed-password"
    rank: "01"
    title: "Production DocumentDB (Managed Password)"
    excerpt: "This preset creates a production-shaped DocumentDB cluster: one writer and one reader instance on shared cluster storage, an AWS-managed master password in Secrets Manager, encrypted storage,..."
  - slug: "02-serverless"
    rank: "02"
    title: "DocumentDB Serverless"
    excerpt: "This preset creates a DocumentDB Serverless cluster: a single `db.serverless` instance that scales between 0.5 and 16 DCUs as demand moves, an AWS-managed master password in Secrets Manager,..."
---

# DocumentDB Presets

Ready-to-deploy configuration presets for DocumentDB. Each preset is a complete manifest you can copy, customize, and deploy.
