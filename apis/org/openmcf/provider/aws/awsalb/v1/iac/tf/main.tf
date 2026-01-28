resource "aws_lb" "this" {
  name                       = local.resource_name
  load_balancer_type         = "application"
  security_groups            = local.security_group_ids
  subnets                    = local.subnet_ids
  internal                   = try(var.spec.internal, false)
  ip_address_type            = "ipv4"
  enable_deletion_protection = try(var.spec.delete_protection_enabled, false)
  idle_timeout               = try(var.spec.idle_timeout_seconds, 60)

  tags = local.tags
}

# HTTP listener that redirects to HTTPS when SSL is enabled
resource "aws_lb_listener" "http_redirect" {
  count = local.is_ssl_enabled ? 1 : 0

  load_balancer_arn = aws_lb.this.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

# HTTPS listener when SSL is enabled
resource "aws_lb_listener" "https" {
  count = local.is_ssl_enabled ? 1 : 0

  load_balancer_arn = aws_lb.this.arn
  port              = 443
  protocol          = "HTTPS"
  certificate_arn   = local.certificate_arn
  ssl_policy        = "ELBSecurityPolicy-2016-08"

  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "OK"
      status_code  = "200"
    }
  }
}

# Optional Route53 records for each hostname when DNS is enabled
resource "aws_route53_record" "this" {
  for_each = local.create_dns_records ? toset(var.spec.dns.hostnames) : []

  zone_id = local.route53_zone_id
  name    = each.value
  type    = "A"

  alias {
    name                   = aws_lb.this.dns_name
    zone_id                = aws_lb.this.zone_id
    evaluate_target_health = false
  }
}


