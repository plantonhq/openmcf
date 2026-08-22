# One Lambda layer version published from an S3 archive, with its
# share grants managed in-line.
#
# Lifecycle facts the render below depends on:
#   - every layer-version argument is ForceNew - a config change
#     publishes a NEW version (functions keep the exact version ARN
#     they pinned, so a replacement never breaks consumers mid-run);
#   - skip_destroy leaves the previous version available in AWS on
#     replacement/destroy (dormant versions bill nothing);
#   - permissions are per-VERSION policy statements keyed by
#     statement_id - they replace alongside the version they grant;
#   - source_code_hash is a local change detector only - AWS never
#     reports it back.
#
# Neither resource is taggable at AWS - hence no tags anywhere here.

resource "aws_lambda_layer_version" "this" {
  layer_name = var.metadata.name

  s3_bucket         = var.spec.code.bucket
  s3_key            = var.spec.code.key
  s3_object_version = var.spec.code.version != "" ? var.spec.code.version : null

  source_code_hash = var.spec.source_code_hash != "" ? var.spec.source_code_hash : null

  description              = var.spec.description != "" ? var.spec.description : null
  compatible_runtimes      = length(var.spec.compatible_runtimes) > 0 ? var.spec.compatible_runtimes : null
  compatible_architectures = length(var.spec.compatible_architectures) > 0 ? var.spec.compatible_architectures : null
  license_info             = var.spec.license_info != "" ? var.spec.license_info : null

  skip_destroy = var.spec.skip_destroy ? true : null
}

resource "aws_lambda_layer_version_permission" "this" {
  for_each = { for permission in var.spec.permissions : permission.statement_id => permission }

  # lambda:GetLayerVersion is the only action AWS supports on layers -
  # pinned here, never spec surface.
  action         = "lambda:GetLayerVersion"
  layer_name     = aws_lambda_layer_version.this.layer_arn
  version_number = aws_lambda_layer_version.this.version
  statement_id   = each.value.statement_id
  principal      = each.value.principal

  organization_id = each.value.organization_id != "" ? each.value.organization_id : null
  skip_destroy    = each.value.skip_destroy ? true : null
}
