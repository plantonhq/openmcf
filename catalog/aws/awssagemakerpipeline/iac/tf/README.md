# AwsSagemakerPipeline — Terraform/OpenTofu module

Deploys an Amazon SageMaker AI pipeline (`aws_sagemaker_pipeline`) —
the ML workflow DAG executions run against.

Module facts worth knowing before editing:

- **Everything except the name updates in place**; creating a pipeline
  is free (executions bill).
- **The component's name IS the pipeline name** — derived from
  `metadata.name`; the display name defaults to it because the provider
  REQUIRES a display name.
- **The definition comes from exactly one place** (spec-validated) —
  the inline arm is `jsonencode`d from the structured spec field; the
  S3 arm renders bucket / object key / optional version.
- **AWS's describe API returns only the RESOLVED definition, never the
  S3 location** — the location is config-only on import and S3-object
  drift is invisible to refresh (taught on the spec field).

Outputs mirror the Pulumi module key-for-key: `pipeline_name`,
`pipeline_arn`.
