---
title: "Presets"
description: "Ready-to-deploy configuration presets for Lambda Event Source Mapping"
type: "preset-list"
componentSlug: "lambda-event-source-mapping"
componentTitle: "Lambda Event Source Mapping"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-sqs-worker"
    rank: "01"
    title: "SQS Queue Worker"
    excerpt: "An SQS-to-Lambda mapping with partial-batch failure reporting -- the right default for independent record processing."
  - slug: "02-kinesis-consumer"
    rank: "02"
    title: "Kinesis Stream Consumer"
    excerpt: "A Kinesis-to-Lambda mapping that processes the backlog from the oldest available record with partial-batch failure reporting."
  - slug: "03-msk-shared-pollers"
    rank: "03"
    title: "MSK Consumer with Shared Provisioned Pollers"
    excerpt: "This preset consumes an Amazon MSK topic with a pinned poller fleet that is SHARED across mappings: every mapping naming the same `pollerGroupName` draws from (and is jointly capped by) one..."
---

# Lambda Event Source Mapping Presets

Ready-to-deploy configuration presets for Lambda Event Source Mapping. Each preset is a complete manifest you can copy, customize, and deploy.
