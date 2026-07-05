---
title: "Presets"
description: "Ready-to-deploy configuration presets for Storage Table"
type: "preset-list"
componentSlug: "storage-table"
componentTitle: "Storage Table"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-app-entities"
    rank: "01"
    title: "Application Entities Table"
    excerpt: "This preset creates a plain entities table -- the serverless key/value store applications reach for when they need cheap, huge, schemaless storage addressed by partition key + row key."
  - slug: "02-audit-trail"
    rank: "02"
    title: "Audit Trail Table"
    excerpt: "This preset creates a table for append-heavy audit/event records -- the workload Table storage is economically unbeatable for: petabytes of small entities, written once, point-read or range-scanned..."
  - slug: "03-policy-anchored-access"
    rank: "03"
    title: "Policy-Anchored Read-Only Table"
    excerpt: "This preset creates a table whose external access rides a read-only stored access policy. SAS tokens issued against the policy inherit its window and query-only permissions -- and revoking the policy..."
---

# Storage Table Presets

Ready-to-deploy configuration presets for Storage Table. Each preset is a complete manifest you can copy, customize, and deploy.
