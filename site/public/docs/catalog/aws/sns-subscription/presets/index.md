---
title: "Presets"
description: "Ready-to-deploy configuration presets for SNS Subscription"
type: "preset-list"
componentSlug: "sns-subscription"
componentTitle: "SNS Subscription"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-sqs-fanout"
    rank: "01"
    title: "SQS Fan-out Consumer"
    excerpt: "Subscribe an SQS queue to a topic with attribute filtering and raw delivery — the standard pub/sub fan-out edge."
  - slug: "02-https-webhook"
    rank: "02"
    title: "HTTPS Webhook with Dead-Letter Queue"
    excerpt: "Push messages to a self-confirming HTTPS endpoint and capture undeliverable messages in an SQS DLQ."
---

# SNS Subscription Presets

Ready-to-deploy configuration presets for SNS Subscription. Each preset is a complete manifest you can copy, customize, and deploy.
