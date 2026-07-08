---
title: "Cross-Region DR Replication"
description: "This preset continuously replicates a container to an account in another region, backfilling everything that already exists -- the blob-level disaster-recovery posture where reads can fail over to..."
type: "preset"
rank: "01"
presetSlug: "01-cross-region-dr"
componentSlug: "storage-object-replication"
componentTitle: "Storage Object Replication"
provider: "azure"
icon: "package"
order: 1
---

# Cross-Region DR Replication

This preset continuously replicates a container to an account in
another region, backfilling everything that already exists -- the
blob-level disaster-recovery posture where reads can fail over to the
replica.

## When to Use

- DR for blob data whose account-level geo-replication (GRS) is not
  enough: object replication gives you a READABLE, independently
  configured replica account
- Migrating a container between regions with a live cutover

## Key Configuration Choices

- **`copyBlobsCreatedAfter: Everything` backfills the container** --
  completion takes time proportional to its size; monitor with
  `az storage account or-policy show --policy-id <policy_id output>`
- **Account prerequisites**: the source account needs
  `blobProperties.versioningEnabled` + `changeFeedEnabled`, the
  destination `versioningEnabled` -- Azure rejects the policy otherwise
- **Replication is asynchronous with no default RPO** -- treat the
  replica as eventually consistent

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<primary-account-resource-name>` | The source AzureStorageAccount's Planton resource name | Your primary-region composition |
| `<dr-account-resource-name>` | The destination AzureStorageAccount's Planton resource name | Your DR-region composition |
| `<source-container-resource-name>` | The replicated AzureStorageContainer | Your primary-region composition |
| `<destination-container-resource-name>` | The replica AzureStorageContainer | Your DR-region composition |
