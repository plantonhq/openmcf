---
title: "Fine-Tune Minimal"
description: "This preset fine-tunes the cheapest fine-tunable Amazon base model (Titan Text Lite) for a single epoch — the pipeline-validation shape that proves your data format, role permissions, and S3 wiring..."
type: "preset"
rank: "01"
presetSlug: "01-fine-tune-minimal"
componentSlug: "bedrock-custom-model"
componentTitle: "Bedrock Custom Model"
provider: "aws"
icon: "package"
order: 1
---

# Fine-Tune Minimal

This preset fine-tunes the cheapest fine-tunable Amazon base model (Titan
Text Lite) for a single epoch — the pipeline-validation shape that proves
your data format, role permissions, and S3 wiring before a real training
budget is spent.

## When to Use

- The first customization run in an account
- Validating a new training dataset's format end-to-end

## Key Configuration Choices

- **One epoch, batch size one** — the minimum honest training pass.
- **Titan Text Lite** — the smallest fine-tunable base model; swap the
  ARN once the pipeline is proven.
- **Training data is JSONL** prompt/completion pairs
  (`{"prompt": "...", "completion": "..."}` per line).

## After Deployment

Watch `job_status` until Completed, review metrics in the output S3
location, then scale up epochs/data — remember every change starts a NEW
job and needs a fresh `job_name`.
