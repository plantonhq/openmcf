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
    excerpt: "SQS-to-Lambda mapping with ReportBatchItemFailures and composed AwsLambda and AwsSqsQueue references."
  - slug: "02-kinesis-consumer"
    rank: "02"
    title: "Kinesis Stream Consumer"
    excerpt: "Kinesis stream consumer from TRIM_HORIZON with partial-batch failures and parallelization_factor 2."
---

# Lambda Event Source Mapping Presets

Ready-to-deploy configuration presets for Lambda Event Source Mapping. Each preset is a complete manifest you can copy, customize, and deploy.
