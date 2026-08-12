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

# -----------------------------------------------------------------------------
# Account-default designation
#
# Claims this configuration as the account/region default for new App Runner
# services created without an explicit configuration. One default exists per
# account per region: claiming it silently displaces the previous holder
# (AWS marks it non-default), and only services created AFTERWARDS are
# affected. One-way at AWS -- destroying this resource is a provider no-op
# (AWS has no restore API), so dropping the flag leaves the designation in
# place until another configuration claims it.
# -----------------------------------------------------------------------------
resource "aws_apprunner_default_auto_scaling_configuration_version" "this" {
  count = var.spec.set_as_account_default ? 1 : 0

  auto_scaling_configuration_arn = aws_apprunner_auto_scaling_configuration_version.this.arn
}
