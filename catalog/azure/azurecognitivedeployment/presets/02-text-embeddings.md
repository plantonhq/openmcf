# Text Embeddings

This preset deploys `text-embedding-3-large` with a pinned upgrade policy -- the shape retrieval and RAG pipelines want: vectors stored today must stay comparable with vectors computed next month, so the model version never moves on its own.

## When to Use

- Embedding pipelines feeding a vector store or search index
- RAG applications where index and query embeddings must match
- Batch document processing

## Key Configuration Choices

- **NO_AUTO_UPGRADE** -- embedding outputs differ across model versions; an automatic upgrade would silently make old vectors incomparable with new ones
- **Standard SKU** -- regional per-token billing; raise `capacity` in place when indexing jobs need more throughput
- **Re-embed deliberately** -- when you DO change the version, plan the re-index; that is a data migration, not a config change

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-cognitive-account-id>` | ARM ID of the parent account (kind OpenAI or AIServices) | `AzureCognitiveAccount` status outputs (`cognitive_account_id`), or reference it with valueFrom |
