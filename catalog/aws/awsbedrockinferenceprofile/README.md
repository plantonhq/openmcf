<p align="center">
  <img src="logo.svg" alt="AWS Bedrock Inference Profile" width="80"/>
</p>

# AWS Bedrock Inference Profile

Create and manage [Amazon Bedrock application inference profiles](https://docs.aws.amazon.com/bedrock/latest/userguide/inference-profiles.html) —
named, taggable handles over a foundation model (or an AWS cross-region
profile) so each application's or tenant's model usage can be tracked,
cost-allocated, and IAM-scoped independently.

## What Gets Created

- **An application inference profile** whose ARN applications pass as the
  `modelId` to `InvokeModel`/`Converse`. The profile itself is free;
  invocations bill at the underlying model's rates.

## Why Profiles

One foundation model, many consumers: without profiles, all usage lands in
one undifferentiated bucket. A profile per application/tenant gives each
its own ARN — cost-allocation tags, CloudWatch metrics, and IAM policies
then work per consumer.

## Everything Is Create-Time-Immutable

AWS provides no update for application inference profiles (beyond tags):
changing the source or description replaces the profile AND ITS ARN —
update consumers pinning the old ARN.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockInferenceProfile
metadata:
  name: checkout-nova
spec:
  region: us-west-2
  description: Cost tracking for the checkout service
  sourceArn: arn:aws:bedrock:us-west-2::foundation-model/amazon.nova-micro-v1:0
```

To inherit cross-region routing, point `sourceArn` at the AWS
system-defined geo profile instead
(`arn:aws:bedrock:<region>:<account>:inference-profile/us.<model-id>`).

## Spec Reference

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
