<p align="center">
  <img src="logo.svg" alt="AWS Bedrock Knowledge Base" width="80"/>
</p>

# AWS Bedrock Knowledge Base

Create and manage [Amazon Bedrock knowledge bases](https://docs.aws.amazon.com/bedrock/latest/userguide/knowledge-base.html) —
the RAG (retrieval-augmented generation) store that ingests documents,
embeds them, and answers retrieval queries for agents and applications.

## What Gets Created

- **A knowledge base** of exactly one type:
  - **vector** — embeddings-based retrieval over a vector store you own
    (eight backends: OpenSearch Serverless, managed OpenSearch, S3
    Vectors, Aurora/pgvector, Pinecone, MongoDB Atlas, Neptune
    Analytics GraphRAG, Redis Enterprise Cloud);
  - **managed** — embeddings-based retrieval with the vector store fully
    managed by AWS (no storage block at all);
  - **kendra** — retrieval delegated to an existing Amazon Kendra index;
  - **sql** — natural-language-to-SQL over Amazon Redshift.
- **Data sources** — document connectors, one per `data_sources` entry:
  S3 buckets, website crawls, Confluence, Salesforce, SharePoint, or
  AWS-managed connectors — each with its own chunking and parsing
  pipeline.

Knowledge bases are free to create — embedding and retrieval bill per
use.

## Type and Storage Pairing

The `vector` type pairs with a `storage` backend; `managed`, `kendra`,
and `sql` bring (or delegate) their own storage and take no `storage`
block. Nearly everything except the name, description, and role is
create-time-only — changing the type, storage, or ingestion pipeline
replaces the knowledge base or data source, and AWS re-ingests.

## Prerequisites

- An AWS provider connection in Planton.
- An IAM role trusting `bedrock.amazonaws.com` (`role_arn`) with
  data-source read, embedding-model invoke, and vector-store access.
- For the `vector` type: the vector store, pre-created (S3 Vectors is the
  lowest-friction self-contained option).

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockKnowledgeBase
metadata:
  name: product-docs
spec:
  region: us-west-2
  roleArn:
    valueFrom:
      kind: AwsIamRole
      name: bedrock-kb-role
      fieldPath: status.outputs.role_arn
  managed: {}
  dataSources:
    - name: manuals
      s3:
        bucketArn:
          valueFrom:
            kind: AwsS3Bucket
            name: product-manuals
            fieldPath: status.outputs.bucket_arn
```

## Spec Reference

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
