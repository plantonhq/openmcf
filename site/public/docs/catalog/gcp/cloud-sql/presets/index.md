---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud SQL"
type: "preset-list"
componentSlug: "cloud-sql"
componentTitle: "Cloud SQL"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-postgres-production-private"
    rank: "01"
    title: "Production PostgreSQL (Private IP)"
    excerpt: "This preset deploys a production-grade PostgreSQL 16 instance reachable only over private IP inside your VPC — no public address at all — with regional high availability, point-in-time recovery,..."
  - slug: "02-mysql-high-availability"
    rank: "02"
    title: "High-Availability MySQL (Auth Proxy Access)"
    excerpt: "This preset deploys a MySQL 8.0 instance with regional high availability and binary logging, exposed through a public IP that has **zero authorized networks** — so the only way in is the..."
  - slug: "03-postgres-read-replica"
    rank: "03"
    title: "PostgreSQL Read Replica"
    excerpt: "This preset attaches a read replica to an existing PostgreSQL primary. A replica is a full `GcpCloudSql` node of its own — same kind, own manifest — that references its primary through..."
---

# Cloud SQL Presets

Ready-to-deploy configuration presets for Cloud SQL. Each preset is a complete manifest you can copy, customize, and deploy.
