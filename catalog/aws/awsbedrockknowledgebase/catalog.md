# AWS Bedrock Knowledge Base

The RAG store for the Bedrock AI stack — ingest documents from S3,
websites, or SaaS connectors, embed them, and serve retrieval queries to
agents and applications, with the vector store owned or fully managed by
AWS.

## What Gets Created

- A Bedrock knowledge base of the type you declare: vector (eight vector
  store backends), managed (AWS runs the store), Kendra delegation, or
  natural-language-to-SQL over Redshift.
- Its data sources, keyed by your stable entry names: S3, web crawler,
  Confluence, Salesforce, SharePoint, or AWS-managed connectors — each
  with fixed-size/hierarchical/semantic chunking, multimodal parsing, and
  optional Lambda transformation.

Creating a knowledge base is free; embedding at ingestion and retrieval
at query time bill per use (a managed-type store bills for what AWS
provisions behind it).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Bedrock agent permissions
  (`bedrock:CreateKnowledgeBase`, `bedrock:CreateDataSource` and their
  read/update/delete siblings, plus `iam:PassRole` on the role).

### AWS Account

- Bedrock knowledge bases available in the target region.
- An IAM role trusting `bedrock.amazonaws.com` with data-source read,
  embedding-model invoke, and vector-store data-plane access.
- For the vector type: the store pre-created (an S3 Vectors index, an
  OpenSearch Serverless collection WITH a vector index, an Aurora
  pgvector table, ...).

## Deploy

### Console

Create the resource from the AWS catalog, pick the region and type,
reference the role and store, add data sources, and deploy.

### CLI

```bash
planton apply -f knowledge-base.yaml
```

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
  vector:
    embeddingModelArn: arn:aws:bedrock:us-west-2::foundation-model/amazon.titan-embed-text-v2:0
    embeddingModel:
      dimensions: 256
  storage:
    s3Vectors:
      indexArn: arn:aws:s3vectors:us-west-2:123456789012:bucket/docs-vectors/index/docs
  dataSources:
    - name: manuals
      dataDeletionPolicy: DELETE
      s3:
        bucketArn:
          value: arn:aws:s3:::product-manuals
      vectorIngestion:
        chunking:
          strategy: FIXED_SIZE
          fixedSize:
            maxTokens: 300
            overlapPercentage: 20
```

## Operational Notes

- **Almost everything replaces.** Only the name, description, role, and
  tags update in place; type, storage, and ingestion changes replace the
  resource and AWS re-ingests from the sources afterwards.
- **Match embedding dimensions to the store.** The embedding model's
  `dimensions` must equal the vector index's dimension — a mismatch
  fails at ingestion, not at create.
- **OpenSearch Serverless needs a pre-existing vector index.** AWS does
  not create it; S3 Vectors and the managed type are the self-contained
  paths.
- **`data_deletion_policy: DELETE` keeps teardowns clean** — ingested
  vectors are purged when the data source is removed; RETAIN leaves them
  queryable.
