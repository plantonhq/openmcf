# ---------------------------------------------------------------------------
# Global Accelerator
#
# A global (non-regional) service: the provider transparently homes all API
# calls in us-west-2 regardless of the configured region. Every create/update
# on the accelerator, its listeners, or its endpoint groups is followed by the
# provider waiting for the accelerator to return to the DEPLOYED state, so
# applies take minutes. On destroy the provider first disables the accelerator
# (an AWS requirement), waits for that to deploy, then deletes.
# ---------------------------------------------------------------------------

resource "aws_globalaccelerator_accelerator" "this" {
  name = local.name

  # Presence-honest optional scalars pass through as null when unset in the
  # manifest, letting the provider defaults apply (enabled=true, IPV4).
  enabled         = var.spec.enabled
  ip_address_type = var.spec.ip_address_type

  # BYOIP static addresses (ForceNew). Empty means AWS-allocated anycast IPs.
  ip_addresses = length(var.spec.ip_addresses) > 0 ? var.spec.ip_addresses : null

  # The attributes block is ALWAYS materialized with an explicit
  # flow_logs_enabled value. Flow-log settings live on a separate accelerator-
  # attributes API, and the provider diff-suppresses a missing attributes
  # block — so omitting the block after flow logs were enabled would silently
  # leave AWS logging forever. Sending the explicit disabled state makes the
  # manifest the single source of truth.
  attributes {
    flow_logs_enabled   = var.spec.flow_logs != null ? var.spec.flow_logs.enabled : false
    flow_logs_s3_bucket = var.spec.flow_logs != null && var.spec.flow_logs.enabled ? var.spec.flow_logs.s3_bucket : null
    flow_logs_s3_prefix = var.spec.flow_logs != null && var.spec.flow_logs.enabled && var.spec.flow_logs.s3_prefix != "" ? var.spec.flow_logs.s3_prefix : null
  }

  tags = local.tags
}

# ---------------------------------------------------------------------------
# Listeners
#
# One provider resource per named listener. accelerator_arn is ForceNew on
# the listener, which is irrelevant here — the accelerator and its listeners
# live and die together in this module.
# ---------------------------------------------------------------------------

resource "aws_globalaccelerator_listener" "this" {
  for_each = local.listeners_map

  # The accelerator resource's id IS its ARN.
  accelerator_arn = aws_globalaccelerator_accelerator.this.id
  protocol        = each.value.protocol

  # null falls through to the provider default (NONE).
  client_affinity = each.value.client_affinity

  dynamic "port_range" {
    for_each = each.value.port_ranges
    content {
      from_port = port_range.value.from_port
      to_port   = port_range.value.to_port
    }
  }
}

# ---------------------------------------------------------------------------
# Endpoint Groups
#
# One provider resource per "listener_name/group_name" composite key.
# endpoint_group_region is ForceNew; everything else updates in place. The
# provider always sends the full endpoint and port-override sets on update,
# so removing the last endpoint from a group works.
# ---------------------------------------------------------------------------

resource "aws_globalaccelerator_endpoint_group" "this" {
  for_each = local.endpoint_groups_map

  # The listener resource's id IS its ARN.
  listener_arn = aws_globalaccelerator_listener.this[each.value.listener_name].id

  # Empty string means "inherit the provider region" (the spec's region).
  endpoint_group_region = each.value.group.endpoint_group_region != "" ? each.value.group.endpoint_group_region : null

  # Presence-honest health-check dials: null falls through to the provider/AWS
  # defaults (listener port, TCP, interval 30, threshold 3). health_check_path
  # is meaningful only for HTTP/HTTPS checks (CEL-enforced in the spec).
  health_check_port             = each.value.group.health_check_port
  health_check_protocol         = each.value.group.health_check_protocol
  health_check_path             = each.value.group.health_check_path != "" ? each.value.group.health_check_path : null
  health_check_interval_seconds = each.value.group.health_check_interval_seconds
  threshold_count               = each.value.group.threshold_count

  # null means 100% (the AWS default). An explicit 0 is a real value — it
  # drains the region while keeping its endpoints registered.
  traffic_dial_percentage = each.value.group.traffic_dial_percentage

  dynamic "endpoint_configuration" {
    for_each = each.value.group.endpoints
    content {
      endpoint_id = endpoint_configuration.value.endpoint_id

      # null lets AWS assign the default weight (128); an explicit 0 stops
      # routing to the endpoint without removing it.
      weight = endpoint_configuration.value.weight

      # Tri-state: null lets AWS apply its per-endpoint-type default; an
      # explicit value pins it. Only meaningful for ALB and EC2 endpoints.
      client_ip_preservation_enabled = endpoint_configuration.value.client_ip_preservation_enabled

      # Cross-account endpoints authorize through a Global Accelerator
      # cross-account attachment created in the endpoint-owning account.
      attachment_arn = endpoint_configuration.value.attachment_arn != "" ? endpoint_configuration.value.attachment_arn : null
    }
  }

  dynamic "port_override" {
    for_each = each.value.group.port_overrides
    content {
      listener_port = port_override.value.listener_port
      endpoint_port = port_override.value.endpoint_port
    }
  }
}
