locals {
  # HTTPS and TLS listeners terminate TLS and therefore require a default
  # server certificate; every other protocol must not carry one.
  is_tls_protocol = contains(["HTTPS", "TLS"], var.spec.protocol)

  # The http_headers sub-objects, null-collapsed once here so main.tf reads
  # cleanly (each attribute below is consulted a single time).
  request_headers  = var.spec.http_headers != null ? var.spec.http_headers.request : null
  response_headers = var.spec.http_headers != null ? var.spec.http_headers.response : null

  # Resource-identity tags, matching the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsLbListener"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
