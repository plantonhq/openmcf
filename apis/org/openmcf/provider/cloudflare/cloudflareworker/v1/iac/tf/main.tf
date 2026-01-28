# Fetch worker script bundle from R2
# Uses AWS S3 provider configured for R2
data "aws_s3_object" "worker_bundle" {
  provider = aws.r2
  bucket   = local.r2_bucket
  key      = local.r2_path
}

# Cloudflare Worker Script
resource "cloudflare_workers_script" "main" {
  account_id = var.spec.account_id
  name       = local.script_name

  # Worker content from R2 bundle
  content = data.aws_s3_object.worker_bundle.body

  # Module format (as opposed to service worker format)
  module = true

  # Compatibility settings
  compatibility_date  = local.compatibility_date
  compatibility_flags = ["nodejs_compat"]

  # Plain text bindings (environment variables)
  dynamic "plain_text_binding" {
    for_each = local.env_variables
    content {
      name = plain_text_binding.key
      text = plain_text_binding.value
    }
  }

  # KV namespace bindings
  dynamic "kv_namespace_binding" {
    for_each = local.kv_bindings
    content {
      name         = kv_namespace_binding.value.name
      namespace_id = kv_namespace_binding.value.field_path
    }
  }

  # Secret bindings
  # Note: Terraform doesn't support uploading secrets directly
  # Secrets must be managed separately via Cloudflare API or dashboard
  dynamic "secret_text_binding" {
    for_each = local.env_secrets
    content {
      name = secret_text_binding.key
      text = secret_text_binding.value
    }
  }

  # Enable observability (logs)
  logpush = true
}

# DNS Record for custom domain (if DNS is enabled)
resource "cloudflare_record" "worker_dns" {
  count = local.dns_enabled ? 1 : 0

  zone_id = local.dns_zone_id
  name    = local.dns_hostname
  type    = "AAAA"
  value   = "100::"  # Dummy IPv6 address for Workers routes
  proxied = true      # Orange cloud - required for Workers
}

# Worker Route (if DNS is enabled)
resource "cloudflare_worker_route" "main" {
  count = local.dns_enabled ? 1 : 0

  zone_id     = local.dns_zone_id
  pattern     = local.route_pattern
  script_name = cloudflare_workers_script.main.name

  # Ensure DNS record exists before creating route
  depends_on = [cloudflare_record.worker_dns]
}

