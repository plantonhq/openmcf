# A target group is the composition point of ELBv2 load balancing: listeners
# and listener rules forward to it, ECS services deploy into it, and
# auto-scaling groups register instances with it. Name, port, protocol,
# protocol_version, vpc_id, target_type, and ip_address_type are create-only
# in AWS -- the provider replaces the group (and re-creates registrations)
# when they change; everything else updates in place.
resource "aws_lb_target_group" "this" {
  name = local.target_group_name

  # Port, protocol, and VPC apply to every target type except lambda -- a
  # Lambda function is invoked directly, never addressed over the network.
  port     = local.is_lambda ? null : var.spec.port
  protocol = local.is_lambda ? null : var.spec.protocol
  vpc_id   = local.is_lambda ? null : var.spec.vpc_id

  target_type      = var.spec.target_type != "" ? var.spec.target_type : null
  protocol_version = var.spec.protocol_version != "" ? var.spec.protocol_version : null
  ip_address_type  = var.spec.ip_address_type != "" ? var.spec.ip_address_type : null

  # ALB Target Optimizer: setting the agent port enables per-target readiness
  # routing. Create-only, like the group's other identity fields.
  target_control_port = var.spec.target_control_port > 0 ? var.spec.target_control_port : null

  # 0 means "keep the AWS default" (300s) -- the manifest zero value is not
  # distinguishable from unset, so immediate deregistration is expressed as 1.
  deregistration_delay = var.spec.deregistration_delay_seconds > 0 ? var.spec.deregistration_delay_seconds : null
  slow_start           = var.spec.slow_start_seconds > 0 ? var.spec.slow_start_seconds : null

  load_balancing_algorithm_type     = var.spec.load_balancing_algorithm_type != "" ? var.spec.load_balancing_algorithm_type : null
  load_balancing_anomaly_mitigation = var.spec.load_balancing_anomaly_mitigation != "" ? var.spec.load_balancing_anomaly_mitigation : null
  load_balancing_cross_zone_enabled = var.spec.load_balancing_cross_zone_enabled != "" ? var.spec.load_balancing_cross_zone_enabled : null

  # preserve_client_ip is a nullable tri-state at AWS (the default depends on
  # the target type): null keeps the AWS default, an explicit value overrides.
  preserve_client_ip                 = var.spec.preserve_client_ip == null ? null : tostring(var.spec.preserve_client_ip)
  proxy_protocol_v2                  = var.spec.proxy_protocol_v2 ? true : null
  connection_termination             = var.spec.connection_termination ? true : null
  lambda_multi_value_headers_enabled = var.spec.lambda_multi_value_headers_enabled ? true : null

  # Only explicitly set health-check fields are sent, so AWS keeps its
  # protocol-appropriate defaults for the rest.
  dynamic "health_check" {
    for_each = var.spec.health_check != null ? [var.spec.health_check] : []
    content {
      enabled             = health_check.value.enabled
      protocol            = health_check.value.protocol != "" ? health_check.value.protocol : null
      port                = health_check.value.port != "" ? health_check.value.port : null
      path                = health_check.value.path != "" ? health_check.value.path : null
      healthy_threshold   = health_check.value.healthy_threshold > 0 ? health_check.value.healthy_threshold : null
      unhealthy_threshold = health_check.value.unhealthy_threshold > 0 ? health_check.value.unhealthy_threshold : null
      interval            = health_check.value.interval_seconds > 0 ? health_check.value.interval_seconds : null
      timeout             = health_check.value.timeout_seconds > 0 ? health_check.value.timeout_seconds : null
      matcher             = health_check.value.matcher != "" ? health_check.value.matcher : null
    }
  }

  # Configuring stickiness implies enabling it unless the spec says otherwise
  # -- the same semantics AWS applies.
  dynamic "stickiness" {
    for_each = var.spec.stickiness != null ? [var.spec.stickiness] : []
    content {
      type            = stickiness.value.type
      enabled         = stickiness.value.enabled == null ? true : stickiness.value.enabled
      cookie_duration = stickiness.value.cookie_duration_seconds > 0 ? stickiness.value.cookie_duration_seconds : null
      cookie_name     = stickiness.value.cookie_name != "" ? stickiness.value.cookie_name : null
    }
  }

  # Group-level health policy. The DNS-failover count is a string at AWS (it
  # accepts "off"); the unhealthy-state-routing count is a plain integer -- an
  # AWS asymmetry the spec mirrors on purpose.
  dynamic "target_group_health" {
    for_each = var.spec.target_group_health != null ? [var.spec.target_group_health] : []
    content {
      dynamic "dns_failover" {
        for_each = target_group_health.value.dns_failover != null ? [target_group_health.value.dns_failover] : []
        content {
          minimum_healthy_targets_count      = dns_failover.value.minimum_healthy_targets_count != "" ? dns_failover.value.minimum_healthy_targets_count : null
          minimum_healthy_targets_percentage = dns_failover.value.minimum_healthy_targets_percentage != "" ? dns_failover.value.minimum_healthy_targets_percentage : null
        }
      }
      dynamic "unhealthy_state_routing" {
        for_each = target_group_health.value.unhealthy_state_routing != null ? [target_group_health.value.unhealthy_state_routing] : []
        content {
          minimum_healthy_targets_count      = unhealthy_state_routing.value.minimum_healthy_targets_count > 0 ? unhealthy_state_routing.value.minimum_healthy_targets_count : null
          minimum_healthy_targets_percentage = unhealthy_state_routing.value.minimum_healthy_targets_percentage != "" ? unhealthy_state_routing.value.minimum_healthy_targets_percentage : null
        }
      }
    }
  }

  # Established-connection handling while a target is unhealthy (NLB TCP/TLS).
  dynamic "target_health_state" {
    for_each = var.spec.target_health_state != null ? [var.spec.target_health_state] : []
    content {
      enable_unhealthy_connection_termination = target_health_state.value.enable_unhealthy_connection_termination
      unhealthy_draining_interval             = target_health_state.value.unhealthy_draining_interval_seconds
    }
  }

  tags = local.aws_tags

  lifecycle {
    # The spec cannot enforce VPC requiredness per target type (proto
    # validation cannot inspect reference fields), so the module is the
    # enforcement point: fail at plan time with a clear message instead of
    # letting AWS reject a half-created apply.
    precondition {
      condition     = local.is_lambda || var.spec.vpc_id != ""
      error_message = "vpc_id is required for every target_type except 'lambda'."
    }
  }
}

# Static registrations. Most architectures leave targets empty (ECS/ASG/EKS
# register their own); attachments are keyed by index because a target id may
# be a reference resolved only at apply time.
resource "aws_lb_target_group_attachment" "this" {
  count = length(var.spec.targets)

  target_group_arn  = aws_lb_target_group.this.arn
  target_id         = var.spec.targets[count.index].target_id
  port              = var.spec.targets[count.index].port > 0 ? var.spec.targets[count.index].port : null
  availability_zone = var.spec.targets[count.index].availability_zone != "" ? var.spec.targets[count.index].availability_zone : null
  quic_server_id    = var.spec.targets[count.index].quic_server_id != "" ? var.spec.targets[count.index].quic_server_id : null
}
