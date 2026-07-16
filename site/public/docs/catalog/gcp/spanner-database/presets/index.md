---
title: "Presets"
description: "Ready-to-deploy configuration presets for Spanner Database"
type: "preset-list"
componentSlug: "spanner-database"
componentTitle: "Spanner Database"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-basic-database"
    rank: "01"
    title: "Basic Database"
    excerpt: "Provisions a Spanner database on an existing instance with Google's SQL dialect and a 24-hour point-in-time recovery window. The database is named after `metadata.name`; schema management belongs to..."
  - slug: "02-postgresql-database"
    rank: "02"
    title: "PostgreSQL-Dialect Database"
    excerpt: "Provisions a Spanner database with the PostgreSQL interface — PostgreSQL syntax and tooling on Spanner's globally distributed, strongly consistent storage. The dialect choice is permanent."
  - slug: "03-cmek-encrypted"
    rank: "03"
    title: "CMEK-Encrypted Database"
    excerpt: "Provisions a compliance-grade Spanner database: customer-managed encryption (CMEK) by reference to a `GcpKmsKey`, GCP API-side drop protection, point-in-time recovery, and an explicit UTC time zone."
---

# Spanner Database Presets

Ready-to-deploy configuration presets for Spanner Database. Each preset is a complete manifest you can copy, customize, and deploy.
