---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cosmos DB SQL Role Definition"
type: "preset-list"
componentSlug: "cosmos-db-sql-role-definition"
componentTitle: "Cosmos DB SQL Role Definition"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-read-only-role"
    rank: "01"
    title: "Read-Only Role"
    excerpt: "This preset creates a custom read-only role covering the full read surface -- point reads, SQL queries, and the change feed -- plus the metadata access every Cosmos SDK needs before its first data..."
  - slug: "02-writer-no-delete"
    rank: "02"
    title: "Writer Without Delete"
    excerpt: "This preset creates the role the built-ins cannot express: full read access plus create/replace/upsert -- but never delete. Ingest pipelines, event processors, and application workloads write..."
  - slug: "03-database-scoped-role"
    rank: "03"
    title: "Database-Scoped Role"
    excerpt: "This preset narrows WHERE a role may ever be granted, not just what it allows: its single assignable scope is one database's path, so an assignment of this role at the account level -- or in any..."
---

# Cosmos DB SQL Role Definition Presets

Ready-to-deploy configuration presets for Cosmos DB SQL Role Definition. Each preset is a complete manifest you can copy, customize, and deploy.
