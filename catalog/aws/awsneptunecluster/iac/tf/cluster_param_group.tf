# A managed cluster parameter group exists only when the spec carries
# inline parameters -- a named parameter list is configuration owned by
# exactly this cluster, not a composable node. The family is derived
# from the pinned engine_version (locals.tf); a version change that
# crosses families forces a new group, which is exactly AWS's own
# constraint (parameter families are engine-major-scoped).
resource "aws_neptune_cluster_parameter_group" "this" {
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
      # must be "pending-reboot". Empty defers to the provider default,
      # which is pending-reboot at the pinned provider -- set "immediate"
      # explicitly when a dynamic parameter should land right away.
      apply_method = parameter.value.apply_method != "" ? parameter.value.apply_method : null
    }
  }
}

# A managed INSTANCE parameter group exists only when the spec carries
# inline instance_parameters -- the instance-level twin of the cluster
# group above, applied to every folded instance that does not bring its
# own group. Same family derivation, same ownership reasoning.
resource "aws_neptune_parameter_group" "instance" {
  count = local.manage_instance_parameter_group ? 1 : 0

  name   = "${local.cluster_identifier}-instance"
  family = local.engine_family
  tags   = local.aws_tags

  dynamic "parameter" {
    for_each = var.spec.instance_parameters
    content {
      name  = parameter.value.name
      value = parameter.value.value
      # Same apply_method semantics as the cluster group: empty defers to
      # the provider default (pending-reboot at the pinned provider).
      apply_method = parameter.value.apply_method != "" ? parameter.value.apply_method : null
    }
  }
}
