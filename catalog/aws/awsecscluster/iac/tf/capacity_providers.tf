# The cluster's capacity, materialized from the folded spec entries: one
# aws_ecs_capacity_provider per EC2 entry (keyed by name, so adding or
# removing an entry never disturbs its siblings) plus ONE association
# that PUTs the union of Fargate built-ins and EC2 provider names onto
# the cluster together with the default strategy.
#
# The association is a whole-set PUT -- exactly why it is a single
# resource: two association resources on one cluster would fight each
# other on every apply. FARGATE and FARGATE_SPOT are AWS-managed
# built-ins: never created, only associated.

resource "aws_ecs_capacity_provider" "this" {
  for_each = { for p in var.spec.ec2_capacity_providers : p.name => p }

  name = each.value.name

  # The wrapped auto-scaling group and the provider name are create-time
  # (ForceNew); the managed scaling/draining/protection knobs update in
  # place.
  auto_scaling_group_provider {
    auto_scaling_group_arn = each.value.auto_scaling_group_arn

    dynamic "managed_scaling" {
      for_each = each.value.managed_scaling != null ? [each.value.managed_scaling] : []
      content {
        status                    = managed_scaling.value.status != "" ? managed_scaling.value.status : null
        target_capacity           = managed_scaling.value.target_capacity != 0 ? managed_scaling.value.target_capacity : null
        minimum_scaling_step_size = managed_scaling.value.minimum_scaling_step_size != 0 ? managed_scaling.value.minimum_scaling_step_size : null
        maximum_scaling_step_size = managed_scaling.value.maximum_scaling_step_size != 0 ? managed_scaling.value.maximum_scaling_step_size : null
        instance_warmup_period    = managed_scaling.value.instance_warmup_period_seconds != 0 ? managed_scaling.value.instance_warmup_period_seconds : null
      }
    }

    # Requires the group's own new-instance scale-in protection when
    # ENABLED -- AWS validates the pairing at create.
    managed_termination_protection = each.value.managed_termination_protection != "" ? each.value.managed_termination_protection : null
    managed_draining               = each.value.managed_draining != "" ? each.value.managed_draining : null
  }

  tags = local.aws_tags
}

resource "aws_ecs_cluster_capacity_providers" "this" {
  # A bare cluster (services name a launch_type directly) legitimately
  # associates nothing at all.
  count = length(local.associated_capacity_providers) > 0 ? 1 : 0

  cluster_name       = aws_ecs_cluster.this.name
  capacity_providers = local.associated_capacity_providers

  dynamic "default_capacity_provider_strategy" {
    for_each = var.spec.default_capacity_provider_strategy
    content {
      capacity_provider = default_capacity_provider_strategy.value.capacity_provider
      base              = default_capacity_provider_strategy.value.base
      weight            = default_capacity_provider_strategy.value.weight
    }
  }

  # The custom providers must exist before the PUT names them.
  depends_on = [aws_ecs_capacity_provider.this]
}
