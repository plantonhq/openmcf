---
title: "Inline Training Pipeline"
description: "This preset carries the pipeline definition inline in the manifest — the definition lives next to the rest of the spec as diffable YAML, and drift on it is visible like any other field."
type: "preset"
rank: "01"
presetSlug: "01-inline-training-pipeline"
componentSlug: "sagemaker-pipeline"
componentTitle: "SageMaker Pipeline"
provider: "aws"
icon: "package"
order: 1
---

# Inline Training Pipeline

This preset carries the pipeline definition inline in the manifest —
the definition lives next to the rest of the spec as diffable YAML, and
drift on it is visible like any other field.

## When to Use

- Pipelines whose definitions are small enough to live in the manifest
- Teams that want the definition versioned and reviewed with the
  manifest

## What You Get

- A pipeline created from an inline placeholder DAG (one Fail step — a
  legal single-node graph AWS accepts at create)
- Server-side DAG validation on every apply — a malformed definition
  fails the apply, not the first execution

## Customize

- Replace the placeholder `Steps` with your real graph: generate it
  with the SageMaker Python SDK's `pipeline.definition()` and paste the
  output — don't hand-write step JSON
- Point `roleArn` at a role trusting `sagemaker.amazonaws.com` that can
  run the steps (training jobs, processing jobs, model registration)
- Add `parallelismMaxSteps` to cap parallel step execution across this
  pipeline's executions
