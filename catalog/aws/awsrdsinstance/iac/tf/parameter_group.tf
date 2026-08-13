# The DB parameter group managed for inline spec.parameters -- pure
# glue (a named parameter list), which is why it stays inside this
# module instead of being its own node. The family is derived from
# engine + engine_version (locals.tf); a family change replaces the
# group and re-associates the instance in the same apply.
resource "aws_db_parameter_group" "this" {
  count = local.manage_parameter_group ? 1 : 0

  name   = local.instance_identifier
  family = local.parameter_group_family

  dynamic "parameter" {
    for_each = var.spec.parameters
    content {
      name         = parameter.value.name
      value        = parameter.value.value
      apply_method = parameter.value.apply_method != "" ? parameter.value.apply_method : null
    }
  }

  tags = local.aws_tags
}
