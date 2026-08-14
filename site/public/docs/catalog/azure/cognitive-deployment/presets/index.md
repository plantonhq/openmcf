---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cognitive Deployment"
type: "preset-list"
componentSlug: "cognitive-deployment"
componentTitle: "Cognitive Deployment"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-gpt-4o-chat"
    rank: "01"
    title: "GPT-4o Chat"
    excerpt: "This preset deploys `gpt-4o-mini` on the GlobalStandard SKU -- the standard starting point for chat and completion workloads: per-token billing, no idle cost, and the widest regional model..."
  - slug: "02-text-embeddings"
    rank: "02"
    title: "Text Embeddings"
    excerpt: "This preset deploys `text-embedding-3-large` with a pinned upgrade policy -- the shape retrieval and RAG pipelines want: vectors stored today must stay comparable with vectors computed next month, so..."
---

# Cognitive Deployment Presets

Ready-to-deploy configuration presets for Cognitive Deployment. Each preset is a complete manifest you can copy, customize, and deploy.
