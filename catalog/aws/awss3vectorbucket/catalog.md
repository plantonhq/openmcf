# AWS S3 Vector Bucket

Vector search at S3 economics: store billions of embeddings and pay for storage plus queries — not for an always-on vector database cluster. The storage layer for Bedrock knowledge bases, semantic search, and RAG applications.

## What Gets Managed

- The vector bucket: encryption at rest, force-destroy posture, and its resource policy.
- Its indexes: dimension (must match the embedding model), distance metric (cosine/euclidean), non-filterable metadata keys, and per-index encryption.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with S3 Vectors permissions.

### AWS Prerequisites

- None — vector buckets are standalone. Know your embedding model's output dimension first (the index cannot change it later).

## After You Deploy

- Point a Bedrock knowledge base at an index (the `index_arns` output), or put/query vectors directly through the s3vectors API.
- Store bulky payloads under the non-filterable keys — filterable metadata has a hard per-vector budget.

## Common Changes

- New use case: a new index (a different model dimension is ALWAYS a new index).
- Rotate encryption: encryption is fixed for life of the bucket — a key change means a new bucket and a re-embed or copy.
