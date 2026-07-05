---
title: "Presets"
description: "Ready-to-deploy configuration presets for DynamoDB"
type: "preset-list"
componentSlug: "dynamodb"
componentTitle: "DynamoDB"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-on-demand-simple"
    rank: "01"
    title: "On-Demand Simple Table"
    excerpt: "This preset creates a DynamoDB table with on-demand (pay-per-request) billing and a simple partition key. On-demand pricing automatically scales to handle any traffic level without capacity planning..."
  - slug: "02-provisioned-production"
    rank: "02"
    title: "Provisioned Production Table"
    excerpt: "This preset creates a production DynamoDB table with provisioned capacity, a composite primary key (partition + sort), a global secondary index for an alternate query pattern, and the..."
  - slug: "03-global-table"
    rank: "03"
    title: "Global Table (Multi-Region Active-Active)"
    excerpt: "This preset creates a DynamoDB Global Tables v2 deployment: the table in its home region plus one full read/write replica in a second region. Applications write to whichever region is closest and..."
---

# DynamoDB Presets

Ready-to-deploy configuration presets for DynamoDB. Each preset is a complete manifest you can copy, customize, and deploy.
