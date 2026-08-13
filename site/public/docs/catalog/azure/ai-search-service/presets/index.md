---
title: "Presets"
description: "Ready-to-deploy configuration presets for AI Search Service"
type: "preset-list"
componentSlug: "ai-search-service"
componentTitle: "AI Search Service"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-production-search"
    rank: "01"
    title: "Production Search Service"
    excerpt: "The production posture: a `standard` service with three replicas (the 99.9% read-write SLA floor), Entra RBAC enabled alongside API keys, and a system identity for keyless indexer connections."
  - slug: "02-dev-basic-search"
    rank: "02"
    title: "Development Basic Search"
    excerpt: "The development shape: a single-replica `basic` service -- the cheapest dedicated tier, with defaults everywhere else (API-key auth, public endpoint, one partition)."
  - slug: "03-semantic-rag-search"
    rank: "03"
    title: "Semantic RAG Search"
    excerpt: "The retrieval backend for RAG applications: a `standard` service with semantic ranking enabled, so Azure OpenAI answers ground on semantically re-ranked results rather than plain keyword scores."
---

# AI Search Service Presets

Ready-to-deploy configuration presets for AI Search Service. Each preset is a complete manifest you can copy, customize, and deploy.
