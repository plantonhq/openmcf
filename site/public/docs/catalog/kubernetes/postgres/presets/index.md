---
title: "Presets"
description: "Ready-to-deploy configuration presets for Postgres"
type: "preset-list"
componentSlug: "postgres"
componentTitle: "Postgres"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-single-instance"
    rank: "01"
    title: "Dev Single Instance"
    excerpt: "This preset declares the smallest useful PostgreSQL cluster: one instance, small storage, a fresh `app` database with an operator-generated password, no backups. It is a single point of failure by..."
  - slug: "02-production-ha"
    rank: "02"
    title: "Production HA"
    excerpt: "This preset declares the production PostgreSQL posture on EKS: three instances with quorum synchronous replication (zero data loss on failover), a dedicated WAL volume, hard anti-affinity, and..."
  - slug: "03-s3-compatible-backups"
    rank: "03"
    title: "S3-Compatible Backups"
    excerpt: "This preset declares a highly available PostgreSQL cluster whose backups land in an S3-COMPATIBLE object store — in-cluster MinIO, Cloudflare R2, Ceph RGW, DigitalOcean Spaces, anything speaking the..."
---

# Postgres Presets

Ready-to-deploy configuration presets for Postgres. Each preset is a complete manifest you can copy, customize, and deploy.
