---
title: "Presets"
description: "Ready-to-deploy configuration presets for Vertex AI Index"
type: "preset-list"
componentSlug: "vertex-ai-index"
componentTitle: "Vertex AI Index"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-streaming-tree-ah"
    rank: "01"
    title: "Streaming Tree-AH Index"
    excerpt: "The workhorse shape for RAG retrieval and live semantic search: a streaming index with approximate (tree-AH) search and cosine-equivalent ranking."
  - slug: "02-batch-from-gcs"
    rank: "02"
    title: "Batch Index from Cloud Storage"
    excerpt: "A bulk-built index seeded from a Cloud Storage directory of embedding files — the economical shape for corpora refreshed on a schedule."
  - slug: "03-brute-force-eval"
    rank: "03"
    title: "Brute-Force Evaluation Index"
    excerpt: "An exact-search index over a sample corpus — the ground truth for measuring an approximate index's recall."
---

# Vertex AI Index Presets

Ready-to-deploy configuration presets for Vertex AI Index. Each preset is a complete manifest you can copy, customize, and deploy.
