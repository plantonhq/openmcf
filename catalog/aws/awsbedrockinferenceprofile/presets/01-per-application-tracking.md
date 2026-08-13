# Per-Application Tracking

This preset creates one application inference profile routing straight to
a foundation model — the per-consumer attribution unit: give each service
its own copy of this preset (named after the service) and every service's
Bedrock usage becomes separately meterable, taggable, and IAM-scopeable.

## When to Use

- Any account where more than one application invokes Bedrock models
- Chargeback/showback models that need per-team AI cost lines

## Key Configuration Choices

- **The profile name IS the attribution key** — name it after the
  consumer, not the model.
- **Direct foundation-model source** — the single-region shape; swap the
  source to an AWS geo profile ARN to inherit cross-region routing.

## After Deployment

Grant the consuming service `bedrock:InvokeModel` on
`inference_profile_arn` (from the outputs) and have it pass that ARN as
the modelId.
