---
title: "Public CDN-Origin Container"
description: "This preset creates a container serving objects anonymously by direct URL -- the public-website and CDN-origin pattern. Listing stays disabled: a client must know an object's URL to fetch it."
type: "preset"
rank: "02"
presetSlug: "02-public-cdn-origin"
componentSlug: "storage-container"
componentTitle: "Storage Container"
provider: "azure"
icon: "package"
order: 2
---

# Public CDN-Origin Container

This preset creates a container serving objects anonymously by direct
URL -- the public-website and CDN-origin pattern. Listing stays
disabled: a client must know an object's URL to fetch it.

## When to Use

- Static assets (images, scripts, downloads) fronted by a CDN or
  Front Door
- Public artifact distribution where per-object URLs are shared
  deliberately

## Key Configuration Choices

- **`containerAccessType: BLOB`** -- anonymous object reads, no
  enumeration (CONTAINER would add listing; rarely appropriate)
- **The account gate**: the parent account's
  `allowNestedItemsToBePublic` must be true (Azure's default) --
  accounts hardened with `false` make this preset undeployable by
  design

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `public-assets` | 3-63 lowercase letters/digits/hyphens | Your naming convention |

## Downstream Wiring

Point a CDN origin at the ACCOUNT's blob host; object URLs compose as
`{primary_blob_endpoint}{container_name}/{blob-path}`:

```yaml
# The origin hostname comes from the account, not the container
status.outputs.primary_blob_host  # {account}.blob.core.windows.net
```
