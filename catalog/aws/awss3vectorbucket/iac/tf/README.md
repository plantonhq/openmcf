# AwsS3VectorBucket — Terraform/OpenTofu module

Manages one S3 vector bucket with its indexes and policy (`aws_s3vectors_vector_bucket` + `aws_s3vectors_index` + `aws_s3vectors_vector_bucket_policy`) - three provider resources under one declarative owner.

Module facts worth knowing before editing:

- **Every index property is create-only** — an index is replaced, not edited; dimension must match the embedding model BEFORE the first vector lands (the spec field teaches the common model dimensions).
- **Indexes key by name** — stable across list reorders; a rename replaces the index (and its vectors).
- **`data_type` is module-pinned to float32** — the provider's DataType enum holds exactly that one value (recorded parity exclusion).
- **Bucket encryption is fixed for life** (RequiresReplaceIfConfigured at the provider).
- **The policy is JSON-normalized by AWS** (importIgnore at the provider); **`force_destroy` is config-only** and never round-trips on import (declared in the import catalog).

Outputs mirror the Pulumi module key-for-key: `vector_bucket_arn` (import ID), `index_arns` (keyed by index name).
