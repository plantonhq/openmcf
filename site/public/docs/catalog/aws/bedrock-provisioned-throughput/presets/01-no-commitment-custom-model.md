---
title: "No-Commitment Custom Model Capacity"
description: "This preset buys one model unit of no-commitment capacity for a fine-tuned custom model — the integration-validation shape: hourly billing, deletable the moment testing ends, referencing the custom..."
type: "preset"
rank: "01"
presetSlug: "01-no-commitment-custom-model"
componentSlug: "bedrock-provisioned-throughput"
componentTitle: "Bedrock Provisioned Throughput"
provider: "aws"
icon: "package"
order: 1
---

# No-Commitment Custom Model Capacity

This preset buys one model unit of no-commitment capacity for a fine-tuned
custom model — the integration-validation shape: hourly billing, deletable
the moment testing ends, referencing the custom model by name so the chart
wires the ARN automatically.

## When to Use

- First serving capacity for a freshly fine-tuned model
- Load and integration testing before choosing a commitment term

## Key Configuration Choices

- **No `commitmentDuration`** — hourly billing, delete any time; the safe
  default while sizing.
- **One model unit** — check the model's per-unit tokens-per-minute quote
  in the console, and your account's no-commitment quota (often 0 —
  raise it via Service Quotas first).

## After Deployment

Point applications at `provisioned_model_arn`; watch utilization metrics,
then re-create with more units or a commitment term when the workload
proves steady (both are replacements).
