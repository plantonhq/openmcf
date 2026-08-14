<p align="center">
  <img src="logo.svg" alt="AWS Bedrock Provisioned Throughput" width="80"/>
</p>

# AWS Bedrock Provisioned Throughput

Purchase and manage [Amazon Bedrock Provisioned Throughput](https://docs.aws.amazon.com/bedrock/latest/userguide/prov-throughput.html) —
dedicated, guaranteed model serving capacity in model units. Provisioned
Throughput is the REQUIRED serving path for fine-tuned custom models and
an option for high-volume foundation-model workloads.

## What Gets Created

- **A provisioned model** — dedicated capacity for the referenced model,
  addressable by its own ARN as the `modelId` in
  `InvokeModel`/`Converse`.

## Cost Warning

This resource bills FROM THE MOMENT IT IS CREATED:

- **No commitment** (the default): hourly billing, delete any time.
- **OneMonth / SixMonths commitments**: discounted rates, but the FULL
  term bills regardless and the purchase CANNOT be deleted until the term
  lapses.

Model units for large models run to thousands of dollars per month — size
deliberately, and check your account's no-commitment model-unit quota
(often 0 by default) before the first purchase.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockProvisionedThroughput
metadata:
  name: support-model-capacity
spec:
  region: us-east-1
  modelArn:
    valueFrom:
      kind: AwsBedrockCustomModel
      name: support-titan-ft
  modelUnits: 1
```

## Spec Reference

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
