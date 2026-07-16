---
title: "Vector Nearest-Neighbor Index"
description: "A composite index with a filter field and a vector field last — the enabler for Firestore nearest-neighbor (embedding similarity) queries."
type: "preset"
rank: "02"
presetSlug: "02-vector-neighbors"
componentSlug: "firestore-index"
componentTitle: "Firestore Index"
provider: "gcp"
icon: "package"
order: 2
---

# Vector Nearest-Neighbor Index

A composite index with a filter field and a vector field last — the
enabler for Firestore nearest-neighbor (embedding similarity) queries.

## When to use

Semantic search, recommendation, or RAG retrieval where queries filter
on metadata then rank by embedding distance.

## What to customize

- `fields[].vectorConfig.dimension` — must match your embedding model's
  output size (e.g. 768).
- Preceding `order` fields — any equality filters your query applies
  before vector search.

## Composes with

`GcpFirestoreDatabase` upstream. Vector fields must be last in the
index (Firestore appends `__name__` automatically when needed).
