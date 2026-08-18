# AwsS3VectorBucket — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Dimension is decided by the model, forever

An index's dimension must equal the embedding model's output exactly (Titan Text v2: 1024/512/256; Cohere: 1024) — mismatched puts are rejected vector-by-vector, and the dimension never changes. Model migration = new index + re-embed; budget for it when choosing models.

## The distance metric is part of the model choice

Cosine for normalized text embeddings (the common case), euclidean when the model was trained for it — a wrong metric returns plausible-looking but degraded results, the worst failure mode. Match the model card, not intuition.

## Non-filterable keys are the cost/latency lever

Every filterable metadata byte rides the index's query path and counts against a hard per-vector budget; payloads (source text chunks, URIs) belong in `non_filterable_metadata_keys`. Getting this wrong doesn't error — it makes queries slower and puts fail on oversized metadata.

## Bedrock is the main consumer — wire outputs, not literals

A knowledge base's s3_vectors arm takes the index ARN; hardcoding it survives nothing. The `index_arns` map keyed by index name is the stable reference charts should pass.

## Deleting embeddings deletes real money

Embeddings cost compute to produce. `force_destroy: false` (the default) makes a non-empty bucket refuse teardown — keep it that way outside scratch environments, and treat re-embedding cost as part of any bucket-replacing change.
