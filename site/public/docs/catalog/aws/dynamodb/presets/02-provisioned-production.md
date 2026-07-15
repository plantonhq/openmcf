---
title: "Provisioned Production Table"
description: "This preset creates a production DynamoDB table with provisioned capacity, a composite primary key (partition + sort), a global secondary index for an alternate query pattern, and the..."
type: "preset"
rank: "02"
presetSlug: "02-provisioned-production"
componentSlug: "dynamodb"
componentTitle: "DynamoDB"
provider: "aws"
icon: "package"
order: 2
---

# Provisioned Production Table

This preset creates a production DynamoDB table with provisioned capacity, a composite primary key (partition + sort), a global secondary index for an alternate query pattern, and the production-safety trio: point-in-time recovery, AWS-managed encryption, and deletion protection. Contributor insights is enabled so hot-key and throttling diagnostics are available from day one.

## When to Use

- Sustained, predictable traffic where reserved capacity pricing beats on-demand
- Single-table-design applications that query by composite key and at least one alternate shape
- Workloads planning to purchase reserved capacity (which applies only to provisioned tables)

## Key Configuration Choices

- **Provisioned billing** (`billingMode: PROVISIONED` + `provisionedThroughput`) -- 25 RCU / 25 WCU baseline; each global secondary index carries its own capacity (10/10 here)
- **Composite primary key** (`pk` HASH + `sk` RANGE) -- The single-table-design shape: entity type and hierarchy encoded in the two key attributes
- **Global secondary index** (`gsi1` on `gsi1pk`, projecting ALL) -- One alternate query pattern; add more GSIs (and per-index capacity) as access patterns emerge -- GSIs edit in place on a live table
- **Contributor insights** (`contributorInsights.enabled: true`) -- CloudWatch per-key access profiling that answers "which partition keys are hot or throttled"
- **Production safety** -- Point-in-time recovery, the AWS-managed `aws/dynamodb` encryption key (reference an `AwsKmsKey` via `serverSideEncryption.kmsKeyArn` to hold your own key), and deletion protection

## Placeholders to Replace

- `<aws-region>` -- The AWS region for the table (e.g. `us-west-2`)

Rename `pk`, `sk`, and `gsi1pk` to match your data model, and size `provisionedThroughput` for your sustained baseline -- if traffic is spiky, prefer the on-demand preset instead.

## Related Presets

- **01-on-demand-simple** -- Use instead for variable or unpredictable traffic
- **03-global-table** -- Use instead for multi-region active-active deployments
