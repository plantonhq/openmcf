---
title: "SageMaker Pipeline"
description: "SageMaker Pipeline deployment documentation"
icon: "package"
order: 100
componentName: "awssagemakerpipeline"
---

# AWS SageMaker Pipeline

An ML workflow DAG as managed infrastructure — the processing,
training, evaluation, and registration step graph that SageMaker
executions run against, declared in SageMaker's pipeline-definition
JSON and validated server-side at create. Free to create; only
executions bill.

## What Gets Created

- A pipeline named after `metadata.name`, with a Studio display name
  (defaults to the pipeline name) and an execution role.
- Its definition — inline as structured YAML/JSON, or read from an S3
  object (optionally pinned to a version) — exactly one of the two.
- Optionally, a default cap on parallel step execution.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SageMaker control-plane permissions
  (`sagemaker:CreatePipeline` and its siblings).

### AWS Account

- An IAM role trusting `sagemaker.amazonaws.com` that can run the
  pipeline's steps — training jobs, processing jobs, model registration
  (`role_arn`, referenceable from an AwsIamRole).
- For the S3 arm: the bucket and object holding the definition JSON.

## Deploy

### Console

Create the resource from the AWS catalog, pick the region and execution
role, paste the definition (or point at its S3 object), and deploy.

### CLI

```bash
planton apply -f pipeline.yaml
```

## After Deploy

- `pipeline_name` / `pipeline_arn` identify the pipeline; executions
  start against the name — imperatively, from schedules, event rules,
  or SDK calls.
- Everything except the name updates in place — iterate on the
  definition freely; each apply re-validates the DAG server-side.
- With the S3 arm, remember the blind spot: describe never returns the
  S3 location, so a changed S3 object is invisible drift — pin
  `version_id` or re-apply to pick up changes.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
