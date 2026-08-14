---
title: "Presets"
description: "Ready-to-deploy configuration presets for RDS Instance"
type: "preset-list"
componentSlug: "rds-instance"
componentTitle: "RDS Instance"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-postgresql-production"
    rank: "01"
    title: "PostgreSQL (Production Multi-AZ)"
    excerpt: "This preset creates a production-shaped PostgreSQL instance: Multi-AZ with a synchronous standby and automatic failover, gp3 storage with autoscaling headroom, an AWS-managed master password in..."
  - slug: "02-mysql-production"
    rank: "02"
    title: "MySQL (Production Multi-AZ)"
    excerpt: "This preset creates a production-shaped MySQL instance: Multi-AZ with a synchronous standby, gp3 storage with autoscaling headroom, an AWS-managed master password, encrypted storage, deletion..."
  - slug: "03-read-replica"
    rank: "03"
    title: "Read Replica"
    excerpt: "This preset creates a read replica of an existing RDS instance. The replica inherits engine, version, storage, and credentials from its source -- the manifest carries only what is genuinely the..."
  - slug: "04-oracle-enterprise"
    rank: "04"
    title: "Oracle Enterprise Edition (BYOL)"
    excerpt: "This preset creates a Multi-AZ Oracle EE instance under bring-your-own-license, with the two engine-configuration surfaces the instance owns itself: a parameter group built from inline parameters and..."
  - slug: "05-mysql-s3-migration"
    rank: "05"
    title: "MySQL Migration via S3 Import"
    excerpt: "This preset creates a new Multi-AZ RDS MySQL instance by restoring a Percona XtraBackup stored in S3 -- no logical dump, no replication bootstrap. The instance lands with an AWS-managed master..."
---

# RDS Instance Presets

Ready-to-deploy configuration presets for RDS Instance. Each preset is a complete manifest you can copy, customize, and deploy.
