# The Front Door custom domain -- your own hostname served by Front
# Door. Creation returns with the domain in a PENDING-VALIDATION state
# (deployment does not block on DNS proof): publish the
# validation_token output as a TXT record at _dnsauth.<host_name> and
# Azure flips the domain to approved; then CNAME the hostname to the
# endpoint's host_name so traffic arrives. Routes attach the domain via
# their custom_domain_ids -- the route side owns the attachment.
#
# The minimum TLS version is NOT an input: Azure retired TLS 1.0/1.1
# (March 2025) and the provider accepts exactly one value, so the
# constant TLS12 is sent unconditionally (minimum_version, never the
# deprecated minimum_tls_version alias). No Azure tags: ARM does not
# support tags on custom domains.
resource "azurerm_cdn_frontdoor_custom_domain" "main" {
  name                     = var.spec.domain_name
  cdn_frontdoor_profile_id = var.spec.profile_id
  host_name                = var.spec.host_name
  dns_zone_id              = var.spec.dns_zone_id

  tls {
    minimum_version = "TLS12"

    # certificate_type absent means ManagedCertificate -- Azure issues,
    # hosts, and auto-rotates a free DV certificate. With
    # CustomerCertificate the referenced Front Door secret carries the
    # key material (pairing spec-enforced).
    certificate_type        = var.spec.tls.certificate_type != null ? local.certificate_type_map[var.spec.tls.certificate_type] : "ManagedCertificate"
    cdn_frontdoor_secret_id = var.spec.tls.secret_id

    # A cipher-suite policy is sent only when configured -- absence
    # serves Azure's default suite set. With CUSTOMIZED, tls13 is sent
    # only when the user pinned it (empty means Azure's TLS 1.3
    # defaults; when set, the spec guarantees both mandatory suites).
    dynamic "cipher_suite" {
      for_each = var.spec.tls.cipher_suite != null ? [var.spec.tls.cipher_suite] : []
      content {
        type = local.cipher_suite_set_type_map[cipher_suite.value.type]

        dynamic "custom_ciphers" {
          for_each = cipher_suite.value.custom_ciphers != null ? [cipher_suite.value.custom_ciphers] : []
          content {
            tls12 = custom_ciphers.value.tls12
            tls13 = custom_ciphers.value.tls13 != null && length(custom_ciphers.value.tls13) > 0 ? custom_ciphers.value.tls13 : null
          }
        }
      }
    }
  }
}
