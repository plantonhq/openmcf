---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock Knowledge Base"
type: "preset-list"
componentSlug: "bedrock-knowledge-base"
componentTitle: "Bedrock Knowledge Base"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-s3-docs-on-s3-vectors"
    rank: "01"
    title: "S3 Docs on S3 Vectors"
    excerpt: "This preset creates a vector knowledge base over Amazon S3 Vectors — the pay-per-use vector store with no infrastructure to run — ingesting a documentation bucket's `manuals/` prefix with fixed-size..."
  - slug: "02-managed-store-web-crawl"
    rank: "02"
    title: "Managed Store with Web Crawl"
    excerpt: "This preset creates a MANAGED-type knowledge base — AWS provisions and runs the vector store, so there is nothing to size or operate — crawling a documentation site (host-only scope, 1000 pages, 60..."
---

# Bedrock Knowledge Base Presets

Ready-to-deploy configuration presets for Bedrock Knowledge Base. Each preset is a complete manifest you can copy, customize, and deploy.
