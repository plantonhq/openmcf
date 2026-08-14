---
title: "Archive Cost-Tiering ECR Repository"
description: "This preset creates an ECR repository with the canonical storage-tiering policy: images nobody pulls for 90 days move to the archive storage class (materially cheaper than standard ECR storage), and..."
type: "preset"
rank: "03"
presetSlug: "03-archive-cost-tiering"
componentSlug: "ecr-repo"
componentTitle: "ECR Repo"
provider: "aws"
icon: "package"
order: 3
---

# Archive Cost-Tiering ECR Repository

This preset creates an ECR repository with the canonical storage-tiering
policy: images nobody pulls for 90 days move to the archive storage class
(materially cheaper than standard ECR storage), and images that have sat in
archive for a year are deleted. Archived images are pulled back
transparently, with a retrieval delay — the right trade for rollback
insurance you will almost never use.

## When to Use

- High-volume registries where storage cost is a real line item
- Teams that keep long release histories "just in case" but rarely pull
  anything older than a quarter
- Compliance postures that want old images retained (cheaply) rather than
  deleted outright

## Key Configuration Choices

- **Archive-by-disuse rule** (priority 2, `sinceImagePulled` 90 days,
  `actionType: transition` + `targetStorageClass: archive`) — moves images
  to the cheaper tier based on when they were last PULLED, not pushed, so
  actively-used old releases stay in standard storage
- **Expire-from-archive rule** (priority 3, `sinceImageTransitioned` 365
  days on `storageClass: archive`) — the archive tier's own retention
  clock; together the pair implements "archive after a quarter of disuse,
  delete after a year in archive"
- **Untagged cleanup rule** (priority 1, `sinceImagePushed` 7 days) —
  removes orphaned layers and failed builds before they ever reach the
  tiering rules
- **Immutable tags** (`imageTagMutability: IMMUTABLE`) — archived history
  is only meaningful if tags cannot be rewritten

## Placeholders to Replace

The example `repositoryName` (`team-blue/checkout-service`) is a realistic
registry path — replace it with your own naming convention.

| Placeholder | Description | Where to Find |
|---|---|---|
| `<aws-region>` | The AWS region for the repository | Your region strategy |

## Related Presets

- **01-production-immutable** — exclusion-filtered immutability with
  keep-last-N retention, no storage tiering
- **02-development** — mutable tags and aggressive cleanup for dev velocity
