# DigitalOcean certificate.
#
# The spec's certificate_source oneof picks the branch, and the branch derives
# DigitalOcean's `type` argument: lets_encrypt (issued and auto-renewed by
# DigitalOcean; every domain must be managed by DigitalOcean DNS in this
# account) or custom (user-provided PEM material).
#
# Every argument is create-only, so any change replaces the certificate.
# create_before_destroy makes that replacement zero-downtime for consumers
# that reference the certificate by its stable name.
resource "digitalocean_certificate" "certificate" {
  name = var.spec.certificate_name
  type = local.cert_type

  # Let's Encrypt branch. domains conflicts with the PEM arguments, so it must
  # stay null on the custom branch.
  domains = local.is_lets_encrypt ? var.spec.lets_encrypt.domains : null

  # Custom branch. The provider stores only hashes of the PEM material
  # (the DigitalOcean API never returns it).
  leaf_certificate  = local.is_custom ? var.spec.custom.leaf_certificate : null
  private_key       = local.is_custom ? var.spec.custom.private_key : null
  certificate_chain = local.is_custom && var.spec.custom.certificate_chain != "" ? var.spec.custom.certificate_chain : null

  lifecycle {
    create_before_destroy = true
  }
}
