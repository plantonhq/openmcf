# A Backup Audit Manager framework: controls continuously evaluating
# the account's backup posture on AWS Config.
#
# Lifecycle facts the render below depends on:
#   - AWS framework names forbid hyphens, so the name is
#     spec.framework_name (an explicit field), never metadata.name;
#   - evaluations run on AWS Config: without an ACTIVE recorder
#     recording the backup types, deployment lands FAILED - and the
#     provider treats FAILED as a completed apply (the failure shows
#     in deployment_status, not as an error);
#   - controls with an empty name are skipped by the provider's
#     expander (a set-diff workaround) - the spec forbids them.

resource "aws_backup_framework" "this" {
  name        = var.spec.framework_name
  description = var.spec.description != "" ? var.spec.description : null

  dynamic "control" {
    for_each = { for c in var.spec.controls : c.name => c }
    content {
      name = control.value.name

      dynamic "input_parameter" {
        for_each = control.value.input_parameters
        content {
          name  = input_parameter.value.name
          value = input_parameter.value.value
        }
      }

      dynamic "scope" {
        for_each = control.value.scope != null ? [control.value.scope] : []
        content {
          compliance_resource_ids   = length(scope.value.compliance_resource_ids) > 0 ? scope.value.compliance_resource_ids : null
          compliance_resource_types = length(scope.value.compliance_resource_types) > 0 ? scope.value.compliance_resource_types : null
          tags                      = length(scope.value.tags) > 0 ? scope.value.tags : null
        }
      }
    }
  }

  tags = local.aws_tags
}
