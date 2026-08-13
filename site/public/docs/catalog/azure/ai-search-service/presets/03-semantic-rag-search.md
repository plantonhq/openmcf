---
title: "Semantic RAG Search"
description: "The retrieval backend for RAG applications: a `standard` service with semantic ranking enabled, so Azure OpenAI answers ground on semantically re-ranked results rather than plain keyword scores."
type: "preset"
rank: "03"
presetSlug: "03-semantic-rag-search"
componentSlug: "ai-search-service"
componentTitle: "AI Search Service"
provider: "azure"
icon: "package"
order: 3
---

# Semantic RAG Search

The retrieval backend for RAG applications: a `standard` service with
semantic ranking enabled, so Azure OpenAI answers ground on
semantically re-ranked results rather than plain keyword scores.

## When to Use

- Retrieval-augmented generation (RAG) with Azure OpenAI
- Search relevance that keyword scoring alone cannot reach
- Vector + semantic hybrid retrieval patterns

## Key Configuration Choices

- `semanticSearchSku: standard` -- unlimited semantic ranking, billed
  per use; start with `free` (1000 requests/month) to evaluate, flip
  in place when it caps out. Not available on the `free` SERVICE sku.
- `identity.type: SYSTEM_ASSIGNED` -- lets indexers pull source data
  by managed identity while the RAG application queries with a query
  key or Entra.
- Pair with an AzureCognitiveAccount (Azure OpenAI) -- the account's
  custom-question-answering integration references a search service.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-resource-group-name>` | The resource group to create the service in | Portal -> Resource groups |
| `acme-search-rag` | Your globally-unique service name | It becomes {name}.search.windows.net |

## Related Presets

- `01-production-search` -- add the SLA/RBAC posture on top of this.
