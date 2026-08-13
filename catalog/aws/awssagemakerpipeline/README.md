<p align="center">
  <img src="logo.svg" alt="AWS SageMaker Pipeline" width="80"/>
</p>

# AWS SageMaker Pipeline

Create and manage [Amazon SageMaker AI pipelines](https://docs.aws.amazon.com/sagemaker/latest/dg/pipelines.html)
— the ML workflow DAGs (processing, training, evaluation, registration
steps) that executions run against. Creating a pipeline is free; only
executions bill.

## What Gets Created

- **A pipeline** whose AWS name derives from `metadata.name` (the
  Studio display name defaults to it), running as the execution role in
  `role_arn`.
- **The definition** — SageMaker's pipeline-definition JSON schema,
  normally produced by the SageMaker Python SDK's
  `pipeline.definition()` — provided inline as structured data
  (`definition`) or read from an S3 object (`definition_s3_location`):
  exactly one.
- Optional: `parallelism_max_steps`, a default cap on steps executed in
  parallel across this pipeline's executions.

## AWS Validates the DAG at Create

The definition's step graph is validated server-side at create — a
green apply IS the definition-validity claim. But note the S3 arm's
blind spot: AWS's describe API returns only the RESOLVED definition,
never the S3 location, so drift on the S3 object is invisible to
refresh.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
