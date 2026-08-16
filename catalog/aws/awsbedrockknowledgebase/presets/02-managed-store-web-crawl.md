# Managed Store with Web Crawl

This preset creates a MANAGED-type knowledge base — AWS provisions and
runs the vector store, so there is nothing to size or operate — crawling
a documentation site (host-only scope, 1000 pages, 60 URLs/minute) with
hierarchical chunking suited to long structured pages.

## When to Use

- Teams that want RAG without owning a vector store at all
- Public documentation or marketing sites as the knowledge source

## What You Get

- A managed knowledge base with AWS's default embedding model (set
  `managed.embeddingModelArn` to pick one explicitly)
- A polite web crawler: host-scoped, page-capped, and rate-limited
- Parent/child chunking: children are retrieved for precision, parents
  supply context for generation

## Customize

- Add `inclusionFilters`/`exclusionFilters` (regex) to narrow the crawl
- Add a second data source (e.g. an S3 bucket) — sources sync
  independently into the same store
- Set `managed.kmsKeyArn` to encrypt the managed store under your own key
