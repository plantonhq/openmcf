---
title: "Presets"
description: "Ready-to-deploy configuration presets for RDS Cluster"
type: "preset-list"
componentSlug: "rds-cluster"
componentTitle: "RDS Cluster"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-aurora-postgresql"
    rank: "01"
    title: "Aurora PostgreSQL (Provisioned)"
    excerpt: "This preset creates a production-shaped Aurora PostgreSQL cluster: one writer and one reader instance on shared cluster storage, an AWS-managed master password in Secrets Manager, encrypted storage,..."
  - slug: "02-aurora-mysql"
    rank: "02"
    title: "Aurora MySQL (Provisioned)"
    excerpt: "This preset creates a production-shaped Aurora MySQL cluster: one writer and one reader on shared cluster storage, an AWS-managed master password, encrypted storage, deletion protection, seven days..."
  - slug: "03-aurora-serverless-v2"
    rank: "03"
    title: "Aurora Serverless v2 (Scale-to-Zero)"
    excerpt: "This preset creates an Aurora PostgreSQL Serverless v2 cluster: one `db.serverless` instance that scales between 0 and 16 ACUs with demand, automatic pause after five idle minutes (compute cost drops..."
---

# RDS Cluster Presets

Ready-to-deploy configuration presets for RDS Cluster. Each preset is a complete manifest you can copy, customize, and deploy.
