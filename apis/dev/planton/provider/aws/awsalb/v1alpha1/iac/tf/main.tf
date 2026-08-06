# The Application Load Balancer carries no routing configuration by design:
# listeners, listener rules, and target groups are separate resources that
# attach to it by ARN, so this module owns only what is truly
# load-balancer-wide -- placement, security groups, and the HTTP behavior
# attributes. Changing "internal" replaces the load balancer; the attributes
# update in place.
resource "aws_lb" "this" {
  name               = local.alb_name
  load_balancer_type = "application"
  subnets            = var.spec.subnets
  security_groups    = var.spec.security_groups
  internal           = var.spec.internal

  enable_deletion_protection = var.spec.delete_protection_enabled

  # Only explicitly set attributes are sent, so AWS keeps its own defaults
  # for the rest -- the module never bakes in opinions the spec does not
  # express.
  ip_address_type   = var.spec.ip_address_type != "" ? var.spec.ip_address_type : null
  idle_timeout      = var.spec.idle_timeout_seconds > 0 ? var.spec.idle_timeout_seconds : null
  client_keep_alive = var.spec.client_keep_alive_seconds > 0 ? var.spec.client_keep_alive_seconds : null

  # http2_enabled is a tri-state: null keeps the AWS default (true); an
  # explicit false downgrades clients to HTTP/1.1.
  enable_http2                                = var.spec.http2_enabled
  enable_waf_fail_open                        = var.spec.waf_fail_open_enabled ? true : null
  enable_zonal_shift                          = var.spec.zonal_shift_enabled ? true : null
  drop_invalid_header_fields                  = var.spec.drop_invalid_header_fields ? true : null
  preserve_host_header                        = var.spec.preserve_host_header ? true : null
  enable_xff_client_port                      = var.spec.xff_client_port_enabled ? true : null
  xff_header_processing_mode                  = var.spec.xff_header_processing_mode != "" ? var.spec.xff_header_processing_mode : null
  desync_mitigation_mode                      = var.spec.desync_mitigation_mode != "" ? var.spec.desync_mitigation_mode : null
  enable_tls_version_and_cipher_suite_headers = var.spec.tls_version_and_cipher_suite_headers_enabled ? true : null

  # The three S3 log streams share one shape; "enabled" is implied by the
  # block's presence in the spec (a bucket with logging off is meaningless).
  # The bucket must carry the ELB log-delivery bucket policy or delivery
  # fails silently.
  dynamic "access_logs" {
    for_each = var.spec.access_logs != null ? [var.spec.access_logs] : []
    content {
      bucket  = access_logs.value.bucket
      prefix  = access_logs.value.prefix != "" ? access_logs.value.prefix : null
      enabled = true
    }
  }

  dynamic "connection_logs" {
    for_each = var.spec.connection_logs != null ? [var.spec.connection_logs] : []
    content {
      bucket  = connection_logs.value.bucket
      prefix  = connection_logs.value.prefix != "" ? connection_logs.value.prefix : null
      enabled = true
    }
  }

  dynamic "health_check_logs" {
    for_each = var.spec.health_check_logs != null ? [var.spec.health_check_logs] : []
    content {
      bucket  = health_check_logs.value.bucket
      prefix  = health_check_logs.value.prefix != "" ? health_check_logs.value.prefix : null
      enabled = true
    }
  }

  tags = local.aws_tags
}

# Optional WAFv2 protection: an ALB has at most one web ACL, and the binding
# is a setting of the PROTECTED resource, so the association lives here (the
# web ACL itself never knows its consumers). Both fields are ForceNew on the
# association -- re-pointing the ALB at a different ACL replaces only this
# glue resource, never the load balancer.
resource "aws_wafv2_web_acl_association" "this" {
  count = var.spec.web_acl_arn != "" ? 1 : 0

  resource_arn = aws_lb.this.arn
  web_acl_arn  = var.spec.web_acl_arn
}

# Optional Route53 records for each hostname when DNS is enabled.
# allow_overwrite adopts an existing alias record (e.g. left by a prior partial apply,
# or one already pointing at this ALB) instead of failing the apply on a CREATE
# collision -- this alias record is owned by the ALB module.
resource "aws_route53_record" "this" {
  # toset([]) rather than [] in the false arm: both arms of a conditional
  # must have one type, and a bare [] is a tuple -- for_each rejects it.
  for_each = local.create_dns_records ? toset(var.spec.dns.hostnames) : toset([])

  allow_overwrite = true
  zone_id         = var.spec.dns.route53_zone_id
  name            = each.value
  type            = "A"

  alias {
    name    = aws_lb.this.dns_name
    zone_id = aws_lb.this.zone_id
    # false on purpose: target-health evaluation only changes behavior under
    # failover/weighted routing policies, and a simple alias should not pay
    # for health evaluation. Must stay identical in the Pulumi module
    # (cross-engine parity).
    evaluate_target_health = false
  }
}
