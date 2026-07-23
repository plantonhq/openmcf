---
title: "Presets"
description: "Ready-to-deploy configuration presets for MySQL"
type: "preset-list"
componentSlug: "mysql"
componentTitle: "MySQL"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "02-production-ha"
    rank: "02"
    title: "Production HA preset"
    excerpt: "Three-node Galera cluster (quorum-safe synchronous replication) with three HAProxy replicas, nightly XtraBackup to S3, dedicated PITR storage, zone-spread anti-affinity, and cert-manager TLS via a..."
---

# MySQL Presets

Ready-to-deploy configuration presets for MySQL. Each preset is a complete manifest you can copy, customize, and deploy.
