# A managed cluster parameter group exists only when the spec carries
# inline parameters -- a named parameter list is configuration owned by
# exactly this cluster, not a composable node. The family is derived
# from the pinned engine + engine_version (locals.tf); a version change
# that crosses families forces a new group, which is exactly AWS's own
# constraint (parameter families are engine-major-scoped).
resource "aws_rds_cluster_parameter_group" "this" {
  count = local.manage_parameter_group ? 1 : 0

  name   = local.cluster_identifier
  family = local.engine_family
  tags   = local.aws_tags

  dynamic "parameter" {
    for_each = var.spec.parameters
    content {
      name  = parameter.value.name
      value = parameter.value.value
      # "immediate" applies dynamic parameters now; static parameters
      # must be "pending-reboot". Empty defers to the provider default
      # (immediate).
      apply_method = parameter.value.apply_method != "" ? parameter.value.apply_method : null
    }
  }
}
