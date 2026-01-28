locals {
  # Convenience locals
  metadata        = var.metadata
  spec            = var.spec
  is_ssl_enabled  = try(var.spec.ssl.enabled, false)
  certificate_arn = try(var.spec.ssl.certificate_arn.value, null)

  # resource name and tags
  resource_name = coalesce(try(var.metadata.name, null), "aws-alb")
  tags          = merge({
    "Name" = local.resource_name
  }, try(var.metadata.labels, {}))

  # flatten foreign-key types to primitive lists
  subnet_ids = [for s in try(var.spec.subnets, []) : coalesce(try(s.value, null), try(s.value_from.name, null))]

  security_group_ids = [for sg in try(var.spec.security_groups, []) : coalesce(try(sg.value, null), try(sg.value_from.name, null))]

  # dns helpers
  create_dns_records = try(var.spec.dns.enabled, false) && length(try(var.spec.dns.hostnames, [])) > 0
  route53_zone_id    = try(var.spec.dns.route53_zone_id.value, null)
}


