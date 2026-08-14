# AwsBedrockKnowledgeBase — Terraform/OpenTofu module

Deploys an Amazon Bedrock knowledge base
(`aws_bedrockagent_knowledge_base`) with its folded data sources
(`aws_bedrockagent_data_source`), keyed by their stable spec entry names.

Module facts worth knowing before editing:

- **Type discriminators are derived, never spec surface.** Exactly one of
  `spec.vector` / `spec.managed` / `spec.kendra` / `spec.sql` is set (CEL
  guarded); the module derives AWS's `type` strings from the configured
  arm — the same pattern for the eight storage backends and the six data
  source connectors.
- **Create-time-only surface.** Nearly everything except name,
  description, role, and tags is `RequiresReplace` upstream — the
  ingestion pipeline (`vector_ingestion_configuration`) replaces the DATA
  SOURCE, the type/storage arms replace the KNOWLEDGE BASE. AWS
  re-ingests after replacement.
- **Propagation retries are the provider's.** IAM assume-role, Kendra
  access, and OpenSearch data-access-policy propagation are retried
  upstream for ~2 minutes at create; the module adds none.
- **One-value vocabularies are module constants**: SQL engine `REDSHIFT`,
  crawl filter `PATTERN`, Confluence host `SAAS`, SharePoint host
  `ONLINE`, Salesforce auth `OAUTH2_CLIENT_CREDENTIALS`, parsing modality
  `MULTIMODAL`, transformation step `POST_CHUNKING`, supplemental storage
  location `S3`.
- **The OpenSearch Serverless backend needs a PRE-EXISTING vector index**
  in the collection — no AWS-provider resource creates one (upstream's own
  tests use the third-party opensearch provider). The self-contained
  backends are S3 Vectors and the MANAGED type.

Outputs mirror the Pulumi module key-for-key: `knowledge_base_id`,
`knowledge_base_arn`, `data_source_ids` (keyed by entry name).
