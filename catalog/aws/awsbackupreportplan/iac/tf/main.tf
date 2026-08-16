# A Backup Audit Manager report plan: scheduled backup/copy/restore
# job or compliance reports delivered as CSV/JSON to an S3 bucket.
#
# Lifecycle facts the render below depends on:
#   - AWS report plan names forbid hyphens, so the name is
#     spec.report_plan_name (an explicit field), never metadata.name;
#   - report_setting.report_template is ForceNew from INSIDE the
#     nested block - changing it replaces the whole report plan;
#   - number_of_frameworks is sent only when positive (AWS computes it
#     otherwise) - the spec's zero-sentinel mirrors that contract;
#   - the destination bucket needs a policy allowing the Backup report
#     service to write (taught on the spec field).

resource "aws_backup_report_plan" "this" {
  name        = var.spec.report_plan_name
  description = var.spec.description != "" ? var.spec.description : null

  report_delivery_channel {
    s3_bucket_name = var.spec.delivery_channel.s3_bucket_name
    s3_key_prefix  = var.spec.delivery_channel.s3_key_prefix != "" ? var.spec.delivery_channel.s3_key_prefix : null
    formats        = length(var.spec.delivery_channel.formats) > 0 ? var.spec.delivery_channel.formats : null
  }

  report_setting {
    report_template      = var.spec.report_setting.report_template
    framework_arns       = length(var.spec.report_setting.framework_arns) > 0 ? var.spec.report_setting.framework_arns : null
    number_of_frameworks = var.spec.report_setting.number_of_frameworks != 0 ? var.spec.report_setting.number_of_frameworks : null
    accounts             = length(var.spec.report_setting.accounts) > 0 ? var.spec.report_setting.accounts : null
    organization_units   = length(var.spec.report_setting.organization_units) > 0 ? var.spec.report_setting.organization_units : null
    regions              = length(var.spec.report_setting.regions) > 0 ? var.spec.report_setting.regions : null
  }

  tags = local.aws_tags
}
