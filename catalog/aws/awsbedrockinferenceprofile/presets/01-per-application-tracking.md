# Per-Application Tracking

This preset creates one application inference profile — the per-consumer
attribution unit: give each service its own copy of this preset (named
after the service) and every service's Bedrock usage becomes separately
meterable, taggable, and IAM-scopeable.

## When to Use

- Any account where more than one application invokes Bedrock models
- Chargeback/showback models that need per-team AI cost lines

## Key Configuration Choices

- **The profile name IS the attribution key** — name it after the
  consumer, not the model.
- **The source is the AWS geo profile ARN for Nova-class models** —
  replace `123456789012` with your account id. Models like Nova support
  only INFERENCE_PROFILE invocation, so a direct foundation-model ARN is
  rejected at create ("does not support On Demand inference"); point
  directly at `arn:...::foundation-model/<model-id>` only for
  ON_DEMAND-capable models
  (`aws bedrock list-foundation-models --by-inference-type ON_DEMAND`).

## After Deployment

Grant the consuming service `bedrock:InvokeModel` on
`inference_profile_arn` (from the outputs) and have it pass that ARN as
the modelId.
