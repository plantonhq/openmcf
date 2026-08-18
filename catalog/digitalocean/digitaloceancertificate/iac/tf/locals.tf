locals {
  # The certificate_source branch fully determines DigitalOcean's `type`
  # argument -- there is no separate type field to keep consistent.
  is_lets_encrypt = var.spec.lets_encrypt != null
  is_custom       = var.spec.custom != null

  cert_type = local.is_lets_encrypt ? "lets_encrypt" : "custom"
}
