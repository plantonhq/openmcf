---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bigtable Table"
type: "preset-list"
componentSlug: "bigtable-table"
componentTitle: "Bigtable Table"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-time-series"
    rank: "01"
    title: "Time-Series Table"
    excerpt: "The classic Bigtable shape: a wide measurements family with age-based retention, a small metadata family capped by versions, and key-prefix pre-splits so initial load distributes across tablets."
  - slug: "02-cdc-enabled"
    rank: "02"
    title: "CDC-Enabled Table with Backups"
    excerpt: "An operational table wired for downstream processing: a change-stream feed for Dataflow CDC pipelines, daily automated backups, and a combined age-plus-versions retention policy."
  - slug: "03-aggregate-counters"
    rank: "03"
    title: "Aggregate Counter Table"
    excerpt: "A metering/analytics table using Bigtable's server-side aggregate cells: the `intsum` family increments atomically at write time, eliminating read-modify-write races and application-side counter..."
---

# Bigtable Table Presets

Ready-to-deploy configuration presets for Bigtable Table. Each preset is a complete manifest you can copy, customize, and deploy.
