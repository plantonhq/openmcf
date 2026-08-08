# The cluster's capacity, materialized from the folded spec entries: one
# aws_ecs_capacity_provider per EC2 entry and per managed-instances entry
# (each keyed by name, so adding or removing an entry never disturbs its
# siblings) plus ONE association that PUTs the union of Fargate built-ins
# and folded provider names onto the cluster together with the default
# strategy.
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

# Managed-instances capacity providers: ECS launches and retires the EC2
# instances itself, so unlike the ASG-backed providers there is no group
# to wrap -- the provider is bound to its cluster at create (the AWS API
# requires the pairing) and carries the whole launch template inline.
# Creating one launches nothing; instances appear only when a service's
# strategy schedules tasks onto it.
resource "aws_ecs_capacity_provider" "managed_instances" {
  for_each = { for p in var.spec.managed_instances_capacity_providers : p.name => p }

  name    = each.value.name
  cluster = aws_ecs_cluster.this.name

  managed_instances_provider {
    infrastructure_role_arn = each.value.infrastructure_role_arn

    # -1 disables scale-in entirely; null keeps AWS's default
    # optimization -- distinct values, so the unset sentinel is null,
    # never a number.
    dynamic "infrastructure_optimization" {
      for_each = each.value.scale_in_after_seconds != null ? [each.value.scale_in_after_seconds] : []
      content {
        scale_in_after = infrastructure_optimization.value
      }
    }

    instance_launch_template {
      ec2_instance_profile_arn = each.value.instance_launch_template.ec2_instance_profile_arn

      # Changing the purchase model replaces the whole capacity provider
      # (ForceNew); the rest of the launch template updates in place.
      capacity_option_type = each.value.instance_launch_template.capacity_option_type != "" ? each.value.instance_launch_template.capacity_option_type : null

      dynamic "capacity_reservations" {
        for_each = each.value.instance_launch_template.capacity_reservations != null ? [each.value.instance_launch_template.capacity_reservations] : []
        content {
          reservation_preference = capacity_reservations.value.reservation_preference != "" ? capacity_reservations.value.reservation_preference : null
          reservation_group_arn  = capacity_reservations.value.reservation_group_arn != "" ? capacity_reservations.value.reservation_group_arn : null
        }
      }

      network_configuration {
        subnets         = each.value.instance_launch_template.network_configuration.subnets
        security_groups = length(each.value.instance_launch_template.network_configuration.security_groups) > 0 ? each.value.instance_launch_template.network_configuration.security_groups : null
      }

      dynamic "instance_requirements" {
        for_each = each.value.instance_launch_template.instance_requirements != null ? [each.value.instance_launch_template.instance_requirements] : []
        content {
          # The two required dimensions; max 0 means "no upper bound".
          memory_mib {
            min = instance_requirements.value.memory_mib.min
            max = instance_requirements.value.memory_mib.max != 0 ? instance_requirements.value.memory_mib.max : null
          }
          vcpu_count {
            min = instance_requirements.value.vcpu_count.min
            max = instance_requirements.value.vcpu_count.max != 0 ? instance_requirements.value.vcpu_count.max : null
          }

          allowed_instance_types  = length(instance_requirements.value.allowed_instance_types) > 0 ? instance_requirements.value.allowed_instance_types : null
          excluded_instance_types = length(instance_requirements.value.excluded_instance_types) > 0 ? instance_requirements.value.excluded_instance_types : null
          instance_generations    = length(instance_requirements.value.instance_generations) > 0 ? instance_requirements.value.instance_generations : null
          cpu_manufacturers       = length(instance_requirements.value.cpu_manufacturers) > 0 ? instance_requirements.value.cpu_manufacturers : null
          bare_metal              = instance_requirements.value.bare_metal != "" ? instance_requirements.value.bare_metal : null
          burstable_performance   = instance_requirements.value.burstable_performance != "" ? instance_requirements.value.burstable_performance : null

          require_hibernate_support = instance_requirements.value.require_hibernate_support ? true : null

          spot_max_price_percentage_over_lowest_price             = instance_requirements.value.spot_max_price_percentage_over_lowest_price != 0 ? instance_requirements.value.spot_max_price_percentage_over_lowest_price : null
          max_spot_price_as_percentage_of_optimal_on_demand_price = instance_requirements.value.max_spot_price_as_percentage_of_optimal_on_demand_price != 0 ? instance_requirements.value.max_spot_price_as_percentage_of_optimal_on_demand_price : null
          on_demand_max_price_percentage_over_lowest_price        = instance_requirements.value.on_demand_max_price_percentage_over_lowest_price != 0 ? instance_requirements.value.on_demand_max_price_percentage_over_lowest_price : null

          local_storage       = instance_requirements.value.local_storage != "" ? instance_requirements.value.local_storage : null
          local_storage_types = length(instance_requirements.value.local_storage_types) > 0 ? instance_requirements.value.local_storage_types : null

          dynamic "total_local_storage_gb" {
            for_each = instance_requirements.value.total_local_storage_gb != null ? [instance_requirements.value.total_local_storage_gb] : []
            content {
              min = total_local_storage_gb.value.min != 0 ? total_local_storage_gb.value.min : null
              max = total_local_storage_gb.value.max != 0 ? total_local_storage_gb.value.max : null
            }
          }
          dynamic "memory_gib_per_vcpu" {
            for_each = instance_requirements.value.memory_gib_per_vcpu != null ? [instance_requirements.value.memory_gib_per_vcpu] : []
            content {
              min = memory_gib_per_vcpu.value.min != 0 ? memory_gib_per_vcpu.value.min : null
              max = memory_gib_per_vcpu.value.max != 0 ? memory_gib_per_vcpu.value.max : null
            }
          }
          dynamic "network_interface_count" {
            for_each = instance_requirements.value.network_interface_count != null ? [instance_requirements.value.network_interface_count] : []
            content {
              min = network_interface_count.value.min != 0 ? network_interface_count.value.min : null
              max = network_interface_count.value.max != 0 ? network_interface_count.value.max : null
            }
          }
          dynamic "network_bandwidth_gbps" {
            for_each = instance_requirements.value.network_bandwidth_gbps != null ? [instance_requirements.value.network_bandwidth_gbps] : []
            content {
              min = network_bandwidth_gbps.value.min != 0 ? network_bandwidth_gbps.value.min : null
              max = network_bandwidth_gbps.value.max != 0 ? network_bandwidth_gbps.value.max : null
            }
          }
          dynamic "baseline_ebs_bandwidth_mbps" {
            for_each = instance_requirements.value.baseline_ebs_bandwidth_mbps != null ? [instance_requirements.value.baseline_ebs_bandwidth_mbps] : []
            content {
              min = baseline_ebs_bandwidth_mbps.value.min != 0 ? baseline_ebs_bandwidth_mbps.value.min : null
              max = baseline_ebs_bandwidth_mbps.value.max != 0 ? baseline_ebs_bandwidth_mbps.value.max : null
            }
          }
          dynamic "accelerator_count" {
            for_each = instance_requirements.value.accelerator_count != null ? [instance_requirements.value.accelerator_count] : []
            content {
              min = accelerator_count.value.min != 0 ? accelerator_count.value.min : null
              max = accelerator_count.value.max != 0 ? accelerator_count.value.max : null
            }
          }
          accelerator_manufacturers = length(instance_requirements.value.accelerator_manufacturers) > 0 ? instance_requirements.value.accelerator_manufacturers : null
          accelerator_names         = length(instance_requirements.value.accelerator_names) > 0 ? instance_requirements.value.accelerator_names : null
          accelerator_types         = length(instance_requirements.value.accelerator_types) > 0 ? instance_requirements.value.accelerator_types : null
          dynamic "accelerator_total_memory_mib" {
            for_each = instance_requirements.value.accelerator_total_memory_mib != null ? [instance_requirements.value.accelerator_total_memory_mib] : []
            content {
              min = accelerator_total_memory_mib.value.min != 0 ? accelerator_total_memory_mib.value.min : null
              max = accelerator_total_memory_mib.value.max != 0 ? accelerator_total_memory_mib.value.max : null
            }
          }
        }
      }

      dynamic "local_storage_configuration" {
        for_each = each.value.instance_launch_template.use_local_storage != null ? [each.value.instance_launch_template.use_local_storage] : []
        content {
          use_local_storage = local_storage_configuration.value
        }
      }

      monitoring = each.value.instance_launch_template.monitoring != "" ? each.value.instance_launch_template.monitoring : null

      dynamic "storage_configuration" {
        for_each = each.value.instance_launch_template.storage_size_gib != 0 ? [each.value.instance_launch_template.storage_size_gib] : []
        content {
          storage_size_gib = storage_configuration.value
        }
      }
    }

    propagate_tags = each.value.propagate_tags != "" ? each.value.propagate_tags : null
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
  depends_on = [
    aws_ecs_capacity_provider.this,
    aws_ecs_capacity_provider.managed_instances,
  ]
}
