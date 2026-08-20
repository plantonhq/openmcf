# Amazon SageMaker AI pipeline: the ML workflow DAG executions run
# against.
#
# Lifecycle facts the renders below depend on:
#   - everything except the name updates in place; creating a pipeline
#     is free (executions bill);
#   - the definition comes from exactly one place (spec-validated) -
#     inline JSON or an S3 object;
#   - AWS's describe API returns only the RESOLVED definition, never the
#     S3 location - the location is config-only on import and S3-object
#     drift is invisible to refresh (taught on the spec field).

resource "aws_sagemaker_pipeline" "this" {
  # The component's name IS the pipeline name.
  pipeline_name = local.pipeline_name

  # Required by the provider - defaults to the pipeline name.
  pipeline_display_name = local.display_name

  pipeline_description = var.spec.description != "" ? var.spec.description : null

  role_arn = var.spec.role_arn

  # The inline definition arm.
  pipeline_definition = var.spec.definition != null ? jsonencode(var.spec.definition) : null

  # The S3-location arm.
  dynamic "pipeline_definition_s3_location" {
    for_each = var.spec.definition_s3_location != null ? [var.spec.definition_s3_location] : []
    content {
      bucket     = pipeline_definition_s3_location.value.bucket
      object_key = pipeline_definition_s3_location.value.object_key
      version_id = pipeline_definition_s3_location.value.version_id != "" ? pipeline_definition_s3_location.value.version_id : null
    }
  }

  dynamic "parallelism_configuration" {
    for_each = var.spec.parallelism_max_steps != null ? [var.spec.parallelism_max_steps] : []
    content {
      max_parallel_execution_steps = parallelism_configuration.value
    }
  }

  tags = local.aws_tags
}
