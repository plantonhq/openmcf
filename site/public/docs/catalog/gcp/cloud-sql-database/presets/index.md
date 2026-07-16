---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud SQL Database"
type: "preset-list"
componentSlug: "cloud-sql-database"
componentTitle: "Cloud SQL Database"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-postgres-app-database"
    rank: "01"
    title: "PostgreSQL Application Database"
    excerpt: "This preset creates one application database on an existing PostgreSQL instance, referencing the instance by name. PostgreSQL databases are created as UTF8; the engine defaults handle collation, so..."
  - slug: "02-mysql-utf8mb4-database"
    rank: "02"
    title: "MySQL Database (utf8mb4)"
    excerpt: "This preset creates a MySQL application database with the modern `utf8mb4` character set — full 4-byte UTF-8 that stores emoji and astral-plane characters correctly, unlike MySQL's legacy 3-byte..."
---

# Cloud SQL Database Presets

Ready-to-deploy configuration presets for Cloud SQL Database. Each preset is a complete manifest you can copy, customize, and deploy.
