locals {
  # AWS limits load balancer names to 32 characters; truncate deterministically
  # so the same manifest always yields the same name.
  nlb_name = substr(var.metadata.name, 0, 32)

  # DNS records are created only when DNS is enabled AND hostnames exist --
  # an enabled block with no hostnames is a no-op, not an error.
  create_dns_records = var.spec.dns != null && try(var.spec.dns.enabled, false) && length(try(var.spec.dns.hostnames, [])) > 0

  # Resource-identity tags, matching the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = local.nlb_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsNlb"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
