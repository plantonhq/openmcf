---
title: "Presets"
description: "Ready-to-deploy configuration presets for BigQuery Table"
type: "preset-list"
componentSlug: "bigquery-table"
componentTitle: "BigQuery Table"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-partitioned-analytics"
    rank: "01"
    title: "Preset: Partitioned Analytics Table"
    excerpt: "Use this preset for the workhorse analytics pattern: an append-only event or fact table that grows without bound. Day partitioning plus clustering keeps query cost proportional to the data actually..."
  - slug: "02-authorized-view"
    rank: "02"
    title: "Preset: Authorized View"
    excerpt: "Use this preset to expose a filtered or aggregated slice of sensitive data without granting readers any access to the raw tables. The view lives in a reader-facing dataset; the source dataset..."
  - slug: "03-external-gcs-table"
    rank: "03"
    title: "Preset: External GCS Table"
    excerpt: "Use this preset to query data-lake files in GCS without loading them into BigQuery-managed storage: parquet/CSV/JSON exports, hive-partitioned lake layouts, or data shared with Spark and other..."
---

# BigQuery Table Presets

Ready-to-deploy configuration presets for BigQuery Table. Each preset is a complete manifest you can copy, customize, and deploy.
