# Bedrock Knowledge Base Backend

The RAG starter: one index sized for Titan Text v2 (1024, cosine) with the bulky payload fields — chunk text and source URI — kept non-filterable. Point an AwsBedrockKnowledgeBase's s3_vectors arm at the `index_arns` output and ingest.
