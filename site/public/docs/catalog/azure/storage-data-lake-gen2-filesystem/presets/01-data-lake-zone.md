---
title: "Data Lake Zone Filesystem"
description: "This preset creates one zone of a medallion-style data lake (raw, curated, gold) as its own filesystem -- the grain that gives each zone its own POSIX posture and its own RBAC scope."
type: "preset"
rank: "01"
presetSlug: "01-data-lake-zone"
componentSlug: "storage-data-lake-gen2-filesystem"
componentTitle: "Storage Data Lake Gen2 Filesystem"
provider: "azure"
icon: "package"
order: 1
---

# Data Lake Zone Filesystem

This preset creates one zone of a medallion-style data lake (raw,
curated, gold) as its own filesystem -- the grain that gives each zone
its own POSIX posture and its own RBAC scope.

## When to Use

- Standing up the classic raw/curated/gold lake layout (deploy this
  preset once per zone, changing the name)
- Any analytics storage where a zone's write access and read access
  belong to different teams

## Key Configuration Choices

- **DEFAULT-scope entries are the inheritance template** -- files and
  directories landing in the zone inherit group-read/other-none without
  per-object management
- **The owning group reads, everyone else is shut out** -- grant teams
  by adding qualified `objectId` entries (USER/GROUP with an Entra
  object ID), or grant at the filesystem scope with an
  AzureRoleAssignment against the `filesystem_id` output
- **The account must be HNS-enabled** (`isHnsEnabled: true`) -- the
  POSIX surface only exists on Data Lake Gen2 accounts

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The HNS AzureStorageAccount's Planton resource name | Your lake composition |
