---
title: "Single Environment Queue"
description: "The standard starting point: one queue mapped onto one compute environment, with stuck-job protection so a mis-sized job cannot block the queue forever."
type: "preset"
rank: "01"
presetSlug: "01-single-environment"
componentSlug: "batch-job-queue"
componentTitle: "Batch Job Queue"
provider: "aws"
icon: "package"
order: 1
---

# Single Environment Queue

The standard starting point: one queue mapped onto one compute environment, with stuck-job protection so a mis-sized job cannot block the queue forever.

## When to Use

- Every first Batch deployment — a compute environment cannot accept jobs without a queue
- Single-team workloads with one compute profile

## What It Configures

- **One environment at order 1** — referenced by ARN through the compute environment's output
- **Stuck-job fuse** — jobs sitting in RUNNABLE for over an hour are cancelled with an operator-readable reason
- **Priority 10** — leaves headroom to add lower-priority backfill queues on the same environment later

## What to Customize

- Replace `<aws-region>` and the referenced compute environment name
- Add a second `computeEnvironmentOrder` entry for the Spot-overflow pattern (see the next preset)
- Attach an `AwsBatchSchedulingPolicy` via `schedulingPolicy` for fair-share ordering — noting it can never be removed once set
