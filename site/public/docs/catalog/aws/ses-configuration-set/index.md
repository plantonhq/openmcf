---
title: "SES Configuration Set"
description: "SES Configuration Set deployment documentation"
icon: "package"
order: 100
componentName: "awssesconfigurationset"
---

# AWS SES Configuration Set

Define an SES sending posture once — TLS enforcement, dedicated IP pool, open/click tracking, suppression, deliverability dashboards — and publish email events (sends, bounces, complaints, opens, clicks) to CloudWatch, EventBridge, Kinesis Firehose, SNS, or Pinpoint. Email identities reference a configuration set as their default rules, and any SendEmail call can name one explicitly.

## What Gets Created

- An SESv2 configuration set carrying delivery, tracking, suppression, reputation, sending, and VDM options.
- One SESv2 event destination per named `eventDestinations[]` entry, each publishing a chosen slice of email events into exactly one AWS destination.

## Prerequisites

- None — the configuration set is a dependency-free leaf. Event destinations optionally reference an SNS topic, EventBridge bus, or Firehose delivery stream + IAM role.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsSesConfigurationSet
metadata:
  name: transactional-set
spec:
  region: us-west-2
  reputationMetricsEnabled: true
  deliveryOptions:
    tlsPolicy: REQUIRE
  suppressedReasons:
    - BOUNCE
    - COMPLAINT
```

## Configuration Reference

### Common Fields

| Field | Description |
|---|---|
| `region` | AWS region; identities can only reference sets in their own region. |
| `deliveryOptions.tlsPolicy` | `REQUIRE` (fail sends the receiver cannot TLS) or `OPTIONAL` (the SMTP default). |
| `deliveryOptions.maxDeliverySeconds` | Retry window before bouncing, 300–50400 seconds (AWS default 14 hours). |
| `deliveryOptions.sendingPoolName` | Dedicated IP pool to send from; unset uses the shared SES IP space. |
| `reputationMetricsEnabled` | Publish bounce/complaint rates to CloudWatch — recommended for any production sender. |
| `sendingEnabled` | The per-set kill switch (default true). |
| `suppressedReasons` | Account-suppression reasons honored by this set: `BOUNCE`, `COMPLAINT`. |
| `trackingOptions.customRedirectDomain` | Keep open/click tracking URLs on your own domain (CNAME to the regional SES tracking endpoint). |
| `vdmOptions` | Virtual Deliverability Manager engagement metrics + optimized shared delivery overrides. |
| `eventDestinations[]` | Named event publishers — exactly one destination arm each: `cloudWatch`, `eventBus`, `firehose`, `snsTopic`, or `pinpointApplicationArn`. |

## Stack Outputs

| Output | Description |
|---|---|
| `configuration_set_arn` | ARN for IAM policies scoping who may send under the set. |
| `configuration_set_name` | The name identities reference and SendEmail calls name. |

## Related Components

- [AWS SES Email Identity](/docs/catalog/aws/ses-email-identity) — references a configuration set as its default sending rules.
- [AWS SNS Topic](/docs/catalog/aws/sns-topic) — the classic bounce/complaint feedback-loop destination.
- [AWS Kinesis Firehose](/docs/catalog/aws/kinesis-firehose) — durable event analytics into S3/Redshift/OpenSearch.
- [AWS EventBridge Bus](/docs/catalog/aws/eventbridge-bus) — rule-based routing of email events.
