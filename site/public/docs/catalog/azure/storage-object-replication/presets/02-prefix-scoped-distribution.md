---
title: "Prefix-Scoped Content Distribution"
description: "This preset replicates only one namespace prefix of a container to a regional account, forward-looking from creation -- the read-local distribution pattern that keeps consumers close to their data..."
type: "preset"
rank: "02"
presetSlug: "02-prefix-scoped-distribution"
componentSlug: "storage-object-replication"
componentTitle: "Storage Object Replication"
provider: "azure"
icon: "package"
order: 2
---

# Prefix-Scoped Content Distribution

This preset replicates only one namespace prefix of a container to a
regional account, forward-looking from creation -- the read-local
distribution pattern that keeps consumers close to their data without
copying the whole container.

## When to Use

- Fanning published/approved content out to consumer-region accounts
  while drafts stay local to the origin
- Selective archival: only `exports/` or `final/` trees leave the
  origin account

## Key Configuration Choices

- **`prefixMatch` entries are INCLUDE filters** -- only blobs whose
  names start with a listed prefix replicate (matching ARM's own
  prefixMatch semantics)
- **The default `OnlyNewObjects` skips backfill** -- existing blobs
  stay put; set `copyBlobsCreatedAfter: Everything` (or an RFC 3339
  instant) to pull history too
- One policy carries up to 1000 rules -- add a rule per container pair
  rather than a policy per pair

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<origin-account-resource-name>` | The source AzureStorageAccount's Planton resource name | Your origin composition |
| `<regional-account-resource-name>` | The destination AzureStorageAccount's Planton resource name | Your regional composition |
| `<source-container-resource-name>` | The origin AzureStorageContainer | Your origin composition |
| `<destination-container-resource-name>` | The regional AzureStorageContainer | Your regional composition |
