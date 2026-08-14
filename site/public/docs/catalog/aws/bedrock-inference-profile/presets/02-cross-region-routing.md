---
title: "Cross-Region Routing"
description: "This preset creates an application profile sourced from an AWS system-defined geo inference profile (`us.` prefix) — invocations then ride AWS's cross-region capacity pools for higher availability..."
type: "preset"
rank: "02"
presetSlug: "02-cross-region-routing"
componentSlug: "bedrock-inference-profile"
componentTitle: "Bedrock Inference Profile"
provider: "aws"
icon: "package"
order: 2
---

# Cross-Region Routing

This preset creates an application profile sourced from an AWS
system-defined geo inference profile (`us.` prefix) — invocations then
ride AWS's cross-region capacity pools for higher availability and burst
headroom, while your application still gets its own attribution ARN.

## When to Use

- Production assistants that must absorb regional capacity fluctuations
- Models whose on-demand quotas are tight in a single region

## Key Configuration Choices

- **Replace the account id (123456789012) with your own** — AWS surfaces
  its system-defined geo profiles under your account's namespace.
- **The geo prefix must match your region's geography** (`us.` for US
  regions, `eu.` for European ones).
- The underlying model may require access first — see
  `AwsBedrockModelAccess`.

## After Deployment

Same consumption pattern as any profile: applications pass
`inference_profile_arn` as the modelId.
