---
title: "Presets"
description: "Ready-to-deploy configuration presets for OpenSearch Serverless Collection"
type: "preset-list"
componentSlug: "opensearch-serverless-collection"
componentTitle: "OpenSearch Serverless Collection"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-dev-search"
    rank: "01"
    title: "Dev Search Collection"
    excerpt: "This preset creates the cheapest usable search collection: standby replicas disabled (half the OCU floor), public network reachability, and a data-access rule granting one application role read/write."
  - slug: "02-production-timeseries"
    rank: "02"
    title: "Production Time-Series (Log Analytics)"
    excerpt: "This preset creates a production TIMESERIES collection with standby replicas, a customer-managed KMS key, least-privilege data access split between an ingest role and an analyst role, and retention..."
  - slug: "03-bedrock-vector-store"
    rank: "03"
    title: "Bedrock Vector Store"
    excerpt: "This preset creates a VECTORSEARCH collection shaped as the vector store for an Amazon Bedrock knowledge base, with the knowledge base's service role granted the index access it needs."
---

# OpenSearch Serverless Collection Presets

Ready-to-deploy configuration presets for OpenSearch Serverless Collection. Each preset is a complete manifest you can copy, customize, and deploy.
