---
title: "Presets"
description: "Ready-to-deploy configuration presets for Firestore Index"
type: "preset-list"
componentSlug: "firestore-index"
componentTitle: "Firestore Index"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-composite-filter-sort"
    rank: "01"
    title: "Composite Filter-and-Sort Index"
    excerpt: "The standard multi-field index: equality filter on one field, sort on another — the shape Firestore error-message links usually suggest."
  - slug: "02-vector-neighbors"
    rank: "02"
    title: "Vector Nearest-Neighbor Index"
    excerpt: "A composite index with a filter field and a vector field last — the enabler for Firestore nearest-neighbor (embedding similarity) queries."
---

# Firestore Index Presets

Ready-to-deploy configuration presets for Firestore Index. Each preset is a complete manifest you can copy, customize, and deploy.
