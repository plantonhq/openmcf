# App Runner auto scaling configuration (version).
#
# AWS versions this resource: every settable value is create-time immutable,
# so any change destroys this revision and registers the NEXT revision under
# the same configuration name. That is the intended lifecycle, not an
# accident -- referencing services pick up the new revision-carrying ARN
# through the resource graph on their next deployment.
resource "aws_apprunner_auto_scaling_configuration_version" "this" {
  auto_scaling_configuration_name = local.resource_name

  # Nulls fall through to AWS's own defaults (100 concurrency / 25 max /
  # 1 min) -- the platform normally materializes the spec defaults before
  # the module runs, so these are explicit in practice.
  max_concurrency = var.spec.max_concurrency
  max_size        = var.spec.max_size
  min_size        = var.spec.min_size

  tags = local.aws_tags
}
