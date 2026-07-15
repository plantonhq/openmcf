---
title: "CloudWatch Event Destination"
description: "Transactional posture plus a CloudWatch event destination for zero-infrastructure bounce and complaint metrics — alarm on `AWS/SES` without standing up SNS or Firehose."
type: "preset"
rank: "02"
presetSlug: "02-cloudwatch-events"
componentSlug: "ses-configuration-set"
componentTitle: "SES Configuration Set"
provider: "aws"
icon: "package"
order: 2
---

# CloudWatch Event Destination

Transactional posture plus a CloudWatch event destination for zero-infrastructure
bounce and complaint metrics — alarm on `AWS/SES` without standing up SNS or Firehose.

## When to Use

- You want bounce-rate alarms without additional AWS plumbing
- Campaign or tenant slicing via `MESSAGE_TAG` dimensions at send time

## What It Configures

- TLS REQUIRE + BOUNCE/COMPLAINT suppression (same as the minimal preset)
- A **`bounce-metrics`** destination publishing SEND, BOUNCE, COMPLAINT, and DELIVERY
  events to CloudWatch, dimensioned by a `campaign` message tag

## What to Customize

- Replace `<aws-region>`
- Tag sends at runtime with `X-SES-MESSAGE-TAGS` so the `campaign` dimension populates
- Add SNS or EventBridge destinations for operational feedback loops
