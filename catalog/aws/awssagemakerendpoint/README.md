<p align="center">
  <img src="logo.svg" alt="AWS SageMaker Endpoint" width="80"/>
</p>

# AWS SageMaker Endpoint

Create and manage [Amazon SageMaker AI real-time inference endpoints](https://docs.aws.amazon.com/sagemaker/latest/dg/realtime-endpoints.html)
together with their endpoint configuration — the immutable
capacity/variant definition the endpoint points at, folded into one
resource.

## What Gets Created

- **An endpoint** plus its **endpoint configuration** — variants are
  [serverless](https://docs.aws.amazon.com/sagemaker/latest/dg/serverless-endpoints.html)
  (no idle charge, billed per inference) or instance-backed, never both on
  one variant.
- **Production variants** (1–10) splitting traffic by weight, and
  optional **shadow variants** receiving a copy of production traffic
  without answering callers (exactly one variant on each side when
  shadow testing).
- Optional: data capture to S3 (the Model Monitor feed), async
  inference with SNS notifications, KMS-encrypted storage volumes, and
  [deployment guardrails](https://docs.aws.amazon.com/sagemaker/latest/dg/deployment-guardrails.html)
  — blue/green or rolling updates with CloudWatch-alarm auto-rollback.

## Configurations Roll, the Endpoint Repoints

The endpoint configuration is immutable upstream, so any
variant/capture/async change mints a NEW name-suffixed configuration
(created before the old one is destroyed) and the endpoint repoints at
it via UpdateEndpoint — AWS's own documented pattern. The endpoint
never references a deleted configuration.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
