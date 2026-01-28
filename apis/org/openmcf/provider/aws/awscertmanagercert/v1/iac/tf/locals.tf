locals {
  # Safe handling of Route53 hosted zone ID with fallback
  safe_route53_zone_id = coalesce(
    try(var.spec.route53_hosted_zone_id.value, null),
    try(var.spec.route53_hosted_zone_id.value_from.name, null)
  )

  # Boolean for DNS validation method
  is_dns_validation = upper(try(var.spec.validation_method, "DNS")) == "DNS"

  # Boolean for email validation method
  is_email_validation = upper(try(var.spec.validation_method, "DNS")) == "EMAIL"

  # Check if alternate domain names are provided
  has_alternate_domains = length(try(var.spec.alternate_domain_names, [])) > 0
}


