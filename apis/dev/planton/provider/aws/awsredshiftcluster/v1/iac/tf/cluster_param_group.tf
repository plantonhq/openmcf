# A managed parameter group exists only when the spec carries inline
# parameters -- a named parameter list is configuration owned by exactly
# this cluster, not a composable node. The family defaults to
# redshift-1.0 (accepted on every cluster); parameter_group_family
# selects the redshift-2.0 generation when the group should track it.
resource "aws_redshift_parameter_group" "this" {
  count = local.manage_parameter_group ? 1 : 0

  name   = local.cluster_identifier
  family = var.spec.parameter_group_family != "" ? var.spec.parameter_group_family : "redshift-1.0"
  tags   = local.aws_tags

  dynamic "parameter" {
    for_each = var.spec.parameters
    content {
      name  = parameter.value.name
      value = parameter.value.value
    }
  }
}
