---
title: "Bedrock Vector Store"
description: "This preset creates a VECTORSEARCH collection shaped as the vector store for an Amazon Bedrock knowledge base, with the knowledge base's service role granted the index access it needs."
type: "preset"
rank: "03"
presetSlug: "03-bedrock-vector-store"
componentSlug: "opensearch-serverless-collection"
componentTitle: "OpenSearch Serverless Collection"
provider: "aws"
icon: "package"
order: 3
---

# Bedrock Vector Store

This preset creates a VECTORSEARCH collection shaped as the vector store for an Amazon Bedrock knowledge base, with the knowledge base's service role granted the index access it needs.

## When to Use

- Retrieval-augmented generation (RAG) with Bedrock knowledge bases
- Any vector-embedding workload (semantic search, recommendations) on managed OpenSearch

## Key Configuration Choices

- **`type: VECTORSEARCH`** — the collection type Bedrock requires; wire the `collection_arn` output into the knowledge base's vector store configuration
- **Knowledge-base service role granted `aoss:*` on indexes** — Bedrock creates and queries the vector index through its service role; without this data-access rule ingestion fails with authorization errors that read like Bedrock problems
- **`standbyReplicas: DISABLED`** — start at the half floor; switch to ENABLED (a new collection — the setting is ForceNew) when the knowledge base goes production

## Scaling Note

Vector workloads with high query volume can enable `serverlessVectorAcceleration` (GPU-accelerated capacity, VECTORSEARCH only) — leave it unset until profiling shows the need.
