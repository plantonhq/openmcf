# Streaming Tree-AH Index

The workhorse shape for RAG retrieval and live semantic search: a streaming
index with approximate (tree-AH) search and cosine-equivalent ranking.

## What this preset creates

An empty `STREAM_UPDATE` index named `RAG Document Chunks` in `us-central1`.
Vectors are upserted via the Vertex AI API (`indexes.upsertDatapoints`) and
become searchable within seconds once the index is deployed onto an index
endpoint. Ranking uses dot-product distance over unit-normalized vectors —
equivalent to cosine similarity, the standard pairing for text embeddings.

## When to use

- RAG document retrieval where chunks are added and removed continuously
- Live recommendation or personalization corpora
- Any index where waiting for batch rebuilds is unacceptable

## Remix ideas

- Set `dimensions` to your embedding model's actual output size (e.g. 768
  for many sentence encoders).
- Raise `treeAhConfig.leafNodesToSearchPercent` (up to 100) for better
  recall at the cost of query latency.
- Pin `shardSize: SHARD_SIZE_MEDIUM` when you know the corpus will outgrow
  the small 2 GB shards — it is immutable, so choose ahead of growth.
