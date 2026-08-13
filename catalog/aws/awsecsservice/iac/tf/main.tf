# AWS ECS Service Terraform Module
#
# The service only SCHEDULES: the task definition, cluster, target groups,
# subnets, and security groups are all referenced resources resolved before
# this module runs -- the module never creates or mutates a resource it
# merely references. Traffic reaches the service through first-class
# AwsLbTargetGroup / AwsLbListener / AwsLbListenerRule nodes; this module
# only registers task IPs into the referenced groups.

resource "aws_ecs_service" "this" {
  name            = local.service_name
  cluster         = var.spec.cluster_arn
  task_definition = var.spec.task_definition
  desired_count   = local.desired_count

  # Exactly one of these is set (spec CEL): a named launch type, or a
  # capacity-provider blend.
  launch_type = local.launch_type

  dynamic "capacity_provider_strategy" {
    for_each = var.spec.capacity_provider_strategy
    content {
      capacity_provider = capacity_provider_strategy.value.capacity_provider
      base              = capacity_provider_strategy.value.base
      weight            = capacity_provider_strategy.value.weight
    }
  }

  platform_version    = var.spec.platform_version != "" ? var.spec.platform_version : null
  scheduling_strategy = var.spec.scheduling_strategy != "" ? var.spec.scheduling_strategy : null

  # network is a pruned optional message: guard with a ternary (HCL's &&
  # does not short-circuit on null). Required for awsvpc task definitions --
  # every Fargate task and the modern EC2 posture.
  dynamic "network_configuration" {
    for_each = var.spec.network != null ? [var.spec.network] : []
    content {
      subnets          = network_configuration.value.subnets
      security_groups  = network_configuration.value.security_groups
      assign_public_ip = network_configuration.value.assign_public_ip
    }
  }

  # Load-balancer wiring registers task IPs into referenced target groups.
  # AWS requires each group to already be associated with a listener when
  # the service is created -- the graph's FK ordering guarantees it.
  dynamic "load_balancer" {
    for_each = var.spec.load_balancers
    content {
      target_group_arn = load_balancer.value.target_group_arn
      container_name   = load_balancer.value.container_name
      container_port   = load_balancer.value.container_port

      # The blue/green pair: ECS swaps the production listener rule between
      # the two target groups as deployments bake.
      dynamic "advanced_configuration" {
        for_each = load_balancer.value.advanced_configuration != null ? [load_balancer.value.advanced_configuration] : []
        content {
          alternate_target_group_arn = advanced_configuration.value.alternate_target_group_arn
          production_listener_rule   = advanced_configuration.value.production_listener_rule
          test_listener_rule         = advanced_configuration.value.test_listener_rule != "" ? advanced_configuration.value.test_listener_rule : null
          role_arn                   = advanced_configuration.value.role_arn
        }
      }
    }
  }

  health_check_grace_period_seconds = local.health_check_grace_period_seconds

  # Proto-optional bounds: null lets AWS default (200/100), an explicit
  # value -- including one that matches the default -- is sent as-is.
  deployment_maximum_percent         = var.spec.deployment_maximum_percent
  deployment_minimum_healthy_percent = var.spec.deployment_minimum_healthy_percent

  dynamic "deployment_circuit_breaker" {
    for_each = var.spec.deployment_circuit_breaker != null ? [var.spec.deployment_circuit_breaker] : []
    content {
      enable   = deployment_circuit_breaker.value.enable
      rollback = deployment_circuit_breaker.value.rollback
    }
  }

  # Alarm gating watches CloudWatch alarms BY NAME during deployments -- the
  # referenced AwsCloudwatchAlarm nodes publish their names as outputs
  # precisely for consumers like this.
  dynamic "alarms" {
    for_each = var.spec.alarms != null ? [var.spec.alarms] : []
    content {
      alarm_names = alarms.value.alarm_names
      enable      = alarms.value.enable
      rollback    = alarms.value.rollback
    }
  }

  dynamic "deployment_configuration" {
    for_each = var.spec.deployment_configuration != null ? [var.spec.deployment_configuration] : []
    content {
      strategy             = deployment_configuration.value.strategy != "" ? deployment_configuration.value.strategy : null
      bake_time_in_minutes = deployment_configuration.value.bake_time_in_minutes

      dynamic "canary_configuration" {
        for_each = deployment_configuration.value.canary_configuration != null ? [deployment_configuration.value.canary_configuration] : []
        content {
          canary_percent              = canary_configuration.value.canary_percent
          canary_bake_time_in_minutes = canary_configuration.value.canary_bake_time_in_minutes
        }
      }

      dynamic "linear_configuration" {
        for_each = deployment_configuration.value.linear_configuration != null ? [deployment_configuration.value.linear_configuration] : []
        content {
          step_percent              = linear_configuration.value.step_percent
          step_bake_time_in_minutes = linear_configuration.value.step_bake_time_in_minutes
        }
      }

      dynamic "lifecycle_hook" {
        for_each = deployment_configuration.value.lifecycle_hooks
        content {
          hook_target_arn  = lifecycle_hook.value.hook_target_arn
          role_arn         = lifecycle_hook.value.role_arn
          lifecycle_stages = lifecycle_hook.value.lifecycle_stages
          hook_details     = lifecycle_hook.value.hook_details != "" ? lifecycle_hook.value.hook_details : null
        }
      }
    }
  }

  dynamic "deployment_controller" {
    for_each = var.spec.deployment_controller != "" ? [var.spec.deployment_controller] : []
    content {
      type = deployment_controller.value
    }
  }

  dynamic "service_connect_configuration" {
    for_each = var.spec.service_connect != null ? [var.spec.service_connect] : []
    content {
      enabled   = service_connect_configuration.value.enabled
      namespace = service_connect_configuration.value.namespace != "" ? service_connect_configuration.value.namespace : null

      # Per-request access logs from the Service Connect proxy, emitted to
      # the proxy's log destination configured below.
      dynamic "access_log_configuration" {
        for_each = service_connect_configuration.value.access_log_configuration != null ? [service_connect_configuration.value.access_log_configuration] : []
        content {
          format                   = access_log_configuration.value.format
          include_query_parameters = access_log_configuration.value.include_query_parameters != "" ? access_log_configuration.value.include_query_parameters : null
        }
      }

      dynamic "log_configuration" {
        for_each = service_connect_configuration.value.log_configuration != null ? [service_connect_configuration.value.log_configuration] : []
        content {
          log_driver = log_configuration.value.log_driver
          options    = log_configuration.value.options

          # Secret options are name -> ARN pairs the ECS agent resolves at
          # task start; sorted iteration keeps the plan deterministic.
          dynamic "secret_option" {
            for_each = { for option_name in sort(keys(log_configuration.value.secret_options)) : option_name => log_configuration.value.secret_options[option_name] }
            content {
              name       = secret_option.key
              value_from = secret_option.value
            }
          }
        }
      }

      dynamic "service" {
        for_each = service_connect_configuration.value.services
        content {
          port_name             = service.value.port_name
          discovery_name        = service.value.discovery_name != "" ? service.value.discovery_name : null
          ingress_port_override = service.value.ingress_port_override

          dynamic "client_alias" {
            for_each = service.value.client_alias != null ? [service.value.client_alias] : []
            content {
              port     = client_alias.value.port
              dns_name = client_alias.value.dns_name != "" ? client_alias.value.dns_name : null

              # Requests carrying a matching header route to the TEST
              # revision during a blue/green deployment.
              dynamic "test_traffic_rules" {
                for_each = client_alias.value.test_traffic_rules
                content {
                  dynamic "header" {
                    for_each = test_traffic_rules.value.header != null ? [test_traffic_rules.value.header] : []
                    content {
                      name = header.value.name
                      value {
                        exact = header.value.value.exact
                      }
                    }
                  }
                }
              }
            }
          }

          dynamic "timeout" {
            for_each = service.value.timeout != null ? [service.value.timeout] : []
            content {
              idle_timeout_seconds        = timeout.value.idle_timeout_seconds > 0 ? timeout.value.idle_timeout_seconds : null
              per_request_timeout_seconds = timeout.value.per_request_timeout_seconds > 0 ? timeout.value.per_request_timeout_seconds : null
            }
          }

          dynamic "tls" {
            for_each = service.value.tls != null ? [service.value.tls] : []
            content {
              issuer_cert_authority {
                aws_pca_authority_arn = tls.value.aws_pca_authority_arn
              }
              kms_key  = tls.value.kms_key != "" ? tls.value.kms_key : null
              role_arn = tls.value.role_arn != "" ? tls.value.role_arn : null
            }
          }
        }
      }
    }
  }

  dynamic "service_registries" {
    for_each = var.spec.service_registries != null ? [var.spec.service_registries] : []
    content {
      registry_arn   = service_registries.value.registry_arn
      container_name = service_registries.value.container_name != "" ? service_registries.value.container_name : null
      container_port = service_registries.value.container_port
      port           = service_registries.value.port
    }
  }

  dynamic "volume_configuration" {
    for_each = var.spec.volume_configuration != null ? [var.spec.volume_configuration] : []
    content {
      name = volume_configuration.value.name

      managed_ebs_volume {
        role_arn         = volume_configuration.value.managed_ebs_volume.role_arn
        size_in_gb       = volume_configuration.value.managed_ebs_volume.size_in_gb > 0 ? volume_configuration.value.managed_ebs_volume.size_in_gb : null
        volume_type      = volume_configuration.value.managed_ebs_volume.volume_type != "" ? volume_configuration.value.managed_ebs_volume.volume_type : null
        iops             = volume_configuration.value.managed_ebs_volume.iops > 0 ? volume_configuration.value.managed_ebs_volume.iops : null
        throughput       = volume_configuration.value.managed_ebs_volume.throughput > 0 ? volume_configuration.value.managed_ebs_volume.throughput : null
        encrypted        = volume_configuration.value.managed_ebs_volume.encrypted
        kms_key_id       = volume_configuration.value.managed_ebs_volume.kms_key_id != "" ? volume_configuration.value.managed_ebs_volume.kms_key_id : null
        snapshot_id      = volume_configuration.value.managed_ebs_volume.snapshot_id != "" ? volume_configuration.value.managed_ebs_volume.snapshot_id : null
        file_system_type = volume_configuration.value.managed_ebs_volume.file_system_type != "" ? volume_configuration.value.managed_ebs_volume.file_system_type : null

        # Snapshot hydration speed; only meaningful with snapshot_id
        # (CEL/AWS both ignore it otherwise).
        volume_initialization_rate = volume_configuration.value.managed_ebs_volume.volume_initialization_rate > 0 ? volume_configuration.value.managed_ebs_volume.volume_initialization_rate : null

        # Creation-time tags on each per-task volume -- without these the
        # volumes carry no cost-allocation tags at all.
        dynamic "tag_specifications" {
          for_each = volume_configuration.value.managed_ebs_volume.tag_specifications
          content {
            resource_type  = tag_specifications.value.resource_type
            tags           = length(tag_specifications.value.tags) > 0 ? tag_specifications.value.tags : null
            propagate_tags = tag_specifications.value.propagate_tags != "" ? tag_specifications.value.propagate_tags : null
          }
        }
      }
    }
  }

  # VPC Lattice target-group attachments: ECS registers each task's named
  # port with the target group as tasks start and stop, assuming the
  # infrastructure role to do it.
  dynamic "vpc_lattice_configurations" {
    for_each = var.spec.vpc_lattice_configurations
    content {
      role_arn         = vpc_lattice_configurations.value.role_arn
      target_group_arn = vpc_lattice_configurations.value.target_group_arn
      port_name        = vpc_lattice_configurations.value.port_name
    }
  }

  dynamic "ordered_placement_strategy" {
    for_each = var.spec.ordered_placement_strategy
    content {
      type  = ordered_placement_strategy.value.type
      field = ordered_placement_strategy.value.field != "" ? ordered_placement_strategy.value.field : null
    }
  }

  dynamic "placement_constraints" {
    for_each = var.spec.placement_constraints
    content {
      type       = placement_constraints.value.type
      expression = placement_constraints.value.expression != "" ? placement_constraints.value.expression : null
    }
  }

  # Unset lets AWS decide (new services default to ENABLED where supported)
  # -- the spec deliberately has no default here because the provider
  # dropped its own.
  availability_zone_rebalancing = var.spec.availability_zone_rebalancing != "" ? var.spec.availability_zone_rebalancing : null

  propagate_tags          = var.spec.propagate_tags != "" ? var.spec.propagate_tags : null
  enable_ecs_managed_tags = var.spec.enable_ecs_managed_tags
  enable_execute_command  = var.spec.enable_execute_command
  force_delete            = var.spec.force_delete

  tags = local.aws_tags

  # desired_count is runtime state once the service is live: the autoscaler
  # owns it when configured, and operators may scale out of band. The module
  # seeds the initial count and then leaves it alone. ignore_changes must be
  # a static literal list (OpenTofu forbids a conditional expression here),
  # so this always ignores -- matching the house convention for
  # autoscaling-managed counts (the Pulumi module ignores the same path).
  lifecycle {
    ignore_changes = [desired_count]
  }
}

# ---------------------------------------------------------------------------
# Folded autoscaling: the scaler's identity IS this service (one scalable
# target per service), so it lives here rather than as its own kind. Each
# target-tracking policy creates its own AWS-managed CloudWatch alarms.
# ---------------------------------------------------------------------------

resource "aws_appautoscaling_target" "this" {
  count              = local.autoscaling_enabled ? 1 : 0
  max_capacity       = var.spec.autoscaling.max_tasks
  min_capacity       = var.spec.autoscaling.min_tasks
  resource_id        = "service/${local.cluster_name}/${aws_ecs_service.this.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "cpu" {
  # HCL's && does not short-circuit on null, so the autoscaling guard and
  # the per-policy guard must nest as ternaries.
  count              = local.autoscaling_enabled ? (var.spec.autoscaling.cpu != null ? 1 : 0) : 0
  name               = "${local.service_name}-cpu-scaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.this[0].resource_id
  scalable_dimension = aws_appautoscaling_target.this[0].scalable_dimension
  service_namespace  = aws_appautoscaling_target.this[0].service_namespace

  target_tracking_scaling_policy_configuration {
    target_value       = var.spec.autoscaling.cpu.target_percent
    scale_in_cooldown  = coalesce(var.spec.autoscaling.cpu.scale_in_cooldown_seconds, local.default_scale_in_cooldown)
    scale_out_cooldown = coalesce(var.spec.autoscaling.cpu.scale_out_cooldown_seconds, local.default_scale_out_cooldown)
    disable_scale_in   = var.spec.autoscaling.cpu.disable_scale_in

    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
  }
}

resource "aws_appautoscaling_policy" "memory" {
  count              = local.autoscaling_enabled ? (var.spec.autoscaling.memory != null ? 1 : 0) : 0
  name               = "${local.service_name}-memory-scaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.this[0].resource_id
  scalable_dimension = aws_appautoscaling_target.this[0].scalable_dimension
  service_namespace  = aws_appautoscaling_target.this[0].service_namespace

  target_tracking_scaling_policy_configuration {
    target_value       = var.spec.autoscaling.memory.target_percent
    scale_in_cooldown  = coalesce(var.spec.autoscaling.memory.scale_in_cooldown_seconds, local.default_scale_in_cooldown)
    scale_out_cooldown = coalesce(var.spec.autoscaling.memory.scale_out_cooldown_seconds, local.default_scale_out_cooldown)
    disable_scale_in   = var.spec.autoscaling.memory.disable_scale_in

    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageMemoryUtilization"
    }
  }
}

resource "aws_appautoscaling_policy" "requests" {
  count              = local.autoscaling_enabled ? (var.spec.autoscaling.requests_per_target != null ? 1 : 0) : 0
  name               = "${local.service_name}-requests-scaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.this[0].resource_id
  scalable_dimension = aws_appautoscaling_target.this[0].scalable_dimension
  service_namespace  = aws_appautoscaling_target.this[0].service_namespace

  target_tracking_scaling_policy_configuration {
    target_value       = var.spec.autoscaling.requests_per_target.target_requests_per_target
    scale_in_cooldown  = coalesce(var.spec.autoscaling.requests_per_target.scale_in_cooldown_seconds, local.default_scale_in_cooldown)
    scale_out_cooldown = coalesce(var.spec.autoscaling.requests_per_target.scale_out_cooldown_seconds, local.default_scale_out_cooldown)
    disable_scale_in   = var.spec.autoscaling.requests_per_target.disable_scale_in

    predefined_metric_specification {
      predefined_metric_type = "ALBRequestCountPerTarget"
      # ALBRequestCountPerTarget is scoped by "<lb-arn-suffix>/<tg-arn-suffix>"
      # -- both halves come from the referenced load balancer's and target
      # group's arn_suffix outputs.
      resource_label = "${var.spec.autoscaling.requests_per_target.load_balancer_arn_suffix}/${var.spec.autoscaling.requests_per_target.target_group_arn_suffix}"
    }
  }
}
